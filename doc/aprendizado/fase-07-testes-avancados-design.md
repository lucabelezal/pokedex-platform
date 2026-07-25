# Fase 07 — Testes Avançados & Filosofia de Design

Esta fase fecha o ciclo. Você vai entender **como Go constrói soluções** — a filosofia
por trás da linguagem — e **como testar cada camada** de uma aplicação backend:
handlers, services, repositories, HTTP clients e banco de dados.

---

# Parte 1 — Filosofia de Design em Go

## Go não é orientado a objetos

Go **não tem classes, herança, nem construtores**. Se você tentar escrever Go como
se fosse Java ou Swift com classes, vai lutar contra a linguagem.

O que Go tem:

| Conceito OOP | Equivalente em Go |
|-------------|-------------------|
| Classe | `struct` + métodos com receiver |
| Herança | Embedding (composição) |
| Polimorfismo | Interfaces (implícitas) |
| Construtor | Função `NewTipo(...)` (convenção) |
| `this` / `self` | Receiver com 1-2 letras (`p`, `s`, `c`) |
| `private` / `public` | minúscula / Maiúscula |
| `abstract class` | Interface + embedding de struct |

Go segue o princípio: **composição sobre herança**. Você não estende tipos — você
compõe comportamentos.

```go
// Não é: class Pokemon extends Entity implements Describable
// É:

type Pokemon struct {
    ID    string
    Nome  string
    Level int
}

// Comportamento via métodos
func (p *Pokemon) LevelUp() { p.Level++ }

// Polimorfismo via interfaces — implementação implícita
type Descritor interface {
    Descreve() string
}

func (p Pokemon) Descreve() string {
    return fmt.Sprintf("%s (lvl %d)", p.Nome, p.Level)
}
// Pokemon implementa Descritor automaticamente. Sem declaração.
```

## Pacotes como unidade de design

Em Go, o **pacote** é a unidade de abstração, não a classe. Um bom pacote:

- Tem **um propósito claro** (`package user`, não `package utils`)
- Expõe **poucos tipos públicos**
- Agrupa **dados + comportamento relacionado**
- É pequeno e focado

```go
// Pacote pokemon — um propósito: gerenciar pokémons
package pokemon

// Tudo que o mundo externo precisa saber sobre Pokémon está aqui.
// Detalhes internos (conexão com banco, cache, etc.) são privados.

type Pokemon struct {
    ID    string
    Nome  string
    Level int
}

type Repository interface {
    FindByID(ctx context.Context, id string) (*Pokemon, error)
    Save(ctx context.Context, p *Pokemon) error
}

type Service struct {
    repo Repository
}

func NewService(repo Repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) LevelUp(ctx context.Context, id string) error {
    p, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return fmt.Errorf("level up: %w", err)
    }
    p.Level++
    return s.repo.Save(ctx, p)
}
```

## Interfaces implícitas = inversão de dependência natural

A grande sacada de Go: interfaces são **implícitas**. Isso significa que:

- Quem **consome** define a interface (não quem implementa)
- Interfaces são **pequenas** (1-3 métodos)
- Você pode criar uma interface para código de terceiros sem modificar o original

Isso produz **inversão de dependência natural** — sem framework, sem annotation,
sem ceremony:

```
Handler → depende de → inbound.PokemonUseCase (interface)
                            ↑ implementada por
                       pokemon.Service (struct)
                            ↓ depende de
                       outbound.PokemonRepository (interface)
                            ↑ implementada por
                       postgres.Repository (struct)
                       http.Client (struct)
                       memory.Repository (struct)
```

```go
// Porta inbound — definida por quem consome (handler)
// arquivo: ports/inbound/pokemon_usecase.go
type PokemonUseCase interface {
    List(ctx context.Context, params SearchParams) (*PokemonPage, error)
    GetByID(ctx context.Context, id string) (*PokemonDetail, error)
}

// Porta outbound — definida por quem consome (service)
// arquivo: ports/outbound/pokemon_repository.go
type PokemonRepository interface {
    GetByID(ctx context.Context, id string) (*Pokemon, error)
    List(ctx context.Context, params SearchParams) ([]Pokemon, error)
}
```

