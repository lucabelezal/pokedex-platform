# Notas De Arquitetura

Esta seção documenta o pensamento arquitetural por trás do projeto em um formato mais leve, próximo de um texto de blog.

## O Que Vamos Ver

- Por que times saem de uma organização em estilo MVC para uma visão de ports and adapters.
- Como SOLID, DDD e arquitetura hexagonal ajudam a criar fronteiras mais claras.
- Como essas ideias se aplicam a este repositório hoje.
- Quais trade-offs, pontos de resistência e decisões práticas vêm junto com esse estilo.

## Ordem De Leitura

1. [HEXAGONAL-ARCHITECTURE-01.md](./HEXAGONAL-ARCHITECTURE-01.md)
2. [HEXAGONAL-ARCHITECTURE-02.md](./HEXAGONAL-ARCHITECTURE-02.md)
3. [HEXAGONAL-ARCHITECTURE-03.md](./HEXAGONAL-ARCHITECTURE-03.md)

## Complemento Técnico

Para um mapeamento técnico mais direto das camadas do BFF, veja:

- [../bff/HEXAGONAL-ARCHITECTURE.md](../bff/HEXAGONAL-ARCHITECTURE.md)

## Referências Que Basearam Esta Seção

- https://www.freecodecamp.org/news/solid-principles-explained-in-plain-english/
- https://johnfercher.medium.com/go-arquitetura-hexagonal-dbcd2e986b55
- https://skoredin.pro/blog/golang/hexagonal-architecture-go
- https://blog.masteringbackend.com/software-architecture-with-golang
- https://prabogo.com/docs/architecture.html

## Recapitulação

O objetivo desta seção não é defender arquitetura como dogma. O objetivo é explicar por que este projeto usa uma estrutura mais intencional, quais benefícios esperamos dela e em quais pontos estamos deliberadamente mantendo a abordagem pragmática.
