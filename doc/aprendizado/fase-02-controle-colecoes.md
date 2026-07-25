# Fase 02 — Controle de Fluxo & Coleções

Nesta fase você vai dominar as estruturas de controle (`if`, `for`, `switch`, `defer`)
e as coleções de dados (`array`, `slice`, `map`, `string`). Ao final, você conseguirá
escrever funções que manipulam dados com a idiomática Go.

---

## `if` e `else`

O `if` em Go é similar ao de outras linguagens, com duas diferenças importantes:
**não usa parênteses** e **permite declaração inline** (init statement).

```go
package main

import "fmt"

func main() {
    nivel := 15

    // if simples — sem parênteses!
    if nivel > 10 {
        fmt.Println("Pokémon forte")
    }

    // if/else
    if nivel > 50 {
        fmt.Println("Muito forte")
    } else if nivel > 20 {
        fmt.Println("Forte")
    } else {
        fmt.Println("Em treinamento")
    }
}
```

### `if` com init — o "guard let" do Go

O `if` em Go aceita uma declaração antes da condição, separada por `;`:

```go
package main

import (
    "errors"
    "fmt"
)

func buscaPokemon(nome string) (string, error) {
    if nome == "Pikachu" {
        return "⚡ encontrado", nil
    }
    return "", errors.New("não encontrado")
}

func main() {
    if resultado, err := buscaPokemon("Pikachu"); err != nil {
        fmt.Println("Erro:", err)
        return
    } else {
        fmt.Println("Resultado:", resultado)
    }
}
```

A variável `resultado` e `err` só existem **dentro** do bloco `if/else`.
Fora dele, elas não são acessíveis. Isso reduz o escopo das variáveis — um
padrão idiomático em Go.

### Swift vs Go

```swift
// Swift — guard let
guard let resultado = try? buscaPokemon("Pikachu") else {
    print("Erro")
    return
}
print("Resultado:", resultado)
```

```go
// Go — if com init
if resultado, err := buscaPokemon("Pikachu"); err != nil {
    fmt.Println("Erro:", err)
    return
}
// resultado visível aqui (se o if não tiver else)
fmt.Println("Resultado:", resultado)
```

**Atenção:** Se você usa `return` dentro do `if`, o `else` é desnecessário.
Este é um padrão chamado **indent error flow** — o "happy path" do código fica
sem indentação extra:

```go
// Bom — happy path sem indentação
f, err := os.Open("arquivo.txt")
if err != nil {
    return err
}
defer f.Close()
// ... usa f ...

// Ruim — aninhamento desnecessário
f, err := os.Open("arquivo.txt")
if err == nil {
    defer f.Close()
    // ... usa f ...
}
return err
```

---

## `for` — o único loop de Go

Go tem **um** loop: `for`. Mas ele assume três formas diferentes.

### Forma 1 — clássica (init; condição; pós)

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
// 0 1 2 3 4
```

### Forma 2 — como `while` (só condição)

```go
nivel := 1
for nivel < 10 {
    nivel *= 2
}
fmt.Println("Nível final:", nivel) // 16
```

### Forma 3 — loop infinito

```go
for {
    // executa para sempre (use break para sair)
    if condicao {
        break
    }
}
```

### `continue` e `break`

```go
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue  // pula números pares
    }
    if i > 7 {
        break     // para no 7
    }
    fmt.Println(i)  // 1 3 5 7
}
```

### Swift vs Go

```swift
// Swift — multiple loop constructs
for i in 0..<5 { }          // range-based
while condicao { }           // while
repeat { } while condicao    // repeat-while
```

```go
// Go — só for, três formas
for i := 0; i < 5; i++ { }  // clássico
for condicao { }              // como while
for { }                       // infinito
```

**Atenção:** `++` e `--` em Go são **statements**, não expressões.
Você **não pode** escrever `x := i++`. Eles não retornam valor.

---

## `switch`

O `switch` em Go é diferente do `switch` em outras linguagens:
- **Não tem `break` automático** (o case termina automaticamente no fim do bloco)
- Aceita **múltiplos valores** por case
- Pode ser usado **sem expressão** (equivale a uma cadeia de `if/else`)

```go
package main

