CREATE TABLE IF NOT EXISTS payments (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `booking_id` INT UNSIGNED NOT NULL,
    `payment_type` ENUM('deposit','balance','full','refund') NOT NULL,
    `amount` DECIMAL(12, 2) NOT NULL,
    `currency` VARCHAR(10) DEFAULT 'USD',
    `payment_status` ENUM('pending','completed','failed','refunded') DEFAULT 'pending',
    `payment_method` ENUM('card', 'paypack', 'paypal', 'bank_transfer', 'cash') DEFAULT 'card',
    `transaction_reference` VARCHAR(255) DEFAULT NULL,
    `paid_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`booking_id`) REFERENCES bookings(`id`) ON DELETE CASCADE
)ENGINE=InnoDB;
