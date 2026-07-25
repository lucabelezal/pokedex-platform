# Decisões de Estilo

[Visão Geral](README.md) | [Guia](guide.md) | [Decisões](decisions.md) | [Melhores Práticas](best-practices.md) | [Regras do Projeto](regras-projeto.md)

**Status:** `[Normativo]` — regras com justificativa detalhada, usadas como referência por revisores.
Podem evoluir com novas versões do Go.

**Nota:** Este documento é subordinado ao [Guia de Estilo](guide.md). Em caso de
conflito, o Guia prevalece.

---

## Estrutura de código

### Ordem de importação em grupos

Imports devem ser organizados em três grupos, separados por linha em branco:

1. Biblioteca padrão
2. Dependências de terceiros
3. Módulos internos do projeto

| Ruim | Bom |
|------|-----|
| | |

```go
// Bom:
import (
    "context"
    "fmt"
    "time"

    "github.com/sony/gobreaker/v2"

    in "pokedex-platform/core/bff/mobile-bff/internal/ports/inbound"
    "pokedex-platform/core/bff/mobile-bff/internal/domain"
)
```

Ao usar alias de importação, mantenha consistência: o mesmo alias deve ser usado
em todos os arquivos que importam o mesmo pacote.

```go
// Bom — alias padronizado no projeto:
httpadapter "pokedex-platform/core/bff/mobile-bff/internal/adapters/inbound/http"
httpclient  "pokedex-platform/core/bff/mobile-bff/internal/adapters/outbound/http"
outbound    "pokedex-platform/core/bff/mobile-bff/internal/ports/outbound"
```

### Agrupar declarações semelhantes

Agrupe declarações relacionadas usando blocos `var` e `const`.

```go
// Bom:
const (
    defaultTimeout = 30 * time.Second
    maxRetries     = 3
)
```

### Ordem de funções

Funções exportadas devem aparecer primeiro no arquivo, seguidas pelas funções
não exportadas. Dentro de cada grupo, funções que chamam outras devem vir antes
das funções chamadas (ordem top-down de leitura).

```go
// Bom:
func (s *Service) List(ctx context.Context) ([]Pokemon, error) {  // exportada, chama fetchAll
    return s.fetchAll(ctx)
}

func (s *Service) fetchAll(ctx context.Context) ([]Pokemon, error) {  // não exportada
    // ...
}
```

### Reduzir aninhamento

Reduza o aninhamento sempre que possível lidando com casos de erro ou
condições especiais primeiro, retornando cedo ou continuando o loop.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim:
for _, v := range data {
    if v.F1 == 1 {
        v = process(v)
        if err := v.Call(); err == nil {
            v.Send()
        } else {
            return err
        }
    } else {
        log.Printf("Inválido v: %v", v)
    }
}
```

```go
// Bom:
for _, v := range data {
    if v.F1 != 1 {
        log.Printf("Inválido v: %v", v)
        continue
    }
    v = process(v)
    if err := v.Call(); err != nil {
        return err
    }
    v.Send()
}
```

### Else desnecessário

Evite `else` quando o bloco `if` termina com `return`, `break` ou `continue`.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim:
if err != nil {
    return err
} else {
    return process(data)
}
```

```go
// Bom:
if err != nil {
    return err
}
return process(data)
```

### Reduzir o escopo de variáveis

Declare variáveis o mais próximo possível do seu uso.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim:
var user User
user, err = repo.FindByID(ctx, id)
if err != nil {
    return err
}
```

```go
// Bom:
user, err := repo.FindByID(ctx, id)
if err != nil {
    return err
}
```

Em testes paralelos com table-driven, declare `tt := tt` dentro do loop para
evitar captura incorreta da variável do loop.

```go
// Bom:
for _, tt := range tests {
    tt := tt // necessário para t.Parallel()
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel()
        // ...
    })
}
```

### Declarações de variáveis globais

Variáveis globais devem ser declaradas no topo do arquivo, antes das funções.

Variáveis globais não exportadas devem ser prefixadas com `_` para deixar claro
que são globais.

```go
// Bom:
var _defaultTimeout = 30 * time.Second

