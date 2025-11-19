CREATE TABLE IF NOT EXISTS apartment_images (
    id SERIAL PRIMARY KEY,
    apartment_id INT NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    caption VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_apartment
        FOREIGN KEY(apartment_id)
        REFERENCES apartments(id)
        ON DELETE CASCADE
);