O service **não sabe** se o repository é PostgreSQL, HTTP ou memória.
O handler **não sabe** como o service implementa o use case.
As interfaces são o contrato — o resto é detalhe.

## "Accept interfaces, return structs"

Este é o mantra de design de APIs em Go:

```go
// CORRETO — aceita interface, retorna struct concreta
func NewService(repo Repository) *Service {  // aceita interface
    return &Service{repo: repo}              // retorna struct
}

// ERRADO — retornar interface sem necessidade
func NewService(repo Repository) PokemonUseCase {  // retorna interface
    return &Service{repo: repo}
}
```

Retornar interfaces limita quem chama — ela não pode acessar outros métodos
do tipo concreto. Retorne structs. Quem chama extrai a interface se precisar.

## Small interfaces

As interfaces mais poderosas de Go têm **1 método**:

```go
type Reader interface {
    Read(p []byte) (n int, err error)  // 1 método
}

type Writer interface {
    Write(p []byte) (n int, err error)  // 1 método
}

type Closer interface {
    Close() error  // 1 método
}

// Composição de interfaces pequenas
type ReadCloser interface {
    Reader
    Closer
}
```

Quanto menor a interface, mais tipos a implementam automaticamente.
Quanto mais tipos implementam, mais reutilizável é o código que a consome.

## Functional options pattern

Quando uma struct tem muitos campos configuráveis, use o padrão de opções funcionais:

```go
type ServerConfig struct {
    Port    int
    Timeout time.Duration
    MaxConn int
    TLS     *tls.Config
}

type Option func(*ServerConfig)

func WithPort(port int) Option {
    return func(c *ServerConfig) { c.Port = port }
}

func WithTimeout(d time.Duration) Option {
    return func(c *ServerConfig) { c.Timeout = d }
}

func WithMaxConn(n int) Option {
    return func(c *ServerConfig) { c.MaxConn = n }
}

func WithTLS(cfg *tls.Config) Option {
    return func(c *ServerConfig) { c.TLS = cfg }
}

func NewServer(opts ...Option) *Server {
    cfg := &ServerConfig{
        Port:    8080,         // defaults
        Timeout: 30 * time.Second,
        MaxConn: 100,
    }
    for _, opt := range opts {
        opt(cfg)
    }
    return &Server{config: cfg}
}

// Uso — claro e extensível
srv := NewServer(
    WithPort(9090),
    WithTimeout(60 * time.Second),
    WithTLS(myTLSConfig),
)
```

## Como montar uma solução em Go

O ciclo de design em Go é:

```
1. Modele os dados → structs com campos exportados
2. Defina comportamentos → métodos com receivers
3. Extraia contratos → interfaces pequenas
4. Implemente adapters → structs que satisfazem as interfaces
5. Componha no entry point → wiring manual no main.go
```

Exemplo completo — um serviço de Pokémons:

```go
// PASSO 1: Modele os dados
type Pokemon struct {
    ID    string
    Nome  string
    Level int
}

// PASSO 2: Defina comportamentos
func (p *Pokemon) LevelUp() { p.Level++ }
func (p *Pokemon) IsValid() bool { return p.Nome != "" && p.Level > 0 }

// PASSO 3: Extraia contratos (interfaces)
type PokemonRepository interface {
    FindByID(ctx context.Context, id string) (*Pokemon, error)
    Save(ctx context.Context, p *Pokemon) error
}

// PASSO 4: Implemente adapters
type PostgresRepository struct {
    db *sql.DB
}
func (r *PostgresRepository) FindByID(ctx context.Context, id string) (*Pokemon, error) { /* ... */ }
func (r *PostgresRepository) Save(ctx context.Context, p *Pokemon) error { /* ... */ }

// PASSO 5: Componha no entry point
func main() {
    db, _ := sql.Open("postgres", dsn)
    repo := &PostgresRepository{db: db}    // adapter concreto
    svc := pokemon.NewService(repo)         // service depende da interface
    h := http.NewHandler(svc)               // handler depende da interface
    http.ListenAndServe(":8080", h.Routes())
}
```

