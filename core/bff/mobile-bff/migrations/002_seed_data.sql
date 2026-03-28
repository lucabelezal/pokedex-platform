-- Seed data for development and testing

-- Insert Pokemons
INSERT INTO pokemons (id, name, number, types, height, weight, description, image_url, element_color, element_type)
VALUES 
    ('1', 'Bulbasaur', '001', ARRAY['Grass', 'Poison'], 0.7, 6.9, 'Bulbasaur can be seen napping in bright sunlight. There is a seed on its back.', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/1.png', 'green', 'Grass'),
    ('4', 'Charmander', '004', ARRAY['Fire'], 0.6, 8.5, 'The flame on its tail shows the strength of its life force.', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/4.png', 'red', 'Fire'),
    ('7', 'Squirtle', '007', ARRAY['Water'], 0.5, 9.0, 'It hides in its shell when it attacks. It squirts water with high accuracy.', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/7.png', 'blue', 'Water'),
    ('25', 'Pikachu', '025', ARRAY['Electric'], 0.4, 6.0, 'When several of these Pokémon gather, their electricity can build and cause lightning storms.', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/25.png', 'yellow', 'Electric'),
    ('39', 'Jigglypuff', '039', ARRAY['Normal', 'Fairy'], 0.5, 5.5, 'Jigglypuff''s body is soft and rubbery. When angered, it will suck in air and inflate itself to balloon.', 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/39.png', 'pink', 'Normal')
ON CONFLICT (id) DO NOTHING;

-- Insert test users
INSERT INTO users (id, email)
VALUES 
    ('user-001', 'user1@example.com'),
    ('user-002', 'user2@example.com'),
    ('user-003', 'user3@example.com')
ON CONFLICT (id) DO NOTHING;

-- Insert some test favorites
INSERT INTO favorites (user_id, pokemon_id)
VALUES 
    ('user-001', '25'),
    ('user-001', '1'),
    ('user-002', '4'),
    ('user-002', '7')
ON CONFLICT (user_id, pokemon_id) DO NOTHING;
