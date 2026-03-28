#!/usr/bin/env python3
"""
Coleta dados factuais dos Pokemon da Geracao 1 via PokeAPI e gera arquivos JSON
para apoiar o preenchimento dos dados da Pokedex no projeto.

Saidas padrao:
- gen1_facts.json
- gen1_evolution_chains.json
- gen1_description_drafts_ptbr.json
"""

from __future__ import annotations

import argparse
import json
import pathlib
import time
import urllib.error
import urllib.request
from dataclasses import dataclass
from typing import Any

POKEAPI_BASE = "https://pokeapi.co/api/v2"

# Mapeia slug EN de tipo para id canônico no projeto (02_type.json).
TYPE_ID_BY_SLUG = {
    "normal": 1,
    "fire": 2,
    "water": 3,
    "electric": 4,
    "grass": 5,
    "ice": 6,
    "fighting": 7,
    "poison": 8,
    "ground": 9,
    "flying": 10,
    "psychic": 11,
    "bug": 12,
    "rock": 13,
    "ghost": 14,
    "dragon": 15,
    "dark": 16,
    "steel": 17,
    "fairy": 18,
}

# Mapeamento util para montar texto em PT-BR sem copiar descricoes de terceiros.
TYPE_PT_BY_SLUG = {
    "normal": "Normal",
    "fire": "Fogo",
    "water": "Agua",
    "electric": "Eletrico",
    "grass": "Grama",
    "ice": "Gelo",
    "fighting": "Lutador",
    "poison": "Venenoso",
    "ground": "Terrestre",
    "flying": "Voador",
    "psychic": "Psiquico",
    "bug": "Inseto",
    "rock": "Pedra",
    "ghost": "Fantasma",
    "dragon": "Dragao",
    "dark": "Sombrio",
    "steel": "Aco",
    "fairy": "Fada",
}


@dataclass
class HttpClient:
    timeout: float = 20.0
    retries: int = 3
    retry_delay_s: float = 0.6

    def get_json(self, url: str) -> dict[str, Any]:
        last_error: Exception | None = None
        for attempt in range(1, self.retries + 1):
            try:
                req = urllib.request.Request(
                    url,
                    headers={
                        "User-Agent": "pokedex-gen1-facts-script/1.0",
                        "Accept": "application/json",
                    },
                    method="GET",
                )
                with urllib.request.urlopen(req, timeout=self.timeout) as resp:
                    return json.loads(resp.read().decode("utf-8"))
            except (urllib.error.URLError, urllib.error.HTTPError, TimeoutError, json.JSONDecodeError) as exc:
                last_error = exc
                if attempt < self.retries:
                    time.sleep(self.retry_delay_s * attempt)
                    continue
                break

        raise RuntimeError(f"falha ao buscar {url}: {last_error}")


def format_number_4(n: int) -> str:
    return f"{n:04d}"


def english_genus(species_json: dict[str, Any]) -> str:
    for entry in species_json.get("genera", []):
        if entry.get("language", {}).get("name") == "en":
            return str(entry.get("genus", "")).strip()
    return ""


def english_flavor(species_json: dict[str, Any]) -> str:
    for entry in species_json.get("flavor_text_entries", []):
        if entry.get("language", {}).get("name") == "en":
            raw = str(entry.get("flavor_text", ""))
            return " ".join(raw.replace("\n", " ").replace("\f", " ").split())
    return ""


def stat_value(pokemon_json: dict[str, Any], stat_name: str) -> int:
    for stat in pokemon_json.get("stats", []):
        if stat.get("stat", {}).get("name") == stat_name:
            return int(stat.get("base_stat", 0))
    return 0


def extract_chain_node(node: dict[str, Any]) -> dict[str, Any]:
    details = node.get("evolution_details") or []
    condition = {"type": "unknown", "description": "Condicao nao mapeada"}

    if details:
        d0 = details[0]
        trigger = (d0.get("trigger") or {}).get("name")
        if trigger == "level-up" and d0.get("min_level") is not None:
            level = int(d0.get("min_level"))
            condition = {
                "type": "level_up",
                "value": level,
                "description": f"Nivel {level}",
            }
        elif trigger == "use-item" and d0.get("item"):
            item = str(d0.get("item", {}).get("name", "item")).replace("-", " ")
            condition = {
                "type": "stone_or_item",
                "description": f"Uso de {item}",
            }
        elif trigger == "trade":
            condition = {
                "type": "trade",
                "description": "Troca",
            }
        elif trigger == "level-up" and d0.get("min_happiness") is not None:
            condition = {
                "type": "happiness",
                "description": "Subir nivel com felicidade",
            }

    return {
        "pokemon": {
            "id": int((node.get("species", {}).get("url", "").rstrip("/").split("/")[-1]) or 0),
            "name": str(node.get("species", {}).get("name", "")).title(),
        },
        "condition": condition,
        "evolutions_to": [extract_chain_node(child) for child in node.get("evolves_to", [])],
    }