Nenhum framework. Nenhuma annotation. Nenhuma injeção de dependência mágica.
**Wiring manual, explícito, visível.**

---

# Parte 2 — Testes Avançados

## A filosofia de teste em Go

Go não tem frameworks de mock pesados como Mockito (Java) ou Cuckoo (Swift).
A razão é simples: **interfaces resolvem o problema de mocking naturalmente**.

Em vez de gerar mocks com reflection ou codegen, você escreve uma struct que
implementa a interface e a usa no teste. Isso é mais verboso? Sim. Mas é explícito,
debugável e não requer ferramentas externas.

## Test doubles — stub, spy, mock, fake

### Stub — retorna valores fixos

```go
// Stub: implementação simplificada que retorna dados pré-definidos
type StubPokemonRepository struct {
    pokemon *Pokemon
    err     error
}

func (s *StubPokemonRepository) FindByID(ctx context.Context, id string) (*Pokemon, error) {
    return s.pokemon, s.err
}

func (s *StubPokemonRepository) Save(ctx context.Context, p *Pokemon) error {
    return nil
}

// Uso no teste
func TestService_GetByID(t *testing.T) {
    stub := &StubPokemonRepository{
        pokemon: &Pokemon{ID: "25", Nome: "Pikachu", Level: 25},
    }
    svc := NewService(stub)

    p, err := svc.GetByID(context.Background(), "25")

    assert.NoError(t, err)
    assert.Equal(t, "Pikachu", p.Nome)
    assert.Equal(t, 25, p.Level)
}
```

### Spy — registra chamadas

```go
// Spy: registra o que foi chamado, com quais argumentos
type SpyPokemonRepository struct {
    FindByIDCalls []string          // registra todos os IDs buscados
    SaveCalls     []*Pokemon        // registra todos os pokémons salvos
    pokemon       *Pokemon          // valor a retornar
}

func (s *SpyPokemonRepository) FindByID(ctx context.Context, id string) (*Pokemon, error) {
    s.FindByIDCalls = append(s.FindByIDCalls, id)
    return s.pokemon, nil
}

func (s *SpyPokemonRepository) Save(ctx context.Context, p *Pokemon) error {
    s.SaveCalls = append(s.SaveCalls, p)
    return nil
}

// Uso no teste — verifica se o service chamou o repository corretamente
func TestService_LevelUp_CallsSave(t *testing.T) {
    spy := &SpyPokemonRepository{
        pokemon: &Pokemon{ID: "25", Nome: "Pikachu", Level: 25},
    }
    svc := NewService(spy)

    err := svc.LevelUp(context.Background(), "25")

    assert.NoError(t, err)
    assert.Equal(t, 1, len(spy.FindByIDCalls))
    assert.Equal(t, "25", spy.FindByIDCalls[0])
    assert.Equal(t, 1, len(spy.SaveCalls))
    assert.Equal(t, 26, spy.SaveCalls[0].Level)
}
```

### Mock — stub + spy + assertions