func (c *Client) Do(req *http.Request) (*http.Response, error) {
    // usa _defaultTimeout
}
```

---

## Declarações e inicialização

### Declarações de variáveis locais

Use `:=` (declaração curta) para a maioria das variáveis locais. Use `var`
quando:

- O valor zero não for suficiente e nenhum valor inicial for atribuído
- O tipo for necessário para clareza (ex: interface vazia)

| Ruim | Bom |
|------|-----|
| `var s string = "foo"` | `s := "foo"` |
| `var count = 0` | `var count int` (zero value explícito) |

### Inicialização de structs

#### Usar nomes de campos

Sempre use nomes de campos ao inicializar structs. A única exceção são
structs com 1-2 campos onde o significado é óbvio e o contexto é muito pequeno
(ex: `Point{X: 10, Y: 20}` pode ser `Point{10, 20}` em testes).

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim:
user := User{"Maria", 30, true}
```

```go
// Bom:
user := User{
    Name:   "Maria",
    Age:    30,
    Active: true,
}
```

#### Omitir campos de valor zero

Omita campos com valor zero ao inicializar structs, a menos que forneçam contexto
significativo.

```go
// Bom:
user := User{
    Name: "Maria",
    // Age e Active são omitidos (zero values)
}
```

#### Usar `&T{}` vs `new(T)`

Prefira `&T{}` a `new(T)` para consistência e para permitir inicialização com campos.

```go
// Bom:
user := &User{Name: "Maria"}
```

#### Structs de valor zero

Use `var t T` para structs de valor zero em vez de `t := T{}`.

| Ruim | Bom |
|------|-----|
| `t := T{}` | `var t T` |

### Inicialização de mapas

Prefira `make(map[T]U)` ou `map[T]U{}` para mapas vazios. Use `make` com
capacidade quando souber o tamanho aproximado.

```go
// Bom:
m := make(map[string]int)
m := make(map[string]int, expectedSize)
```

### Slices

`nil` é um slice válido. Use `nil` para slices vazios, exceto ao retornar
via API/JSON — nesse caso, use `make([]T, 0)` para garantir serialização como
array vazio `[]` em vez de `null`.

| Ruim | Bom |
|------|-----|
| `return nil, nil` para lista vazia via API | `return make([]Pokemon, 0), nil` |

```go
// Bom:
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
    pokemons := make([]dto.PokemonDTO, 0) // garante JSON "[]" e não "null"
    // ...
}
```

### Literais de string raw

Use literais de string raw (`` ` ``) para evitar escape quando a string contiver
muitas aspas ou caracteres especiais.

```go
// Bom:
const sql = `
    SELECT id, name, type
    FROM pokemons
    WHERE region = $1
`
```

### Evitar parâmetros isolados

Evite parâmetros isolados que não deixam claro seu propósito na chamada.
Prefira tipos nomeados ou opções funcionais.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim:
func NewServer(addr string, bool, *zap.Logger) *Server
```

```go
// Bom:
type ServerOption interface { ... }

func WithTLS(enabled bool) ServerOption { ... }
func WithLogger(log *zap.Logger) ServerOption { ... }

func NewServer(addr string, opts ...ServerOption) *Server
```

---

## Interfaces

### Verificar conformidade de interface

Verifique a conformidade de interface em tempo de compilação usando a declaração
`var _ Interface = (*Impl)(nil)`. Isso garante que um tipo implementa a interface
esperada e falha no build se o contrato for quebrado.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim — sem verificação de conformidade:
type Handler struct { ... }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { ... }
```

```go
// Bom — verificação em tempo de compilação:
type Handler struct { ... }

var _ http.Handler = (*Handler)(nil)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) { ... }
```

Esta verificação deve estar presente em cada adapter e service do projeto:

```go
// Em service/pokemon_service.go:
var _ inbound.PokemonUseCase = (*PokemonService)(nil)

