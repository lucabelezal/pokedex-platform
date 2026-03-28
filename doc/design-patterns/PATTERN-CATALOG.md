# Catálogo De Padrões

## O Que Vamos Ver

- Um catálogo de design patterns comuns em Go.
- Links para exemplos canônicos.
- Uma indicação rápida de quando cada pattern é útil.

## Fonte De Referência

Este catálogo é baseado nos exemplos em Go do Refactoring Guru:

- https://refactoring.guru/design-patterns/go

## Patterns Criacionais

### Abstract Factory

- Link: https://refactoring.guru/design-patterns/abstract-factory/go/example
- Use quando: você precisa criar famílias de objetos relacionados sem acoplar o código a implementações concretas.

### Builder

- Link: https://refactoring.guru/design-patterns/builder/go/example
- Use quando: a construção do objeto tem muitas etapas ou múltiplas representações.

### Factory Method

- Link: https://refactoring.guru/design-patterns/factory-method/go/example
- Use quando: a criação do objeto varia por configuração ou subtipo.

### Prototype

- Link: https://refactoring.guru/design-patterns/prototype/go/example
- Use quando: clonar objetos existentes é mais fácil do que reconstruí-los.

### Singleton

- Link: https://refactoring.guru/design-patterns/singleton/go/example
- Use quando: uma única instância compartilhada é realmente necessária.

Atenção:

- em Go, estado global frequentemente causa mais mal do que bem

## Patterns Estruturais

### Adapter

- Link: https://refactoring.guru/design-patterns/adapter/go/example
- Use quando: duas partes do sistema precisam colaborar por meio de interfaces incompatíveis.

### Bridge

- Link: https://refactoring.guru/design-patterns/bridge/go/example
- Use quando: abstração e implementação precisam evoluir de forma independente.

### Composite

- Link: https://refactoring.guru/design-patterns/composite/go/example
- Use quando: você precisa lidar de forma uniforme com estruturas de objetos em árvore.

### Decorator

- Link: https://refactoring.guru/design-patterns/decorator/go/example
- Use quando: comportamento deve ser adicionado dinamicamente por meio de wrappers.

### Facade

- Link: https://refactoring.guru/design-patterns/facade/go/example
- Use quando: é necessária uma interface mais simples sobre um subsistema mais complexo.

### Flyweight

- Link: https://refactoring.guru/design-patterns/flyweight/go/example
- Use quando: muitos objetos podem compartilhar estado para reduzir uso de memória.

### Proxy

- Link: https://refactoring.guru/design-patterns/proxy/go/example
- Use quando: o acesso a um objeto real precisa ser controlado, adiado, armazenado em cache ou protegido.

## Patterns Comportamentais

### Chain of Responsibility

- Link: https://refactoring.guru/design-patterns/chain-of-responsibility/go/example
- Use quando: uma requisição pode passar por múltiplos handlers até que um a processe.

### Command

- Link: https://refactoring.guru/design-patterns/command/go/example
- Use quando: requisições precisam ser encapsuladas como objetos, enfileiradas, reexecutadas ou desfeitas.

### Iterator

- Link: https://refactoring.guru/design-patterns/iterator/go/example
- Use quando: você quer percorrer coleções sem expor detalhes da estrutura.

### Mediator

- Link: https://refactoring.guru/design-patterns/mediator/go/example
- Use quando: muitos componentes precisam de coordenação indireta sem acoplamento direto.

### Memento

- Link: https://refactoring.guru/design-patterns/memento/go/example
- Use quando: o estado de um objeto precisa ser salvo e restaurado com segurança.

### Observer

- Link: https://refactoring.guru/design-patterns/observer/go/example
- Use quando: muitos listeners reagem a mudanças em um objeto fonte.

### State

- Link: https://refactoring.guru/design-patterns/state/go/example
- Use quando: o comportamento muda dependendo de transições internas de estado.

### Strategy

- Link: https://refactoring.guru/design-patterns/strategy/go/example
- Use quando: múltiplos algoritmos ou implementações precisam ser intercambiáveis.

### Template Method

- Link: https://refactoring.guru/design-patterns/template-method/go/example
- Use quando: o esqueleto de um algoritmo é fixo, mas algumas etapas variam.

## Quais Patterns Mais Importam Aqui

Para este repositório, os patterns mais imediatamente relevantes são:

- Adapter
- Strategy
- Facade
- Factory Method
- Repository

Repository não faz parte do catálogo do Gang of Four, mas é muito relevante para a arquitetura atual.

## Recapitulação

O catálogo completo é útil como vocabulário compartilhado, mas nem todo pattern deve aparecer neste projeto. Os mais valiosos são os que ajudam a isolar infraestrutura, esclarecer fronteiras e manter a lógica de negócio independente de detalhes de transporte e persistência.
