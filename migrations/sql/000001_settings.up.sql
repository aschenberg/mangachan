CREATE TABLE IF NOT EXISTS settings (
    setting_id INT PRIMARY KEY,
    first_voucher BOOLEAN NOT NULL DEFAULT TRUE
);
