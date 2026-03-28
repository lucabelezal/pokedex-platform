import json
from pathlib import Path

base = Path('infra/postgres/source-json')
gen = base / 'generated'

facts = json.loads((gen / 'gen1_facts.json').read_text())
evo = json.loads((gen / 'gen1_evolution_chains.json').read_text())
drafts = json.loads((gen / 'gen1_description_drafts_ptbr.json').read_text())

species = json.loads((base / '06_species.json').read_text())
egg_groups = json.loads((base / '03_egg_group.json').read_text())

species_by_num = {}
for s in species:
    species_by_num.setdefault(s['pokemon_number'], []).append(s)

egg_name_to_id = {e['name']: int(e['id']) for e in egg_groups}
slug_to_egg_pt = {
    'monster': 'Monstro',
    'water1': 'Água 1',
    'bug': 'Inseto',
    'flying': 'Voador',
    'field': 'Campo',
    'ground': 'Campo',
    'fairy': 'Fada',
    'grass': 'Grama',
    'plant': 'Grama',
    'human-like': 'Humanoide',
    'humanshape': 'Humanoide',
    'water3': 'Água 3',
    'mineral': 'Mineral',
    'amorphous': 'Amorfo',
    'indeterminate': 'Amorfo',
    'water2': 'Água 2',
    'ditto': 'Ditto',
    'dragon': 'Dragão',
    'no-eggs': 'Indescoberto',
}

pokemon_to_chain_id = {}


def walk(node, chain_id):
    p = node.get('pokemon') or {}
    pid = p.get('id')
    if isinstance(pid, int):
        pokemon_to_chain_id[pid] = chain_id
    for nxt in node.get('evolutions_to', []) or []:
        walk(nxt, chain_id)


for chain in evo:
    walk(chain['chain'], int(chain['id']))

pt_drafts = {int(d['id']): d['description_ptbr_draft'] for d in drafts}

stats_rows = []
pokemon_rows = []
errors = []

for f in facts:
    pid = int(f['id'])
    number = f['number']
    name = f['name']

    st = f['stats']
    stats_rows.append({
        'id': pid,
        'pokemon_number': number,
        'pokemon_name': name,
        'total': int(st['total']),
        'hp': int(st['hp']),
        'attack': int(st['attack']),
        'defense': int(st['defense']),
        'sp_atk': int(st['sp_atk']),
        'sp_def': int(st['sp_def']),
        'speed': int(st['speed']),
    })

    candidates = species_by_num.get(number, [])
    species_id = None
    for c in candidates:
        if c.get('name') == name:
            species_id = int(c['id'])
            break
    if species_id is None and candidates:
        species_id = int(candidates[0]['id'])
    if species_id is None:
        errors.append(f'species not found for {number} {name}')
        continue

    egg_ids = []
    for slug in f['species'].get('egg_group_slugs_en', []):
        pt_name = slug_to_egg_pt.get(slug)
        if not pt_name or pt_name not in egg_name_to_id:
            errors.append(f'egg group unknown: {slug} ({number} {name})')
            continue
        egg_ids.append(egg_name_to_id[pt_name])

    abilities = []
    for ab in f.get('abilities', []):
        ability_id = int(ab.get('ability_id') or 0)
        if ability_id <= 0:
            slug = ab.get('name_en', '')
            errors.append(f'ability_id ausente: {slug} ({number} {name})')
            continue
        abilities.append({
            'ability_id': ability_id,
            'is_hidden': bool(ab.get('is_hidden', False)),
        })

    gr = int(f['species'].get('gender_rate_value', -1))
    gender = None
    if 0 <= gr <= 8:
        female = round(gr * 12.5, 1)
        male = round((8 - gr) * 12.5, 1)
        gender = {'male': male, 'female': female}

    row = {
        'id': pid,
        'number': number,
        'name': name,
        'description': pt_drafts.get(pid, ''),
        'height': float(f['height_m']),
        'weight': float(f['weight_kg']),
        'stats_id': pid,
        'generation_id': int(f['generation_id']),
        'species_id': species_id,
        'region_id': int(f['region_id']),
        'evolution_chain_id': int(pokemon_to_chain_id.get(pid, 0)),
        'gender_rate_value': gr,
        'egg_cycles': int(f['species'].get('hatch_counter', 0)),
        'egg_group_ids': egg_ids,
        'type_ids': [int(t) for t in f.get('type_ids', [])],
        'abilities': abilities,
        'sprites': f.get('sprites', {}),
    }
    if gender is not None:
        row['gender'] = gender

    pokemon_rows.append(row)

if errors:
    raise SystemExit('Erro(s) de mapeamento:\n- ' + '\n- '.join(errors[:20]) + ('' if len(errors) <= 20 else f'\n... total={len(errors)}'))

(base / '07_stats.json').write_text(json.dumps(stats_rows, ensure_ascii=False, indent=2) + '\n')
(base / '08_evolution_chains.json').write_text(json.dumps(evo, ensure_ascii=False, indent=2) + '\n')
(base / '09_pokemon.json').write_text(json.dumps(pokemon_rows, ensure_ascii=False, indent=2) + '\n')

print('ok: 07_stats.json ->', len(stats_rows), 'registros')
print('ok: 08_evolution_chains.json ->', len(evo), 'cadeias')
print('ok: 09_pokemon.json ->', len(pokemon_rows), 'pokemons')
