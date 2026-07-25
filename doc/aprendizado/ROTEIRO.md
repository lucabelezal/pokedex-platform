# Roteiro de Aprendizado — Go para Desenvolvedores iOS/Swift

## Sua missão

Você é desenvolvedor iOS com experiência em Swift e quer aprender Go para atuar no backend.
Este roteiro foi desenhado para você: cada conceito é explicado do zero, com exemplos funcionais,
e fazendo pontes com Swift sempre que isso acelera o entendimento.

Ao final das 7 fases você terá autonomia para ler, escrever e contribuir com código Go idiomático
— o mesmo padrão usado no mobile-bff e nos demais serviços da Pokedex Platform.

## Como usar

Cada fase é um arquivo `.md` autossuficiente. Você não precisa de links externos, cursos ou vídeos.
Leia na ordem. Execute os exemplos no seu terminal. Faça os exercícios.

| Etapa | O que fazer |
|-------|------------|
| 1 | **Ler** a explicação do conceito — entenda o "por quê" |
| 2 | **Executar** os exemplos de código — `go run exemplo.go` |
| 3 | **Ler** a seção "Atenção" — evite as pegadinhas |
| 4 | **Fazer** o exercício prático — consolide o aprendizado |
| 5 | **Consultar** o [glossário](GLOSSARIO.md) e o [cheatsheet](CHEATSHEET.md) como referência |

## O que esperar em cada fase

| Fase | Arquivo | Conceitos | Linha de chegada |
|------|---------|-----------|-----------------|
| 01 | [fase-01-fundamentos.md](fase-01-fundamentos.md) | Tipos, zero values, variáveis, constantes, ponteiros, funções, `fmt`, `iota` | Você escreve e executa programas Go simples |
| 02 | [fase-02-controle-colecoes.md](fase-02-controle-colecoes.md) | `if`/`else`, `for`, `switch`, `defer`, arrays, slices, maps, strings | Você manipula coleções e controla fluxo com idiomática Go |
| 03 | [fase-03-structs-interfaces.md](fase-03-structs-interfaces.md) | Structs, métodos, embedding, interfaces, generics, erros, JSON, pacotes | Você modela domínios com structs e interfaces; trata erros corretamente |
| 04 | [fase-04-concorrencia.md](fase-04-concorrencia.md) | Goroutines, channels, `select`, `sync`, `context` | Você escreve código concorrente seguro e com cancelamento |
| 05 | [fase-05-testes-web.md](fase-05-testes-web.md) | Testes, `net/http`, `database/sql`, toolchain, graceful shutdown | Você testa, serve HTTP, acessa banco e compila para produção |
| 06 | [fase-06-biblioteca-padrao.md](fase-06-biblioteca-padrao.md) | `time`, arquivos, `flag`, `slog`, `sort`, `panic`/`recover`, SHA256/HMAC, regexp, URL, Base64, templates, rate limiting | Você domina os pacotes utilitários do dia a dia backend |
| 07 | [fase-07-testes-avancados-design.md](fase-07-testes-avancados-design.md) | Filosofia de design Go, stubs/spies/mocks/fakes manuais, `testify/mock`, teste de banco, teste de client HTTP, integration tests | Você testa qualquer camada e entende como Go estrutura soluções |

## Pontes Swift → Go

Alguns conceitos são mais fáceis de entender quando você conhece o equivalente em Swift:

| Swift | Go | Nota |
|-------|-----|------|
| `protocol` | `interface` | Em Go a implementação é **implícita** — não se declara conformidade |
| `struct` (value type) | `struct` | Igual — value type, copiado na atribuição |
| `class` (reference type) | Ponteiro para struct (`*T`) | Go não tem classes. Use `*Struct` quando precisar de referência |
| `enum` com associated values | `iota` + `switch type` | Go não tem enums com valores associados; use constantes + type switch |
| `Result<T, Error>` | `(T, error)` | Retorno múltiplo — o padrão Go para success/failure |
| `guard let` / `if let` | `if err != nil` + comma ok | Padrão "happy path" — trate o erro primeiro |
| `async/await` + `Task` | Goroutines + channels + `context` | Concorrência é mais explícita e poderosa em Go |
| `extension` | Funções no mesmo pacote | Não existe `extension`. Adicione funções ao tipo no mesmo pacote |
| ARC (retain/release) | Garbage Collector | Sem `weak`, `unowned`. O GC cuida de tudo |
| `throws` / `try` / `catch` | `error` como valor de retorno | Não há exceções. Erro é um valor como outro qualquer |
| `Optional<T>` | Ponteiro (`*T`) ou comma ok | Go não tem Optional. Use ponteiros ou o segundo valor de retorno |
| `typealias` | `type MeuTipo TipoBase` | Igual — cria um tipo nomeado distinto |

## Referências rápidas

- [GLOSSARIO.md](GLOSSARIO.md) — dicionário de termos Go com equivalente Swift
- [CHEATSHEET.md](CHEATSHEET.md) — sintaxe rápida para consulta diária

## Progresso

| Fase | Status | Data |
|------|--------|------|
| 01 — Fundamentos | ⬜ Pendente | |
| 02 — Controle & Coleções | ⬜ Pendente | |
| 03 — Structs & Interfaces | ⬜ Pendente | |
| 04 — Concorrência | ⬜ Pendente | |
| 05 — Testes & Web | ⬜ Pendente | |
| 06 — Biblioteca Padrão | ⬜ Pendente | |
| 07 — Testes Avançados & Design | ⬜ Pendente | |