import "fmt"

func main() {
    tipo := "Elétrico"

    // switch com expressão
    switch tipo {
    case "Fogo":
        fmt.Println("🔥 Fraco contra Água")
    case "Água":
        fmt.Println("💧 Fraco contra Elétrico")
    case "Elétrico":
        fmt.Println("⚡ Fraco contra Terra")
    default:
        fmt.Println("Tipo desconhecido")
    }

    // múltiplos valores no mesmo case
    switch tipo {
    case "Fogo", "Elétrico", "Lutador":
        fmt.Println("Tipo ofensivo")
    case "Água", "Planta":
        fmt.Println("Tipo defensivo")
    }

    // switch sem expressão (como if/else if)
    nivel := 42
    switch {
    case nivel >= 80:
        fmt.Println("Lendário")
    case nivel >= 50:
        fmt.Println("Elite")
    case nivel >= 20:
        fmt.Println("Intermediário")
    default:
        fmt.Println("Iniciante")
    }
}
```

### Swift vs Go

```swift
// Swift — switch exaustivo, com fallthrough explícito
switch tipo {
case "Fogo":         print("🔥")
case "Água":         print("💧")
case "Elétrico":     print("⚡")
default:             print("?")
}
```

```go
// Go — break implícito, fallthrough explícito (raro)
switch tipo {
case "Fogo":         fmt.Println("🔥")
case "Água":         fmt.Println("💧")
case "Elétrico":     fmt.Println("⚡")
default:             fmt.Println("?")
}
```

**Atenção:** Em Go, se você **quiser** que o `switch` continue no próximo case
(equivalente ao comportamento padrão de Swift), use `fallthrough`.
Mas isso é raro na prática.

---

## `defer`

`defer` agenda a execução de uma função para **o final da função atual**.
É usado principalmente para limpeza de recursos: fechar arquivos, conexões,
liberar locks.

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    f, err := os.Create("pokemon.txt")
    if err != nil {
        fmt.Println("Erro:", err)
        return
    }
    defer f.Close()   // garante que f.Close() será chamado ao final de main

    fmt.Fprintln(f, "Pikachu")
    fmt.Fprintln(f, "Charizard")
    fmt.Println("Arquivo escrito com sucesso")
}
```

### `defer` executa em ordem LIFO (último a entrar, primeiro a sair)

```go
func main() {
    defer fmt.Println("primeiro defer")
    defer fmt.Println("segundo defer")
    defer fmt.Println("terceiro defer")
    fmt.Println("corpo da função")
}

// Saída:
// corpo da função
// terceiro defer
// segundo defer
// primeiro defer
```

### Swift vs Go

```swift
// Swift — defer idêntico ao Go
func processa() {
    let f = try abrirArquivo()
    defer { f.close() }
    // ... usa f ...
}
```

```go
// Go — mesmo conceito, mesma sintaxe
func processa() {
    f, _ := os.Open("arquivo.txt")
    defer f.Close()
    // ... usa f ...
}
```

**Atenção:** Os argumentos de uma função chamada com `defer` são **avaliados
imediatamente**, não no momento da execução:

```go
x := 1
defer fmt.Println(x)   // x = 1 é avaliado AGORA
x = 2
// Saída: 1 (não 2)
```

---

## Arrays

Arrays em Go têm **tamanho fixo** que faz parte do tipo. `[3]int` e `[4]int`
são tipos diferentes. Na prática, você quase sempre vai usar slices em vez de arrays.

```go
package main

import "fmt"

func main() {
    var times [3]string          // declara array de 3 strings
    times[0] = "Valor"
    times[1] = "Instinct"
    times[2] = "Mystic"

    niveis := [3]int{10, 25, 50} // literal

    fmt.Println(times)           // [Valor Instinct Mystic]
    fmt.Println(len(times))      // 3
    fmt.Println(niveis[1])       // 25
}
```