```go
// Mock: combina stub (valores fixos) + spy (registra chamadas) + assertions
type MockPokemonRepository struct {
    FindByIDFunc func(ctx context.Context, id string) (*Pokemon, error)
    SaveFunc     func(ctx context.Context, p *Pokemon) error

    // tracking
    FindByIDCalls int
    SaveCalls     int
}

func (m *MockPokemonRepository) FindByID(ctx context.Context, id string) (*Pokemon, error) {
    m.FindByIDCalls++
    return m.FindByIDFunc(ctx, id)
}

func (m *MockPokemonRepository) Save(ctx context.Context, p *Pokemon) error {
    m.SaveCalls++
    return m.SaveFunc(ctx, p)
}

// Uso — define comportamento inline no teste
func TestService_LevelUp_NotFound(t *testing.T) {
    mock := &MockPokemonRepository{
        FindByIDFunc: func(ctx context.Context, id string) (*Pokemon, error) {
            return nil, ErrNotFound
        },
        SaveFunc: func(ctx context.Context, p *Pokemon) error {
            return nil
        },
    }
    svc := NewService(mock)

    err := svc.LevelUp(context.Background(), "999")

    assert.ErrorIs(t, err, ErrNotFound)
    assert.Equal(t, 1, mock.FindByIDCalls)
    assert.Equal(t, 0, mock.SaveCalls) // Save nunca foi chamado
}
```

### Fake — implementação funcional simplificada

```go
// Fake: implementação real, mas simplificada (ex: banco em memória)
type FakePokemonRepository struct {
    pokemons map[string]*Pokemon
}

func NewFakeRepository() *FakePokemonRepository {
    return &FakePokemonRepository{pokemons: make(map[string]*Pokemon)}
}

func (f *FakePokemonRepository) FindByID(ctx context.Context, id string) (*Pokemon, error) {
    p, ok := f.pokemons[id]
    if !ok {
        return nil, ErrNotFound
    }
    return p, nil
}

func (f *FakePokemonRepository) Save(ctx context.Context, p *Pokemon) error {
    f.pokemons[p.ID] = p
    return nil
}

// Uso — útil para testes de integração e cenários com estado
func TestService_MultipleOperations(t *testing.T) {
    repo := NewFakeRepository()
    repo.Save(nil, &Pokemon{ID: "25", Nome: "Pikachu", Level: 25})
    svc := NewService(repo)

    svc.LevelUp(nil, "25")
    svc.LevelUp(nil, "25")

    p, _ := svc.GetByID(nil, "25")
    assert.Equal(t, 27, p.Level)
}
```

### Quando usar cada um

| Double | Use quando... |
|--------|--------------|
| **Stub** | O teste só precisa de valores de retorno fixos |
| **Spy** | Você precisa verificar **se** e **como** o dependente foi chamado |
| **Mock** | Você precisa de stubs + spies + assertions inline (comportamento por teste) |
| **Fake** | Você precisa de estado entre operações (ex: múltiplas chamadas no mesmo teste) |

## Mock com `testify/mock`

Para interfaces com muitos métodos, `testify/mock` reduz boilerplate:

```go
package service_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock gerado com testify
type MockPokemonRepository struct {
    mock.Mock
}

func (m *MockPokemonRepository) FindByID(ctx context.Context, id string) (*Pokemon, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*Pokemon), args.Error(1)
}

func (m *MockPokemonRepository) Save(ctx context.Context, p *Pokemon) error {
    args := m.Called(ctx, p)
    return args.Error(0)
}

func TestService_LevelUp_WithTestify(t *testing.T) {
    mockRepo := new(MockPokemonRepository)

    // configura expectativas
    mockRepo.On("FindByID", mock.Anything, "25").
        Return(&Pokemon{ID: "25", Nome: "Pikachu", Level: 25}, nil)
    mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(p *Pokemon) bool {
        return p.ID == "25" && p.Level == 26
    })).Return(nil)

    svc := NewService(mockRepo)
    err := svc.LevelUp(context.Background(), "25")

    assert.NoError(t, err)
    mockRepo.AssertExpectations(t) // verifica que todas as chamadas esperadas ocorreram
}
```

**Atenção:** `testify/mock` é útil, mas use com moderação. Para interfaces pequenas
(1-3 métodos), o mock manual é mais legível e não adiciona dependência.

## `testify/assert` e `testify/require`

