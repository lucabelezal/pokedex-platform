# Fase 05 — Testes, Web & Ferramentas

Fase final. Você vai aprender a testar código Go, construir servidores HTTP,
acessar banco de dados e dominar as ferramentas do ecossistema. Ao final, você terá
a visão completa de como uma aplicação Go funciona de ponta a ponta.

---

## Testes

Testes em Go usam o pacote `testing` da biblioteca padrão. Arquivos de teste
terminam com `_test.go` e funções de teste começam com `Test`.

### Teste simples

```go
// arquivo: pokemon_test.go
package pokemon

import "testing"

func TestLevelUp(t *testing.T) {
    p := Pokemon{Level: 25}
    p.LevelUp()

    if p.Level != 26 {
        t.Errorf("esperado 26, obteve %d", p.Level)
    }
}
```

Execute com:

```bash
go test              # pacote atual
go test ./...         # todos os pacotes
go test -v            # verbose
go test -run TestLevelUp  # filtra por nome
```

### Table-driven tests

O padrão idiomático de Go. Cada caso de teste é uma entrada na tabela:

```go
func TestCalculaDano(t *testing.T) {
    tests := []struct {
        name   string
        ataque int
        defesa int
        want   int
    }{
        {"dano normal", 50, 20, 30},
        {"dano mínimo", 10, 30, 1},
        {"ataque igual defesa", 20, 20, 1},
        {"defesa zero", 100, 0, 100},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := CalculaDano(tt.ataque, tt.defesa)
            if got != tt.want {
                t.Errorf("CalculaDano(%d, %d) = %d; want %d",
                    tt.ataque, tt.defesa, got, tt.want)
            }
        })
    }
}
```

### Swift vs Go — testes

```swift
// Swift — XCTest
func testCalculaDano() {
    XCTAssertEqual(calculaDano(50, 20), 30)
    XCTAssertEqual(calculaDano(10, 30), 1)
}
```

```go
// Go — table-driven tests
func TestCalculaDano(t *testing.T) {
    tests := []struct {
        name   string
        ataque int
        defesa int
        want   int
    }{{"normal", 50, 20, 30}, {"mínimo", 10, 30, 1}}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := CalculaDano(tt.ataque, tt.defesa); got != tt.want {
                t.Errorf("= %d; want %d", got, tt.want)
            }
        })
    }
}
```

### Testes paralelos

```go
func TestOperacao(t *testing.T) {
    t.Parallel()  // roda em paralelo com outros testes paralelos

    tests := []struct{ name string }{{"a"}, {"b"}}
    for _, tt := range tests {
        tt := tt  // captura a variável para uso na closure
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel()
            // ...
        })
    }
}
```

**Atenção:** Sempre faça `tt := tt` ao usar `t.Parallel()` dentro de um loop.
Isso garante que cada subteste capture o valor correto da iteração.

### Helpers

```go
func assertEqual[T comparable](t *testing.T, got, want T) {
    t.Helper()  // marca como helper — stack trace aponta o caller, não aqui
    if got != want {
        t.Errorf("got %v; want %v", got, want)
    }
}
```

### `httptest` — testando handlers HTTP

```go
package main

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func handler(w http.ResponseWriter, r *http.Request) {
    nome := r.URL.Query().Get("nome")
    if nome == "" {
        http.Error(w, "nome obrigatório", http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Olá, " + nome))
}

func TestHandler(t *testing.T) {
    tests := []struct {
        name       string
        query      string
        wantStatus int
        wantBody   string
    }{
        {"com nome", "?nome=Pikachu", http.StatusOK, "Olá, Pikachu"},
        {"sem nome", "", http.StatusBadRequest, "nome obrigatório\n"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/"+tt.query, nil)
            rec := httptest.NewRecorder()

            handler(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d; want %d", rec.Code, tt.wantStatus)
            }
            if rec.Body.String() != tt.wantBody {
                t.Errorf("body = %q; want %q", rec.Body.String(), tt.wantBody)
            }
        })
    }
}
```

