CREATE TABLE IF NOT EXISTS trailers (
    id bigserial primary key,
    trailer_name text not null,
    duration int not null,
    premier_date text not null
)