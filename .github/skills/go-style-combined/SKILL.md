---
name: go-style-combined
description: >-
  Guia de estilo Go completo para a Pokedex Platform, combinando o Uber Go Style Guide e o Google Go Style Guide
  com as decisões específicas já tomadas no projeto. Cobre naming, receivers, initialisms, imports de 3 grupos,
  asserções de interface em tempo de compilação, make() para slices, time.Duration para constantes de tempo,
  errors (wrapping, sentinels, strings), goroutines com context e ciclo de vida, mutex, slices/maps,
  struct init, declarações agrupadas, testes orientados a tabela, functional options, panic, init(),
  variáveis globais mutáveis, atomic, embed e linting.
  Use para QUALQUER revisão ou escrita de código Go na plataforma. Esta skill é a referência canônica.
  Trigger examples: "estilo Go", "naming", "receiver", "import", "erro", "goroutine", "mutex",
  "slice nil", "interface compliance", "time.Duration", "table test", "functional options",
  "código idiomático Go", "code review", "revisar código".
applyTo:
  - "core/bff/mobile-bff/**/*.go"
  - "core/app/**/*.go"
  - "core/gateway/**/*.go"
---

# Go Style Guide — Pokedex Platform

> Síntese canônica do **Uber Go Style Guide** + **Google Go Style Guide**  
> com decisões específicas do projeto documentadas.
>
> Idioma dos comentários no código: **Português Brasil**  
> Quando houver conflito entre os guias, prevalece a **consistência local** (Google principle §Consistency).

---

## Conteúdo

