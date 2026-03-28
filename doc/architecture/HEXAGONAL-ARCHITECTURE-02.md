# Arquitetura Hexagonal 02

## O Que Vamos Ver

- Como a arquitetura hexagonal se mapeia neste repositório.
- Quais responsabilidades pertencem a `mobile-bff`, `pokemon-catalog-service` e `auth-service`.
- Quais decisões nós já tomamos.
- Em quais pontos a implementação atual é intencionalmente pragmática.

## Visão Do Repositório

No nível do repositório, a plataforma está organizada assim:

```text
core/
├── app/
│   ├── auth-service/
│   └── pokemon-catalog-service/
├── bff/
│   └── mobile-bff/
├── gateway/
├── infra/
└── docker-compose.yml
```

Isto não é arquitetura hexagonal para o monorepo inteiro. É uma plataforma orientada a serviços na qual um dos componentes, `mobile-bff`, se beneficia diretamente de um estilo de ports and adapters.

## Responsabilidades Dos Serviços

### `mobile-bff`

Possui a experiência do cliente e a composição das respostas.

Exemplos:

- endpoints orientados ao frontend
- orquestração de favoritos
- exposição dos fluxos de auth para os clientes
- transporte e formatação de resposta para necessidades da UI

### `pokemon-catalog-service`

Possui o catálogo canônico de Pokémons.

Exemplos:

- listar, buscar e filtrar Pokémons
- devolver dados canônicos do catálogo
- centralizar o comportamento de leitura do catálogo

### `auth-service`

Possui autenticação e ciclo de vida de tokens.

Exemplos:

- signup
- login
- refresh
- logout

Essa separação importa porque evita que o BFF vire um monólito disfarçado.

## Como A Hexagonal Se Aplica Ao BFF

Dentro de `core/bff/mobile-bff`, a estrutura é intencionalmente próxima de ports and adapters:

```text
internal/
├── domain/
├── ports/
├── service/
└── adapters/
```

### Domínio

Guarda conceitos locais de negócio e erros em nível de domínio.

Regra importante:

- o domínio não deve conhecer HTTP, PostgreSQL nem detalhes de implementação de serviços externos

### Ports

Descrevem do que a aplicação precisa e quais casos de uso ela oferece.

Exemplos:

- `PokemonUseCase`
- `FavoriteUseCase`
- `AuthUseCase`
- ports de repository e provider

### Services

Implementam os casos de uso e orquestram o fluxo de negócio.

É aqui que o BFF decide coisas como:

- paginação padrão
- orquestração de favoritos
- orquestração dos fluxos de auth por meio de um provider port

### Adapters

Traduzem detalhes externos para a linguagem esperada pelo core do BFF.

Exemplos:

- HTTP handlers
- PostgreSQL repositories
- auth client
- catalog client

## Decisões Arquiteturais Já Refletidas No Código

### Decisão 1: Manter O Código De Runtime Em `core/`

Movemos os artefatos executáveis da plataforma para `core/` para que a raiz do repositório pudesse ficar mais limpa para documentação e arquivos de colaboração.

Por que isso importa:

- deixa a área de runtime explícita
- separa documentação de arquitetura do código de implementação
- melhora a descoberta para novos contribuidores

### Decisão 2: Renomear `pokedex-service` Para `pokemon-catalog-service`

Renomeamos o serviço de catálogo para refletir sua responsabilidade real.

Por que isso importa:

- melhora a linguagem ubíqua
- reduz ambiguidade
- deixa a posse da responsabilidade do serviço mais clara

### Decisão 3: Remover Mocks De Teste Da Composição De Runtime

O runtime do BFF não importa mais implementações mockadas a partir de `tests/`.

Por que isso importa:

- a composição de produção não depende mais de pacotes de teste
- adapters de fallback agora vivem no código interno de runtime

### Decisão 4: Colocar Auth Atrás De Uma Fronteira De Caso De Uso

O handler HTTP agora conversa com `AuthUseCase`, e não diretamente com o client concreto de auth.

Por que isso importa:

- o handler não está mais acoplado ao contrato externo de auth
- o client adapter normaliza erros externos para erros locais de domínio
- o BFF ganhou uma fronteira de aplicação mais limpa

## Onde Ainda Estamos Sendo Pragmáticos

Este projeto não está tentando modelar todo conceito com rigor profundo de DDD.

Por exemplo:

- nem todo tipo é um aggregate rico
- nem todo fluxo precisa de múltiplas camadas de indireção
- parte da lógica de mapeamento ainda é simples e próxima do adapter

Isso é aceitável. O ponto não é maximizar pureza teórica. O ponto é manter a arquitetura compreensível e útil para o time.

## Recapitulação

Neste repositório, a arquitetura hexagonal é principalmente uma disciplina aplicada ao `mobile-bff`: contratos orientados ao negócio no centro, adapters específicos de tecnologia nas bordas e uma separação mais clara entre as responsabilidades de catálogo, auth e BFF. A estrutura é intencional, mas continua pragmática.
