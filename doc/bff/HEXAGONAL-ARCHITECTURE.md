# Arquitetura Hexagonal

## Objetivo

Este documento explica como a arquitetura hexagonal está aplicada atualmente dentro de `core/bff/mobile-bff`.

## Modelo De Camadas

```text
Adapters de entrada
  -> ports de entrada
    -> services de aplicação
      -> ports de saída
        -> adapters de saída
```

## Mapeamento Atual No Projeto

### Domínio

Caminho: `core/bff/mobile-bff/internal/domain`

Responsabilidade:

- Definir modelos centrais de negócio.
- Guardar erros em nível de domínio e regras de validação.
- Permanecer independente de HTTP, PostgreSQL e serviços externos.

Exemplos:

- `pokemon.go`
- `errors.go`

### Ports De Entrada

Caminho: `core/bff/mobile-bff/internal/ports`

Responsabilidade:

- Definir o que a aplicação pode fazer do ponto de vista de quem a chama.
- Descrever casos de uso sem acoplá-los a HTTP.

Exemplos:

- `PokemonUseCase`
- `FavoriteUseCase`

### Services De Aplicação

Caminho: `core/bff/mobile-bff/internal/service`

Responsabilidade:

- Implementar casos de uso.
- Coordenar objetos de domínio e ports de saída.
- Aplicar regras de aplicação como validação, orquestração e paginação padrão.

Exemplos:

- `PokemonService`
- `FavoriteService`

### Ports De Saída

Caminho: `core/bff/mobile-bff/internal/ports`

Responsabilidade:

- Definir quais capacidades de infraestrutura a aplicação precisa.
- Esconder detalhes de persistência e comunicação remota atrás de interfaces.

Exemplos:

- `PokemonRepository`
- `FavoriteRepository`

### Adapters De Entrada

Caminho: `core/bff/mobile-bff/internal/adapters/http`

Responsabilidade:

- Receber requisições HTTP.
- Fazer o parse de dados específicos do transporte.
- Chamar casos de uso.
- Converter resultados da aplicação em DTOs e respostas HTTP.

Exemplos:

- `handlers.go`
- `middleware.go`
- `response_builder.go`
- `dto/`

### Adapters De Saída

Caminho: `core/bff/mobile-bff/internal/adapters/repository`

Responsabilidade:

- Implementar ports usando PostgreSQL ou serviços HTTP remotos.
- Traduzir detalhes de infraestrutura para um comportamento voltado à aplicação.

Exemplos:

- `favorite_repository.go`
- `pokemon_repository.go`
- `pokemon_catalog_service_repository.go`
- `auth_service_client.go`

## O Que Está Funcionando Bem

- O código já separa domínio, ports, services e adapters.
- Os casos de uso estão representados como interfaces.
- Repositories são abstrações em vez de chamadas diretas ao banco dentro de handlers.
- DTOs HTTP permanecem dentro do pacote do adapter HTTP.
- Adapters de fallback de runtime vivem em código interno de runtime em vez de `tests/`.
- Os fluxos de auth agora passam por `AuthUseCase` em vez de acoplar o adapter HTTP ao client externo de auth.

## Onde O Hexágono Está Mais Fraco

### Mapeamento de negócio duplicado entre camadas

O mapeamento de cor por tipo de Pokémon existe em mais de um pacote.

Por que isso importa:

- O comportamento do domínio pode se desalinhar entre as camadas de adapter e service.

## Próxima Refatoração Recomendada

### Passo 1

Centralizar as regras de apresentação de tipo de Pokémon usadas tanto pelo service quanto pelos response builders.

### Passo 2

Revisar se parte da lógica de enriquecimento de favoritos deve ir para um mapper voltado à aplicação ou para um query service dedicado.

### Passo 3

Continuar expandindo os testes ao redor das fronteiras arquiteturais sempre que novos adapters ou integrações forem introduzidos.

## Avaliação Final

O BFF agora aplica arquitetura hexagonal de forma mais consistente do que antes. O trabalho restante tem menos a ver com corrigir grandes violações de fronteira e mais com refinar regras de mapeamento compartilhadas e manter a arquitetura leve à medida que o projeto evolui.