### Benchmarks

```go
func BenchmarkLevelUp(b *testing.B) {
    p := Pokemon{Level: 1}
    for i := 0; i < b.N; i++ {
        p.LevelUp()
    }
}
```

Execute com:

```bash
go test -bench=.      # roda benchmarks
go test -bench=. -benchmem  # com alocação de memória
```

### Cobertura

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out  # abre no navegador
```

---

## `net/http` — Servidor e Cliente

Go tem um servidor HTTP completo na biblioteca padrão — sem dependências externas.

### Servidor HTTP

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type Pokemon struct {
    Nome  string `json:"nome"`
    Level int    `json:"level"`
}

func pokemonHandler(w http.ResponseWriter, r *http.Request) {
    // só aceita GET
    if r.Method != http.MethodGet {
        http.Error(w, "método não permitido", http.StatusMethodNotAllowed)
        return
    }

    pokemon := Pokemon{Nome: "Pikachu", Level: 25}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(pokemon)
}

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("GET /api/pokemon", pokemonHandler)

    mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status":"ok"}`))
    })

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    log.Println("Servidor rodando em :8080")
    log.Fatal(srv.ListenAndServe())
}
```

No Go 1.22+, as rotas suportam **métodos HTTP** e **path parameters** no `ServeMux` padrão:

```go
mux.HandleFunc("GET /api/pokemon/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    // ...
})
```

### Cliente HTTP

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type Pokemon struct {
    Nome  string `json:"nome"`
    Level int    `json:"level"`
}

func main() {
    client := &http.Client{
        Timeout: 5 * time.Second,
    }

    resp, err := client.Get("http://localhost:8080/api/pokemon")
    if err != nil {
        fmt.Println("Erro:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        fmt.Println("Status:", resp.Status)
        return
    }

    var pokemon Pokemon
    if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
        fmt.Println("Erro ao decodificar:", err)
        return
    }

    fmt.Printf("Pokémon: %s (lvl %d)\n", pokemon.Nome, pokemon.Level)
}
```

### Swift vs Go — HTTP

```swift
// Swift — URLSession
let url = URL(string: "http://localhost:8080/api/pokemon")!
let (data, response) = try await URLSession.shared.data(from: url)
let pokemon = try JSONDecoder().decode(Pokemon.self, from: data)
```

```go
// Go — net/http
resp, err := http.Get("http://localhost:8080/api/pokemon")
defer resp.Body.Close()
var pokemon Pokemon
json.NewDecoder(resp.Body).Decode(&pokemon)
```

---

## Middleware

Middleware é uma função que envolve um handler, adicionando comportamento
(logging, auth, rate limiting):

```go
package main

import (
    "log"
    "net/http"
    "time"
)

// middleware de logging
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        inicio := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(inicio))
    })
}

// middleware de autenticação
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "não autorizado", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// middleware de recovery (panic → 500)
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v", err)
                http.Error(w, "erro interno", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /api/pokemon", pokemonHandler)

    // empilha middlewares (o último é o mais externo)
    handler := recoveryMiddleware(
        loggingMiddleware(
            authMiddleware(mux),
        ),
    )

    http.ListenAndServe(":8080", handler)
}
```

---

## `database/sql` — Acesso a banco de dados

Go acessa bancos relacionais via `database/sql` + um driver específico
(ex: `pgx` para PostgreSQL):

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/lib/pq"  // driver PostgreSQL (import com _ = side-effect import)
)

type Pokemon struct {
    ID    int
    Nome  string
    Level int
}

