# Cheatsheet — Sintaxe Go

Referência rápida de sintaxe Go. Imprima e mantenha ao lado enquanto programa.

## Estrutura mínima

```go
package main

import "fmt"

func main() {
    fmt.Println("hello")
}
```

Executar: `go run main.go`

## Variáveis e constantes

```go
var nome string = "Pikachu"      // declaração completa
var nivel int                     // zero value = 0
pokemon := "Charizard"            // declaração curta (dentro de funções)
const maxHP = 100                 // constante compile-time
const (
    StatusAtivo  = iota           // 0
    StatusInativo                 // 1
    StatusBanido                  // 2
)
```

## Tipos básicos

```go
var (
    a int     = 42
    b int64   = 9999999999
    c float64 = 3.14
    d string  = "hello"
    e bool    = true
    f byte    = 'A'           // alias para uint8
    g rune    = '世'          // alias para int32 (code point Unicode)
)
```

## Conversão de tipos

```go
var i int = 42
var f float64 = float64(i)   // conversão explícita sempre obrigatória
var s string = strconv.Itoa(i)
n, err := strconv.Atoi("42")
```

## Ponteiros

```go
x := 42
p := &x          // endereço de x
*p = 99          // desreferencia e modifica x
// Go NÃO tem aritmética de ponteiros
```

## Funções

```go
func soma(a, b int) int {
    return a + b
}

// retorno múltiplo — o padrão Go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("divisao por zero")
    }
    return a / b, nil
}

// função variádica
func concat(sep string, partes ...string) string {
    return strings.Join(partes, sep)
}
```

## Controle de fluxo

```go
// if — parênteses NUNCA, chaves SEMPRE
if x > 10 {
    // ...
} else if x > 5 {
    // ...
}

// if com init — o "guard let" do Go
if err := doSomething(); err != nil {
    return err
}

// for — o único loop de Go (3 formas)
for i := 0; i < 10; i++ { }             // clássico
for x < 10 { x++ }                       // como while
for { }                                  // loop infinito

// for range
for i, v := range slice { }              // índice + valor
for k, v := range map { }                // chave + valor
for _, v := range slice { }              // só valor (descarta índice)

// switch — sem break automático
switch status {
case 200:
    // ...
case 404, 410:                           // múltiplos valores
    // ...
default:
    // ...
}

// switch sem expressão (equivale a if-else if)
switch {
case x > 10:
    // ...
case x > 5:
    // ...
}

// defer — executa ao final da função (LIFO)
f, _ := os.Open("arquivo.txt")
defer f.Close()
```

## Slices

```go
s := []int{1, 2, 3}                      // literal
s := make([]int, 0, 10)                  // len=0, cap=10
s = append(s, 4)                         // adiciona
s = append(s, 5, 6, 7)                   // adiciona vários
s = append(s, outroSlice...)             // spread
copia := make([]int, len(s))
copy(copia, s)                           // cópia independente
fatia := s[1:3]                          // [lo:hi], hi é exclusivo
// API/JSON: use make([]T, 0) em vez de nil
```

## Maps

```go
m := map[string]int{"a": 1, "b": 2}
m := make(map[string]int)
m["c"] = 3
delete(m, "a")
v, ok := m["c"]                          // comma ok idiom
```

## Structs

```go
type Pokemon struct {
    Nome  string
    Level int
    Tipos []string
}

p := Pokemon{Nome: "Pikachu", Level: 25} // com nomes de campo
p := Pokemon{}                            // zero value
p2 := &Pokemon{Nome: "Charizard"}        // ponteiro para struct

// método com pointer receiver
func (p *Pokemon) LevelUp() {
    p.Level++
}
```

## Interfaces

```go
type Descritor interface {
    Descreve() string
}

// Pokemon implementa Descritor automaticamente (implícito!)
func (p Pokemon) Descreve() string {
    return fmt.Sprintf("%s (lvl %d)", p.Nome, p.Level)
}

// verificação em tempo de compilação
var _ Descritor = (*Pokemon)(nil)
var _ Descritor = Pokemon{}              // se método tem value receiver

// type assertion
d, ok := obj.(Descritor)

// type switch
switch v := obj.(type) {
case Pokemon:
    // v é Pokemon
case *Pokemon:
    // v é *Pokemon
}
```

## Erros

```go
// sentinel errors
var ErrNotFound = errors.New("nao encontrado")

// erro customizado
type ValidationError struct {
    Campo string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("campo %s invalido", e.Campo)
}

// wrapping
if err != nil {
    return fmt.Errorf("buscar pokemon: %w", err)   // %w = wrap (inspecionável)
}

// verificação
if errors.Is(err, ErrNotFound) { }
var vErr *ValidationError
if errors.As(err, &vErr) { }
```

## Goroutines e channels

```go
go fn()                                  // inicia goroutine

ch := make(chan int)                     // unbuffered
ch := make(chan int, 10)                 // buffered

ch <- 42                                 // envia
v := <-ch                                // recebe
v, ok := <-ch                            // ok=false se fechado
close(ch)                                // fecha (envia zero value + ok=false)

// select
select {
case v := <-ch1:
    // ...
case ch2 <- 42:
    // ...
case <-time.After(5 * time.Second):
    // timeout
}

// range sobre channel (até ser fechado)
for v := range ch { }

// channel directions
func envia(ch chan<- int) { ch <- 42 }   // send-only
func recebe(ch <-chan int) { v := <-ch } // receive-only
```

## Context

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

ctx, cancel := context.WithCancel(context.Background())
ctx = context.WithValue(ctx, chave, valor)

select {
case <-ctx.Done():
    return ctx.Err()
case result := <-resultCh:
    return result, nil
}
```

## Testes

```go
// arquivo: pokemon_test.go
func TestLevelUp(t *testing.T) {
    tests := []struct {
        name     string
        levelIn  int
        levelOut int
    }{
        {"sobe nivel", 25, 26},
        {"nivel maximo", 100, 100},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            p := Pokemon{Level: tt.levelIn}
            p.LevelUp()
            if p.Level != tt.levelOut {
                t.Errorf("esperado %d, obteve %d", tt.levelOut, p.Level)
            }
        })
    }
}

func BenchmarkLevelUp(b *testing.B) {
    p := Pokemon{Level: 1}
    for i := 0; i < b.N; i++ {
        p.LevelUp()
    }
}
```

## JSON

```go
type PokemonDTO struct {
    Nome  string   `json:"nome"`
    Level int      `json:"nivel"`
    Tipos []string `json:"tipos,omitempty"`
    interno string `json:"-"`            // omitido
}

data, _ := json.Marshal(p)
json.Unmarshal(data, &p)
```

## Pacotes e módulos

```go
// go.mod
module github.com/usuario/projeto
go 1.22

// imports
import (
    "fmt"
    "github.com/usuario/lib"
)

// visibilidade: Maiúscula = público, minúscula = privado
func Publica() {}    // exportada
func privada() {}    // não exportada
```

## Ferramentas CLI

```bash
go mod init modulo           # criar módulo
go mod tidy                  # limpar dependências
go build                     # compilar
go build -o bin/app ./cmd/   # compilar com nome
go run .                     # compilar e executar
go run main.go               # executar arquivo
go test ./...                 # todos os testes
go test -v -run TestNome    # teste específico
go test -race ./...          # detector de data races
go test -cover ./...         # cobertura
go vet ./...                  # análise estática
go fmt ./...                  # formatar código
go doc fmt.Println            # documentação
GOOS=linux GOARCH=amd64 go build  # cross-compile
```
