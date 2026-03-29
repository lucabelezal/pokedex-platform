INSERT INTO pokemons (
    id, number, name, height, weight, description, sprites,
    gender_male, gender_female, generation_id, region_id
) VALUES
    (245, '0245', 'Suicune', 2.0, 187.0, '', '{"other":{"home":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/home/245.png"}}}'::jsonb, NULL, NULL, 2, 2),
    (306, '0306', 'Aggron', 2.1, 360.0, '', '{"other":{"home":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/home/306.png"}}}'::jsonb, 50.0, 50.0, 3, 3),
    (384, '0384', 'Rayquaza', 7.0, 206.5, '', '{"other":{"home":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/home/384.png"}}}'::jsonb, NULL, NULL, 3, 3),
    (448, '0448', 'Lucario', 1.2, 54.0, '', '{"other":{"home":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/home/448.png"}}}'::jsonb, 87.5, 12.5, 4, 4),
    (497, '0497', 'Serperior', 3.3, 60.3, '', '{"other":{"home":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/home/497.png"}}}'::jsonb, 87.5, 12.5, 5, 5),
    (571, '0571', 'Zoroark', 1.6, 81.1, '', '{"other":{"home":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/home/571.png"}}}'::jsonb, 87.5, 12.5, 5, 5),
    (609, '0609', 'Chandelure', 1.0, 34.3, '', '{"other":{"home":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/home/609.png"}}}'::jsonb, 50.0, 50.0, 5, 5),
    (613, '0613', 'Cubchoo', 0.5, 8.5, '', '{"other":{"home":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/home/613.png"}}}'::jsonb, 50.0, 50.0, 5, 5),
    (733, '0733', 'Toucannon', 1.1, 26.0, '', '{"other":{"home":{"front_default":"https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/home/733.png"}}}'::jsonb, 50.0, 50.0, 7, 7)
ON CONFLICT (id) DO NOTHING;

INSERT INTO pokemon_types (pokemon_id, type_id) VALUES
    (245, 3),
    (306, 17), (306, 9),
    (384, 15),
    (448, 7), (448, 17),
    (497, 5),
    (571, 16),
    (609, 14), (609, 2),
    (613, 6),
    (733, 10), (733, 1)
ON CONFLICT (pokemon_id, type_id) DO NOTHING;

INSERT INTO pokemon_weaknesses (pokemon_id, type_id) VALUES
    (245, 4), (245, 5),
    (306, 2), (306, 3), (306, 7), (306, 9),
    (384, 6), (384, 13), (384, 15), (384, 18),
    (448, 2), (448, 7), (448, 9),
    (497, 2), (497, 6), (497, 8), (497, 10),
    (571, 7), (571, 12), (571, 18),
    (609, 3), (609, 9), (609, 13), (609, 14), (609, 16),
    (613, 2), (613, 7), (613, 13), (613, 17),
    (733, 4), (733, 6), (733, 13)
ON CONFLICT (pokemon_id, type_id) DO NOTHING;
