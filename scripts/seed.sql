-- Seed data for Go Inventory System
-- Sample inventory items for development and testing

INSERT INTO items (name, sku, description, quantity, price, category, created_at, updated_at)
VALUES
    ('Laptop - Dell XPS 15', 'LAPTOP-XPS15-001', 'High-performance laptop with 16GB RAM and 512GB SSD', 25, 1299.99, 'Electronics', NOW(), NOW()),
    ('Wireless Mouse - Logitech MX Master 3', 'MOUSE-MX3-001', 'Ergonomic wireless mouse with customizable buttons', 150, 99.99, 'Accessories', NOW(), NOW()),
    ('Mechanical Keyboard - Keychron K2', 'KEYBOARD-K2-001', 'Wireless mechanical keyboard with RGB backlight', 75, 89.99, 'Accessories', NOW(), NOW()),
    ('Monitor - LG 27" 4K UHD', 'MONITOR-LG27-001', '27-inch 4K UHD monitor with HDR support', 40, 449.99, 'Electronics', NOW(), NOW()),
    ('USB-C Hub - Anker 7-in-1', 'HUB-ANKER7-001', '7-in-1 USB-C hub with HDMI, USB 3.0, and SD card reader', 200, 49.99, 'Accessories', NOW(), NOW()),
    ('Webcam - Logitech C920', 'WEBCAM-C920-001', 'Full HD 1080p webcam with auto-focus', 60, 79.99, 'Electronics', NOW(), NOW()),
    ('Headphones - Sony WH-1000XM4', 'HEADPHONE-SONY-001', 'Wireless noise-cancelling headphones', 45, 349.99, 'Audio', NOW(), NOW()),
    ('External SSD - Samsung T7 1TB', 'SSD-T7-1TB-001', 'Portable external SSD with USB 3.2 Gen 2', 100, 159.99, 'Storage', NOW(), NOW()),
    ('Docking Station - CalDigit TS3 Plus', 'DOCK-TS3-001', 'Thunderbolt 3 docking station with 15 ports', 30, 299.99, 'Accessories', NOW(), NOW()),
    ('Standing Desk - FlexiSpot E7', 'DESK-E7-001', 'Electric height-adjustable standing desk', 15, 599.99, 'Furniture', NOW(), NOW())
ON CONFLICT (sku) DO NOTHING;