A diferença entre `assert` e `require`:

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestExemplo(t *testing.T) {
    // assert — marca o teste como falho, mas CONTINUA executando
    assert.Equal(t, 25, pokemon.Level)
    assert.Equal(t, "Pikachu", pokemon.Nome) // executa mesmo se o anterior falhou

    // require — marca o teste como falho e PARA imediatamente
    pokemon, err := repo.FindByID(ctx, "25")
    require.NoError(t, err)  // se falhar, as próximas linhas não executam
    require.NotNil(t, pokemon)
    assert.Equal(t, 25, pokemon.Level)  // seguro: sabemos que pokemon não é nil
}
```

## Testando handlers HTTP

`httptest.NewRecorder` já foi coberto na Fase 05. Aqui vai o padrão completo com
service mockado:

```go
func TestPokemonHandler_GetByID(t *testing.T) {
    tests := []struct {
        name       string
        id         string
        setupMock  func(*MockPokemonUseCase)
        wantStatus int
        wantBody   string
    }{
        {
            name: "pokemon encontrado",
            id:   "25",
            setupMock: func(m *MockPokemonUseCase) {
                m.On("GetByID", mock.Anything, "25").
                    Return(&PokemonDetail{Nome: "Pikachu", Level: 25}, nil)
            },
            wantStatus: http.StatusOK,
            wantBody:   `{"nome":"Pikachu","level":25}`,
        },
        {
            name: "pokemon nao encontrado",
            id:   "999",
            setupMock: func(m *MockPokemonUseCase) {
                m.On("GetByID", mock.Anything, "999").
                    Return(nil, ErrNotFound)
            },
            wantStatus: http.StatusNotFound,
            wantBody:   `{"erro":"nao encontrado"}`,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockUseCase := new(MockPokemonUseCase)
            tt.setupMock(mockUseCase)

            handler := NewPokemonHandler(mockUseCase)
            req := httptest.NewRequest("GET", "/api/pokemon/"+tt.id, nil)
            rec := httptest.NewRecorder()

            handler.GetByID(rec, req)

            assert.Equal(t, tt.wantStatus, rec.Code)
            assert.JSONEq(t, tt.wantBody, rec.Body.String())
            mockUseCase.AssertExpectations(t)
        })
    }
}
```

## Testando clientes HTTP — `httptest.Server` como fake upstream

Quando seu código faz chamadas HTTP para serviços externos, use `httptest.Server`
para simular o serviço:

```go
func TestPokemonCatalogClient_GetByID(t *testing.T) {
    // servidor fake que simula o pokemon-catalog-service
    fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // verifica que a requisição está correta
        assert.Equal(t, "/api/v1/pokemons/25", r.URL.Path)
        assert.Equal(t, "GET", r.Method)

        // responde como o serviço real responderia
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "id":    "25",
            "nome":  "Pikachu",
            "level": 25,
        })
    }))
    defer fakeServer.Close()

    // o client aponta para o servidor fake
    client := NewPokemonCatalogClient(fakeServer.URL)

    pokemon, err := client.GetByID(context.Background(), "25")

    assert.NoError(t, err)
    assert.Equal(t, "Pikachu", pokemon.Nome)
    assert.Equal(t, 25, pokemon.Level)
}

func TestPokemonCatalogClient_ServerError(t *testing.T) {
    fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusInternalServerError)
    }))
    defer fakeServer.Close()

    client := NewPokemonCatalogClient(fakeServer.URL)
    _, err := client.GetByID(context.Background(), "25")

    assert.Error(t, err)
}
```

## Testando acesso a banco de dados

### Opção 1 — `sqlmock` (testes de unidade do repository)

`sqlmock` simula o driver SQL — não precisa de banco real. Ideal para testar
a lógica do repository (queries, scanning, transações):

```go
import (
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
)

