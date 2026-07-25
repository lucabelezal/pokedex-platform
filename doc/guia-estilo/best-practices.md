# Melhores Práticas

[Visão Geral](README.md) | [Guia](guide.md) | [Decisões](decisions.md) | [Melhores Práticas](best-practices.md) | [Regras do Projeto](regras-projeto.md)

**Status:** `[Informativo]` — padrões recomendados que resolvem problemas comuns.
Não são canônicos, mas seu uso é encorajado para manter a base de código uniforme.

---

## Testes

### Testes baseados em tabela (table-driven)

Testes table-driven com [subtestes] são o padrão para escrever testes quando a
lógica central é repetitiva com múltiplas entradas e saídas.

[subtestes]: https://blog.golang.org/subtests

Use a convenção: slice de structs chamada `tests` e cada caso chamado `tt`.
Explicite entradas e saídas com os prefixos `give` e `want`.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim — repetição manual:
host, port, err := net.SplitHostPort("192.0.2.0:8000")
require.NoError(t, err)
assert.Equal(t, "192.0.2.0", host)

host, port, err = net.SplitHostPort(":8000")
require.NoError(t, err)
assert.Equal(t, "", host)
```

```go
// Bom — table-driven:
tests := []struct {
    give     string
    wantHost string
    wantPort string
}{
    {give: "192.0.2.0:8000", wantHost: "192.0.2.0", wantPort: "8000"},
    {give: ":8000", wantHost: "", wantPort: "8000"},
    {give: "192.0.2.0:http", wantHost: "192.0.2.0", wantPort: "http"},
}

for _, tt := range tests {
    t.Run(tt.give, func(t *testing.T) {
        host, port, err := net.SplitHostPort(tt.give)
        require.NoError(t, err)
        assert.Equal(t, tt.wantHost, host)
        assert.Equal(t, tt.wantPort, port)
    })
}
```

### Evitar complexidade desnecessária em tabelas

Testes de tabela **não** devem ser usados quando houver lógica condicional
complexa ou ramificação dentro dos subtestes. Nesses casos, prefira funções
`Test...` separadas.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim — tabela complexa com flags condicionais:
tests := []struct {
    give     string
    want     string
    wantErr  error
    shouldCallX bool
    shouldCallY bool
    giveXResponse string
    giveYResponse string
}{ ... }
```

```go
// Bom — testes separados e focados:
func TestShouldCallX(t *testing.T) {
    ctrl := gomock.NewController(t)
    xMock := xmock.NewMockX(ctrl)
    xMock.EXPECT().Call().Return("XResponse", nil)

    got, err := DoComplexThing("inputX", xMock, nil)
    require.NoError(t, err)
    assert.Equal(t, "want", got)
}

func TestShouldCallYAndFail(t *testing.T) {
    ctrl := gomock.NewController(t)
    yMock := ymock.NewMockY(ctrl)
    yMock.EXPECT().Call().Return("", errors.New("Y failed"))

    _, err := DoComplexThing("inputY", nil, yMock)
    assert.EqualError(t, err, "Y failed")
}
```

Alguns ideais para tabelas de teste:

- Focar na unidade mais estreita de comportamento
- Minimizar a "profundidade do teste" e evitar asserções condicionais
- Garantir que todos os campos da tabela sejam usados em todos os testes
- Se o corpo do teste for curto, um único campo `shouldErr` é aceitável

### Nomes de funções de teste

Funções de teste devem seguir o padrão `Test<NomeDaFuncao>`. Use underscore
para separar conceitos em nomes de teste.

```go
// Bom:
func TestSplitHostPort(t *testing.T) { ... }
func TestAuthService_Login(t *testing.T) { ... }
func TestPokemonService_List_WithFilters(t *testing.T) { ... }
```

### Cobertura

- Cobertura mínima: **75%**
- Cobertura ideal: **90%**
- Execute `make coverage` para gerar relatório

### Testes paralelos

Testes paralelos devem declarar `tt := tt` dentro do loop para evitar captura
incorreta da variável do loop:

```go
for _, tt := range tests {
    tt := tt
    t.Run(tt.give, func(t *testing.T) {
        t.Parallel()
        // ...
    })
}
```

### Exemplos testáveis

Use funções `Example` em arquivos `*_test.go` para documentar o uso de APIs
públicas. Exemplos são executados como testes e aparecem na documentação.

```go
// Bom:
func ExamplePokemonService_List() {
    svc := NewPokemonService(repo)
    page, err := svc.List(context.Background(), DefaultParams)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(len(page.Items))
    // Output: 0
}
```

---

## Desempenho

### Preferir `strconv` a `fmt`

Para conversões de tipos primitivos para string e vice-versa, prefira funções
do pacote `strconv`. Elas são mais rápidas e não alocam tanto quanto `fmt.Sprintf`.

| Ruim | Bom |
|------|-----|
| `fmt.Sprintf("%d", num)` | `strconv.Itoa(num)` |
| `fmt.Sscanf(s, "%d", &num)` | `strconv.Atoi(s)` |

### Especificar capacidade de containers

Sempre que souber o tamanho aproximado, especifique a capacidade de maps e slices
com `make` para evitar realocações.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim:
items := []string{}
for _, v := range data {
    items = append(items, v.Name)
}
```

```go
// Bom:
items := make([]string, 0, len(data))
for _, v := range data {
    items = append(items, v.Name)
}
```

```go
// Bom:
m := make(map[string]int, expectedSize)
```

### Evitar conversões repetidas de string para bytes

Evite converter repetidamente entre `string` e `[]byte` em loops. Converta uma
vez e reutilize.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim:
for _, s := range items {
    data := []byte(s)
    w.Write(data)
}
```