def build_sprites_block(pokemon_id: int, pokemon_json: dict[str, Any]) -> dict[str, Any]:
    s = pokemon_json.get("sprites", {})
    # Mantem os links mais estaveis em raw.githubusercontent (padrao adotado no projeto).
    raw = "https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon"

    return {
        "back_default": f"{raw}/back/{pokemon_id}.png",
        "back_female": s.get("back_female"),
        "back_shiny": f"{raw}/back/shiny/{pokemon_id}.png",
        "back_shiny_female": s.get("back_shiny_female"),
        "front_default": f"{raw}/{pokemon_id}.png",
        "front_female": s.get("front_female"),
        "front_shiny": f"{raw}/shiny/{pokemon_id}.png",
        "front_shiny_female": s.get("front_shiny_female"),
        "other": {
            "dream_world": {
                "front_default": s.get("other", {}).get("dream_world", {}).get("front_default"),
                "front_female": s.get("other", {}).get("dream_world", {}).get("front_female"),
            },
            "home": {
                "front_default": f"{raw}/other/home/{pokemon_id}.png",
                "front_female": s.get("other", {}).get("home", {}).get("front_female"),
                "front_shiny": f"{raw}/other/home/shiny/{pokemon_id}.png",
                "front_shiny_female": s.get("other", {}).get("home", {}).get("front_shiny_female"),
            },
            "official-artwork": {
                "front_default": f"{raw}/other/official-artwork/{pokemon_id}.png",
                "front_shiny": f"{raw}/other/official-artwork/shiny/{pokemon_id}.png",
            },
            "showdown": {
                "back_default": s.get("other", {}).get("showdown", {}).get("back_default"),
                "back_female": s.get("other", {}).get("showdown", {}).get("back_female"),
                "back_shiny": s.get("other", {}).get("showdown", {}).get("back_shiny"),
                "back_shiny_female": s.get("other", {}).get("showdown", {}).get("back_shiny_female"),
                "front_default": s.get("other", {}).get("showdown", {}).get("front_default"),
                "front_female": s.get("other", {}).get("showdown", {}).get("front_female"),
                "front_shiny": s.get("other", {}).get("showdown", {}).get("front_shiny"),
                "front_shiny_female": s.get("other", {}).get("showdown", {}).get("front_shiny_female"),
            },
        },
    }


def create_ptbr_description_draft(name: str, type_slugs: list[str], genus_en: str, habitat: str | None) -> str:
    type_pt = [TYPE_PT_BY_SLUG.get(t, t.title()) for t in type_slugs]
    type_text = " e ".join(type_pt) if len(type_pt) == 2 else type_pt[0]
    habitat_text = habitat.replace("-", " ") if habitat else "ambientes variados"

    return (
        f"{name} e um Pokemon do tipo {type_text}. "
        f"E conhecido na categoria {genus_en or 'Pokemon'} e costuma aparecer em {habitat_text}. "
        "Este texto e um rascunho original para revisao editorial em portugues do Brasil."
    )