func TestPostgresRepository_FindByID(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    // configura o mock: "quando essa query chegar, retorne essas linhas"
    rows := sqlmock.NewRows([]string{"id", "nome", "level"}).
        AddRow("25", "Pikachu", 25)

    mock.ExpectQuery(`SELECT id, nome, level FROM pokemons WHERE id = \$1`).
        WithArgs("25").
        WillReturnRows(rows)

    repo := &PostgresRepository{db: db}
    pokemon, err := repo.FindByID(context.Background(), "25")

    assert.NoError(t, err)
    assert.Equal(t, "Pikachu", pokemon.Nome)
    assert.Equal(t, 25, pokemon.Level)
    assert.NoError(t, mock.ExpectationsWereMet()) // verifica que a query foi executada
}
```

### Opção 2 — Transação com rollback (testes de integração)

Usa um banco real (PostgreSQL/MySQL) mas faz rollback ao final — o banco
nunca é alterado:

```go
func TestPostgresRepository_Save(t *testing.T) {
    if testing.Short() {
        t.Skip("pula teste de integração em modo short")
    }

    db, err := sql.Open("postgres", "host=localhost dbname=pokedex_test sslmode=disable")
    require.NoError(t, err)
    defer db.Close()

    // inicia transação — tudo será desfeito no rollback
    tx, err := db.Begin()
    require.NoError(t, err)
    defer tx.Rollback() // rollback garante que o banco volta ao estado original

    repo := &PostgresRepository{db: tx} // usa a transação, não o db direto

    p := &Pokemon{ID: "25", Nome: "Pikachu", Level: 25}
    err = repo.Save(context.Background(), p)
    assert.NoError(t, err)

    // verifica que foi salvo (dentro da transação)
    found, err := repo.FindByID(context.Background(), "25")
    assert.NoError(t, err)
    assert.Equal(t, "Pikachu", found.Nome)
    // rollback automático — o banco real não é afetado
}
```

### Opção 3 — Testcontainers (testes de integração isolados)

Sobe um container Docker com banco real por teste:

```go
import (
    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestWithTestcontainers(t *testing.T) {
    if testing.Short() {
        t.Skip()
    }

    ctx := context.Background()

    // sobe PostgreSQL em container
    pgContainer, err := postgres.Run(ctx,
        "postgres:16-alpine",
        postgres.WithDatabase("pokedex_test"),
        postgres.WithUsername("test"),
        postgres.WithPassword("test"),
    )
    require.NoError(t, err)
    defer pgContainer.Terminate(ctx)

    connStr, _ := pgContainer.ConnectionString(ctx)
    db, _ := sql.Open("postgres", connStr)
    defer db.Close()

    // roda migrations
    runMigrations(db)

    // testa com banco real, isolado, descartável
    repo := &PostgresRepository{db: db}
    // ...
}
```

## Testando middleware

```go
func TestAuthMiddleware(t *testing.T) {
    tests := []struct {
        name       string
        token      string
        wantStatus int
    }{
        {"com token", "Bearer valid-token", http.StatusOK},
        {"sem token", "", http.StatusUnauthorized},
        {"token invalido", "Bearer invalid", http.StatusUnauthorized},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // handler real que o middleware protege
            next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
            })

            handler := AuthMiddleware(next)

            req := httptest.NewRequest("GET", "/protegido", nil)
            if tt.token != "" {
                req.Header.Set("Authorization", tt.token)
            }
            rec := httptest.NewRecorder()

            handler.ServeHTTP(rec, req)

            assert.Equal(t, tt.wantStatus, rec.Code)
        })
    }
}
```

## Testes de integração — cenário completo

Testa handler → service → repository real (fake em memória):

```go
func TestIntegration_PokemonFlow(t *testing.T) {
    // wiring manual — sem mocks, com fake repository
    repo := NewFakePokemonRepository()
    svc := NewPokemonService(repo)
    handler := NewPokemonHandler(svc)

    // PASSO 1: cria um pokémon
    body := strings.NewReader(`{"nome":"Pikachu","level":25}`)
    req := httptest.NewRequest("POST", "/api/pokemons", body)
    rec := httptest.NewRecorder()
    handler.Create(rec, req)
    assert.Equal(t, http.StatusCreated, rec.Code)

    var created Pokemon
    json.Unmarshal(rec.Body.Bytes(), &created)

    // PASSO 2: busca o pokémon criado
    req = httptest.NewRequest("GET", "/api/pokemons/"+created.ID, nil)
    rec = httptest.NewRecorder()
    handler.GetByID(rec, req)
    assert.Equal(t, http.StatusOK, rec.Code)

    var found Pokemon
    json.Unmarshal(rec.Body.Bytes(), &found)
    assert.Equal(t, "Pikachu", found.Nome)
    assert.Equal(t, 25, found.Level)

    // PASSO 3: faz level up
    req = httptest.NewRequest("POST", "/api/pokemons/"+created.ID+"/levelup", nil)
    rec = httptest.NewRecorder()
    handler.LevelUp(rec, req)
    assert.Equal(t, http.StatusOK, rec.Code)

    // PASSO 4: verifica que o level subiu
    req = httptest.NewRequest("GET", "/api/pokemons/"+created.ID, nil)
    rec = httptest.NewRecorder()
    handler.GetByID(rec, req)
    json.Unmarshal(rec.Body.Bytes(), &found)
    assert.Equal(t, 26, found.Level)
}
```

## Test helpers e fixtures

Extraia lógica repetitiva de setup:

```go
// helper: cria um service com repository fake populado
func setupTestService(t *testing.T, pokemons ...*Pokemon) (*PokemonService, *FakePokemonRepository) {
    t.Helper()
    repo := NewFakePokemonRepository()
    for _, p := range pokemons {
        repo.Save(context.Background(), p)
    }
    return NewPokemonService(repo), repo
}