```go
// Bom:
for _, s := range items {
    w.Write([]byte(s)) // já é eficiente se w.Write não mantém referência
}

// Melhor ainda — use io.WriteString quando possível:
for _, s := range items {
    io.WriteString(w, s)
}
```

### Preferir `time.Time` e `time.Duration`

Use os tipos do pacote `time` para manipular tempo. Não use inteiros como
timestamps ou durações em milissegundos.

```go
// Bom:
const timeout = 30 * time.Second
var createdAt time.Time
```

---

## Padrões

### Opções funcionais

Use o padrão de opções funcionais para argumentos opcionais em construtores e
APIs públicas que você prevê precisar expandir, especialmente se já houver
três ou mais argumentos.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim:
db.Open(addr, true, zap.NewNop())
db.Open(addr, false /* cache */, log)
```

```go
// Bom:
db.Open(addr)
db.Open(addr, db.WithLogger(log))
db.Open(addr, db.WithCache(false), db.WithLogger(log))
```

Implementação recomendada com interface `Option` e método não exportado:

```go
type options struct {
    cache  bool
    logger *zap.Logger
}

type Option interface {
    apply(*options)
}

type cacheOption bool

func (c cacheOption) apply(opts *options) {
    opts.cache = bool(c)
}

func WithCache(c bool) Option {
    return cacheOption(c)
}

type loggerOption struct {
    Log *zap.Logger
}

func (l loggerOption) apply(opts *options) {
    opts.logger = l.Log
}

func WithLogger(log *zap.Logger) Option {
    return loggerOption{Log: log}
}

func Open(addr string, opts ...Option) (*Connection, error) {
    options := options{
        cache:  defaultCache,
        logger: zap.NewNop(),
    }
    for _, o := range opts {
        o.apply(&options)
    }
    // ...
}
```

Este padrão permite que opções sejam comparadas em testes e implementem outras
interfaces como `fmt.Stringer`.

### `defer` para limpeza

Use `defer` para garantir que recursos sejam liberados, independentemente de
como a função termine. Posicione o `defer` logo após a aquisição do recurso.

```go
// Bom:
f, err := os.Open(file)
if err != nil {
    return fmt.Errorf("abrir arquivo: %w", err)
}
defer f.Close()
```

### Enums com `iota`

Inicie enums em 1 (não em 0) para que o valor zero seja distinguível de um
valor intencional. Isso evita bugs onde uma variável não inicializada parece
ter um valor válido.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim:
const (
    TypeNormal   = iota // 0
    TypeFire            // 1
    TypeWater           // 2
)
```

```go
// Bom:
const (
    TypeNormal = iota + 1 // 1
    TypeFire              // 2
    TypeWater             // 3
)
```

### Evitar `init()`

Evite usar funções `init()` sempre que possível. Elas tornam o código difícil
de entender, testar e depurar. Prefira inicialização explícita.

Se `init()` for inevitável:
- Mantenha-o simples e determinístico
- Nunca inicie goroutines em `init()`
- Nunca acesse o sistema de arquivos ou rede
- Documente claramente os efeitos colaterais

### Sair em `main`

Apenas a função `main` (ou `run` chamada por `main`) deve chamar `os.Exit` ou
`log.Fatal`. Bibliotecas e serviços nunca devem encerrar o programa.

```go
// Bom — apenas em main.go:
func main() {
    if err := run(); err != nil {
        slog.Error("servidor encerrado com erro", "error", err)
        os.Exit(1)
    }
}
```

### Sair apenas uma vez

Garanta que `os.Exit` seja chamado no máximo uma vez. Múltiplas chamadas podem
causar comportamento imprevisível com `defer`.

### `nil` é um slice válido

Não retorne um slice vazio explícito (`[]T{}`) a menos que necessário (ex: JSON).
`nil` é um slice válido em Go e funciona com `range`, `len`, `append`.

```go
// Bom (uso interno):
func findItems() []string {
    return nil // aceitável
}

// Bom (retorno via API/JSON):
func listItems() []string {
    return make([]string, 0) // garante "[]" no JSON
}
```

### Tags de campos em structs serializadas

Use tags de struct consistentes para serialização. Mantenha o formato padrão com
espaços entre tags.

```go
// Bom:
type Pokemon struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Type string `json:"type,omitempty"`
}
```

### Evitar incorporação de tipos em structs públicas

Evite incorporar tipos em structs públicas. A incorporação vaza detalhes de
implementação e dificulta a evolução da API.

| Ruim | Bom |
|------|-----|
| | |

```go
// Ruim — incorpora sync.Mutex nos métodos públicos:
type Counter struct {
    sync.Mutex
    value int
}
```

```go
// Bom — mutex é campo privado:
type Counter struct {
    mu    sync.Mutex
    value int
}
```

### Evitar usar nomes embutidos

Não use nomes que colidem com identificadores pré-declarados da linguagem Go
(`string`, `error`, `int`, `len`, `cap`, `make`, `copy`, `append`, `close`,
`delete`, `new`, `panic`, `recover`, `real`, `imag`, `complex`).

| Ruim | Bom |
|------|-----|
| `var error string` | `var errorMessage string` |
| `var string string` | `var str string` |
