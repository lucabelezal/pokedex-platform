# DDD Neste Projeto

## O Que Vamos Ver

- Onde ideias de DDD já aparecem neste repositório.
- Quais partes da plataforma sugerem subdomínios ou bounded contexts diferentes.
- Em quais pontos o modelo atual ainda é leve.
- O que devemos ter em mente se o projeto crescer.

## Visão De Domínio Da Plataforma

Este repositório não é um único modelo de negócio monolítico. Ele está mais próximo de uma plataforma com múltiplas capacidades:

- catálogo de Pokémons
- autenticação
- favoritos e experiência do cliente
- gateway e suporte de infraestrutura

Isso já sugere subdomínios diferentes e, muito provavelmente, bounded contexts diferentes.

## Candidatos Atuais A Contexto

### `pokemon-catalog-service`

Este é o candidato mais claro a bounded context de catálogo.

Por quê:

- ele possui a informação canônica do catálogo
- sua linguagem é sobre dados e recuperação de Pokémons
- ele não deve ser moldado por preocupações de frontend

### `auth-service`

Este é um bounded context claro de autenticação.

Por quê:

- ele possui identidade e ciclo de vida de tokens
- sua linguagem é sobre signup, login, refresh e logout
- ele deve permanecer focado e transversal

### `mobile-bff`

Este está mais para um contexto de experiência de aplicação do que para um domínio canônico de negócio.

Por quê:

- ele compõe comportamento voltado ao cliente
- ele enriquece respostas para consumo da UI
- ele coordena favoritos e a exposição de auth para os clientes

Isso não torna o BFF "menos importante". Apenas significa que seu modelo é orientado à experiência, e não canônico.

## Sinais De DDD Já Presentes

### Nomenclatura Melhor De Serviços

Renomear `pokedex-service` para `pokemon-catalog-service` melhorou a linguagem ubíqua da plataforma.

Por que isso importa:

- o nome agora descreve responsabilidade, e não apenas tema

### Ports Explícitos

Os ports no BFF definem contratos na linguagem da aplicação em vez de acoplamento técnico direto.

### Erros De Domínio

Erros em nível de domínio ajudam a manter o comportamento da aplicação alinhado ao significado de negócio em vez de ficar preso a detalhes brutos de infraestrutura.

## Onde O Projeto Ainda É Leve

A base de código atual não é uma implementação pesada de DDD, e isso está tudo bem.

O que isso significa:

- os modelos ainda são relativamente simples
- os aggregates não estão profundamente modelados
- parte da lógica continua mais orientada à aplicação do que profundamente orientada ao domínio
- o BFF é mais centrado em orquestração do que em domínio rico

Isso não é necessariamente uma fraqueza. Pode simplesmente refletir o nível atual de complexidade do produto.

## O Que Manter Em Mente Daqui Para Frente

### 1. Proteger A Linguagem Ubíqua

Manter a nomenclatura alinhada com a responsabilidade real de cada serviço e caso de uso.

### 2. Evitar Misturar Modelos Canônicos E De Experiência

O contexto de catálogo não deve derivar para formatações específicas de frontend.

### 3. Manter Contextos Claros

As preocupações de auth, catálogo e BFF devem continuar distintas mesmo colaborando de perto.

### 4. Modelar Comportamento Mais Rico Apenas Onde A Complexidade Exigir

Se favoritos, regras de auth ou fluxos de catálogo ficarem mais ricos, então um modelamento mais profundo pode passar a valer a pena.

### 5. Não Forçar Uniformidade De DDD Em Todos Os Serviços

Contextos diferentes podem ter níveis diferentes de complexidade interna.

## Uma Checklist Prática De DDD Para Este Repositório

- O nome do serviço está alinhado com o que ele realmente possui?
- Esta regra faz parte de um modelo canônico de negócio ou apenas de orquestração para o cliente?
- Esta mudança pertence a catálogo, auth ou BFF?
- Estamos usando uma linguagem que o time de negócio ou produto reconheceria?
- Estamos protegendo as invariantes importantes na camada certa?

## Recapitulação

DDD neste projeto deve ser aplicado como uma lente, e não como uma cerimônia rígida. Os resultados mais valiosos são nomes de serviço mais claros, posse mais clara de responsabilidade, linguagem mais clara e uma separação melhor entre contextos canônicos e orquestração orientada à experiência. Isso por si só já é suficiente para melhorar a arquitetura de forma relevante.