// fixture: dados de teste reutilizáveis
var (
    pikachu   = &Pokemon{ID: "25", Nome: "Pikachu", Level: 25}
    charizard = &Pokemon{ID: "6", Nome: "Charizard", Level: 36}
)

func TestService_GetByID(t *testing.T) {
    svc, _ := setupTestService(t, pikachu, charizard)

    p, err := svc.GetByID(context.Background(), "25")
    assert.NoError(t, err)
    assert.Equal(t, "Pikachu", p.Nome)
}
```

## Swift vs Go — testes

```swift
// Swift — XCTest + protocolos para mocking
protocol PokemonRepository {
    func findById(_ id: String) throws -> Pokemon
}

class MockRepository: PokemonRepository {
    var findByIdResult: Result<Pokemon, Error> = .failure(NSError())
    var findByIdCalled = false
    func findById(_ id: String) throws -> Pokemon {
        findByIdCalled = true
        return try findByIdResult.get()
    }
}
```

```go
// Go — testing + interfaces para mocking (sem frameworks mágicos)
type PokemonRepository interface {
    FindByID(ctx context.Context, id string) (*Pokemon, error)
}

type MockRepository struct {
    FindByIDFunc func(ctx context.Context, id string) (*Pokemon, error)
}
func (m *MockRepository) FindByID(ctx context.Context, id string) (*Pokemon, error) {
    return m.FindByIDFunc(ctx, id)
}
```

**A diferença fundamental:** Em Swift/iOS, você depende de protocolos + XCTest +
possivelmente Cuckoo/Sourcery para gerar mocks. Em Go, a interface é suficiente —
sem codegen, sem framework, sem reflection.

---

## Exercícios da Fase 07

### 1. Crie um spy

Implemente um `SpyPokemonRepository` que registra quantas vezes `FindByID` foi
chamado e com quais IDs. Use-o para testar que `GetByID` chama o repository
exatamente uma vez com o ID correto.

<details>
<summary>Gabarito</summary>

```go
type SpyPokemonRepository struct {
    CalledWith []string
    Pokemon    *Pokemon
}

func (s *SpyPokemonRepository) FindByID(ctx context.Context, id string) (*Pokemon, error) {
    s.CalledWith = append(s.CalledWith, id)
    return s.Pokemon, nil
}

