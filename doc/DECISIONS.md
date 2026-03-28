# Decisões

## Objetivo

Este arquivo registra as principais decisões arquiteturais que já estão visíveis na base de código.

## Decisão 1: Usar Um BFF Para A API Voltada Ao Cliente

### Decisão

Expor endpoints orientados ao cliente por meio do `mobile-bff`, em vez de expor diretamente os serviços internos.

### Por Que

- As respostas podem ser moldadas para necessidades de UX.
- Os serviços internos permanecem mais focados.
- A orquestração entre serviços fica fora do cliente.

### Consequência

O BFF se torna uma camada importante de composição e deve evitar crescer como um monólito genérico.

## Decisão 2: Manter O Catálogo Canônico Em `pokemon-catalog-service`

### Decisão

Usar o `pokemon-catalog-service` como fonte canônica de leitura para informações do catálogo de Pokémon.

### Por Que

- As regras do catálogo permanecem centralizadas.
- O BFF pode continuar focado em apresentação e orquestração.

### Consequência

O BFF deve evitar reimplementar regras do catálogo além de formatações específicas de apresentação.

## Decisão 3: Manter Favoritos No Contexto Do BFF Por Enquanto

### Decisão

Os favoritos são atualmente tratados no contexto do BFF, e não em um serviço dedicado.

### Por Que

- Implementação mais simples para o estágio atual do projeto.
- Favoritos estão muito ligados à experiência do usuário autenticado.

### Consequência

Isso é aceitável agora, mas pode se tornar um candidato à extração futura se a lógica de favoritos crescer, exigir escala independente ou precisar ser compartilhada por mais clientes.

## Decisão 4: Usar Arquitetura Hexagonal No BFF

### Decisão

Organizar o BFF em torno de domínio, portas, serviços e adaptadores.

### Por Que

- Melhora a separação de responsabilidades.
- Facilita testes.
- Reduz o acoplamento a detalhes de transporte e persistência.

### Consequência

A base de código deve continuar reforçando a direção das dependências. Adaptadores concretos devem depender de portas, e adaptadores de entrada não devem contornar a camada de aplicação.

## Decisão 5: Preferir Geração Determinística De Seed Em Vez De Setup Dinâmico

### Decisão

Manter arquivos JSON como fonte de dados e gerar seeds SQL de forma determinística.

### Por Que

- Ambientes locais reproduzíveis.
- Facilidade para revisar mudanças nos dados de origem.
- Pipeline claro do conteúdo de origem até a inicialização do banco.

### Consequência

Qualquer mudança no conteúdo do catálogo deve respeitar o fluxo de geração JSON para SQL.