### Swift vs Go

```swift
// Swift — Array é sempre dinâmico, tamanho não faz parte do tipo
var times: [String] = ["Valor", "Instinct", "Mystic"]
```

```go
// Go — array tem tamanho fixo no tipo
var times [3]string
times = [3]string{"Valor", "Instinct", "Mystic"}
```

**Atenção:** Em Go, arrays são **value types**. Quando você atribui um array a
outra variável, todo o conteúdo é copiado:

```go
a := [3]int{1, 2, 3}
b := a       // cópia completa
b[0] = 99
fmt.Println(a[0])  // 1 (não foi alterado)
```

---

## Slices

Slice é a estrutura de dados mais importante em Go. É uma **visão flexível**
sobre um array interno (backing array). Diferente do array, o slice tem tamanho
dinâmico.

```go
package main

import "fmt"

func main() {
    // literal de slice (note a ausência de tamanho nos colchetes)
    pokemons := []string{"Pikachu", "Charizard", "Blastoise"}

    // make — cria slice com tamanho e capacidade
    numeros := make([]int, 3, 5)  // len=3, cap=5
    fmt.Println(len(numeros))     // 3
    fmt.Println(cap(numeros))     // 5

    // append — adiciona elementos (pode realocar)
    pokemons = append(pokemons, "Gengar")
    pokemons = append(pokemons, "Dragonite", "Mewtwo")  // múltiplos

    // spread de outro slice
    novos := []string{"Snorlax", "Lapras"}
    pokemons = append(pokemons, novos...)

    fmt.Println(pokemons)
}
```

### Anatomia de um slice

Internamente, um slice é uma struct com três campos:

```
type slice struct {
    pointer *Elemento  // aponta para o array interno
    length  int        // número de elementos visíveis
    capacity int       // capacidade total do array interno
}
```

Isso tem consequências importantes:

### Slice operator `[lo:hi]`

```go
s := []int{10, 20, 30, 40, 50}

a := s[1:3]   // [20, 30] — índice 1 até 3 (exclusivo)
b := s[:3]    // [10, 20, 30] — do início até 3
c := s[2:]    // [30, 40, 50] — do 2 até o final
```

### Cuidado: slices compartilham o array interno

```go
s := []int{10, 20, 30, 40, 50}
fatia := s[0:3]       // [10, 20, 30]
fatia[0] = 999        // modifica o array interno compartilhado!

fmt.Println(s)         // [999, 20, 30, 40, 50] — s também foi alterado!
fmt.Println(fatia)     // [999, 20, 30]
```

### `copy` — cria uma cópia independente

```go
s := []int{10, 20, 30, 40, 50}
copia := make([]int, len(s))
copy(copia, s)         // cópia profunda

copia[0] = 999
fmt.Println(s[0])      // 10 (não foi alterado)
```

### `append` e capacidade

```go
s := make([]int, 0, 4)   // len=0, cap=4
s = append(s, 1)         // len=1, cap=4 (sem realocação)
s = append(s, 2, 3, 4)   // len=4, cap=4 (cheio)
s = append(s, 5)         // len=5, cap=8 (dobrou a capacidade!)
```

Quando a capacidade se esgota, `append` aloca um novo array com o dobro da
capacidade (aproximadamente) e copia os elementos. Por isso, `append` **sempre
retorna o slice**, e você deve **sempre** reatribuir:

```go
s = append(s, elemento)   // correto
append(s, elemento)        // errado — perde o novo slice
```

### Swift vs Go

```swift
// Swift — Array é similar ao slice de Go
var pokemons = ["Pikachu", "Charizard"]
pokemons.append("Gengar")

let fatia = pokemons[0..<2]  // ArraySlice — cuidado com compartilhamento
```

