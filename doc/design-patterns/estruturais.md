# Padrões Estruturais

Como Go compõe estruturas para conectar componentes — adaptando interfaces,
simplificando subsistemas e adicionando comportamento sem modificar código existente.

---

## Adapter

### Propósito

Permitir que dois componentes com interfaces incompatíveis trabalhem juntos.
O Adapter traduz uma interface para outra.

### Filosofia Go

Adapter é um dos padrões mais naturais em Go. Como interfaces são implícitas,
um adapter é simplesmente uma struct que **wrapeia** o adaptee e implementa a
interface alvo. Sem herança, sem declaração de conformidade.

Na arquitetura hexagonal, **todo adapter outbound é um Adapter**: traduz a
interface do repository para chamadas HTTP, SQL ou gRPC.

### Código idiomático

```go
package main

import "fmt"

// Interface que o cliente espera
type PokemonRepository interface {
    Buscar(id string) (string, error)
}

// API externa com interface incompatível
type PokeAPI struct{}

func (p PokeAPI) FetchByID(id int) (string, error) {
    return fmt.Sprintf("dados do pokemon %d", id), nil
}

// Adapter — traduz FetchByID(int) → Buscar(string)
type PokeAPIAdapter struct {
    api PokeAPI
}

func (a PokeAPIAdapter) Buscar(id string) (string, error) {
    // converte string → int, chama a API, retorna no formato esperado
    var intID int
    fmt.Sscanf(id, "%d", &intID)
    return a.api.FetchByID(intID)
}

func main() {
    var repo PokemonRepository = PokeAPIAdapter{api: PokeAPI{}}
    resultado, _ := repo.Buscar("25")
    fmt.Println(resultado)
}
```

### Onde aparece no projeto

```go
// Adapter que traduz Pokémon Catalog Service (HTTP) → PokemonRepository (interface)
type PokemonCatalogServiceRepository struct {
    baseURL    string
    httpClient *http.Client
}

func (r *PokemonCatalogServiceRepository) GetByID(ctx context.Context, id string) (*domain.Pokemon, error) {
    // traduz: chamada HTTP → retorno da interface
    resp, err := r.httpClient.Get(r.baseURL + "/pokemons/" + id)
    // ... parse JSON → domain.Pokemon
}
```

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/adapter && go run .
```

---

## Bridge

### Propósito

Separar uma abstração da sua implementação para que ambas possam variar
independentemente.

### Filosofia Go

Em Go, Bridge se resolve com **duas hierarquias de interfaces** que se comunicam.
A abstração (`Computer`) contém uma referência à implementação (`Printer`), e
você pode trocar qualquer lado em runtime.

Bridge é menos comum em Go do que simplesmente injetar uma interface — que já
resolve 80% dos casos de "separar abstração de implementação".

### Código idiomático

```go
package main

import "fmt"

// Implementação — interface
type TipoPokemon interface {
    Multiplicador() float64
}

type Fogo struct{}
func (Fogo) Multiplicador() float64 { return 1.5 }

type Agua struct{}
func (Agua) Multiplicador() float64 { return 1.0 }

// Abstração — struct que contém a implementação
type Ataque struct {
    nome string
    tipo TipoPokemon   // bridge: abstração separada da implementação
}

func (a Ataque) Dano(poder int) float64 {
    return float64(poder) * a.tipo.Multiplicador()
}

func main() {
    lancachamas := Ataque{nome: "Lança-chamas", tipo: Fogo{}}
    jatoagua := Ataque{nome: "Jato d'Água", tipo: Agua{}}

    fmt.Printf("%s: %.0f\n", lancachamas.nome, lancachamas.Dano(90)) // 135
    fmt.Printf("%s: %.0f\n", jatoagua.nome, jatoagua.Dano(90))       // 90
}
```

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/bridge && go run .
```

---

## Composite

### Propósito

Compor objetos em estruturas de árvore para representar hierarquias todo-parte.
Clientes tratam objetos individuais e composições uniformemente.

### Filosofia Go

Composite é natural em Go: defina uma **interface comum** para folhas e
composites. O composite contém um **slice da interface** (`[]Component`).

A recursão é natural — cada método do composite itera sobre os filhos e delega.

### Código idiomático

