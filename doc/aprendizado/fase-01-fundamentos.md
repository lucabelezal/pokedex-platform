# Fase 01 — Fundamentos

O objetivo desta fase é você conseguir escrever e executar programas Go simples
que usam variáveis, constantes, tipos básicos, funções e ponteiros.

## Ambiente Go

Antes de escrever código, você precisa do compilador Go instalado. O mais comum
é baixar do site oficial e seguir as instruções do seu sistema operacional.

Depois de instalado, verifique que está funcionando:

```bash
go version
# go version go1.22.0 darwin/arm64
```

O comando `go` é a ferramenta central do ecossistema. Com ele você compila, testa,
formata, gerencia dependências e muito mais. Por enquanto, você só precisa saber:

| Comando | O que faz |
|---------|-----------|
| `go run arquivo.go` | Compila e executa (útil para desenvolvimento) |
| `go build` | Compila e gera um binário (útil para produção) |
| `go fmt ./...` | Formata todo o código (use sempre antes de commitar) |

Go é **opinativo** sobre formatação. Não existe discussão sobre estilo — `gofmt`
resolve tudo. Configure seu editor para formatar ao salvar o arquivo.

---

## Estrutura de um programa Go

Crie um arquivo `main.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Pokémon!")
}
```

Execute com `go run main.go`.

O que está acontecendo:

- **`package main`** — Todo arquivo `.go` pertence a um pacote. O pacote `main`
  é especial: ele define um programa executável (não uma biblioteca).
- **`import "fmt"`** — Importa o pacote `fmt` da biblioteca padrão, que fornece
  funções de formatação e impressão.
- **`func main()`** — O ponto de entrada do programa. Quando você executa o binário,
  esta é a primeira função chamada. Não recebe argumentos nem retorna nada.
- **Chaves `{` na mesma linha** — Go não aceita quebra de linha antes da chave.
  Isso é imposto pelo compilador, não é preferência estilística.

### Swift vs Go

```swift
// Swift — estrutura mínima
print("Hello, Pokémon!")
// pode escrever solto no Playground ou em um arquivo .swift
```

```go
// Go — estrutura obrigatória
package main
import "fmt"
func main() { fmt.Println("Hello, Pokémon!") }
// TODO programa executável precisa de package main e func main()
```

**Atenção:** Go é mais verboso que Swift para o mínimo. Mas essa estrutura
explícita escala bem para projetos grandes. Não há "magia" — tudo está declarado.

---

## Variáveis

Go oferece duas formas de declarar variáveis:

### `var` — declaração explícita

```go
package main

import "fmt"

func main() {
    var nome string          // declara com tipo; valor inicial = "" (zero value)
    var nivel int = 25       // declara com tipo e valor
    var ativo = true         // declara sem tipo; Go infere bool
    var x, y int = 10, 20    // múltiplas variáveis do mesmo tipo

    fmt.Println(nome, nivel, ativo, x, y)
}
```

### `:=` — declaração curta (short variable declaration)

```go
package main

import "fmt"

func main() {
    pokemon := "Pikachu"      // infere string
    hp := 100                 // infere int
    taxa := 0.75              // infere float64

    fmt.Println(pokemon, hp, taxa)
}
```

### Quando usar cada uma

| `var` | `:=` |
|-------|------|
| Declaração em nível de pacote (fora de funções) | Só funciona **dentro de funções** |
| Quando o tipo precisa ser explícito | Quando o tipo é óbvio pelo valor |
| Quando você quer o zero value sem valor inicial | Quando você tem um valor inicial |
| `var handler slog.Handler` (zero value útil) | `pokemon := "Pikachu"` (óbvio) |

**Atenção:** `:=` **não** funciona fora de funções. Isto dá erro de compilação:

```go
package main

pokemon := "Pikachu"  // ERRO: short declaration não permitida aqui

func main() {
    pokemon := "Pikachu"  // OK: dentro de função
}
```

### Swift vs Go

```swift
// Swift
var nome: String            // precisa ser inicializada antes de usar
var nome = "Pikachu"        // inferência de tipo
let nome = "Pikachu"        // constante (let)
```

```go
// Go
var nome string             // zero value = "" — pode usar imediatamente
var nome = "Pikachu"        // inferência de tipo (raramente usado assim)
nome := "Pikachu"           // forma idiomática
nome := "Pikachu"           // variável, não constante. Go usa const para constantes
```

**Atenção:** Em Go, `:=` cria uma variável **mutável**. Diferente de Swift onde
`let` cria constantes. Para constantes em Go, use `const`.

---

## Tipos básicos

Go tem um conjunto pequeno de tipos fundamentais:

