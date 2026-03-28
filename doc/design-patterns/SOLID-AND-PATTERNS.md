# SOLID E Padrões

## O Que Vamos Ver

- Como design patterns podem ajudar a aplicar SOLID na prática.
- Quais patterns são mais úteis para cada princípio do SOLID.
- Como isso se conecta com Go e com este repositório.
- Onde patterns ajudam e onde viram cerimônia desnecessária.

## Introdução

SOLID nos dá princípios de design. Design patterns nos dão estratégias recorrentes de implementação. Eles não são a mesma coisa, mas funcionam bem juntos.

Pensando de forma simples:

- SOLID nos diz que tipo de pressão de design devemos respeitar.
- patterns frequentemente sugerem uma forma possível de resolver essa pressão.

Então, quando dizemos "usar patterns para resolver SOLID", isso não significa que cada princípio tem um pattern obrigatório. Significa que alguns patterns são especialmente bons em apoiar princípios específicos.

## Princípio Da Responsabilidade Única

### A Pressão De Design

Uma unidade de código deve ter uma razão clara para mudar.

### Patterns Que Frequentemente Ajudam

#### Adapter

Adapter ajuda a separar preocupações de transporte e infraestrutura das preocupações de aplicação e domínio.

Por que ajuda:

- o parse de HTTP fica no adapter
- detalhes de banco ficam no adapter
- a tradução de contratos de serviços externos fica no adapter

Isso evita que um componente mude por múltiplas preocupações não relacionadas.

#### Facade

Facade pode ajudar quando uma superfície voltada ao cliente deve continuar simples enquanto a coordenação acontece por trás.

Por que ajuda:

- clientes dependem de uma interface menor
- a orquestração pode ficar em um lugar dedicado em vez de vazar para todo lado

### Neste Repositório

Bons exemplos:

- HTTP handlers as adapters
- repository clients as adapters
- `mobile-bff` behaving as a facade for clients

## Princípio Aberto/Fechado

### A Pressão De Design

Queremos estender comportamento sem reescrever código estável a todo momento.

### Patterns Que Frequentemente Ajudam

#### Strategy

Strategy é uma das formas práticas mais claras de apoiar OCP.

Por que ajuda:

- o comportamento pode variar atrás de uma interface
- novas implementações podem ser adicionadas sem modificar quem chama

#### Factory Method

Factory Method ajuda quando a escolha da implementação precisa mudar conforme configuração ou ambiente.

Por que ajuda:

- a construção do objeto muda sem afetar o fluxo de negócio

#### Abstract Factory

Útil quando famílias de dependências relacionadas precisam ser criadas em conjunto.

Por que ajuda:

- a construção em grupo permanece consistente
- quem chama não precisa conhecer tipos concretos

### Neste Repositório

Bons exemplos:

- repository ports with swappable implementations
- auth provider abstraction
- runtime wiring in `main` selecting concrete adapters

## Princípio Da Substituição De Liskov

### A Pressão De Design

Uma implementação deve poder substituir com segurança a abstração que ela afirma implementar.

### Patterns Que Frequentemente Ajudam

#### Strategy

Strategy só funciona bem se as implementações forem realmente substituíveis. Então ela naturalmente nos força a pensar em LSP.

Por que ajuda:

- contratos de interface ficam explícitos
- implementações precisam preservar o comportamento esperado

#### Adapter

Adapter ajuda porque mantém modelos externos incompatíveis fora do core da aplicação.

Por que ajuda:

- em vez de forçar o domínio a aceitar semânticas estranhas, o adapter as traduz primeiro

### Neste Repositório

Sempre que uma implementação de repository ou provider é injetada por um port, LSP importa. Se uma implementação se comportar de forma surpreendente, o caso de uso fica frágil.

## Princípio Da Segregação De Interface

### A Pressão De Design

Clientes não devem depender de métodos que não precisam.

### Patterns Que Frequentemente Ajudam

#### Strategy

Strategies menores costumam ser mais fáceis de trocar e entender do que interfaces gigantes e multipropósito.

#### Repository

O design de repository fica mais limpo quando as interfaces estão focadas em uma capacidade real em vez de virarem um balde genérico de persistência.

### Neste Repositório

É por isso que ports focados importam mais do que interfaces gigantes de service. Um port estreito costuma nos dar testes mais limpos, responsabilidade mais clara e menos dependências acidentais.

## Princípio Da Inversão De Dependência

### A Pressão De Design

O fluxo de negócio de alto nível deve depender de abstrações, e não de detalhes técnicos concretos.

### Patterns Que Frequentemente Ajudam

#### Adapter

Adapter é uma das formas mais diretas de apoiar DIP.

Por que ajuda:

- o core depende de um port
- o adapter implementa esse port e esconde detalhes técnicos

#### Strategy

Strategy apoia DIP porque quem chama depende da interface, e não da implementação concreta.

#### Repository

Repository é uma fronteira orientada a DDD que apoia fortemente DIP para persistência.

### Neste Repositório

Este princípio é central para a arquitetura atual:

- handlers depend on use cases
- services depend on ports
- adapters implement ports

Esse é exatamente o tipo de direção de dependência que DIP tenta proteger.

## Um Mapeamento Prático

### SOLID -> Patterns

- SRP -> Adapter, Facade
- OCP -> Strategy, Factory Method, Abstract Factory
- LSP -> Strategy, Adapter
- ISP -> Strategy, Repository
- DIP -> Adapter, Strategy, Repository

Isto não é uma fórmula rígida. É um guia prático.

## Aviso Importante

Patterns podem apoiar SOLID, mas também podem violá-lo quando usados em excesso.

Exemplos:

- interfaces demais podem violar SRP e ISP ao criar ruído em vez de clareza
- uma Strategy falsa com apenas uma implementação pode adicionar cerimônia de OCP sem valor real
- uma Facade que continua absorvendo regras de negócio pode virar um depósito de tudo

Então a pergunta certa nunca é:

"Which pattern can I force here?"

A pergunta certa é:

"Which design pressure am I dealing with, and does a known pattern help me handle it clearly?"

## Recapitulação

SOLID e design patterns se reforçam quando usados com intenção. SOLID nos diz que tipo de qualidade de design queremos. Patterns oferecem formas repetíveis de caminhar nessa direção. Neste projeto, Adapter, Strategy, Facade e Repository são as pontes práticas mais fortes entre esses princípios e o código Go do dia a dia.
