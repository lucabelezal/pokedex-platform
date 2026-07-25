# Fase 05 — Testes, Web & Deploy

**Aulas do curso:** 85 a 98

Fase final. Você vai conectar testes, banco de dados, HTTP e as ferramentas do ecossistema Go.
Ao final, você terá a visão completa de como a Pokedex Platform funciona de ponta a ponta.

## Sumário

| Passo | Recurso | Tempo est. |
|-------|---------|-----------|
| 1 | Assistir aulas do curso | ~3h |
| 2 | Executar exemplos Go by Example | ~40min |
| 3 | Ler Effective Go | ~25min |
| 4 | Conferir roadmap.sh | ~5min |
| 5 | Ler styleguide | ~25min |
| 6 | Fazer exercício prático | ~30min |

---

## 1. Aulas do curso

| # | Aula | Conceito |
|---|------|----------|
| 85 | Teste Unitário Básico | `func TestXxx(t *testing.T)`, `go test` |
| 86 | Criando Dataset para os Testes | Dados de teste, table-driven implícito |
| 87 | Tipo de Arquitetura e os Testes | Testabilidade e design |
| 88 | Gerando Relatório de Cobertura de Testes | `go test -cover`, `go test -coverprofile` |
| 89 | Orientações para Instalação do MySQL | Setup de banco |
| 90 | Criando o Schema e a Tabela | DDL |
| 91 | Executando Inserts | `INSERT` |
| 92 | Executando Inserts em uma Transação | `BEGIN`, `COMMIT`, `ROLLBACK` |
| 93 | Executando Update e Delete | `UPDATE`, `DELETE` |
| 94 | Executando Select e Mapeando p/ Struct | `SELECT`, `rows.Scan` |
| 95 | Criando um Servidor Estático | `http.FileServer` |
| 96 | Gerando Conteúdo Dinâmico | `http.HandleFunc`, templates |
| 97 | Integrando Http e SQL (2 Serviços REST) | API REST completa com banco |
| 98 | Obrigado e Até Breve | Encerramento |

---

## 2. Go by Example — exemplos desta fase

| Exemplo | Link | O que observar |
|---------|------|---------------|
| Testing and Benchmarking | https://gobyexample.com/testing-and-benchmarking | `func TestXxx`, `func BenchmarkXxx`, table-driven |
| HTTP Client | https://gobyexample.com/http-client | `http.Get`, `http.Post`, `json.NewDecoder` |
| HTTP Server | https://gobyexample.com/http-server | `http.HandleFunc`, `http.ListenAndServe` |
| TCP Server | https://gobyexample.com/tcp-server | `net.Listen`, `net.Conn` — base do HTTP |
| Sorting | https://gobyexample.com/sorting | `sort.Ints`, `sort.Strings`, `sort.Slice` |
| Sorting by Functions | https://gobyexample.com/sorting-by-functions | `sort.Slice` com função de comparação customizada |
| Panic | https://gobyexample.com/panic | Revisão: `panic()` interrompe a goroutine |
| Defer | https://gobyexample.com/defer | Revisão: LIFO, limpeza de recursos |
| Recover | https://gobyexample.com/recover | Revisão: `recover()` dentro de `defer` |
| Context | https://gobyexample.com/context | Revisão: cancelamento, timeout, valores |
| Command-Line Arguments | https://gobyexample.com/command-line-arguments | `os.Args` — argumentos da linha de comando |
| Command-Line Flags | https://gobyexample.com/command-line-flags | `flag.String`, `flag.Int`, `flag.Parse` |
| Command-Line Subcommands | https://gobyexample.com/command-line-subcommands | `flag.NewFlagSet` |
| Environment Variables | https://gobyexample.com/environment-variables | `os.Getenv`, `os.Setenv`, `os.Environ` |
| Logging | https://gobyexample.com/logging | `log.Println`, `log.Fatalln`, `log.Panicln` |
| Reading Files | https://gobyexample.com/reading-files | `os.ReadFile`, `os.Open`, `bufio.Scanner` |
| Writing Files | https://gobyexample.com/writing-files | `os.WriteFile`, `os.Create` |
| Line Filters | https://gobyexample.com/line-filters | `bufio.Scanner` para processar stdin linha a linha |
| File Paths | https://gobyexample.com/file-paths | `filepath.Base`, `Dir`, `Join`, `Ext` |
| Directories | https://gobyexample.com/directories | `os.MkdirAll`, `os.ReadDir` |
| Temporary Files/Dirs | https://gobyexample.com/temporary-files-and-directories | `os.CreateTemp`, `os.MkdirTemp` |
| Spawning Processes | https://gobyexample.com/spawning-processes | `exec.Command` — executar comandos externos |
| Exec'ing Processes | https://gobyexample.com/execing-processes | `syscall.Exec` — substituir o processo atual |
| Signals | https://gobyexample.com/signals | `os.Signal` — SIGINT, SIGTERM, graceful shutdown |
| Exit | https://gobyexample.com/exit | `os.Exit(code)` — saída com código de status |

---

## 3. Effective Go — seções para ler

