-- Add category_id column to products table
ALTER TABLE products ADD COLUMN IF NOT EXISTS category_id INTEGER;

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    code VARCHAR(32) UNIQUE NOT NULL,
    name VARCHAR(256) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Add foreign key constraint
ALTER TABLE products 
ADD CONSTRAINT fk_products_category_id 
FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;