```go
package main

import (
    "fmt"
    "strings"
)

// Interface comum
type ItemPokedex interface {
    Peso() int
    String() string
}

// Folha
type Pokemon struct {
    Nome    string
    PesoKg  int
}

func (p Pokemon) Peso() int       { return p.PesoKg }
func (p Pokemon) String() string  { return fmt.Sprintf("%s (%dkg)", p.Nome, p.PesoKg) }

// Composite
type TimePokemon struct {
    Nome    string
    Membros []ItemPokedex
}

func (t TimePokemon) Peso() int {
    total := 0
    for _, m := range t.Membros {
        total += m.Peso()
    }
    return total
}

func (t TimePokemon) String() string {
    var nomes []string
    for _, m := range t.Membros {
        nomes = append(nomes, m.String())
    }
    return fmt.Sprintf("[%s: %s]", t.Nome, strings.Join(nomes, ", "))
}

func main() {
    pikachu := Pokemon{"Pikachu", 6}
    charizard := Pokemon{"Charizard", 90}

    time := TimePokemon{
        Nome:    "Time Ash",
        Membros: []ItemPokedex{pikachu, charizard},
    }

    fmt.Println(time.String())       // [Time Ash: Pikachu (6kg), Charizard (90kg)]
    fmt.Println("Peso total:", time.Peso(), "kg")  // 96kg
}
```

A beleza: `TimePokemon` pode conter outros `TimePokemon`, formando hierarquias profundas
— tudo através da mesma interface `ItemPokedex`.

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/composite && go run .
```

---

## Decorator

### Propósito

Adicionar comportamento a um objeto dinamicamente, sem modificar sua estrutura.
Decorators envolvem o objeto original e adicionam funcionalidade antes/depois.

### Filosofia Go

Decorator é **onipresente** em Go. Todo middleware HTTP é um decorator.
Toda função que wrapped um `io.Reader` com buffer é um decorator.

A implementação é trivial: uma struct que implementa a mesma interface do
objeto decorado e contém uma referência a ele.

### Código idiomático — middleware HTTP

```go
package main

import (
    "fmt"
    "net/http"
    "time"
)

// Handler é o componente base
func pokemonHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte(`{"nome": "Pikachu"}`))
}

// Decorator de logging — wrapped o handler, adiciona log
func loggingDecorator(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        inicio := time.Now()
        next(w, r)
        fmt.Printf("%s %s levou %v\n", r.Method, r.URL.Path, time.Since(inicio))
    }
}

// Decorator de autenticação — wrapped o handler, verifica auth
func authDecorator(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Authorization") == "" {
            http.Error(w, "não autorizado", http.StatusUnauthorized)
            return
        }
        next(w, r)
    }
}

func main() {
    // Empilha decorators: auth → logging → handler
    handler := authDecorator(loggingDecorator(pokemonHandler))
    http.HandleFunc("/pokemon", handler)
    http.ListenAndServe(":8080", nil)
}
```

### Decorator com `io.Reader`

```go
// Decorator que conta bytes lidos — wrapped io.Reader, mesma interface
type CountingReader struct {
    reader    io.Reader
    bytesRead int64
}

func (c *CountingReader) Read(p []byte) (int, error) {
    n, err := c.reader.Read(p)
    c.bytesRead += int64(n)
    return n, err
}

// Uso: wrapped um reader comum com contagem
reader := &CountingReader{reader: os.Stdin}
```

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/decorator && go run .
```

---

## Facade

### Propósito

Fornecer uma interface simplificada para um subsistema complexo. A Facade
esconde a complexidade interna atrás de uma API de alto nível.

### Filosofia Go

Em Go, Facade é uma **struct que agrega múltiplos componentes** e expõe métodos
que orquestram chamadas entre eles. É o ponto de entrada para um subsistema.

Na arquitetura hexagonal, um **Service** frequentemente atua como Facade sobre
múltiplos repositories e serviços externos.

### Código idiomático

