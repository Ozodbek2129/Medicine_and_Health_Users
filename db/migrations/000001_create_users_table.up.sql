CREATE TYPE IF NOT EXISTS gender AS ENUM ('male', 'female');

CREATE TYPE IF NOT EXISTS role AS ENUM ('patient', 'doctor', 'admin');

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(35) NOT NULL,
    last_name VARCHAR(55) NOT NULL,
    date_of_birth DATE NOT NULL,
    gender gender NOT NULL,
    role role NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);