```go
// Go — slice
pokemons := []string{"Pikachu", "Charizard"}
pokemons = append(pokemons, "Gengar")

fatia := pokemons[0:2]  // compartilha array interno — cuidado!
```

**Atenção para APIs JSON:** Ao retornar slices via API/JSON, use `make([]T, 0)`
em vez de `nil`:

```go
// JSON: [] vs null
var s1 []string           // nil → JSON: null
s2 := make([]string, 0)   // vazio → JSON: []
```

---

## Maps

Map é a estrutura chave-valor de Go. Similar ao `Dictionary` de Swift.

```go
package main

import "fmt"

func main() {
    // literal
    pokedex := map[string]int{
        "Pikachu":   25,
        "Charizard": 6,
        "Mewtwo":    150,
    }

    // make
    times := make(map[string]string)
    times["Ash"] = "Pallet"
    times["Misty"] = "Cerulean"

    // leitura
    nivel := pokedex["Pikachu"]    // 25
    fmt.Println(nivel)

    // chave inexistente → zero value do tipo
    desconhecido := pokedex["MissingNo"]  // 0 (int zero value)
    fmt.Println(desconhecido)

    // comma ok idiom — verifica se a chave existe
    nivel, existe := pokedex["Mewtwo"]
    if existe {
        fmt.Println("Mewtwo nível:", nivel)
    }

    // delete
    delete(pokedex, "Charizard")
    fmt.Println(len(pokedex))   // 2
}
```

### Comma ok idiom

O "comma ok" é um dos idioms mais importantes de Go. Ele aparece em maps,
type assertions, leitura de channels:

```go
// map
v, ok := m["chave"]

// type assertion
v, ok := x.(T)

// channel
v, ok := <-ch
```

Em todos os casos, o segundo valor booleano indica sucesso.

### Swift vs Go

```swift
// Swift — Dictionary com Optional
var pokedex: [String: Int] = ["Pikachu": 25, "Charizard": 6]
let nivel = pokedex["Pikachu"]  // Int? (Optional)

if let nivel = pokedex["Mewtwo"] {
    print("Nível:", nivel)
}
```

```go
// Go — map com zero value + comma ok
pokedex := map[string]int{"Pikachu": 25, "Charizard": 6}
nivel := pokedex["Pikachu"]  // int (zero value se não existir)

if nivel, ok := pokedex["Mewtwo"]; ok {
    fmt.Println("Nível:", nivel)
}
```

**Atenção:** Maps em Go **não são thread-safe**. Se você precisar de acesso
concorrente, use `sync.Map` ou proteja com `sync.Mutex` (veremos na Fase 04).

---

## `range` — iteração

`range` itera sobre arrays, slices, maps, strings e channels.

```go
package main

import "fmt"

func main() {
    pokemons := []string{"Pikachu", "Charizard", "Blastoise"}

    // slice: índice + valor
    for i, nome := range pokemons {
        fmt.Printf("%d: %s\n", i, nome)
    }

    // só valor (descarta índice com _)
    for _, nome := range pokemons {
        fmt.Println(nome)
    }

    // só índice
    for i := range pokemons {
        fmt.Println(i)
    }

    // map: chave + valor
    pokedex := map[string]int{"Pikachu": 25, "Charizard": 6}
    for nome, nivel := range pokedex {
        fmt.Printf("%s: nível %d\n", nome, nivel)
    }

    // string: índice do byte + rune
    for i, r := range "Pokémon" {
        fmt.Printf("%d: %c\n", i, r)
    }
}
```

**Atenção:** A ordem de iteração sobre maps é **aleatória**. Não confie na ordem.
Se precisar de ordem, extraia as chaves, ordene-as, e itere sobre as chaves ordenadas.

---

## Strings e runes

Strings em Go são **sequências de bytes** UTF-8. Isso significa que indexar uma
string retorna bytes, não caracteres:

