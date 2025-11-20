CREATE TABLE IF NOT EXISTS orders (
                                      id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
                                      amount INT NOT NULL,
                                      status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE=InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE = utf8mb4_unicode_ci;