| Seção | Link | O que aprender |
|-------|------|---------------|
| **Errors** | https://go.dev/doc/effective_go#errors | Revisão: `error` interface, type assertion em erros, `PathError` |
| **Panic** | https://go.dev/doc/effective_go#panic | `panic()` para erros irrecuperáveis |
| **Recover** | https://go.dev/doc/effective_go#recover | `recover()` em `defer`, graceful degradation |
| **A web server** | https://go.dev/doc/effective_go#web_server | Exemplo completo: servidor HTTP + templates + QR code |

---

## 4. Roadmap.sh — tópicos desta fase

| Categoria | Tópicos |
|-----------|---------|
| Testing & Benchmarking | `testing` package basics, Table-driven Tests, Mocks and Stubs, `httptest` for HTTP Tests, Benchmarks, Coverage |
| Web Development | `net/http` (standard), Frameworks: `gin`, `echo`, `fiber`, `beego`, gRPC & Protocol Buffers |
| ORMs & DB Access | `pgx`, `GORM` |
| Logging | `slog` (Go 1.21+), Zerolog, Zap |
| Core Go Commands | `go run`, `go build`, `go install`, `go fmt`, `go mod`, `go test`, `go clean`, `go doc`, `go version` |
| Code Quality | `go vet`, `goimports`, Linters (`revive`, `staticcheck`, `golangci-lint`) |
| Security | `govulncheck` |
| Performance | `pprof`, `trace`, Race Detector |
| Deployment | Cross-compilation, Docker, Kubernetes |
| Standard Library | `os`, `bufio`, `slog`, `regexp`, `go:embed` |
| CLI | `cobra`, `urfave/cli`, `bubbletea` |

---

## 5. Guia de Estilo — regras desta fase

> Leia [`best-practices.md`](../guia-estilo/best-practices.md) e [`regras-projeto.md`](../guia-estilo/regras-projeto.md) completos.

### best-practices.md

