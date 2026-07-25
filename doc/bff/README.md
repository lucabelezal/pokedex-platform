# Backend for Frontend (BFF)

## Estrutura

```
doc/bff/
├── README.md                 ← este arquivo
└── 01-padrao-bff.md          ← padrão arquitetural BFF sob a ótica de System Design
```

## Contexto no projeto

A implementação prática do BFF na Pokedex Platform está documentada em [doc/BFF.md](../BFF.md), que cobre:

- Responsabilidades atuais do `mobile-bff`
- Contratos de API (home, detalhe, regiões, favoritos, perfil)
- Modelagem UI-oriented
- Diagrama de comunicação e arquitetura hexagonal
- Pontos de melhoria identificados

O arquivo `01-padrao-bff.md` complementa com a **teoria arquitetural** — os conceitos por trás das decisões de design já aplicadas no projeto.

## Referência

| Fonte | Conteúdo |
|-------|----------|
| [fidelissauro.dev/bffs](https://fidelissauro.dev/bffs/) | Padrão BFF sob a ótica de System Design, por Matheus Fidelis |

## Conexões com a Pokedex Platform

O `mobile-bff` já implementa vários conceitos do artigo:

| Conceito do artigo | Onde aparece no projeto |
|-------------------|------------------------|
| **API Composition** | `PokemonService` agrega dados do catálogo + favoritos do PostgreSQL em uma única resposta |
| **Segregação de canais** | O BFF é dedicado ao cliente mobile — contratos, payloads e estados de tela são mobile-first |
| **Resiliência** | Circuit breaker com retry em todos os clients HTTP (`circuit_breaker.go`) |
| **Contratos UI-oriented** | Estados de tela (`unauthenticated`, `empty`, `has_data`), labels, CTAs nos payloads |
| **Desacoplamento de métricas** | Observabilidade com OpenTelemetry (`infrastructure/observability/`) |

### Próximos passos alinhados ao artigo

- **Blast radius:** isolar grupos de chamadas com bulkheads — evitar que falha no catálogo derrube autenticação
- **Versionamento:** feature toggles para experimentar novas versões de contrato sem quebrar clientes existentes
- **Segregação futura:** se houver um cliente web ou IoT, criar BFFs dedicados em vez de encher o `mobile-bff` de condicionais
