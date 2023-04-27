package main

import (
	"bytes"
	"encoding/json"
	"greenlight.aslan/internal/data"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// #1
func TestCreateMovieHandler(t *testing.T) {

	app := &application{}

	db, err := dbConnection()
	if err != nil {
		t.Errorf("Database is not working correctly: %s", err)
	}

	app.models.Movies.DB = db

	movie := struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}{
		"Harry Potter and the Philosopher's Stone",
		2004,
		154,
		[]string{"drama", "fantasy"},
	}

	jsMovie, err := json.Marshal(movie)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/movies", bytes.NewBuffer(jsMovie))
	if err != nil {
		t.Errorf("Error in request creation: %s", err)
	}

	rr := httptest.NewRecorder()

	app.routes().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("CreateMovie handler returned wrong status code: got - \"%v\", expected - \"%v\"", status, http.StatusOK)
	}
}

// #2
func TestUpdateMovieHandler(t *testing.T) {

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

	movieJs, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PATCH", "/v1/movies/2", bytes.NewBuffer(movieJs))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	app.routes().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("UpdateMovie handler returned wrong status code: got - \"%v\", expected - \"%v\"", status, http.StatusOK)
	}
}

// #3
func TestDeleteMovieHandler(t *testing.T) {

	app := &application{}

	db, err := dbConnection()
	if err != nil {
		t.Errorf("Database is not working correctly: %s", err)
	}

	app.models.Movies.DB = db

	req, err := http.NewRequest("DELETE", "/v1/movies/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	app.routes().ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DeleteMovie handler returned wrong status code: got - \"%v\", expected - \"%v\"", status, http.StatusOK)
	}
}

// #4
func TestShowMovieHandler(t *testing.T) {

	app := &application{}

	db, err := dbConnection()
	if err != nil {
		t.Errorf("Database is not working correctly: %s", err)
	}

	app.models.Movies.DB = db

	type Movie struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	expected := struct {
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

	req, err := http.NewRequest("GET", "/v1/movies/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	app.routes().ServeHTTP(rr, req)

	m := Movie{}

	last := len(rr.Body.String()) - 2
	s := rr.Body.String()[9:last]

	err = json.Unmarshal([]byte(s), &m)
	if err != nil {
		t.Fatal(err)
	}

	if m.Title != expected.Title {
		t.Errorf("Title has an error insertion within handler: got - %s, expected - %s", m.Title, expected.Title)
	}
	if m.Year != expected.Year {
		t.Errorf("Year has an error insertion within handler: got - %d, expected - %d", m.Year, expected.Year)
	}
	if m.Runtime != expected.Runtime {
		t.Errorf("Runtime has an error insertion within handler: got - %d, expected - %d", m.Runtime, expected.Runtime)
	}
	gate := true
	for i := 0; i < len(m.Genres); i++ {
		if m.Genres[i] != expected.Genres[i] {
			gate = false
		}
	}
	if !gate {
		t.Errorf("Genres has an error insertion within handler: got - %v, expected - %v", m.Genres, expected.Genres)
	}

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("DeleteMovie handler returned wrong status code: got - \"%v\", expected - \"%v\"", status, http.StatusOK)
	}
}

// #5
func TestCreateAuthenticationTokenHandle(t *testing.T) {

	app := &application{}

	db, err := dbConnection()
	if err != nil {
		t.Errorf("Database is not working correctly: %s", err)
	}

	app.models.Tokens.DB = db
	app.models.Users.DB = db

	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		"nissanfordo03@gmail.com",
		"12345678",
	}

	type Token struct {
		Token  string `json:"token"`
		Expiry string `json:"expiry"`
	}

	jsInput, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/tokens/authentication", bytes.NewBuffer(jsInput))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	app.routes().ServeHTTP(rr, req)

	tokenObj := Token{}

	last := len(rr.Body.String()) - 2
	s := rr.Body.String()[24:last]

	err = json.Unmarshal([]byte(s), &tokenObj)
	if err != nil {
		t.Fatal(err)
	}

	tt, err := time.Parse(time.RFC3339, tokenObj.Expiry)
	if err != nil {
		t.Fatal(err)
	}

	expected := time.Now().AddDate(0, 0, 2)

	if tt.Before(time.Now()) {
		t.Errorf("Token creation is working with expired date: got - %v, expected - %v", tt.UTC().Format("2006-01-02"), expected.UTC().Format("2006-01-02"))
	}

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("CreateToken handler returned wrong status code: got - \"%v\", expected - \"%v\"", status, http.StatusOK)
	}
}

// #6
func TestRegisterUserHandler(t *testing.T) {

	app := &application{}

	db, err := dbConnection()
	if err != nil {
		t.Errorf("Database is not working correctly: %s", err)
	}

	app.models.Users.DB = db
	app.models.Tokens.DB = db

	input := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}{
		"Kusainov Aslan",
		"nissanfordo03@gmail.com",
		"12345678",
		"User",
	}

	userObj, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/v1/users", bytes.NewBuffer(userObj))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	app.routes().ServeHTTP(rr, req)

	last := len(rr.Body.String()) - 2
	body := rr.Body.String()[8:last]

	userInput := struct {
		Id                int32  `json:"id"`
		Activated         bool   `json:"activated"`
		Role              string `json:"role"`
		expectedId        int32
		expectedActivated bool
		expectedRole      string
	}{
		expectedId:        1,
		expectedActivated: false,
		expectedRole:      "User",
	}

	err = json.Unmarshal([]byte(body), &userInput)
	if err != nil {
		t.Fatal(err)
	}

	if userInput.Id != userInput.expectedId {
		t.Errorf("ID field does not set right: got - %d, expected - %d", userInput.Id, userInput.expectedId)
	}
	if userInput.Activated != userInput.expectedActivated {
		t.Errorf("Activated field does not set right: got - %v, expected - %v", userInput.Activated, userInput.Activated)
	}
	if userInput.Role != userInput.Role {
		t.Errorf("Role field does not set right: got - %s, expected - %s", userInput.Role, userInput.expectedRole)
	}

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("User registration handler returned wrong status code: got - \"%v\", expected - \"%v\"", status, http.StatusOK)
	}
}

// #7
func TestListMoviesHandler(t *testing.T) {

	app := &application{}

	db, err := dbConnection()
	if err != nil {
		t.Errorf("Database is not working correctly: %s", err)
	}

	app.models.Movies.DB = db

	md := struct {
		CurrentPage  int `json:"current_page"`
		PageSize     int `json:"page_size"`
		FirstPage    int `json:"first_page"`
		LastPage     int `json:"last_page"`
		TotalRecords int `json:"total_records"`
	}{}

	req, err := http.NewRequest("GET", "/v1/movies", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	app.routes().ServeHTTP(rr, req)

	last := strings.Index(rr.Body.String(), "movies") - 2
	mdBody := rr.Body.String()[12:last]

	err = json.Unmarshal([]byte(mdBody), &md)
	if err != nil {
		t.Fatal(err)
	}

	expectedTotalRecords := 0

	rows, err := db.Query("select count(*) from movies")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&expectedTotalRecords)
		if err != nil {
			t.Fatal(err)
		}
	}

	if md.TotalRecords != expectedTotalRecords {
		t.Errorf("The number of users in METADATA is not equal to DB: got - %d, expected - %d", md.TotalRecords, expectedTotalRecords)
	}
}
