CREATE TABLE IF NOT EXISTS apartments (
    id SERIAL PRIMARY KEY,
    code VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    rooms INT NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'available' CHECK (status IN ('available', 'booked')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

