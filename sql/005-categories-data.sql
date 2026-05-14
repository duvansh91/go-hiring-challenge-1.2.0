-- Insert categories
INSERT INTO categories (code, name) VALUES
('CLOTHING', 'Clothing'),
('SHOES', 'Shoes'),
('ACCESSORIES', 'Accessories');

-- Add category_id column to products if it doesn't exist
-- Update products with category associations
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'CLOTHING') WHERE code IN ('PROD001', 'PROD004', 'PROD007');
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'SHOES') WHERE code IN ('PROD002', 'PROD006');
UPDATE products SET category_id = (SELECT id FROM categories WHERE code = 'ACCESSORIES') WHERE code IN ('PROD003', 'PROD005', 'PROD008');
