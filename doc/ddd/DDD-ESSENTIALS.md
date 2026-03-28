# Fundamentos De DDD

## O Que Vamos Ver

- O que DDD está tentando atacar.
- Por que complexidade de domínio é diferente de complexidade técnica.
- Quais conceitos são essenciais antes de falar de entities ou aggregates.
- Por que muitos times aplicam DDD de forma errada.

## Introdução

Domain-Driven Design é uma abordagem de desenvolvimento de software que coloca o domínio de negócio no centro das decisões de design.

Isso pode soar abstrato, mas a ideia central é simples:

Quando a parte mais difícil do software é entender o próprio negócio, o código deve ser moldado em torno desse entendimento.

DDD é especialmente útil quando o problema real não é "como eu exponho um endpoint?", mas sim:

- quais regras realmente definem uma operação válida?
- quais conceitos importam neste negócio?
- onde uma área de conhecimento termina e outra começa?
- qual linguagem o código deve usar para continuar alinhado ao negócio?

## Os Três Tipos De Complexidade

Uma forma útil de pensar nisso é:

- complexidade técnica
- complexidade de legado
- complexidade de domínio

DDD ajuda principalmente com a complexidade de domínio.

Ele não resolve magicamente:

- espalhamento de frameworks
- deriva de infraestrutura
- bagunça herdada do legado

Ele ajuda o time a modelar melhor o problema pelo qual o negócio realmente paga.

## Conceitos Essenciais

### Domain

O domínio é a área de conhecimento que o software está tentando apoiar.

Exemplos:

- catalog management
- authentication
- logistics
- payments

### Modelo De Domínio

O modelo de domínio é a representação, no software, dos conceitos, regras e comportamentos de negócio que são relevantes.

Ele não é apenas um diagrama ou um artefato de documentação. Em DDD, o modelo vive no código e evolui com o produto.

### Subdomínios

Subdomínios são recortes significativos do domínio de negócio mais amplo.

Eles ajudam a evitar a pretensão de que o negócio inteiro seja um único modelo coeso.

### Linguagem Ubíqua

Esta é uma das ideias mais importantes de DDD.

O time deve usar uma linguagem compartilhada entre especialistas do negócio e desenvolvedores, e essa linguagem deve moldar o código.

Se o negócio fala em "favorito", o código não deveria chamar isso de `bookmarkRecordManager`.

### Bounded Context

Um bounded context é a fronteira dentro da qual um modelo e uma linguagem permanecem consistentes.

A mesma palavra pode ter significados diferentes em contextos diferentes. DDD não força um modelo universal para tudo.

## Por Que Muitos Times Aplicam DDD De Forma Errada

DDD frequentemente é reduzido a uma lista de artefatos técnicos:

- entity
- value object
- aggregate
- repository
- domain service

Essas coisas importam, mas não são o ponto de partida.

Quando os times pulam linguagem, subdomínios e fronteiras de contexto, frequentemente acabam com:

- modelos de domínio anêmicos
- nomes genéricos
- camadas técnicas fingindo ser "domínio"
- muitos patterns com pouco significado de negócio

Essa é uma das razões pelas quais alguns desenvolvedores acham que DDD é só arquitetura cheia de buzzwords.

## O Que Devemos Manter Em Mente

Se quisermos aplicar DDD bem, estas são as perguntas essenciais:

- qual problema de negócio estamos realmente modelando?
- onde está a complexidade real do domínio?
- quais termos importam para os especialistas?
- quais partes do sistema pertencem a contextos diferentes?
- o que deve permanecer invariável independentemente de transporte ou persistência?

Essas perguntas são mais importantes do que perguntar "onde eu coloco meu repository?".

## Recapitulação

DDD tem mais valor quando ajuda o time a modelar a complexidade do domínio com uma linguagem compartilhada e fronteiras claras. Antes de pensar em entities e aggregates, precisamos entender o domínio, o modelo, os subdomínios e os contextos nos quais essa linguagem é válida.