def collect_gen1(client: HttpClient, start: int, end: int) -> tuple[list[dict[str, Any]], list[dict[str, Any]], list[dict[str, Any]]]:
    facts: list[dict[str, Any]] = []
    chains: list[dict[str, Any]] = []
    desc_drafts: list[dict[str, Any]] = []

    seen_chain_ids: set[int] = set()

    for pokemon_id in range(start, end + 1):
        pokemon = client.get_json(f"{POKEAPI_BASE}/pokemon/{pokemon_id}")
        species = client.get_json(f"{POKEAPI_BASE}/pokemon-species/{pokemon_id}")

        type_slugs = [t.get("type", {}).get("name", "") for t in sorted(pokemon.get("types", []), key=lambda x: x.get("slot", 0))]
        type_ids = [TYPE_ID_BY_SLUG[t] for t in type_slugs if t in TYPE_ID_BY_SLUG]

        abilities = []
        for ability in sorted(pokemon.get("abilities", []), key=lambda x: x.get("slot", 0)):
            ability_url = str(ability.get("ability", {}).get("url", ""))
            ability_id = 0
            if ability_url:
                try:
                    ability_id = int(ability_url.rstrip("/").split("/")[-1])
                except ValueError:
                    ability_id = 0
            abilities.append(
                {
                    "ability_id": ability_id,
                    "name_en": ability.get("ability", {}).get("name"),
                    "is_hidden": bool(ability.get("is_hidden", False)),
                    "slot": int(ability.get("slot", 0)),
                }
            )

        egg_group_slugs = [e.get("name") for e in species.get("egg_groups", [])]

        fact = {
            "id": pokemon_id,
            "number": format_number_4(pokemon_id),
            "name": str(pokemon.get("name", "")).title(),
            "generation_id": 1,
            "region_id": 1,
            "type_slugs_en": type_slugs,
            "type_ids": type_ids,
            "height_m": float(pokemon.get("height", 0)) / 10.0,
            "weight_kg": float(pokemon.get("weight", 0)) / 10.0,
            "stats": {
                "total": stat_value(pokemon, "hp")
                + stat_value(pokemon, "attack")
                + stat_value(pokemon, "defense")
                + stat_value(pokemon, "special-attack")
                + stat_value(pokemon, "special-defense")
                + stat_value(pokemon, "speed"),
                "hp": stat_value(pokemon, "hp"),
                "attack": stat_value(pokemon, "attack"),
                "defense": stat_value(pokemon, "defense"),
                "sp_atk": stat_value(pokemon, "special-attack"),
                "sp_def": stat_value(pokemon, "special-defense"),
                "speed": stat_value(pokemon, "speed"),
            },
            "abilities": abilities,
            "species": {
                "genus_en": english_genus(species),
                "flavor_en": english_flavor(species),
                "capture_rate": species.get("capture_rate"),
                "base_happiness": species.get("base_happiness"),
                "is_legendary": species.get("is_legendary"),
                "is_mythical": species.get("is_mythical"),
                "gender_rate_value": species.get("gender_rate"),
                "hatch_counter": species.get("hatch_counter"),
                "egg_group_slugs_en": egg_group_slugs,
                "habitat_en": (species.get("habitat") or {}).get("name"),
            },
            "sprites": build_sprites_block(pokemon_id, pokemon),
        }
        facts.append(fact)

        desc_drafts.append(
            {
                "id": pokemon_id,
                "number": format_number_4(pokemon_id),
                "name": fact["name"],
                "description_ptbr_draft": create_ptbr_description_draft(
                    name=fact["name"],
                    type_slugs=type_slugs,
                    genus_en=fact["species"]["genus_en"],
                    habitat=fact["species"]["habitat_en"],
                ),
            }
        )

        chain_url = str((species.get("evolution_chain") or {}).get("url", "")).strip()
        if chain_url:
            chain_id = int(chain_url.rstrip("/").split("/")[-1])
            if chain_id not in seen_chain_ids:
                seen_chain_ids.add(chain_id)
                chain_json = client.get_json(chain_url)
                chain_root = chain_json.get("chain", {})
                chains.append(
                    {
                        "id": chain_id,
                        "chain": {
                            "pokemon": {
                                "id": int((chain_root.get("species", {}).get("url", "").rstrip("/").split("/")[-1]) or 0),
                                "name": str(chain_root.get("species", {}).get("name", "")).title(),
                            },
                            "evolutions_to": [extract_chain_node(child) for child in chain_root.get("evolves_to", [])],
                        },
                    }
                )

        print(f"[ok] coletado #{format_number_4(pokemon_id)} {fact['name']}")

    chains.sort(key=lambda x: x["id"])
    return facts, chains, desc_drafts


def main() -> int:
    parser = argparse.ArgumentParser(description="Coleta fatos da Gen 1 na PokeAPI para apoiar o dataset local")
    parser.add_argument("--start", type=int, default=1, help="ID inicial (padrao: 1)")
    parser.add_argument("--end", type=int, default=151, help="ID final (padrao: 151)")
    parser.add_argument(
        "--output-dir",
        type=str,
        default="infra/postgres/source-json/generated",
        help="Diretorio de saida para os JSONs gerados",
    )
    parser.add_argument("--timeout", type=float, default=20.0, help="Timeout HTTP em segundos")
    parser.add_argument("--retries", type=int, default=3, help="Tentativas por requisicao")
    args = parser.parse_args()

    if args.start < 1 or args.end < args.start:
        raise SystemExit("intervalo invalido: use start>=1 e end>=start")

    repo_root = pathlib.Path(__file__).resolve().parents[4]
    out_dir = repo_root / args.output_dir
    out_dir.mkdir(parents=True, exist_ok=True)

    client = HttpClient(timeout=args.timeout, retries=args.retries)
    facts, chains, desc_drafts = collect_gen1(client, args.start, args.end)

    (out_dir / "gen1_facts.json").write_text(json.dumps(facts, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    (out_dir / "gen1_evolution_chains.json").write_text(
        json.dumps(chains, ensure_ascii=False, indent=2) + "\n",
        encoding="utf-8",
    )
    (out_dir / "gen1_description_drafts_ptbr.json").write_text(
        json.dumps(desc_drafts, ensure_ascii=False, indent=2) + "\n",
        encoding="utf-8",
    )

    print("\nArquivos gerados:")
    print(f"- {out_dir / 'gen1_facts.json'}")
    print(f"- {out_dir / 'gen1_evolution_chains.json'}")
    print(f"- {out_dir / 'gen1_description_drafts_ptbr.json'}")

    return 0


if __name__ == "__main__":
    raise SystemExit(main())
