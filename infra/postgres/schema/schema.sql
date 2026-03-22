-- init.sql

-- Drop tables if they exist to ensure a clean slate
-- CASCADE é usado para remover automaticamente objetos que dependem das tabelas que estão sendo descartadas (por exemplo, chaves estrangeiras)
DROP TABLE IF EXISTS pokemon_types, pokemon_abilities, pokemon_egg_groups, pokemon_weaknesses, evolution_chains, pokemons, species, stats, abilities, generations, regions, types, egg_groups CASCADE;

-- Table: regions
CREATE TABLE regions (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Table: types
CREATE TABLE types (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    color VARCHAR(7)
);

-- Table: egg_groups
CREATE TABLE egg_groups (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

-- Table: generations
CREATE TABLE generations (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    region_id BIGINT,
    FOREIGN KEY (region_id) REFERENCES regions(id)
);

-- Table: abilities
CREATE TABLE abilities (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    introduced_generation_id BIGINT,
    FOREIGN KEY (introduced_generation_id) REFERENCES generations(id)
);

-- Table: species
CREATE TABLE species (
    id BIGINT PRIMARY KEY,
    pokemon_number VARCHAR(10),
    name VARCHAR(255) NOT NULL,
    species_en VARCHAR(255),
    species_pt VARCHAR(255)
);

-- Table: stats (Agora independente de dados específicos de Pokémon, apenas as estatísticas)
CREATE TABLE stats (
    id BIGINT PRIMARY KEY,
    total INT,
    hp INT,
    attack INT,
    defense INT,
    sp_atk INT,
    sp_def INT,
    speed INT
);

-- Table: evolution_chains (Agora armazena os dados da cadeia diretamente como JSONB)
CREATE TABLE evolution_chains (
    id BIGINT PRIMARY KEY,
    chain_data JSONB NOT NULL -- Adicionado: Coluna JSONB para armazenar a estrutura aninhada da cadeia de evolução
);

-- Table: pokemons
CREATE TABLE pokemons (
    id BIGINT PRIMARY KEY,
    number VARCHAR(10),
    name VARCHAR(255) NOT NULL,
    height NUMERIC(5, 2),
    weight NUMERIC(6, 2),
    description TEXT,
    sprites JSONB,
    gender_male NUMERIC(4, 1),
    gender_female NUMERIC(4, 1),
    gender_rate_value INT,
    egg_cycles INT,
    stats_id BIGINT UNIQUE, -- UNIQUE para relação One-to-One
    generation_id BIGINT,
    species_id BIGINT,
    region_id BIGINT,
    evolution_chain_id BIGINT, -- 'pokemons' ainda referencia 'evolution_chains'
    FOREIGN KEY (stats_id) REFERENCES stats(id),
    FOREIGN KEY (generation_id) REFERENCES generations(id),
    FOREIGN KEY (species_id) REFERENCES species(id),
    FOREIGN KEY (region_id) REFERENCES regions(id),
    FOREIGN KEY (evolution_chain_id) REFERENCES evolution_chains(id) -- Mantida a referência de 'pokemons' para 'evolution_chains'
);

-- Join Table: pokemon_types (Many-to-Many)
CREATE TABLE pokemon_types (
    pokemon_id BIGINT,
    type_id BIGINT,
    PRIMARY KEY (pokemon_id, type_id),
    FOREIGN KEY (pokemon_id) REFERENCES pokemons(id),
    FOREIGN KEY (type_id) REFERENCES types(id)
);

-- Join Table: pokemon_abilities (Many-to-Many with extra column - modelada como entidade própria)
CREATE TABLE pokemon_abilities (
    id BIGSERIAL PRIMARY KEY,
    pokemon_id BIGINT NOT NULL,
    ability_id BIGINT NOT NULL,
    is_hidden BOOLEAN NOT NULL,
    FOREIGN KEY (pokemon_id) REFERENCES pokemons(id),
    FOREIGN KEY (ability_id) REFERENCES abilities(id)
);

-- Join Table: pokemon_egg_groups (Many-to-Many)
CREATE TABLE pokemon_egg_groups (
    pokemon_id BIGINT,
    egg_group_id BIGINT,
    PRIMARY KEY (pokemon_id, egg_group_id),
    FOREIGN KEY (pokemon_id) REFERENCES pokemons(id),
    FOREIGN KEY (egg_group_id) REFERENCES egg_groups(id)
);

-- Join Table: pokemon_weaknesses (Many-to-Many)
CREATE TABLE pokemon_weaknesses (
    pokemon_id BIGINT,
    type_id BIGINT,
    PRIMARY KEY (pokemon_id, type_id),
    FOREIGN KEY (pokemon_id) REFERENCES pokemons(id),
    FOREIGN KEY (type_id) REFERENCES types(id)
);