// Em adapters/outbound/http/pokemon_catalog_client.go:
var _ outbound.PokemonRepository = (*PokemonCatalogServiceRepository)(nil)

// Em adapters/outbound/postgres/pokemon_repository.go:
var _ outbound.PokemonRepository = (*PostgresPokemonRepository)(nil)
```

O lado direito deve ser o valor zero do tipo: `nil` para ponteiros, slices e maps;
struct vazia para tipos struct.

### Receptores e interfaces

Ao implementar uma interface com um tipo concreto, use receptor **ponteiro** se
qualquer método da interface precisar modificar o receiver. Use receptor **valor**
apenas quando o tipo for imutável e pequeno.

```go
// Bom — receptor ponteiro (métodos modificam estado):
func (s *PokemonService) List(ctx context.Context) (*PokemonPage, error) { ... }

// Bom — receptor valor (tipo imutável e pequeno):
func (p Point) Distance(other Point) float64 { ... }
```

### Ponteiros para interfaces

**Nunca** use ponteiro para interface. Interfaces já são tipos de referência.
Um ponteiro para interface quase nunca é necessário.

| Ruim | Bom |
|------|-----|
| `func process(p *PokemonRepository)` | `func process(p PokemonRepository)` |
| `var repo *outbound.PokemonRepository` | `var repo outbound.PokemonRepository` |

---

## Erros

### Tipos de erro

Use erros sentinela (variáveis `var ErrXxx = errors.New(...)`) para erros
de domínio que precisam ser inspecionados pelos chamadores com `errors.Is`.

```go
// Bom — em internal/domain/errors.go:
var (
    ErrPokemonNotFound       = errors.New("pokemon nao encontrado")
    ErrUserAlreadyExists     = errors.New("usuario ja existe")
    ErrFavoriteAlreadyExists = errors.New("favorito ja existe")
    ErrInvalidCredentials    = errors.New("credenciais invalidas")
    ErrInvalidToken          = errors.New("token invalido")
    ErrServiceUnavailable    = errors.New("servico temporariamente indisponivel")
)
```

Use tipos de erro customizados (implementando `Error() string`) apenas quando
precisar de metadados adicionais (ex: código de status HTTP, campo de validação).

### Envoltório de erros (wrapping)

Ao propagar erros, use `fmt.Errorf` para adicionar contexto. Escolha entre `%w`
ou `%v`:

- **`%w`**: O chamador deve poder inspecionar o erro subjacente com `errors.Is`/`errors.As`.
  Padrão recomendado para a maioria dos casos.
- **`%v`**: O erro subjacente deve ser obscurecido. O chamador não poderá fazer match.
  Use quando o erro interno for detalhe de implementação.

Mantenha o contexto sucinto, evitando frases como "falha ao" que se acumulam:

| Ruim | Bom |
|------|-----|
| `fmt.Errorf("falha ao criar novo armazenamento: %w", err)` | `fmt.Errorf("novo armazenamento: %w", err)` |

```go
// Bom — %w para erros inspecionáveis:
s, err := store.New()
if err != nil {
    return fmt.Errorf("novo armazenamento: %w", err)
}

// Bom — %v para obscurecer detalhes internos:
data, err := decrypt(raw)
if err != nil {
    return fmt.Errorf("decrypt: %v", err)
}
```

Erros de domínio devem ser documentados como parte do contrato da função quando
expostos via `%w`.

### Nomeação de erros

- Variáveis sentinela exportadas: `ErrXxx` (ex: `ErrPokemonNotFound`)
- Variáveis sentinela não exportadas: `errXxx` (ex: `errInvalidFormat`)
- Tipos de erro customizados exportados: `XxxError` (ex: `ValidationError`)
- Variáveis locais de erro: `err` (sempre)

### Tratar erros uma vez

Um erro deve ser tratado apenas uma vez. Não registre (log) um erro e depois
o retorne — escolha um ou outro.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim — loga e retorna o mesmo erro:
if err != nil {
    log.Printf("erro ao buscar usuário: %v", err)
    return err
}
```