| Regra | Seção | Por que importa agora |
|-------|-------|----------------------|
| Table-driven tests com `t.Run` | [Testes baseados em tabela](../guia-estilo/best-practices.md#testes-baseados-em-tabela) | Padrão do projeto. `tests` + `tt` + `give`/`want` |
| Evitar complexidade em tabelas | [Evitar complexidade](../guia-estilo/best-practices.md#evitar-complexidade-desnecessária-em-tabelas) | Sem `shouldErr`, `shouldCallX` — separe em funções |
| `tt := tt` para paralelismo | [Testes paralelos](../guia-estilo/best-practices.md#testes-paralelos) | Necessário para `t.Parallel()` |
| Cobertura: 75% mín, 90% ideal | [Cobertura](../guia-estilo/best-practices.md#cobertura) | `make coverage` para relatório |
| Opções funcionais | [Opções funcionais](../guia-estilo/best-practices.md#opções-funcionais) | Padrão para construtores com 3+ argumentos |
| Usar `defer` para limpeza | [defer para limpeza](../guia-estilo/best-practices.md#defer-para-limpeza) | Recursos sempre liberados |
| Evitar `init()` | [Evitar init()](../guia-estilo/best-practices.md#evitar-init) | Prefira inicialização explícita |
| Sair em `main` | [Sair em Main](../guia-estilo/best-practices.md#sair-em-main) | Apenas `main()` ou `run()` chamada por `main` |
| Tags de struct | [Tags de campos](../guia-estilo/best-practices.md#tags-de-campos-em-structs-serializadas) | `json:"id"` com espaços |
| Evitar embedding público | [Evitar incorporação](../guia-estilo/best-practices.md#evitar-incorporação-de-tipos-em-structs-públicas) | Não incorpore `sync.Mutex` em structs públicas |

### regras-projeto.md

| Regra | Seção | Por que importa agora |
|-------|-------|----------------------|
| Handlers dependem de use cases | [Regras arquiteturais](../guia-estilo/regras-projeto.md#regras-arquiteturais) | Nunca injete client HTTP direto no handler |
| Services dependem de ports outbound | [Regras arquiteturais](../guia-estilo/regras-projeto.md#regras-arquiteturais-2) | Injeção manual no `main.go` |
| Erros normalizados no adapter/service | [Regras arquiteturais](../guia-estilo/regras-projeto.md#regras-arquiteturais-3) | Handler mapeia erro de domínio → HTTP status |
| DTOs separados de entidades | [DTOs e serialização](../guia-estilo/regras-projeto.md#dtos-e-serialização-json) | `dto/` para resposta HTTP, `domain/` para lógica |
| `context.Context` passado, nunca armazenado | [Regras arquiteturais](../guia-estilo/regras-projeto.md#regras-arquiteturais-6) | Nunca `ctx` em campo de struct |
| Commits: Conventional Commits em pt-BR | [Commits](../guia-estilo/regras-projeto.md#commits) | `feat(bff): adicionar ...` |
| `go vet` + `golangci-lint` | [CI/CD](../guia-estilo/regras-projeto.md#ci-cd) | Build, test, vet, lint |
| Config via variáveis de ambiente | [Configuração](../guia-estilo/regras-projeto.md#configuração) | `LoadConfig()` — nada hardcoded |
| Segurança: cookies, rate limiting, JWT | [Segurança](../guia-estilo/regras-projeto.md#segurança) | HttpOnly, Secure, SameSite, MaxBytesReader |
| Observabilidade: slog + OTel + Prometheus | [Observabilidade](../guia-estilo/regras-projeto.md#observabilidade) | Logs estruturados, tracing, métricas |
| Spec-driven development | [Desenvolvimento spec-driven](../guia-estilo/regras-projeto.md#desenvolvimento-spec-driven) | Specify → Design → Tasks → Execute |

---

## 6. No código do projeto

### Table-driven tests reais

```go
// Em tests/unit/service_test.go:
func TestPokemonService_List(t *testing.T) {
    tests := []struct {
        name    string
        params  domain.SearchParams
        wantLen int
        wantErr bool
    }{
        {name: "lista todos", params: domain.SearchParams{}, wantLen: 5, wantErr: false},
        {name: "filtra por tipo", params: domain.SearchParams{Type: "Fire"}, wantLen: 1, wantErr: false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            svc := service.NewPokemonService(mockRepo)
            page, err := svc.List(context.Background(), tt.params)
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            assert.Len(t, page.Items, tt.wantLen)
        })
    }
}
```

### Entry point com wiring manual

```go
// Em cmd/server/main.go:
func main() {
    logger.Setup("mobile-bff")
    cfg := config.LoadConfig()

    // Ports outbound — implementações concretas
    pokemonRepo := httpclient.NewPokemonCatalogServiceRepository(cfg.PokemonCatalogSvcURL)
    var favoriteRepo outbound.FavoriteRepository  // fallback: postgres ou memory

    // Services — implementam inbound ports
    pokemonSvc := service.NewPokemonService(pokemonRepo)
    favoriteSvc := service.NewFavoriteService(favoriteRepo)
    authSvc := service.NewAuthService(authProvider)

    // Handler — depende apenas de inbound ports
    h := httphandler.NewHandler(pokemonSvc, favoriteSvc, authSvc)

    // Server com timeouts
    srv := &http.Server{
        Addr:         ":" + cfg.Port,
        Handler:      h.RegisterRoutes(mux),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
}
```

### Estrutura completa de diretórios

```
core/bff/mobile-bff/
├── cmd/server/main.go           ← Entry point, wiring
├── internal/
│   ├── domain/                  ← Entidades, value objects, erros
│   ├── ports/
│   │   ├── inbound/             ← Use case interfaces
│   │   └── outbound/            ← Repository/client interfaces
│   ├── service/                 ← Implementação dos use cases
│   ├── adapters/
│   │   ├── inbound/http/        ← Handlers, DTOs, middleware
│   │   └── outbound/            ← HTTP clients, PostgreSQL, memory
│   ├── config/                  ← LoadConfig()
│   └── infrastructure/          ← Logger, tracing, metrics
├── migrations/                  ← SQL migrations
├── tests/                       ← Testes (unit, integration, mocks)
└── Makefile                     ← test, coverage, build, lint
```

---

## 7. Exercício prático

**Objetivo:** Rodar a suíte completa de testes e verificar a cobertura.

1. Navegue até `core/bff/mobile-bff/` e execute:

```bash
make test          # roda todos os testes unitários
make coverage      # gera relatório de cobertura
make lint          # roda golangci-lint
```

2. Veja o relatório de cobertura em `coverage.html`

3. Encontre um arquivo com cobertura abaixo de 75%:

```bash
go tool cover -func=coverage.out | grep -E "coverage: [0-6][0-9]\.[0-9]%"
```

4. Escolha uma função não coberta e escreva um teste table-driven para ela:

```go
func TestSuaFuncao(t *testing.T) {
    tests := []struct {
        name    string
        give    string
        want    string
        wantErr bool
    }{
        {name: "caso feliz", give: "input", want: "expected", wantErr: false},
        {name: "caso erro", give: "invalid", want: "", wantErr: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := SuaFuncao(tt.give)
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

5. Execute `make coverage` novamente e veja a cobertura subir

6. **Desafio extra:** Execute `go test -race ./...` para verificar data races. Se encontrar alguma, investigue com `go test -race -v -run=NomeDoTeste`.

---

## Parabéns!

Você completou as 5 fases do roteiro de aprendizado. Agora você:

- Sabe escrever Go idiomático (Effective Go)
- Conhece os patterns e idioms (Go by Example)
- Tem a visão completa dos tópicos (Roadmap.sh)
- Aplica as regras de estilo do projeto (Guia de Estilo)
- Entende a arquitetura hexagonal da Pokedex Platform (Regras do Projeto)

**Próximos passos:**

- Leia `doc/DECISIONS.md` para decisões arquiteturais
- Leia `doc/SYSTEM-OVERVIEW.md` para visão do sistema
- Explore `doc/architecture/` para diagramas
- Contribua com uma feature seguindo o spec-driven flow em `.specs/`

---

[Voltar ao ROTEIRO](ROTEIRO.md)
