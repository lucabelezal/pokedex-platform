# Padrões De Projeto Em Go

## O Que Vamos Ver

- Como pensar design patterns em Go sem cair em overengineering.
- Quais patterns aparecem naturalmente em arquitetura hexagonal.
- Quais patterns já são visíveis neste projeto.
- Quais patterns são úteis, arriscados ou desnecessários aqui.

## Introdução

Design patterns são soluções recorrentes para problemas recorrentes de design. Eles não são receitas obrigatórias e não devem ser aplicados de forma mecânica. Em Go, isso importa ainda mais porque a linguagem tende a recompensar designs mais simples, interfaces menores e composição explícita.

Isso significa que uma base saudável em Go normalmente não se parece com um sistema clássico de OOP cheio de camadas e hierarquias. Ainda assim, design patterns continuam importando porque ajudam a nomear e entender estruturas que já estamos usando.

## A Regra Mais Importante

Use um pattern quando ele tornar uma fronteira mais clara ou simplificar mudança.

Não use um pattern quando ele só adicionar cerimônia.

## Patterns Que Se Encaixam Naturalmente Com Arquitetura Hexagonal

### Adapter

Este é um dos patterns mais importantes para este projeto.

Por que isso importa:

- handlers HTTP adaptam requisições de transporte para chamadas de caso de uso.
- implementações de repository adaptam PostgreSQL ou integrações HTTP remotas para ports da aplicação.
- clients de serviços externos traduzem contratos técnicos para comportamento voltado à aplicação.

Exemplos neste repositório:

- `core/bff/mobile-bff/internal/adapters/http`
- `core/bff/mobile-bff/internal/adapters/repository`

### Strategy

Strategy aparece quando o comportamento é escolhido por meio de uma interface e as implementações podem ser trocadas.

Por que isso importa:

- ports como repositories e providers permitem injetar implementações diferentes
- o mesmo caso de uso pode rodar com implementações mockadas, postgres ou serviços remotos

Exemplos neste repositório:

- `PokemonRepository`
- `FavoriteRepository`
- `AuthProvider`

### Facade

Facade é útil quando um componente expõe uma interface mais simples sobre várias partes móveis.

Por que isso importa:

- um BFF muitas vezes se comporta como uma facade para clientes de frontend
- ele esconde topologia interna de serviços e detalhes de orquestração

Neste projeto:

- `mobile-bff` não é apenas um BFF pelo papel de deploy, mas também se aproxima de uma facade do ponto de vista de pattern

### Factory Method

Ele aparece quando a lógica de construção é centralizada e implementações diferentes podem ser escolhidas por configuração.

Por que isso importa:

- a composição de runtime em `main` frequentemente atua como uma fronteira prática de factory
- a construção de adapters pode variar por ambiente

Neste projeto:

- `cmd/server/main.go` funciona como uma composition root leve e como ponto de factory

### Repository

Repository não faz parte do catálogo clássico do Gang of Four, mas é muito relevante para DDD e arquitetura hexagonal.

Por que isso importa:

- ele cria uma fronteira entre lógica de negócio e detalhes de persistência
- ele permite que services de aplicação falem a linguagem do domínio em vez da linguagem de SQL

Exemplos neste repositório:

- `PostgresFavoriteRepository`
- `PostgresPokemonRepository`

## Patterns Que Devem Ser Usados Com Cuidado

### Singleton

Desenvolvedores Go frequentemente recorrem a globais em nível de pacote. Isso pode parecer conveniente, mas cria dependências escondidas e dificulta testes.

Diretriz:

- prefira injeção explícita de dependência em vez de singletons globais

### Decorator

Decorator pode ser útil para logging, cache, tracing e métricas ao redor de interfaces. Mas ele também aumenta a indireção.

Diretriz:

- use quando preocupações transversais se repetem claramente
- evite se um wrapper direto já for suficiente

### Template Method

Este pattern é muito mais orientado à herança. Go normalmente prefere composição, interfaces e funções simples.

Diretriz:

- raramente é a primeira escolha em Go

## Patterns Que Normalmente Importam Menos Neste Projeto

Alguns patterns são úteis de forma geral, mas não são centrais para o formato atual deste repositório:

- Composite
- Flyweight
- Memento
- Bridge

Isso não os torna patterns ruins. Significa apenas que eles não estão guiando a arquitetura principal aqui neste momento.

## Orientação Prática Para Este Repositório

Ao trabalhar neste projeto, estas são as perguntas sobre patterns mais úteis de se fazer:

### Isto É Um Adapter?

Se o código traduz entre um protocolo técnico e a linguagem da aplicação, provavelmente é.

### Isto É Uma Strategy?

Se o código depende de uma interface com implementações trocáveis, provavelmente é.

### Isto É Uma Facade?

Se um componente esconde várias dependências internas atrás de uma superfície mais simples voltada ao cliente, pode ser.

### Um Repository Se Justifica Aqui?

Se o fluxo de negócio precisa de persistência sem conhecer detalhes de armazenamento, sim.

### Um Pattern Melhoraria A Clareza Ou Só Adicionaria Cerimônia?

Se a resposta for cerimônia, pule.

## Mapa De Patterns Orientado Ao Projeto

### `mobile-bff`

Patterns mais relevantes:

- Adapter
- Strategy
- Facade
- Repository

### `pokemon-catalog-service`

Patterns mais relevantes:

- Repository
- Adapter

### `auth-service`

Patterns mais relevantes:

- Adapter
- Strategy

## Recapitulação

Em Go, patterns devem ser usados como vocabulário, e não como decoração. Neste repositório, os patterns mais valiosos são os que fortalecem fronteiras: Adapter, Strategy, Facade e Repository. Eles ajudam o código a continuar testável, componível e mais fácil de evoluir sem empurrar o projeto para abstrações desnecessárias.

## Próxima Leitura

Se você quiser conectar esses patterns de forma mais direta com princípios de design, continue em:

- [SOLID-AND-PATTERNS.md](./SOLID-AND-PATTERNS.md)