```go
package main

import "fmt"

func main() {
    var (
        inteiro   int     = -42          // int (32 ou 64 bits dependendo da arquitetura)
        int8bit   int8    = 127          // -128 a 127
        int16bit  int16   = 32767        // -32768 a 32767
        int32bit  int32   = 2147483647   // ...
        int64bit  int64   = 9223372036854775807
        uinteiro  uint    = 42           // unsigned (0 a ...)
        decimal   float64 = 3.14159      // float de 64 bits (padrão)
        float32b  float32 = 3.14         // float de 32 bits
        texto     string  = "Pokémon"    // string UTF-8
        booleano  bool    = true         // true ou false (não aceita 0/1)
        caractere byte    = 'A'          // byte = uint8 (um caractere ASCII)
        unicode   rune    = '世'         // rune = int32 (um code point Unicode)
    )
    fmt.Println(inteiro, int8bit, int16bit, int32bit, int64bit)
    fmt.Println(uinteiro, decimal, float32b, texto, booleano)
    fmt.Println(caractere, unicode)
}
```

Alguns pontos importantes:

- **`int` vs `int64`** — `int` tem o tamanho da palavra da máquina (32 ou 64 bits).
  Use `int` para contadores e tamanhos. Use `int64` quando o tamanho exato importa
  (ex: IDs de banco de dados).
- **`float64` é o padrão** — Quando você escreve `x := 3.14`, o tipo é `float64`.
  `float32` raramente é usado.
- **`byte` e `rune`** — São aliases. `byte` = `uint8`, `rune` = `int32`.
  Use `byte` para dados binários, `rune` para caracteres Unicode.
- **Go não tem `char`** — Use `byte` para ASCII ou `rune` para Unicode.

### Swift vs Go

```swift
// Swift — muitos tipos numéricos com nomes familiares
let x: Int = 42          // Int (64 bits em plataformas 64-bit)
let y: Float = 3.14      // Float (32 bits)
let z: Double = 3.14159  // Double (64 bits) — padrão para decimais
let c: Character = "A"   // Character (um caractere)
```

```go
// Go — tipos similares, mas float64 é o padrão
var x int = 42           // int
var y float32 = 3.14     // float32
var z float64 = 3.14159  // float64 — padrão para decimais
var c byte = 'A'         // byte (não existe Character como em Swift)
```

---

## Zero values

Toda variável em Go é inicializada com um **valor padrão** se você não definir um.
Isso é intencional — não existe "variável não inicializada" em Go.

```go
package main

import "fmt"

func main() {
    var i int         // 0
    var f float64     // 0.0
    var s string      // "" (string vazia, não nil)
    var b bool        // false
    var p *int        // nil

    fmt.Printf("int: %d\n", i)
    fmt.Printf("float64: %f\n", f)
    fmt.Printf("string: %q\n", s)
    fmt.Printf("bool: %t\n", b)
    fmt.Printf("pointer: %v\n", p)
}
```

Saída:

```
int: 0
float64: 0.000000
string: ""
bool: false
pointer: <nil>
```

Este conceito é fundamental em Go. Significa que:

- Você pode declarar uma variável e usá-la imediatamente.
- Um `sync.Mutex` com zero value já está pronto para uso (não precisa inicializar).
- Um `var buf bytes.Buffer` já aceita escrita sem `make()`.

### Swift vs Go

```swift
// Swift — variáveis NÃO inicializadas precisam ser definidas antes do uso
var x: Int
print(x)  // ERRO: variable 'x' used before being initialized

// Swift — Optionals começam como nil por padrão
var x: Int?  // nil
```

```go
// Go — toda variável tem um valor inicial
var x int
fmt.Println(x)  // 0 — funciona, sem erro
```

**Atenção:** `nil` em Go não é um valor universal como em Swift. Em Go, `nil`
só pode ser atribuído a ponteiros, slices, maps, channels, interfaces e funções.
Você **não pode** escrever `var s string = nil` — string tem zero value `""`.

---

## Constantes

Constantes em Go são valores conhecidos em **tempo de compilação**. Diferente de
Swift, você não pode criar uma constante a partir do retorno de uma função.

```go
package main

import "fmt"

const pokemonInicial = "Bulbasaur"
const maxHP = 100

// bloco de constantes
const (
    statusAtivo  = 1
    statusInativo = 2
)

func main() {
    fmt.Println(pokemonInicial, maxHP, statusAtivo, statusInativo)
}
```

### `iota` — enumerados

Go não tem `enum`. Para criar uma sequência de constantes, use `iota`:

