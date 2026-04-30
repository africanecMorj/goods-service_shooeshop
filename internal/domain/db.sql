CREATE TABLE users (
	id SERIAL PRIMARY KEY,
	email TEXT UNIQUE,
	password TEXT
);

CREATE TABLE refresh_tokens (
	id SERIAL PRIMARY KEY,
	user_id INT,
	token TEXT,
	expires_at TIMESTAMP
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price NUMERIC(10,2),
    image_path TEXT
);