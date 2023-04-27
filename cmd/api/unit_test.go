package main

import (
	"database/sql"
	"flag"
	"greenlight.aslan/internal/data"
	"greenlight.aslan/internal/validator"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func getConfigDB() config {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db_dsn", os.Getenv("GREENLIGHT_DB_DSN"), "PostgresSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Parse()

	return cfg
}

func dbConnection() (*sql.DB, error) {
	db, err := openDB(getConfigDB())

	return db, err
}

// #1
func TestHealthcheckHandler(t *testing.T) {

	app := &application{}

	req, err := http.NewRequest("GET", "/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.healthcheckHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("HealthCheck handler returned wrong status code: got - \"%v\", expected - \"%v\"", status, http.StatusOK)
	}

	expected := "available"

	js := string(rr.Body.Bytes())[11:20]

	if expected != js {
		t.Errorf("HealthCheck handler has a different state structure: got - \"%s\", expected - \"%s\"", js, expected)
	}
}

// #2
func TestInsertTrailer(t *testing.T) {
	app := &application{}

	db, err := dbConnection()
	if err != nil {
		t.Errorf("Database is not working correctly: %s", err)
	}

	app.models.Trailers.DB = db

	tests := struct {
		name   string
		input1 string
		input2 int64
		input3 string
	}{
		"Insert method",
		"test1",
		42,
		"test1",
	}

	trailer := &data.Trailer{
		Trailer_name: tests.input1,
		Duration:     tests.input2,
		Premier_date: tests.input3,
	}

	err = app.models.Trailers.Insert(trailer)
	if err != nil {
		t.Errorf("Insert method has an error: %s", err)
	}
}

// #3
func TestDeleteMovie(t *testing.T) {
	app := &application{}

	db, err := dbConnection()
	if err != nil {
		t.Errorf("Database is not working correctly: %s", err)
	}

	app.models.Movies.DB = db

	err = app.models.Movies.Delete(5)
	if err != nil {
		t.Errorf("Delete method is not working: \"%s\"", err)
	}
}

// #4
func TestGetMovie(t *testing.T) {
	app := &application{}

	db, err := dbConnection()
	if err != nil {
		t.Errorf("Database is not working correctly: %s", err)
	}

	app.models.Movies.DB = db

	gate := false

	m := struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}{
		"Harry Potter and the Philosopher's Stone",
		2002,
		152,
		[]string{"drama", "fantasy"},
	}

	movie, err := app.models.Movies.Get(4)
	if err != nil {
		t.Errorf("Get method is not working: \"%s\"", err)
	}

	for i := 0; i < len(m.Genres); i++ {
		if movie.Genres[i] != m.Genres[i] {
			gate = true
		}
	}

	if movie.Title != m.Title {
		t.Errorf("Movie title: got \"%s\" - expected \"%s\"", movie.Title, m.Title)
	} else if movie.Year != m.Year {
		t.Errorf("Movie year: got \"%d\" - expected \"%d\"", movie.Year, m.Year)
	} else if movie.Runtime != m.Runtime {
		t.Errorf("Movie runtime: got \"%d\" - expected \"%d\"", movie.Runtime, m.Runtime)
	} else if gate {
		t.Errorf("Movie genres: got \"%v\" - expected \"%v\"", movie.Genres, m.Genres)
	}
}

// #5
func TestUpdateMovie(t *testing.T) {
	app := &application{}

	db, err := dbConnection()
	if err != nil {
		t.Errorf("Database is not working correctly: %s", err)
	}

	app.models.Movies.DB = db

	m := struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}{
		"Harry Potter and the Philosopher's Stone",
		2002,
		152,
		[]string{"drama", "fantasy"},
	}

	oldMovie, err := app.models.Movies.Get(4)
	if err != nil {
		t.Fatal(err)
	}

	movie := &data.Movie{
		ID:      oldMovie.ID,
		Title:   m.Title,
		Year:    m.Year,
		Runtime: m.Runtime,
		Genres:  m.Genres,
		Version: oldMovie.Version,
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		t.Errorf("Movie update method is not working well: \"%s\"", err)
	}
}

// #6
func TestCalculateMetadata(t *testing.T) {

	tests := []struct {
		name     string
		input1   int
		input2   int
		expected int
	}{
		{"Calculating LastPage size: Ceil(125/10)=13", 125, 10, 13},
		{"Calculating LastPage size: Ceil(47/13)=4", 47, 13, 4},
		{"Calculating LastPage size: Ceil(53/40)=2", 53, 40, 2},
		{"Calculating LastPage size: Ceil(123/42)=3", 123, 42, 3},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			md := data.CalculateMetadata(tst.input1, 1, tst.input2)
			if md.LastPage != tst.expected {
				t.Errorf("Expected \"%d\" - got \"%d\"", md.LastPage, tst.expected)
			}
		})
	}
}

// #7
func TestValidateMatches(t *testing.T) {

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Validation Match method:", "nissanfordo@gmail.com", true},
		{"Validation Match method:", "alex.family.2202@gmail.com", true},
		{"Validation Match method:", "okda@mail.ru", true},
		{"Validation Match method:", "nissanfordo03@gmail.", false},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			gate := validator.Matches(tst.input, validator.EmailRX)
			if gate != tst.expected {
				g := strconv.FormatBool(gate)
				te := strconv.FormatBool(tst.expected)
				t.Errorf("Validation Mathes is not working well: got \"%s\" - expected \"%s\"", g, te)
			}
		})
	}
}

// #8
func TestValidationPermittedValue(t *testing.T) {

	tests := []struct {
		name     string
		input1   string
		input2   []string
		expected bool
	}{
		{"Validation PermittedValue method:", "test1", []string{"test1", "test2", "test3"}, true},
		{"Validation PermittedValue method:", "test2", []string{"test1", "test2", "test3"}, true},
		{"Validation PermittedValue method:", "test3", []string{"test1", "test2", "test4"}, false},
		{"Validation PermittedValue method:", "test4", []string{"test2", "test4", "test1"}, true},
	}

	for _, tst := range tests {
		t.Run(tst.name, func(t *testing.T) {
			gate := validator.PermittedValue(tst.input1, tst.input2...)
			if gate != tst.expected {
				g := strconv.FormatBool(gate)
				te := strconv.FormatBool(tst.expected)
				t.Errorf("Validation PermittedValue is not working well: got \"%s\" - expected \"%s\"", g, te)
			}
		})
	}
}

// #9
func TestGetByEmail(t *testing.T) {
	app := &application{}

	db, err := dbConnection()
	if err != nil {
		t.Errorf("Database is not working correctly: %s", err)
	}

	app.models.Users.DB = db

	u := struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		Activated bool   `json:"activated"`
		Role      string `json:"role"`
	}{
		"Aslan Kusainov",
		"nissanfordo03@gmail.com",
		false,
		"User",
	}

	user, err := app.models.Users.GetByEmail(u.Email)
	if err != nil {
		t.Errorf("GetByEmail method is not working: %s", err)
	}

	if user.Name != u.Name {
		t.Errorf("User name: got \"%s\" - expected \"%s\"", user.Name, u.Name)
	}
	if user.Email != u.Email {
		t.Errorf("User email: got \"%s\" - expected \"%s\"", user.Email, u.Email)
	}
	if user.Activated != u.Activated {
		uA := strconv.FormatBool(user.Activated)
		uAc := strconv.FormatBool(u.Activated)
		t.Errorf("User activation: got \"%s\" - expected \"%s\"", uA, uAc)
	}
	if user.Role != u.Role {
		t.Errorf("User role: got \"%s\" - expected \"%s\"", user.Role, u.Role)
	}
}