```go
package main

import "fmt"

const (
    domingo = iota     // 0
    segunda            // 1
    terca              // 2
    quarta             // 3
    quinta             // 4
    sexta              // 5
    sabado             // 6
)

// iota reseta a cada bloco const
const (
    _ = iota           // descarta 0
    KB = 1 << (10 * iota)  // 1 << 10 = 1024
    MB                     // 1 << 20 = 1048576
    GB                     // 1 << 30 = 1073741824
)

func main() {
    fmt.Println(domingo, segunda, sabado)
    fmt.Println(KB, MB, GB)
}
```

### Swift vs Go

```swift
// Swift — enum poderoso, com valores associados e raw values
enum DiaSemana: Int {
    case domingo = 0
    case segunda, terca, quarta, quinta, sexta, sabado
}
```

```go
// Go — use iota para sequências simples. Para enums mais complexos, use type + const
type Status int

const (
    StatusAtivo Status = iota + 1   // começa em 1
    StatusInativo
    StatusBanido
)
```

**Atenção:** `iota` só funciona dentro de blocos `const`. Ele incrementa a cada
linha do bloco e reseta no próximo bloco `const`.

---

## Conversão de tipos

Go é rigoroso com tipos. **Toda conversão é explícita.** Não existe conversão
implícita, nem entre tipos numéricos similares.

```go
package main

import "fmt"

func main() {
    var i int = 42
    var f float64 = float64(i)   // conversão explícita: int → float64
    var u uint = uint(f)         // conversão explícita: float64 → uint

    fmt.Println(i, f, u)
}
```

Isto **não compila**:

```go
var i int = 42
var f float64 = i   // ERRO: cannot use i (type int) as type float64
```

Para converter entre strings e números, use o pacote `strconv`:

```go
package main

import (
    "fmt"
    "strconv"
)

func main() {
    // número → string
    s := strconv.Itoa(42)            // "42"
    s2 := strconv.FormatFloat(3.14, 'f', 2, 64)  // "3.14"

    // string → número
    n, err := strconv.Atoi("42")     // 42, nil
    if err != nil {
        fmt.Println("erro:", err)
    }

    fmt.Println(s, s2, n)
}
```

### Swift vs Go

```swift
// Swift — conversões também são explícitas na maioria dos casos
let i: Int = 42
let f: Double = Double(i)   // explícita (similar a Go)
let s: String = String(i)   // explícita
let n: Int? = Int("42")     // retorna Optional (similar ao Atoi)
```

```go
// Go — SEMPRE explícito, sem exceções
var i int = 42
var f float64 = float64(i)  // obrigatório
s := strconv.Itoa(i)        // string
n, err := strconv.Atoi("42") // (int, error)
```

---

## Funções

Funções em Go são declaradas com `func`. Parâmetros e retornos são tipados.

```go
package main

import "fmt"

// função simples
func cumprimenta() {
    fmt.Println("Olá, treinador!")
}

// função com parâmetros e retorno
func soma(a int, b int) int {
    return a + b
}

// parâmetros do mesmo tipo podem compartilhar a declaração
func multiplica(a, b int) int {
    return a * b
}

// retorno múltiplo — o padrão Go para success/failure
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("divisão por zero")
    }
    return a / b, nil
}

func main() {
    cumprimenta()
    fmt.Println(soma(3, 4))
    fmt.Println(multiplica(5, 6))

    resultado, err := divide(10, 3)
    if err != nil {
        fmt.Println("Erro:", err)
        return
    }
    fmt.Println("Resultado:", resultado)
}
```

### Swift vs Go

```swift
// Swift — funções e parâmetros
func soma(_ a: Int, _ b: Int) -> Int { a + b }

// Swift — retorno de erro com throws
func divide(_ a: Double, _ b: Double) throws -> Double {
    guard b != 0 else { throw DivError.zero }
    return a / b
}
```

```go
// Go — funções são mais verbosas, mas consistentes
func soma(a, b int) int { return a + b }

// Go — erro é valor de retorno, não exceção
func divide(a, b float64) (float64, error) { ... }
```

**Atenção:** Go não tem `throws`. Erro é um valor como outro qualquer.
Você retorna `(resultado, error)` e a pessoa que chamou decide o que fazer.
Esta é uma das diferenças mais fundamentais entre Go e Swift.

---

## Ponteiros

Ponteiros armazenam o **endereço de memória** de um valor. Em Go, ponteiros são
seguros: não há aritmética de ponteiros e o garbage collector gerencia a memória.

```go
package main

import "fmt"

func main() {
    x := 42

    var p *int = &x    // p aponta para x. & = "endereço de"

    fmt.Println("x:", x)      // 42
    fmt.Println("p:", p)      // 0xc0000140a8 (endereço)
    fmt.Println("*p:", *p)    // 42  (* = desreferencia)

    *p = 99                  // modifica x através do ponteiro
    fmt.Println("x agora:", x) // 99
}
```

