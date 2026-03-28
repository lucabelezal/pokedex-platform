# Arquitetura Hexagonal 01

## O Que Vamos Ver

- Por que MVC muitas vezes começa bem e fica barulhento com o tempo.
- Por que muitos times caminham para arquitetura hexagonal, DDD e fronteiras inspiradas em SOLID.
- Por que essa conversa costuma gerar resistência entre desenvolvedores.
- Por que uma versão mais leve da abordagem costuma ser mais útil do que uma versão dogmática.

## Introdução

Quando um projeto é pequeno, MVC ou estruturas simples por feature costumam parecer suficientes. Um controller recebe uma requisição, chama um service, lê ou escreve por meio de um repository e devolve uma resposta. Isso costuma ser produtivo no começo.

O problema começa quando a aplicação cresce e o time cresce junto com ela. Nesse ponto, a questão normalmente não é que MVC esteja "errado". A questão é que regras de negócio, preocupações de transporte, detalhes de persistência, APIs externas, lógica de validação e mapeamento de resposta começam a viver perto demais umas das outras.

Esse é o momento em que arquitetura deixa de ser um assunto teórico e passa a ser uma ferramenta de coordenação.

## De MVC Para Fronteiras Melhores

Na prática, muitas empresas não se afastam de MVC por moda. Elas mudam porque a base de código começa a fazer perguntas difíceis:

- Onde uma nova regra deve ficar?
- Um DTO HTTP pode ser reutilizado dentro do domínio?
- Um repository pode devolver um modelo moldado para resposta de API?
- Um handler pode chamar infraestrutura diretamente?
- O que é seguro mudar sem quebrar comportamentos não relacionados?

Sem regras claras, o time responde essas perguntas de forma diferente a cada sprint. O resultado é entropia.

Arquitetura hexagonal, especialmente em Go, é útil porque dá uma resposta menor e mais clara do que muitas arquiteturas pesadas:

- manter os conceitos de negócio no centro
- definir ports como contratos
- implementar adapters para HTTP, banco, filas e serviços externos
- manter as dependências apontando para dentro

É por isso que "ports and adapters" costuma ser um modelo mental melhor do que a palavra "hexagonal" por si só.

## Por Que SOLID E DDD Aparecem Nessa Conversa

Mesmo quando um projeto em Go não é fortemente orientado a objetos, as ideias por trás de SOLID ainda ajudam.

### Responsabilidade Única

Um handler deve mudar porque um contrato HTTP mudou, e não porque uma tabela do banco mudou.

### Aberto/Fechado

Um caso de uso deve ser extensível por meio da conexão de um novo adapter, e não pela reescrita do fluxo de negócio sempre que uma integração mudar.

### Segregação De Interface

Ports pequenos e focados costumam ser mais fáceis de entender do que interfaces gigantes do tipo "manager".

### Inversão De Dependência

O fluxo de negócio deve depender de abstrações que descrevem o que ele precisa, e não de clients técnicos concretos.

DDD também aparece aqui por uma razão parecida: ele nos lembra que a linguagem do negócio importa. Uma base de código fica mais fácil de evoluir quando nomeia seus conceitos pela responsabilidade e não pela tecnologia.

Essa é uma das razões pelas quais `pokemon-catalog-service` é um nome melhor do que o antigo e genérico `pokedex-service`. O nome novo diz o que o serviço realmente possui.

## Por Que Desenvolvedores Resistem A Isso

A resistência é compreensível.

Muitos desenvolvedores já viram discussões de arquitetura produzirem:

- pastas demais
- interfaces demais
- boilerplate demais
- pouco valor de negócio
- debates intermináveis sem entrega

Essa crítica muitas vezes é válida. Muitas bases ditas de "clean architecture" são apenas projetos CRUD com camadas demais e mais cerimônia do que clareza.

A resposta não é rejeitar arquitetura por completo. A resposta é usar arquitetura na medida certa para criar fronteiras mais fortes sem transformar o projeto em um labirinto.

## A Posição Prática Deste Projeto

Este projeto não precisa de uma arquitetura maximalista. Ele se beneficia, sim, de uma arquitetura organizada.

Por quê?

- existe mais de um serviço
- o BFF compõe múltiplas dependências
- a plataforma já tem preocupações de gateway, infraestrutura, auth e catálogo
- o time precisa de regras compartilhadas sobre onde o código deve entrar

Então o objetivo aqui não é "complexidade enterprise". O objetivo é crescimento organizado.

## Por Que Uma Abordagem Leve Faz Sentido

Uma abordagem hexagonal leve em Go normalmente significa:

- um pacote de domínio pequeno
- ports focados
- uma camada de service que implementa casos de uso
- adapters para HTTP e infraestrutura
- `main` como ponto de composição

Ela não exige:

- dezenas de pacotes aninhados
- interfaces para toda struct
- factories genéricas para tudo
- camadas cerimoniais sem responsabilidade clara

Esse equilíbrio é o ponto central. Arquitetura deve remover complexidade acidental, e não adicionar complexidade teatral.

## Recapitulação

MVC não é o inimigo. O problema real são fronteiras fracas em um sistema que está crescendo. Arquitetura hexagonal se torna atraente quando o time precisa de uma resposta mais clara sobre onde vivem as regras de negócio, como as dependências fluem e como evoluir a aplicação sem acoplar tudo a HTTP ou persistência.