1. [Princípios](#1-princípios)
2. [Naming](#2-naming)
3. [Comentários e Godoc](#3-comentários-e-godoc)
4. [Imports](#4-imports)
5. [Declarações e agrupamentos](#5-declarações-e-agrupamentos)
6. [Inicialização de structs, slices e maps](#6-inicialização-de-structs-slices-e-maps)
7. [Interfaces e asserções de compilação](#7-interfaces-e-asserções-de-compilação)
8. [Erros](#8-erros)
9. [Context](#9-context)
10. [Goroutines e concorrência](#10-goroutines-e-concorrência)
11. [Mutex e sincronização](#11-mutex-e-sincronização)
12. [Performance](#12-performance)
13. [Controle de fluxo](#13-controle-de-fluxo)
14. [Testes orientados a tabela](#14-testes-orientados-a-tabela)
15. [Functional Options](#15-functional-options)
16. [Panic e init()](#16-panic-e-init)
17. [Variáveis globais mutáveis](#17-variáveis-globais-mutáveis)
18. [Linting e ferramentas](#18-linting-e-ferramentas)
19. [Decisões específicas do projeto](#19-decisões-específicas-do-projeto)

---

## 1. Princípios

O Google Go Style Guide elenca os atributos do código legível, em ordem de importância:

| Prioridade | Atributo | Descrição |
|:---:|---|---|
| 1 | **Clareza** | O propósito e o raciocínio são óbvios para o leitor |
| 2 | **Simplicidade** | O objetivo é atingido da forma mais simples possível |
| 3 | **Concisão** | Alta relação sinal/ruído — cada linha conta |
| 4 | **Manutenibilidade** | Fácil de manter ao longo do tempo |
| 5 | **Consistência** | Consistente com o restante da base de código |

> **Regra de ouro**: prefira o mecanismo mais simples que resolve o problema.  
> Channel, slice, map ou loop nativos > stdlib > third-party > código próprio.

---

## 2. Naming

### 2.1 Pacotes

- Minúsculas, sem underscore, sem camelCase. Ex.: `pokemon`, `auth`, `http`.
- Sem sufixo info/helper/util/common/shared/lib — esses nomes não comunicam nada.
- O nome do pacote faz parte do símbolo: `auth.NewService()`, não `auth.NewAuthService()`.
- Não plural: `net/url`, não `net/urls`.

```go
// ✅
package pokemon

// ❌
package pokemonUtils
package pokemonHelper
```

### 2.2 Receivers

- 1–2 letras, abreviação do tipo. **Nunca** `this`, `self`, `me`.
- Consistente em todos os métodos do mesmo tipo.

```go
// ✅
func (s *PokemonService) List(ctx context.Context) ([]domain.Pokemon, error)
func (c *AuthServiceClient) Login(ctx context.Context, email, password string) (*domain.Session, error)
func (r *PostgresPokemonRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error)

// ❌
func (this *PokemonService) List(ctx context.Context) ([]domain.Pokemon, error)
func (service *PokemonService) List(ctx context.Context) ([]domain.Pokemon, error)
```

### 2.3 Initialisms

Sempre em caixa completa (ou tudo minúsculo para unexported):

| Termo | Exported | Unexported |
|-------|----------|------------|
| ID | `ID` | `id` |
| URL | `URL` | `url` |
| HTTP | `HTTP` | `http` |
| JSON | `JSON` | `json` |
| DB | `DB` | `db` |
| API | `API` | `api` |
| gRPC (exported) | `GRPC` | `gRPC` |

```go
// ✅
type PokemonID string
func GetByID(id PokemonID) {}
var baseURL = "http://..."
type GRPCClient struct{}   // exported
var gRPCClient *GRPCClient // unexported

// ❌
type PokemonId string
func GetById(id PokemonId) {}
var baseUrl = "http://..."
```

### 2.4 Funções e métodos

- **Sem prefixo `Get`** para getters: `Count()` não `GetCount()`, `Name()` não `GetName()`.
  - Exceção: operação cara (RPC, IO) pode usar `Fetch`, `Compute`, `Load`.
- Construtores: `New` + tipo. Ex.: `NewPokemonService`, `NewHandler`, `NewAuthServiceClient`.
- Predicados booleanos: `Is`, `Has`, `Can`, `Should`. Ex.: `IsFavorite`, `HasAccess`.
- Nomes de funções de teste com underscore para agrupamento: `TestListPokemons_WhenUserNotFound`.

### 2.5 Constantes

- MixedCaps. Ex.: `MaxRetryAttempts`, `DefaultPageSize`.
- **Nunca** ALL_CAPS ou prefixo `K`. Ex.: ❌ `MAX_RETRY`, ❌ `kDefaultPage`.
- Nome pelo papel, não pelo valor: ❌ `Twelve = 12`.

### 2.6 Variáveis

- Comprimento proporcional ao escopo. Loop curto → `i`, `v`, `k`. Escopo grande → `userCount`.
- Omita o tipo no nome: `users` não `userSlice`; `count` não `numUsers`.
- Nomes locais pelo conteúdo no contexto atual, não pela origem.

---

## 3. Comentários e Godoc

### 3.1 Regras do projeto (PT-BR)

Todos os comentários devem estar em **Português Brasil**.

#### Quando comentar (alto valor)

- Explicação de lógica não-óbvia ("por que", não "o quê")
- Razões de design/arquitetura: por que esse pattern e não outro
- Avisos de performance, limitações conhecidas
- Contexto de integração entre serviços

#### Quando NÃO comentar (zero valor)

- Óbvio pelo nome → sem comentário
- Descreve o código literalmente: ❌ `// verifica erro`
- Código autoexplicativo com naming descritivo

### 3.2 Godoc

Todo símbolo **exportado** deve ter comentário godoc começando pelo nome:

```go
// PokemonService gerencia operações de listagem, busca e favoritos de Pokémons.
type PokemonService struct { ... }

// GetByID retorna um Pokémon pelo ID canônico.
// Retorna domain.ErrNotFound se o Pokémon não existir no catálogo.
func (s *PokemonService) GetByID(ctx context.Context, id string) (*domain.Pokemon, error)
```

#### Comentários no package:

```go
// Package pokemon implementa os casos de uso relacionados ao catálogo de Pokémons.
// Depende de PokemonRepository (outbound) e implementa PokemonUseCase (inbound).
package pokemon
```

### 3.3 Parâmetros desnudos

Adicione comentários inline `/* isLocal */` quando boolean/int passado sem nome for ambíguo:

```go
// ✅
printInfo("pikachu", true /* isElectric */, false /* isFavorite */)

// Melhor ainda: use tipo nomeado
type ElectricType bool
```

---

## 4. Imports

### 4.1 Três grupos (decisão do projeto)

O projeto usa **3 grupos** separados por linha em branco. O Uber usa 2, o Google usa 4+.  
Decisão: manter 3 para consistência com a base de código atual.

```go
import (
    // 1. stdlib
    "context"
    "fmt"
    "net/http"
    "time"

    // 2. third-party
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"

    // 3. internos (módulo pokedex-platform)
    "pokedex-platform/core/bff/mobile-bff/internal/domain"
    inbound "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"
    outbound "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
)
```

### 4.2 Alias de import

Use alias quando:
- O nome do pacote não bate com o último elemento do path (ex.: paths versionados `v2`)
- Há colisão entre dois pacotes com o mesmo nome
- O projeto já estabeleceu um alias canônico (ex.: `inbound`, `outbound`)

```go
// ✅
inbound "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"

// ❌ — rename desnecessário
pokemonDomain "pokedex-platform/core/bff/mobile-bff/internal/domain"
```

### 4.3 Import dot e blank

- `import .` → **proibido** (obscurece a origem dos símbolos)
- `import _` → apenas em `main` ou arquivos de teste que precisam de side effects

---

## 5. Declarações e agrupamentos

### 5.1 Agrupe declarações relacionadas

```go
// ✅
const (
    DefaultPageSize = 20
    MaxPageSize     = 100
)

var (
    ErrNotFound     = errors.New("não encontrado")
    ErrInvalidToken = errors.New("token inválido")
)

// ❌ — mistura não-relacionados no mesmo bloco
const (
    DefaultPageSize = 20
    EnvKey          = "APP_ENV"  // conceito diferente → bloco separado
)
```

### 5.2 Declaração de variável

```go
// var x T — zero value explícito e intencional
var count int
var repo outbound.PokemonRepository

// x := value — inicializa com valor não-zero
name := "Pikachu"
client := &http.Client{Timeout: 5 * time.Second}

// var x T quando todos os campos são zero
var cfg Config  // não: cfg := Config{}
```

### 5.3 Escopo mínimo

Reduza o escopo de variáveis sem sacrificar a legibilidade:

```go
// ✅
if err := os.WriteFile(name, data, 0644); err != nil {
    return err
}

// use declaração longa quando o result for usado fora do if
data, err := os.ReadFile(name)
if err != nil {
    return err
}
fmt.Println(data)
```

---

## 6. Inicialização de structs, slices e maps

### 6.1 Structs: sempre com nomes de campos

```go
// ✅
svc := &PokemonService{
    pokemonRepo:  repo,
    favoriteRepo: favoriteRepo,
    authProvider: authProvider,
}

// ❌ — quebra se a ordem dos campos mudar
svc := &PokemonService{repo, favoriteRepo, authProvider}
```

Exceção: tabelas de teste com ≤3 campos simples podem omitir os nomes.

### 6.2 Omita campos zero-value

```go
// ✅
cfg := Config{
    Host: "localhost",
    Port: 8080,
}

// ❌
cfg := Config{
    Host:    "localhost",
    Port:    8080,
    Timeout: 0,     // zero value — desnecessário
    Debug:   false, // zero value — desnecessário
}
```

Exceção: em testes, inclua campos zero-value quando o nome fornece contexto relevante:

```go
tests := []struct {
    give string
    want int
}{
    {give: "0", want: 0},  // want: 0 documenta a intenção
}
```

### 6.3 Slices — decisão do projeto

**REGRA DO PROJETO**: Use `make([]T, 0)` quando o slice é retornado via JSON ou API.

Justificativa: `nil` slice serializa como `null` em JSON; `make([]T, 0)` serializa como `[]`.

```go
// ✅ — em repositories e use cases que retornam via API
pokemons := make([]domain.Pokemon, 0)
favorites := make([]string, 0)

// ✅ — com capacidade conhecida (preferível)
pokemons := make([]domain.Pokemon, 0, len(ids))

// ✅ — var é ok para slices internos que não são retornados via JSON
var nums []int
if condition {
    nums = append(nums, 1)
}

// ✅ — verificar se vazio com len(), não com == nil
if len(pokemons) == 0 { ... }
```

### 6.4 Maps

```go
// ✅ — map vazio programático
m := make(map[string]*domain.Pokemon)

// ✅ — com capacidade conhecida
m := make(map[string]*domain.Pokemon, len(ids))

// ✅ — map com conjunto fixo de elementos
statusMessages := map[int]string{
    http.StatusOK:           "ok",
    http.StatusNotFound:     "não encontrado",
    http.StatusUnauthorized: "não autorizado",
}

// ❌ — map nil (pânico em escrita)
var m map[string]*domain.Pokemon
m["pikachu"] = p // panic!
```

### 6.5 Referências de struct: use `&T{}` não `new(T)`

```go
// ✅
svc := &PokemonService{repo: repo}

// ❌ — inconsistente
svc := new(PokemonService)
svc.repo = repo
```

---

## 7. Interfaces e asserções de compilação

### 7.1 Defina no lado do consumidor

```go
// ports/outbound/repository.go — quem consome define a interface
type PokemonRepository interface {
    GetByID(ctx context.Context, id string) (*domain.Pokemon, error)
    ListAll(ctx context.Context) ([]domain.Pokemon, error)
}

// adapters/outbound/postgres/pokemon_repository.go — implementa
type PostgresPokemonRepository struct { db *sql.DB }
```

### 7.2 Interfaces pequenas

Prefira 1–3 métodos. Interfaces grandes são difíceis de mockar e violam ISP:

```go
// ✅ — focada
type TokenValidator interface {
    ValidateToken(ctx context.Context, token string) (*domain.TokenClaims, error)
}

// ❌ — muito ampla
type AuthEverything interface {
    Login(...)
    Logout(...)
    Register(...)
    ValidateToken(...)
    RefreshToken(...)
    ChangePassword(...)
}
```

### 7.3 Asserção de interface em tempo de compilação — OBRIGATÓRIO no projeto

Para cada struct que implementa uma interface, declare a asserção no **final do arquivo**:

```go
// Ao final de cada arquivo de implementação:

// adapters/outbound/postgres/pokemon_repository.go
var _ outbound.PokemonRepository = (*PostgresPokemonRepository)(nil)

// adapters/outbound/postgres/favorite_repository.go
var _ outbound.FavoriteRepository = (*PostgresFavoriteRepository)(nil)

// adapters/outbound/http/auth_service_client.go (implementa 2 interfaces)
var _ outbound.AuthProvider    = (*AuthServiceClient)(nil)
var _ inbound.TokenValidator   = (*AuthServiceClient)(nil)

// adapters/outbound/http/pokemon_catalog_client.go
var _ outbound.PokemonRepository = (*PokemonCatalogServiceRepository)(nil)

// service/pokemon_service.go
var _ inbound.PokemonUseCase = (*PokemonService)(nil)

// service/favorite_service.go
var _ inbound.FavoriteUseCase = (*FavoriteService)(nil)
```

---

## 8. Erros

### 8.1 Nunca descarte erros silenciosamente

```go
// ❌ — proibido em código de produção
_ = repo.Save(ctx, pokemon)

// ✅ — pelo menos logue se não puder retornar
if err := repo.Save(ctx, pokemon); err != nil {
    log.Printf("falha ao salvar pokemon %s: %v", id, err)
}
```

### 8.2 Strings de erro

- Minúsculas, sem ponto final (serão concatenadas)
- Mensagens de log/UI podem ter capitalização normal

```go
// ✅
return fmt.Errorf("pokemon não encontrado: %w", err)
return errors.New("token inválido")

// ❌
return fmt.Errorf("Pokemon não encontrado.")
return errors.New("Token inválido.")
```

### 8.3 Wrapping: `%w` vs `%v`

| Use | Quando |
|-----|--------|
| `%w` | O chamador pode ou vai usar `errors.Is` / `errors.As` |
| `%v` | Somente log/anotação, sem necessidade de inspeção |

```go
// ✅ — retorno ao chamador
if err := s.repo.GetByID(ctx, id); err != nil {
    return nil, fmt.Errorf("buscar pokemon %s: %w", id, err)
}

// ✅ — só log
log.Printf("operação completada com aviso: %v", err)
```

### 8.4 Sentinel errors

Use `var ErrX = errors.New(...)` para erros esperados que os chamadores verificam:

```go
var (
    ErrNotFound      = errors.New("não encontrado")
    ErrInvalidToken  = errors.New("token inválido")
    ErrAlreadyExists = errors.New("já existe")
    ErrUnauthorized  = errors.New("não autorizado")
)
```

### 8.5 Tipos de erro customizados

Use quando precisar transportar dados estruturados:

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("campo %s: %s", e.Field, e.Message)
}
```

### 8.6 Como tratar erros (4 opções)

1. **Handle e retorne** — o mais comum
2. **Encapsule e retorne** com `%w` ao chamador
3. **Logue e degrade** graciosamente (operação não crítica)
4. **Corresponda e degrade** usando `errors.Is` / `errors.As`

```go
// 1. Handle e retorne
if err := validate(input); err != nil {
    return nil, fmt.Errorf("validar entrada: %w", err)
}

// 3. Logue e degrade
if err := emitMetrics(); err != nil {
    // métricas não devem derrubar a aplicação
    log.Printf("falha ao emitir métricas: %v", err)
}

// 4. Corresponda e degrade
tz, err := getUserTimeZone(id)
if err != nil {
    if errors.Is(err, ErrNotFound) {
        tz = time.UTC
    } else {
        return fmt.Errorf("obter timezone do usuário %s: %w", id, err)
    }
}
```

### 8.7 Asserção de tipo: use "comma ok"

```go
// ✅
t, ok := i.(string)
if !ok {
    return fmt.Errorf("tipo inesperado: %T", i)
}

// ❌ — pânico se o tipo estiver errado
t := i.(string)
```

### 8.8 Erros in-band: evite

```go
// ❌ — força o chamador a checar valor sentinela
func Lookup(key string) int  // retorna -1 se não encontrado

// ✅ — múltiplos retornos
func Lookup(key string) (value string, ok bool)
func GetByID(ctx context.Context, id string) (*Pokemon, error)
```

---

## 9. Context

### 9.1 Primeiro parâmetro, sempre

```go
// ✅
func (s *PokemonService) ListPokemons(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error)

// ❌
func (s *PokemonService) ListPokemons(page, pageSize int, ctx context.Context) (*domain.PokemonPage, error)
```

### 9.2 Nunca armazene context em struct

```go
// ❌ — proibido
type Service struct {
    ctx context.Context
}

// ✅ — passe como parâmetro em cada chamada
func (s *Service) Do(ctx context.Context) error { ... }
```

### 9.3 Use context para cancelamento em goroutines

```go
func (s *Scheduler) Start(ctx context.Context) {
    // encerra quando ctx for cancelado pelo chamador
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case <-s.ticker.C:
                s.processWork(ctx)
            }
        }
    }()
}
```

---

## 10. Goroutines e concorrência

### 10.1 Nunca lance goroutines sem estratégia de finalização

```go
// ❌ — goroutine órfã
func init() {
    go backgroundWorker()
}

// ✅ — ciclo de vida documentado, context usado para cancelamento
func NewScheduler(ctx context.Context) *Scheduler {
    s := &Scheduler{}
    // goroutine encerra quando ctx for cancelado
    go s.run(ctx)
    return s
}
```

### 10.2 Documente quando a goroutine encerra

```go
// Start inicia o worker em background.
// A goroutine encerra quando ctx for cancelado. Aguarde com WaitGroup se necessário.
func (w *Worker) Start(ctx context.Context, wg *sync.WaitGroup) {
    wg.Add(1)
    go func() {
        defer wg.Done()
        w.loop(ctx)
    }()
}
```

### 10.3 Tamanho de canal explícito

```go
// ✅ — buffer justificado com comentário
// canal com buffer 1 para desacoplar produtor e consumidor
results := make(chan *domain.Pokemon, 1)

// ❌ — sem buffer quando o consumidor pode ser lento → deadlock
results := make(chan *domain.Pokemon)
```

### 10.4 Operações atômicas: use `go.uber.org/atomic`

```go
// ❌ — fácil esquecer de usar operação atômica
type Service struct {
    running int32 // atômico
}

// ✅ — tipo seguro
type Service struct {
    running atomic.Bool
}

func (s *Service) start() {
    if s.running.Swap(true) {
        return // já rodando
    }
}
```

---

## 11. Mutex e sincronização

### 11.1 Declare mutex junto ao que ele protege

```go
// ✅
type MockPokemonRepository struct {
    mu       sync.RWMutex
    pokemons map[string]*domain.Pokemon // protegido por mu
}
```

### 11.2 RWMutex quando leituras superam escritas

```go
func (m *MockPokemonRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    p, ok := m.pokemons[id]
    if !ok {
        return nil, domain.ErrNotFound
    }
    return p, nil
}

func (m *MockPokemonRepository) Save(ctx context.Context, p *domain.Pokemon) error {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.pokemons[p.ID] = p
    return nil
}
```

### 11.3 Evite embed de tipos em structs públicas

```go
// ❌ — vaza detalhes de implementação, impede evolução
type ConcreteList struct {
    *AbstractList
}

// ✅ — composição explícita
type ConcreteList struct {
    list *AbstractList
}

func (l *ConcreteList) Add(e Entity) { l.list.Add(e) }
```

---

## 12. Performance

### 12.1 Prefira `strconv` a `fmt` para conversões numéricas

```go
// ✅ — mais rápido, 1 alloc
s := strconv.Itoa(42)
n, _ := strconv.Atoi("42")

// ❌ — mais lento, 2 allocs
s := fmt.Sprintf("%d", 42)
```

### 12.2 Evite conversões repetidas de string→byte

```go
// ✅
separator := []byte(", ")  // converte uma vez
for _, item := range items {
    w.Write(separator)
    w.Write([]byte(item))
}

// ❌ — recria o slice a cada iteração
for range items {
    w.Write([]byte(", "))
}
```

### 12.3 Pré-aloque slices e maps com capacidade conhecida

```go
// ✅
pokemons := make([]domain.Pokemon, 0, len(ids))
index := make(map[string]*domain.Pokemon, len(ids))

// ❌ — realoca conforme cresce
var pokemons []domain.Pokemon
index := map[string]*domain.Pokemon{}
```

### 12.4 `time.Duration` para constantes de tempo — decisão do projeto

**REGRA DO PROJETO**: Constantes e variáveis de duração de tempo devem usar `time.Duration`.

```go
// ✅
const defaultAuthRateLimitWindow = 60 * time.Second

// uso direto — sem cast
time.NewTicker(defaultAuthRateLimitWindow)

// ❌ — tipo int, exige cast a cada uso
const defaultAuthRateLimitWindowSeconds = 60

// uso ruim — cast repetido
time.NewTicker(time.Duration(defaultAuthRateLimitWindowSeconds) * time.Second)
```

### 12.5 String de formato como constante

```go
// ✅ — go vet pode verificar em tempo de compilação
const msgFmt = "pokemon %s não encontrado (status=%d)"
fmt.Printf(msgFmt, id, code)

// ❌ — go vet não alcança
msg := "pokemon %s não encontrado (status=%d)"
fmt.Printf(msg, id, code)
```

---

## 13. Controle de fluxo

### 13.1 Indent error flow — retorne cedo

Retorne ao encontrar erro; o caminho feliz fica sem aninhamento.

```go
// ✅
func (s *PokemonService) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
    if id == "" {
        return nil, domain.ErrInvalidInput
    }

    pokemon, err := s.pokemonRepo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("buscar pokemon %s: %w", id, err)
    }

    return pokemon, nil
}

// ❌ — caminho feliz aninhado no else
func (s *PokemonService) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
    if id != "" {
        pokemon, err := s.pokemonRepo.GetByID(ctx, id)
        if err == nil {
            return pokemon, nil
        } else {
            return nil, fmt.Errorf("buscar pokemon %s: %w", id, err)
        }
    } else {
        return nil, domain.ErrInvalidInput
    }
}
```

### 13.2 Switch vs if-else

Prefira `switch` quando houver 3+ ramos sobre a mesma variável:

```go
// ✅
switch resp.StatusCode {
case http.StatusOK:
    return parseResponse(body)
case http.StatusNotFound:
    return nil, domain.ErrNotFound
case http.StatusUnauthorized:
    return nil, domain.ErrInvalidToken
case http.StatusConflict:
    return nil, domain.ErrAlreadyExists
default:
    return nil, fmt.Errorf("status inesperado do auth-service: %d", resp.StatusCode)
}
```

### 13.3 Reduza o aninhamento

```go
// ✅ — retorno antecipado reduz aninhamento
for _, item := range items {
    if item == nil {
        continue
    }
    if !item.IsValid() {
        continue
    }
    process(item)
}

// ❌ — pirâmide de doom
for _, item := range items {
    if item != nil {
        if item.IsValid() {
            process(item)
        }
    }
}
```

### 13.4 Raw string literals para evitar escapes

```go
// ✅ — legível
wantError := `unknown name:"test"`

// ❌ — difícil de ler
wantError := "unknown name:\"test\""
```

---

## 14. Testes orientados a tabela

### 14.1 Quando usar table-driven

Use quando há múltiplos casos que variam apenas nos inputs/outputs com a mesma lógica de verificação:

```go
func TestAuthServiceClientErrorMapping(t *testing.T) {
    tests := []struct {
        name      string
        method    func(client *httpclient.AuthServiceClient) error
        transport roundTripFunc
        wantErr   error
    }{
        {
            name: "login retorna ErrInvalidCredentials para 401",
            method: func(c *httpclient.AuthServiceClient) error {
                _, err := c.Login(context.Background(), "user@test.com", "wrongpass")
                return err
            },
            transport: mockStatusCode(http.StatusUnauthorized),
            wantErr:   domain.ErrInvalidCredentials,
        },
        // ... mais casos
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client := httpclient.NewAuthServiceClient(
                "http://test",
                httpclient.WithHTTPClient(&http.Client{Transport: tt.transport}),
            )
            err := tt.method(client)
            require.ErrorIs(t, err, tt.wantErr)
        })
    }
}
```

### 14.2 Quando NÃO usar table-driven

- Casos com setup/teardown muito diferente entre si
- Casos com assertions complexas e únicas
- Testes de caminho feliz com muitas verificações de campos (use teste explícito)

```go
// ✅ — teste explícito para caminho feliz com múltiplas assertions
func TestAuthServiceClientSignupSuccess(t *testing.T) {
    // setup específico
    client := setupClientWithMockResponse(t, validSignupResponse)
    
    session, err := client.Signup(context.Background(), validSignupRequest)
    
    require.NoError(t, err)
    assert.Equal(t, "expected-token", session.AccessToken)
    assert.Equal(t, "user@test.com", session.UserEmail)
    assert.False(t, session.ExpiresAt.IsZero())
}
```

### 14.3 Naming dos subtests

- Use `t.Run(tt.name, ...)` com nomes descritivos do comportamento esperado
- Padrão: `"método retorna ErroX para condição Y"`

### 14.4 Campos desnecessários na struct de teste

Nunca declare campos na struct de teste que não são usados no corpo do `t.Run`:

```go
// ❌ — wantPath declarado mas nunca usado
tests := []struct {
    name      string
    wantErr   error
    wantPath  string  // declara mas ninguém usa
}{ ... }

// ✅ — apenas o que é verificado
tests := []struct {
    name    string
    wantErr error
}{ ... }
```

---

## 15. Functional Options

Use quando um construtor tem muitos parâmetros opcionais:

```go
type ClientOption func(*AuthServiceClient)

func WithTimeout(d time.Duration) ClientOption {
    return func(c *AuthServiceClient) {
        c.httpClient.Timeout = d
    }
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
    return func(c *AuthServiceClient) {
        c.httpClient = httpClient
    }
}

func NewAuthServiceClient(baseURL string, opts ...ClientOption) *AuthServiceClient {
    c := &AuthServiceClient{
        baseURL: baseURL,
        httpClient: &http.Client{Timeout: 10 * time.Second},
    }
    for _, opt := range opts {
        opt(c)
    }
    return c
}

// Uso
client := NewAuthServiceClient(
    "http://auth-service:8080",
    WithTimeout(30 * time.Second),
)

// Em testes
client := NewAuthServiceClient(
    "http://test",
    WithHTTPClient(&http.Client{Transport: mockTransport}),
)
```

---

## 16. Panic e init()

### 16.1 Não use panic em código de produção

```go
// ❌
func run(args []string) {
    if len(args) == 0 {
        panic("argumento obrigatório")
    }
}

// ✅ — retorne erro e deixe o chamador decidir
func run(args []string) error {
    if len(args) == 0 {
        return errors.New("argumento obrigatório")
    }
    return nil
}

func main() {
    if err := run(os.Args[1:]); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

Exceção: inicialização do programa com `template.Must`, `regexp.MustCompile`:

```go
// ✅ — pânico apenas na inicialização, não durante requisições
var statusTemplate = template.Must(template.New("status").Parse(statusHTML))
```

Em testes, use `t.Fatal` / `t.FailNow`, nunca panic:

```go
// ✅
f, err := os.CreateTemp("", "test")
if err != nil {
    t.Fatal("falha ao configurar o teste:", err)
}

// ❌
if err != nil {
    panic("falha ao configurar o teste")
}
```

### 16.2 Evite init()

`init()` é executado muito cedo, dificulta testes e esconde dependências.

```go
// ❌
func init() {
    go backgroundWorker()    // goroutine sem controle
    db = openDatabase()      // efeito colateral oculto
}

// ✅ — inicialização explícita no construtor
func New(ctx context.Context, dbURL string) (*Service, error) {
    db, err := openDatabase(dbURL)
    if err != nil {
        return nil, fmt.Errorf("abrir banco de dados: %w", err)
    }
    svc := &Service{db: db}
    go svc.run(ctx) // goroutine com ciclo de vida controlado
    return svc, nil
}
```

O único uso aceitável de `init()`: registrar drivers/codecs que é impossível fazer de outro jeito.

---

## 17. Variáveis globais mutáveis

Evite — use injeção de dependência:

```go
// ❌ — mutação global dificulta testes paralelos
var _timeNow = time.Now

func sign(msg string) string {
    return signWithTime(msg, _timeNow())
}

// ✅ — injetável
type Signer struct {
    now func() time.Time
}

func NewSigner() *Signer {
    return &Signer{now: time.Now}
}

func (s *Signer) Sign(msg string) string {
    return signWithTime(msg, s.now())
}

// Em testes:
func TestSign(t *testing.T) {
    s := NewSigner()
    s.now = func() time.Time { return fixedTime }
    assert.Equal(t, want, s.Sign(give))
}
```

### Nomes canônicos para variáveis globais não-exportadas

Use prefixo `_` para distinguir globais de locais:

```go
// ✅
var (
    _defaultClient = &http.Client{}
    _logger        = log.New(os.Stderr, "", 0)
)
```

---

## 18. Linting e ferramentas

### 18.1 Ferramentas essenciais

| Ferramenta | Propósito |
|------------|-----------|
| `golangci-lint` | Lint agregado — único runner |
| `go vet` | Análise básica do compilador |
| `gofmt` / `goimports` | Formatação + imports automáticos |
| `errcheck` | Detecta erros descartados |
| `staticcheck` | Análise estática avançada |
| `revive` | Regras de estilo adicionais |

### 18.2 Linha de comprimento suave

Limite suave: **99 caracteres** (Uber). Quebre antes, mas não é rígido.  
Exceção: URLs longas em comentários não precisam ser quebradas.

### 18.3 Funções Printf customizadas

Se criar wrapper de Printf, termine com `f`:

```go
// ✅ — go vet detecta automaticamente
func Wrapf(format string, args ...interface{}) error { ... }

// ❌
func Wrap(format string, args ...interface{}) error { ... }
```

---

## 19. Decisões específicas do projeto

Estas são decisões já tomadas na Pokedex Platform que divergem ou complementam os guias:

### 19.1 Imports: 3 grupos (não 2 como Uber, não 4 como Google)

Ver [seção 4.1](#41-três-grupos-decisão-do-projeto).

### 19.2 Slices retornados via API: sempre `make([]T, 0)`

Ver [seção 6.3](#63-slices--decisão-do-projeto).

### 19.3 Interface compliance: `var _ I = (*T)(nil)` ao final do arquivo

Ver [seção 7.3](#73-asserção-de-interface-em-tempo-de-compilação--obrigatório-no-projeto).

### 19.4 Constantes de tempo: sempre `time.Duration`

Ver [seção 12.4](#124-timeduration-para-constantes-de-tempo--decisão-do-projeto).

### 19.5 Comentários: sempre Português Brasil

Ver [seção 3.1](#31-regras-do-projeto-pt-br).

### 19.6 Estrutura hexagonal: direção das dependências

```
adapters/inbound/http  →  ports/inbound (use cases)  ←  service
                                                              ↓
                          ports/outbound  ←  adapters/outbound/{postgres,http}
```

- Handlers HTTP **não** dependem de clients concretos
- Erros normalizados no adapter outbound ou na camada de aplicação
- Novas entidades de domínio em `domain/`, nunca em `ports/`

### 19.7 gRPC exported: `GRPC` não `Grpc`

Seguir Google Style Guide:

```go
// ✅
type GRPCClient struct{}
func NewGRPCHandler() *GRPCHandler {}

// ❌
type GrpcClient struct{}
func NewGrpcHandler() *GrpcHandler {}
```

### 19.8 Nomes oficiais dos serviços

- `mobile-bff` — BFF voltado para experiência do cliente
- `pokemon-catalog-service` — fonte canônica do catálogo
- `auth-service` — autenticação e ciclo de vida de token

Nunca use `pokedex-service` (legado).

---

## Referências

- [Uber Go Style Guide (PT-BR)](https://github.com/uber-go/guide/blob/master/style.md)
- [Google Go Style Guide](https://google.github.io/styleguide/go/guide)
- [Google Go Style Decisions](https://google.github.io/styleguide/go/decisions)
- [Google Go Best Practices](https://google.github.io/styleguide/go/best-practices)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