### Por que usar ponteiros?

**1. Modificar um valor dentro de uma função:**

```go
package main

import "fmt"

func levelUp(hp *int) {
    *hp = *hp + 10             // modifica o valor original
}

func main() {
    hp := 100
    levelUp(&hp)
    fmt.Println("HP:", hp)     // 110
}
```

**2. Evitar cópia de structs grandes:**

```go
func processa(p *Pokemon) {   // passa referência, não cópia
    // ...
}
```

### Swift vs Go

```swift
// Swift — `inout` para modificar parâmetro
func levelUp(hp: inout Int) {
    hp += 10
}
var hp = 100
levelUp(hp: &hp)
```

```go
// Go — ponteiro explícito
func levelUp(hp *int) {
    *hp += 10
}
hp := 100
levelUp(&hp)
```

**Atenção:** Go **não** tem `weak` ou `unowned`. O garbage collector resolve
ciclos de referência sem intervenção do programador. Isso simplifica muito o código.

---

## O pacote `fmt`

O pacote `fmt` é usado em praticamente todo programa Go. Aqui estão as funções
que você mais vai usar:

```go
package main

import "fmt"

func main() {
    nome := "Pikachu"
    nivel := 25

    // Println — imprime com espaço entre argumentos e quebra de linha
    fmt.Println(nome, "nível", nivel)

    // Printf — formatação com placeholders
    fmt.Printf("%s está no nível %d\n", nome, nivel)

    // Sprintf — retorna string formatada (não imprime)
    mensagem := fmt.Sprintf("%s está no nível %d", nome, nivel)
    fmt.Println(mensagem)
}
```

### Placeholders mais comuns

| Verbo | Significado | Exemplo de saída |
|-------|-------------|-----------------|
| `%v` | Valor padrão | `fmt.Printf("%v", p)` → `{Pikachu 25}` |
| `%+v` | Valor com nomes de campo | `{Nome:Pikachu Nivel:25}` |
| `%#v` | Sintaxe Go do valor | `main.Pokemon{Nome:"Pikachu", Nivel:25}` |
| `%T` | Tipo do valor | `main.Pokemon` |
| `%d` | Inteiro decimal | `42` |
| `%s` | String | `Pikachu` |
| `%q` | String com aspas | `"Pikachu"` |
| `%f` | Float | `3.141590` |
| `%t` | Boolean | `true` |
| `%p` | Ponteiro (endereço) | `0xc0000140a8` |
| `%%` | Símbolo de porcentagem | `%` |

---

## Exercícios da Fase 01

Crie um arquivo `exercicios_fase01.go` e resolva:

### 1. Calculadora de dano

Escreva uma função `calculaDano(ataque, defesa int) int` que retorna `ataque - defesa`,
com dano mínimo de 1. Se o resultado for menor que 1, retorne 1.

<details>
<summary>Gabarito</summary>

```go
package main

import "fmt"

func calculaDano(ataque, defesa int) int {
    dano := ataque - defesa
    if dano < 1 {
        return 1
    }
    return dano
}

func main() {
    fmt.Println(calculaDano(50, 20))  // 30
    fmt.Println(calculaDano(10, 30))  // 1
    fmt.Println(calculaDano(5, 5))    // 1
}
```
</details>

### 2. Conversão segura

Escreva uma função `stringParaInt(s string) (int, error)` que converte uma string
para inteiro. Se a conversão falhar, retorne o erro apropriado.

<details>
<summary>Gabarito</summary>

```go
package main

import (
    "fmt"
    "strconv"
)

func stringParaInt(s string) (int, error) {
    n, err := strconv.Atoi(s)
    if err != nil {
        return 0, fmt.Errorf("stringParaInt: %w", err)
    }
    return n, nil
}

func main() {
    fmt.Println(stringParaInt("42"))
    fmt.Println(stringParaInt("abc"))
}
```
</details>

### 3. Ponteiro na prática

Escreva uma função `cura(pokemon *int, quantidade int)` que aumenta o HP de um
pokémon. O HP máximo é 100. Teste a função.

<details>
<summary>Gabarito</summary>

```go
package main

import "fmt"

func cura(pokemon *int, quantidade int) {
    *pokemon += quantidade
    if *pokemon > 100 {
        *pokemon = 100
    }
}

func main() {
    hp := 60
    cura(&hp, 20)
    fmt.Println("HP:", hp)  // 80
    cura(&hp, 50)
    fmt.Println("HP:", hp)  // 100 (não passa de 100)
}
```
</details>

---

**Próxima fase:** [Fase 02 — Controle & Coleções](fase-02-controle-colecoes.md)