```go
package main

import "fmt"

func main() {
    s := "Pokémon"

    fmt.Println(len(s))          // 8 (bytes, não caracteres!)
    fmt.Println(s[0])            // 80 ('P' em ASCII)
    fmt.Printf("%c\n", s[0])     // P

    // range itera sobre runes (code points), não bytes
    for i, r := range s {
        fmt.Printf("pos %d: %c\n", i, r)
    }
    // pos 0: P
    // pos 1: o
    // pos 2: k
    // pos 3: é    ← note que 'é' ocupa 2 bytes!
    // pos 5: m    ← índice pulou para 5
    // pos 6: o
    // pos 7: n

    // converter para []rune para acesso por caractere
    runes := []rune(s)
    fmt.Println(len(runes))       // 7 (caracteres)
    fmt.Printf("%c\n", runes[3])  // é
}
```

### Pacote `strings` — funções essenciais

```go
import "strings"

strings.Contains("Pokémon", "mon")    // true
strings.HasPrefix("Pokémon", "Po")    // true
strings.HasSuffix("Pokémon", "mon")   // true
strings.ToUpper("pokémon")            // "POKÉMON"
strings.ToLower("POKÉMON")            // "pokémon"
strings.Replace("Pika Pika", "Pika", "Chu", 1) // "Chu Pika"
strings.ReplaceAll("Pika Pika", "Pika", "Chu")  // "Chu Chu"
strings.Split("a,b,c", ",")           // ["a", "b", "c"]
strings.Join([]string{"a", "b"}, ",") // "a,b"
strings.TrimSpace("  hello  ")        // "hello"
```

---

## Exercícios da Fase 02

### 1. Filtro de Pokémons

Escreva uma função `filtraPorNivel(pokemons map[string]int, nivelMinimo int) []string`
que retorna os nomes dos pokémons com nível maior ou igual ao mínimo.
A ordem dos nomes não importa.

<details>
<summary>Gabarito</summary>

```go
func filtraPorNivel(pokemons map[string]int, nivelMinimo int) []string {
    resultado := make([]string, 0)
    for nome, nivel := range pokemons {
        if nivel >= nivelMinimo {
            resultado = append(resultado, nome)
        }
    }
    return resultado
}
```
</details>

### 2. Contador de tipos

Escreva uma função `contaTipos(pokemons []string) map[string]int` que recebe uma
lista de pokémons (com repetições) e retorna um map com a contagem de cada um.

```go
contaTipos([]string{"Pikachu", "Charizard", "Pikachu", "Mewtwo", "Pikachu"})
// map[Pikachu:3 Charizard:1 Mewtwo:1]
```

<details>
<summary>Gabarito</summary>

```go
func contaTipos(pokemons []string) map[string]int {
    contagem := make(map[string]int)
    for _, nome := range pokemons {
        contagem[nome]++
    }
    return contagem
}
```
</details>

### 3. Cópia segura

Escreva uma função `copiaSlice(original []int) []int` que retorna uma cópia
independente do slice (modificar um não afeta o outro).

<details>
<summary>Gabarito</summary>

```go
func copiaSlice(original []int) []int {
    copia := make([]int, len(original))
    copy(copia, original)
    return copia
}
```
</details>

### 4. Manipulação de strings

Escreva uma função `iniciaisMaiusculas(frase string) string` que transforma a
primeira letra de cada palavra em maiúscula. Ex: `"pikachu eu escolho voce"` →
`"Pikachu Eu Escolho Voce"`.

<details>
<summary>Gabarito</summary>

```go
func iniciaisMaiusculas(frase string) string {
    palavras := strings.Fields(frase)        // split por whitespace
    for i, p := range palavras {
        if len(p) > 0 {
            palavras[i] = strings.ToUpper(p[:1]) + p[1:]
        }
    }
    return strings.Join(palavras, " ")
}
```
</details>

---

**Próxima fase:** [Fase 03 — Structs, Interfaces & Erros](fase-03-structs-interfaces.md)