```go
package main

import "fmt"

// Subsistema complexo — vários componentes
type PokemonRepo struct{}
func (PokemonRepo) Buscar(id string) string { return "Pikachu" }

type CacheService struct{}
func (CacheService) Get(key string) (string, bool) { return "", false }
func (CacheService) Set(key, value string) {}

type AuditLog struct{}
func (AuditLog) Registrar(acao string) { fmt.Println("LOG:", acao) }

// Facade — API simples sobre o subsistema
type PokemonService struct {
    repo  PokemonRepo
    cache CacheService
    log   AuditLog
}

func (s *PokemonService) BuscarPokemon(id string) string {
    // tenta cache
    if cached, ok := s.cache.Get(id); ok {
        s.log.Registrar("cache hit: " + id)
        return cached
    }

    // busca no repository
    pokemon := s.repo.Buscar(id)
    s.cache.Set(id, pokemon)
    s.log.Registrar("busca: " + id)
    return pokemon
}

func main() {
    svc := &PokemonService{
        repo:  PokemonRepo{},
        cache: CacheService{},
        log:   AuditLog{},
    }
    fmt.Println(svc.BuscarPokemon("25"))
}
```

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/facade && go run .
```

---

## Flyweight

### Propósito

Compartilhar objetos para suportar grandes quantidades de objetos com baixo
consumo de memória. Separa estado intrínseco (compartilhado) de extrínseco
(por instância).

### Filosofia Go

Flyweight é **raro em Go** para APIs backend. É mais útil em jogos ou renderização
gráfica — cenários com milhares de objetos similares.

Em Go, implementa-se com um **map como cache** em uma factory. O estado intrínseco
(mesmo para todos) fica no flyweight; o estado extrínseco (único por instância)
fica no objeto cliente.

### Código idiomático

```go
package main

import "fmt"

// Flyweight — estado intrínseco (compartilhado)
type TipoPokemon struct {
    Nome     string
    Elemento string
    Cor      string
}

// Factory com cache de flyweights
var cacheTipos = make(map[string]*TipoPokemon)

func GetTipoPokemon(nome string) *TipoPokemon {
    if t, ok := cacheTipos[nome]; ok {
        return t
    }
    // cria uma única vez
    t := &TipoPokemon{Nome: nome, Elemento: elemento(nome), Cor: cor(nome)}
    cacheTipos[nome] = t
    return t
}

// Cliente — estado extrínseco (único por instância)
type Pokemon struct {
    ID    string
    Level int
    Tipo  *TipoPokemon  // compartilhado via flyweight
}

func main() {
    // Todos os pokémons elétricos compartilham o mesmo *TipoPokemon
    p1 := Pokemon{ID: "025", Level: 25, Tipo: GetTipoPokemon("Elétrico")}
    p2 := Pokemon{ID: "100", Level: 15, Tipo: GetTipoPokemon("Elétrico")}

    fmt.Println(p1.Tipo == p2.Tipo) // true — mesmo ponteiro!
}
```

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/flyweight && go run .
```

---

## Proxy

### Propósito

Fornecer um substituto para outro objeto, controlando o acesso a ele.
O Proxy implementa a mesma interface do objeto real e adiciona controle
(acesso, cache, lazy loading, rate limiting).

### Filosofia Go

Proxy é natural em Go — mesma técnica do Decorator, mas com propósito diferente.
Enquanto Decorator **adiciona** comportamento, Proxy **controla** acesso.

Exemplos em Go: `http.Transport` com rate limiting, `sql.DB` como proxy para
conexões reais, cache proxy para API externa.

### Código idiomático — Rate limiting proxy

```go
package main

import (
    "fmt"
    "time"
)

// Interface comum
type PokemonAPI interface {
    Buscar(id string) string
}

// Serviço real
type PokeAPIService struct{}

func (PokeAPIService) Buscar(id string) string {
    return fmt.Sprintf("Dados reais de %s", id)
}

// Proxy com rate limiting
type RateLimitedProxy struct {
    real  PokemonAPI
    ultimoAcesso time.Time
    minInterval  time.Duration
}

func (p *RateLimitedProxy) Buscar(id string) string {
    if elapsed := time.Since(p.ultimoAcesso); elapsed < p.minInterval {
        return "limite de requisições excedido"
    }
    p.ultimoAcesso = time.Now()
    return p.real.Buscar(id)
}

func main() {
    api := &RateLimitedProxy{
        real:        PokeAPIService{},
        minInterval: 1 * time.Second,
    }

    fmt.Println(api.Buscar("25"))  // sucesso
    fmt.Println(api.Buscar("6"))   // bloqueado — muito rápido
}
```

### Proxy vs Decorator — qual usar?

| Situação | Padrão |
|----------|--------|
| Adicionar log/metrics a um handler | Decorator |
| Limitar taxa de chamadas a uma API | Proxy |
| Adicionar cache a uma operação | Proxy |
| Adicionar compressão a um writer | Decorator |

### Exemplo executável

```bash
cd ../../design-patterns/design-patterns-go/proxy && go run .
```