func main() {
    // abre conexão com o banco
    db, err := sql.Open("postgres",
        "host=localhost user=pokedex password=secret dbname=pokedex sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // configura o pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)

    // verifica se o banco está acessível
    ctx := context.Background()
    if err := db.PingContext(ctx); err != nil {
        log.Fatal("banco indisponível:", err)
    }

    // Query — múltiplas linhas
    rows, err := db.QueryContext(ctx, "SELECT id, nome, level FROM pokemons WHERE level > $1", 10)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    var pokemons []Pokemon
    for rows.Next() {
        var p Pokemon
        if err := rows.Scan(&p.ID, &p.Nome, &p.Level); err != nil {
            log.Fatal(err)
        }
        pokemons = append(pokemons, p)
    }
    // SEMPRE verifique err após o loop
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }

    // QueryRow — única linha
    var p Pokemon
    err = db.QueryRowContext(ctx,
        "SELECT id, nome, level FROM pokemons WHERE id = $1", 25).
        Scan(&p.ID, &p.Nome, &p.Level)
    if err == sql.ErrNoRows {
        fmt.Println("Pokémon não encontrado")
    } else if err != nil {
        log.Fatal(err)
    }

    // Exec — INSERT, UPDATE, DELETE
    result, err := db.ExecContext(ctx,
        "INSERT INTO pokemons (nome, level) VALUES ($1, $2)", "Pikachu", 25)
    if err != nil {
        log.Fatal(err)
    }
    id, _ := result.LastInsertId()
    rowsAffected, _ := result.RowsAffected()
    fmt.Printf("Inserido id=%d, linhas afetadas=%d\n", id, rowsAffected)
}
```

### Transações

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback()  // rollback se não der commit

_, err = tx.ExecContext(ctx, "UPDATE pokemons SET level = level + 1 WHERE id = $1", 25)
if err != nil {
    return fmt.Errorf("update: %w", err)
}

if err := tx.Commit(); err != nil {
    return fmt.Errorf("commit: %w", err)
}
```

---

## Ferramentas do ecossistema

### `go mod` — gerenciamento de dependências

```bash
go mod init github.com/seu-usuario/seu-projeto   # criar módulo
go mod tidy                                        # adiciona dependências usadas, remove não usadas
go get github.com/lib/pq                           # adicionar dependência
go get github.com/lib/pq@v1.10.9                   # versão específica
go mod download                                    # baixar dependências
```

### `go build` — compilação

```bash
go build                          # compila no diretório atual
go build -o bin/app ./cmd/server  # compila com nome e caminho definidos
go build -ldflags="-s -w" ./...   # binário menor (strip debug info)
```

### `go vet` — análise estática

```bash
go vet ./...     # detecta bugs comuns (formatação errada, código inalcançável, etc.)
```

### Cross-compilation

Go compila para qualquer sistema operacional e arquitetura sem dependências externas:

```bash
GOOS=linux GOARCH=amd64 go build        # Linux 64-bit
GOOS=darwin GOARCH=arm64 go build       # macOS Apple Silicon
GOOS=windows GOARCH=amd64 go build      # Windows 64-bit
```

### Graceful shutdown

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    srv := &http.Server{Addr: ":8080"}

    // canal que escuta sinais do SO
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        log.Println("Servidor rodando em :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    <-quit  // bloqueia até receber SIGINT/SIGTERM
    log.Println("Desligando servidor...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Forced shutdown:", err)
    }

    log.Println("Servidor desligado com sucesso")
}
```

### Estrutura de diretórios de um projeto Go típico

```
meu-projeto/
├── cmd/
│   └── server/
│       └── main.go          # entry point
├── internal/                 # código privado do projeto
│   ├── domain/               # entidades, value objects
│   ├── service/              # lógica de negócio
│   ├── ports/                # interfaces (inbound/outbound)
│   └── adapters/             # implementações concretas
├── migrations/               # SQL
├── tests/                    # testes de integração
├── go.mod
├── go.sum
└── Makefile
```

---

## Exercícios da Fase 05

### 1. Teste table-driven

Escreva testes table-driven para a função `CalculaDano(ataque, defesa int) int`
(exercício da Fase 01). Inclua pelo menos 5 casos.

<details>
<summary>Gabarito</summary>

```go
func TestCalculaDano(t *testing.T) {
    tests := []struct {
        name   string
        ataque int
        defesa int
        want   int
    }{
        {"dano normal", 50, 20, 30},
        {"dano mínimo", 10, 30, 1},
        {"ataque igual defesa", 20, 20, 1},
        {"defesa maior", 5, 50, 1},
        {"defesa zero", 100, 0, 100},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := CalculaDano(tt.ataque, tt.defesa)
            if got != tt.want {
                t.Errorf("got %d; want %d", got, tt.want)
            }
        })
    }
}
```
</details>

### 2. Servidor HTTP simples

Crie um servidor HTTP com duas rotas:
- `GET /api/pokemon/{id}` — retorna um Pokémon em JSON
- `GET /api/health` — retorna `{"status":"ok"}`

Adicione pelo menos um middleware (logging ou recovery).

<details>
<summary>Gabarito</summary>

```go
package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"
)