```go
// Bom — retorna com contexto:
if err != nil {
    return fmt.Errorf("buscar usuário: %w", err)
}

// Bom — loga e retorna um erro diferente:
if err != nil {
    log.Printf("erro inesperado: %v", err)
    return domain.ErrServiceUnavailable
}
```

### Error strings

Mensagens de erro devem ser **minúsculas**, sem ponto final, e não devem
começar com o nome do pacote.

| Ruim | Bom |
|------|-----|
| `errors.New("Usuário não encontrado.")` | `errors.New("usuario nao encontrado")` |
| `errors.New("domain: pokemon nao encontrado")` | `errors.New("pokemon nao encontrado")` |

### Indent error flow

Trate o caminho de erro primeiro e retorne imediatamente. O fluxo principal
deve seguir o caminho feliz com indentação mínima.

```go
// Bom:
func process(data []byte) error {
    if len(data) == 0 {
        return fmt.Errorf("dados vazios")
    }
    if err := validate(data); err != nil {
        return fmt.Errorf("validar dados: %w", err)
    }
    return save(data)
}
```

### Lidar com falhas de asserção de tipo

Sempre use a forma de dois valores (`value, ok`) para asserções de tipo.
Nunca use a forma de um valor que causa panic.

| Ruim | Bom |
|------|-----|
| `user := ctx.Value(key).(User)` | `user, ok := ctx.Value(key).(User)` |

```go
// Bom:
user, ok := ctx.Value(userKey).(User)
if !ok {
    return fmt.Errorf("tipo inesperado no contexto")
}
```

### Não entrar em pânico

Não use `panic` para erros de negócio ou condições esperadas. `panic` é reservado
para erros irrecuperáveis que indicam bug do programador (ex: invariante violada).

```go
// Ruim — não use panic para erro de validação:
if input == "" {
    panic("input vazio")
}

// Bom — retorne o erro:
if input == "" {
    return fmt.Errorf("input nao pode ser vazio")
}
```

---

## Concorrência

### Goroutines

#### Aguardar goroutines finalizarem

Sempre garanta que goroutines sejam finalizadas ou tenham um mecanismo de
cancelamento. Use `context.Context` para propagar cancelamento e
`sync.WaitGroup` para aguardar conclusão.

```go
// Bom:
var wg sync.WaitGroup
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

for _, item := range items {
    wg.Add(1)
    go func(item Item) {
        defer wg.Done()
        process(ctx, item)
    }(item)
}
wg.Wait()
```

#### Sem goroutines em `init()`

Nunca inicie goroutines dentro de funções `init()`. Isso torna o comportamento
do programa imprevisível e difícil de testar.

### Canais

O tamanho de um canal deve ser **um** (bufferizado) ou **nenhum** (não bufferizado).
Evite tamanhos arbitrários. Canais não bufferizados fornecem sincronização
implícita.

| Ruim | Bom |
|------|-----|
| `ch := make(chan int, 10)` | `ch := make(chan int)` ou `ch := make(chan int, 1)` |

### Mutex com valor zero é válido

`sync.Mutex` e `sync.RWMutex` são válidos com valor zero. Não use ponteiros
para mutexes desnecessariamente.

```go
// Bom:
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}
```

### Evitar globais mutáveis

Evite variáveis globais mutáveis. Use injeção de dependência em vez de estado
global compartilhado.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim:
var db *sql.DB

func GetUser(id string) (*User, error) {
    // usa db global
}
```

```go
// Bom:
type UserService struct {
    db *sql.DB
}

func (s *UserService) GetUser(id string) (*User, error) {
    // usa s.db
}
```

### Usar `time` para manipular tempo

Use `time.Time` e `time.Duration` para manipular tempo. Não use inteiros para
representar timestamps ou durações.

| Ruim | Bom |
|------|-----|
| `func wait(seconds int)` | `func wait(d time.Duration)` |
| `var createdAt int64` | `var createdAt time.Time` |