func TestService_CallsRepositoryOnce(t *testing.T) {
    spy := &SpyPokemonRepository{Pokemon: &Pokemon{ID: "25", Nome: "Pikachu"}}
    svc := NewService(spy)

    _, _ = svc.GetByID(context.Background(), "25")

    assert.Equal(t, 1, len(spy.CalledWith))
    assert.Equal(t, "25", spy.CalledWith[0])
}
```
</details>

### 2. Teste um client HTTP com `httptest.Server`

Crie um client que busca pokémons em uma API externa e teste-o com um servidor fake
que simula sucesso (200) e erro (500).

<details>
<summary>Gabarito</summary>

```go
func TestExternalAPIClient(t *testing.T) {
    t.Run("sucesso", func(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            assert.Equal(t, "/api/pokemon/25", r.URL.Path)
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(map[string]interface{}{
                "id": "25", "nome": "Pikachu", "level": 25,
            })
        }))
        defer server.Close()

        client := NewPokemonClient(server.URL)
        p, err := client.GetByID(context.Background(), "25")
        assert.NoError(t, err)
        assert.Equal(t, "Pikachu", p.Nome)
    })

    t.Run("erro servidor", func(t *testing.T) {
        server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusInternalServerError)
        }))
        defer server.Close()

        client := NewPokemonClient(server.URL)
        _, err := client.GetByID(context.Background(), "25")
        assert.Error(t, err)
    })
}
```
</details>

### 3. Functional options

Refatore a struct `ServerConfig` do exemplo para usar o padrão functional options.
Adicione pelo menos 3 opções com defaults sensíveis.

<details>
<summary>Gabarito</summary>

```go
type ServerConfig struct {
    Port    int
    Timeout time.Duration
    Debug   bool
}

type ServerOption func(*ServerConfig)

func WithPort(p int) ServerOption             { return func(c *ServerConfig) { c.Port = p } }
func WithTimeout(d time.Duration) ServerOption { return func(c *ServerConfig) { c.Timeout = d } }
func WithDebug() ServerOption                  { return func(c *ServerConfig) { c.Debug = true } }

func NewServerConfig(opts ...ServerOption) *ServerConfig {
    cfg := &ServerConfig{Port: 8080, Timeout: 30 * time.Second}
    for _, opt := range opts {
        opt(cfg)
    }
    return cfg
}

func TestServerConfig(t *testing.T) {
    cfg := NewServerConfig(WithPort(3000), WithDebug())
    assert.Equal(t, 3000, cfg.Port)
    assert.Equal(t, 30*time.Second, cfg.Timeout)
    assert.True(t, cfg.Debug)
}
```
</details>

### 4. Teste de integração com fake

Implemente um `FakePokemonRepository` (map em memória) e escreva um teste de
integração que cria um pokémon, faz level up duas vezes e verifica o resultado final.

<details>
<summary>Gabarito</summary>

```go
type FakePokemonRepository struct {
    data map[string]*Pokemon
}

func NewFakeRepo() *FakePokemonRepository {
    return &FakePokemonRepository{data: make(map[string]*Pokemon)}
}

func (f *FakePokemonRepository) FindByID(ctx context.Context, id string) (*Pokemon, error) {
    p, ok := f.data[id]
    if !ok {
        return nil, ErrNotFound
    }
    return p, nil
}

func (f *FakePokemonRepository) Save(ctx context.Context, p *Pokemon) error {
    f.data[p.ID] = p
    return nil
}

func TestIntegration_LevelUpTwice(t *testing.T) {
    repo := NewFakeRepo()
    repo.Save(nil, &Pokemon{ID: "25", Nome: "Pikachu", Level: 25})
    svc := NewService(repo)

    err := svc.LevelUp(context.Background(), "25")
    require.NoError(t, err)
    err = svc.LevelUp(context.Background(), "25")
    require.NoError(t, err)

    p, err := svc.GetByID(context.Background(), "25")
    assert.NoError(t, err)
    assert.Equal(t, 27, p.Level)
}
```
</details>

---

**Voltar ao início:** [ROTEIRO](ROTEIRO.md)