type Pokemon struct {
    ID    string `json:"id"`
    Nome  string `json:"nome"`
    Level int    `json:"level"`
}

var pokedex = map[string]Pokemon{
    "25": {ID: "25", Nome: "Pikachu", Level: 25},
    "6":  {ID: "6", Nome: "Charizard", Level: 36},
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        inicio := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(inicio))
    })
}

func pokemonHandler(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    p, ok := pokedex[id]
    if !ok {
        http.Error(w, "nao encontrado", http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(p)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"status":"ok"}`))
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /api/pokemon/{id}", pokemonHandler)
    mux.HandleFunc("GET /api/health", healthHandler)

    handler := loggingMiddleware(mux)

    log.Println("Servidor em :8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
```
</details>

### 3. Teste HTTP com `httptest`

Escreva um teste para o handler `pokemonHandler` usando `httptest.NewRecorder`.
Teste o caso de Pokémon encontrado e o caso de Pokémon não encontrado.

<details>
<summary>Gabarito</summary>

```go
func TestPokemonHandler(t *testing.T) {
    tests := []struct {
        name       string
        id         string
        wantStatus int
    }{
        {"pokemon existe", "25", http.StatusOK},
        {"pokemon nao existe", "999", http.StatusNotFound},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/", nil)
            req.SetPathValue("id", tt.id)
            rec := httptest.NewRecorder()

            pokemonHandler(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d; want %d", rec.Code, tt.wantStatus)
            }
        })
    }
}
```
</details>

### 4. Graceful shutdown

Modifique o servidor do exercício 2 para fazer graceful shutdown: capture SIGINT/SIGTERM
e aguarde conexões ativas terminarem (com timeout de 10 segundos) antes de sair.

<details>
<summary>Gabarito</summary>

```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /api/pokemon/{id}", pokemonHandler)
    mux.HandleFunc("GET /api/health", healthHandler)

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      loggingMiddleware(mux),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        log.Println("Servidor em :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()

    <-quit
    log.Println("Desligando...")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatal("Shutdown forçado:", err)
    }
    log.Println("Servidor desligado")
}
```
</details>

---

## Parabéns, Treinador!

Você completou as 5 fases do roteiro. Recapitule o que você domina agora:

| Fase | Você sabe... |
|------|-------------|
| 01 | Declarar variáveis, usar ponteiros, escrever funções, usar `fmt` e `iota` |
| 02 | Controlar fluxo com `if`/`for`/`switch`, manipular slices e maps, iterar com `range` |
| 03 | Modelar domínios com structs, definir interfaces implícitas, tratar erros com wrapping, serializar JSON |
| 04 | Lançar goroutines, comunicar com channels, usar `select`, gerenciar concorrência com `sync` e `context` |
| 05 | Escrever testes table-driven, construir servidores HTTP, testar com `httptest`, acessar banco, compilar e fazer deploy |

**Próximos passos:**

- Consulte o [GLOSSARIO.md](GLOSSARIO.md) e o [CHEATSHEET.md](CHEATSHEET.md) como referência rápida
- Explore o código real da Pokedex Platform em `core/bff/mobile-bff/`
- Comece com tarefas simples: corrigir um bug, adicionar um campo a uma struct, escrever um teste
- Contribua com uma feature seguindo o fluxo spec-driven

[Voltar ao ROTEIRO](ROTEIRO.md)
