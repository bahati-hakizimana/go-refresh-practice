
CREATE TYPE booking_status_enum AS ENUM (
    'pending',
    'confirmed',
    'checked_in',
    'cancelled'
    
);

CREATE TABLE IF NOT EXISTS bookings (
    id SERIAL PRIMARY KEY,
    apartment_id INT NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL,
    checkin_date DATE NOT NULL,
    checkout_date DATE NOT NULL,
    guest_number INT NOT NULL,
    total_price NUMERIC(12,2) NOT NULL,
    booking_amount NUMERIC(12,2) NOT NULL,
    balance_amount NUMERIC(12,2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    booking_status booking_status_enum DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_apartment FOREIGN KEY (apartment_id)
        REFERENCES apartments(id)
        ON DELETE CASCADE
);



