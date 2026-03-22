-- Migration: Create pokemons table
CREATE TABLE IF NOT EXISTS pokemons (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    number VARCHAR(10) NOT NULL UNIQUE,
    types TEXT[] DEFAULT ARRAY[]::TEXT[],
    height DECIMAL(5, 2),
    weight DECIMAL(5, 2),
    description TEXT,
    image_url VARCHAR(512),
    element_color VARCHAR(50),
    element_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for search queries
CREATE INDEX IF NOT EXISTS idx_pokemons_name ON pokemons (name);
CREATE INDEX IF NOT EXISTS idx_pokemons_number ON pokemons (number);
CREATE INDEX IF NOT EXISTS idx_pokemons_types ON pokemons USING GIN (types);

-- Migration: Create users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(255) PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Migration: Create favorites table
CREATE TABLE IF NOT EXISTS favorites (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id VARCHAR(255) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pokemon_id VARCHAR(255) NOT NULL REFERENCES pokemons(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, pokemon_id)
);

-- Create index for user favorites queries
CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites (user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_pokemon_id ON favorites (pokemon_id);
