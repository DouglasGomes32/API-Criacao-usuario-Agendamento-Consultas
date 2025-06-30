-- docker/init.sql

CREATE TABLE IF NOT EXISTS pacientes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS consultas (
    id SERIAL PRIMARY KEY,
    patient_id INTEGER NOT NULL REFERENCES pacientes(id) ON DELETE CASCADE,
    datetime TIMESTAMP NOT NULL
);
