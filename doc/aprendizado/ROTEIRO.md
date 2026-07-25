# Roteiro de Aprendizado Go — Pokedex Platform

Este roteiro integra **6 recursos** para você aprender Go de forma progressiva, aplicando
as convenções da Pokedex Platform desde o primeiro código.

## Os 6 recursos

| # | Recurso | Tipo | Use para... |
|---|---------|------|------------|
| 1 | **Curso** (98 aulas) | Videoaulas práticas | Aprender sintaxe e conceitos passo a passo |
| 2 | **[Roadmap.sh](https://roadmap.sh/golang)** | Checklist visual | Verificar cobertura e descobrir tópicos avançados |
| 3 | **[Effective Go](https://go.dev/doc/effective_go)** | Documento canônico | Entender os idiomas e a filosofia da linguagem |
| 4 | **[Go by Example](https://gobyexample.com)** | Exemplos anotados | Ver e executar código idiomático de cada conceito |
| 5 | **[golang.com.br/aprenda](https://golang.com.br/aprenda/)** | Hub em português | Tutoriais, trilhas, cheatsheet, comunidade e mercado BR |
| 6 | **[Guia de Estilo](../guia-estilo/README.md)** | Regras do projeto | Saber como escrever código na Pokedex Platform |

## Como usar este roteiro

Para cada fase, siga esta ordem:

| Passo | Ação |
|-------|------|
| 1 | **Assista as aulas do curso** — entenda o conceito |
| 2 | **Execute os exemplos do Go by Example** — veja o código funcionando |
| 3 | **Leia as seções do Effective Go** — entenda os idiomas e o "pensar em Go" |
| 4 | **Confira o roadmap.sh** — veja onde está e o que falta |
| 5 | **Leia as regras do styleguide** — aplique as convenções do projeto |
| 6 | **Faça o exercício prático** — consolide o conhecimento no código real |

## Ordem de leitura do styleguide

| Após a fase... | Leia no styleguide... | Por quê |
|----------------|----------------------|---------|
| 01 — Fundamentos | `guide.md` (completo) | Regras fundamentais: MixedCaps, receivers, initialisms, nomes de pacotes |
| 02 — Controle & Coleções | `decisions.md#estrutura-de-código` + `#declarações-e-inicialização` | Controle de fluxo, slices, maps, structs |
| 03 — Structs & Interfaces | `decisions.md#interfaces` + `#erros` | Interface compliance, wrapping, error flow |
| 04 — Concorrência | `decisions.md#concorrência` | Goroutines, canais, mutex, context |
| 05 — Testes & Web | `best-practices.md` + `regras-projeto.md` | Testes table-driven, arquitetura hexagonal, CI/CD |

## Visão geral das fases

| Fase | Arquivo | Curso | Go by Example | Foco do styleguide |
|------|---------|-------|---------------|-------------------|
| 01 | [fase-01-fundamentos.md](fase-01-fundamentos.md) | Aulas 1–26 | 6 exemplos | `guide.md` |
| 02 | [fase-02-controle-colecoes.md](fase-02-controle-colecoes.md) | Aulas 27–43 | 13 exemplos | `decisions.md` (estrutura + declarações) |
| 03 | [fase-03-structs-interfaces.md](fase-03-structs-interfaces.md) | Aulas 44–71 | 25 exemplos | `decisions.md` (interfaces + erros) |
| 04 | [fase-04-concorrencia.md](fase-04-concorrencia.md) | Aulas 72–84 | 19 exemplos | `decisions.md` (concorrência) |
| 05 | [fase-05-testes-web.md](fase-05-testes-web.md) | Aulas 85–98 | 21 exemplos | `best-practices.md` + `regras-projeto.md` |

## Progresso

| Fase | Status | Data |
|------|--------|------|
| 01 — Fundamentos | ⬜ Pendente | |
| 02 — Controle & Coleções | ⬜ Pendente | |
| 03 — Structs & Interfaces | ⬜ Pendente | |
| 04 — Concorrência | ⬜ Pendente | |
| 05 — Testes & Web | ⬜ Pendente | |

---

## Referências rápidas

- [Tour of Go](https://go.dev/tour/list) — tour interativo oficial
- [Go Playground](https://go.dev/play/) — execute Go no browser
- [Go Spec](https://go.dev/ref/spec) — especificação da linguagem
- [Standard Library](https://pkg.go.dev/std) — documentação dos pacotes padrão
- [Go Wiki — CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments) — checklist de revisão de código
- [golang.com.br/aprenda](https://golang.com.br/aprenda/) — hub da comunidade Go Brasil: tutoriais, trilhas, cheatsheet, vagas e mercado
