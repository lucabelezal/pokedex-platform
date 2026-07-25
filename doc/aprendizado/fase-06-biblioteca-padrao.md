# Fase 06 — Biblioteca Padrão & Utilitários

Esta fase cobre os pacotes da biblioteca padrão que você mais usa no dia a dia de
um backend Go: tempo, arquivos, flags, logging, ordenação e outros utilitários.

---

## Tempo — `time`

O pacote `time` é onipresente em aplicações backend. Diferente de Swift (`Date`),
Go tem um sistema próprio de formatação de datas — e ele é baseado em uma **data de referência**.

### `time.Time` e `time.Duration`

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    agora := time.Now()
    fmt.Println("Agora:", agora)

    // criar data específica
    data := time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC)
    fmt.Println("Data:", data)

    // duração
    duracao := 5 * time.Second
    fmt.Println("Duração:", duracao)

    // operações
    futuro := agora.Add(24 * time.Hour)
    passado := agora.Add(-30 * time.Minute)

    // diferença entre datas
    diff := futuro.Sub(agora)
    fmt.Println("Diferença:", diff)  // 24h0m0s

    // comparação
    fmt.Println(agora.Before(futuro))   // true
    fmt.Println(agora.After(passado))   // true
}
```

**Atenção:** Use sempre `time.Duration` para intervalos de tempo. Nunca passe
tempo como `int` (milissegundos, segundos). Isso é propenso a erros e não idiomático.

### Formatação e parsing

Go usa uma **data de referência** para formatação, não tokens como `YYYY-MM-DD`.
A referência é:

```
Mon Jan 2 15:04:05 MST 2006
```

Que em formato numérico é: `01/02 03:04:05PM '06 -0700`

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    t := time.Date(2024, time.January, 15, 10, 30, 0, 0, time.UTC)

    // formatação (use a data de referência como molde)
    fmt.Println(t.Format("2006-01-02"))                 // 2024-01-15
    fmt.Println(t.Format("02/01/2006"))                  // 15/01/2024
    fmt.Println(t.Format("15:04:05"))                    // 10:30:00
    fmt.Println(t.Format("2006-01-02 15:04:05"))         // 2024-01-15 10:30:00
    fmt.Println(t.Format("Monday, January 2, 2006"))      // Monday, January 15, 2024
    fmt.Println(t.Format(time.RFC3339))                  // 2024-01-15T10:30:00Z

    // parsing — mesmo molde
    parsed, err := time.Parse("2006-01-02", "2024-03-20")
    if err != nil {
        fmt.Println("Erro:", err)
        return
    }
    fmt.Println("Parseado:", parsed)
}
```

**Formatos pré-definidos:**

| Constante | Formato | Exemplo |
|-----------|---------|---------|
| `time.RFC3339` | ISO 8601 | `2024-01-15T10:30:00Z` |
| `time.RFC3339Nano` | ISO 8601 com nanosegundos | `2024-01-15T10:30:00.123456789Z` |
| `time.DateTime` | Data e hora | `2024-01-15 10:30:00` |
| `time.DateOnly` | Só data | `2024-01-15` |
| `time.TimeOnly` | Só hora | `10:30:00` |

### Epoch (Unix timestamp)

```go
agora := time.Now()

// time → unix
fmt.Println(agora.Unix())           // 1705319400 (segundos)
fmt.Println(agora.UnixMilli())      // milissegundos
fmt.Println(agora.UnixNano())       // nanossegundos

// unix → time
t := time.Unix(1705319400, 0)
fmt.Println(t)  // 2024-01-15 10:30:00 +0000 UTC
```

### Sleep e Timers

```go
// sleep — pausa a goroutine atual
time.Sleep(2 * time.Second)

// timer — dispara uma vez no futuro
timer := time.NewTimer(5 * time.Second)
<-timer.C  // bloqueia até o timer disparar
timer.Stop() // cancela se ainda não disparou

// ticker — dispara repetidamente
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()
for i := 0; i < 3; i++ {
    <-ticker.C
    fmt.Println("Tick")
}
```

### Swift vs Go — tempo

```swift
// Swift
let agora = Date()
let futuro = agora.addingTimeInterval(86400)
let formatter = DateFormatter()
formatter.dateFormat = "yyyy-MM-dd"
let str = formatter.string(from: agora)
```

```go
// Go
agora := time.Now()
futuro := agora.Add(24 * time.Hour)
str := agora.Format("2006-01-02")          // formato de referência
parsed, _ := time.Parse("2006-01-02", str)
```

**Atenção:** A data de referência parece mágica no início, mas você rapidamente
memoriza: `Mon Jan 2 15:04:05 2006` = `1 2 3 4 5 6`. Os números são sequenciais.

---

## Arquivos & I/O

### Leitura de arquivos

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    // os.ReadFile — lê arquivo inteiro de uma vez (arquivos pequenos)
    data, err := os.ReadFile("pokemon.txt")
    if err != nil {
        fmt.Println("Erro:", err)
        return
    }
    fmt.Println(string(data))

    // os.Open + bufio.Scanner — linha por linha (arquivos grandes)
    f, err := os.Open("pokemon.txt")
    if err != nil {
        fmt.Println("Erro:", err)
        return
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        linha := scanner.Text()
        fmt.Println("Linha:", linha)
    }
    if err := scanner.Err(); err != nil {
        fmt.Println("Erro ao ler:", err)
    }
}
```

### Escrita de arquivos

```go
// os.WriteFile — escreve arquivo inteiro de uma vez
linhas := []byte("Pikachu\nCharizard\nBlastoise\n")
err := os.WriteFile("pokemon.txt", linhas, 0644)

// os.Create + fmt.Fprintln — escreve linha por linha
f, err := os.Create("pokemon.txt")
defer f.Close()

fmt.Fprintln(f, "Pikachu")
fmt.Fprintln(f, "Charizard")
fmt.Fprintf(f, "%s (lvl %d)\n", "Mewtwo", 150)

// os.OpenFile — abre com flags (append, read/write, etc.)
f, err = os.OpenFile("pokemon.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
fmt.Fprintln(f, "nova linha")
```

### Caminhos de arquivos

```go
import "path/filepath"

caminho := "data/pokemons/pikachu.json"

filepath.Base(caminho)               // "pikachu.json"
filepath.Dir(caminho)                // "data/pokemons"
filepath.Ext(caminho)                // ".json"
filepath.Join("data", "pokemons")    // "data/pokemons"
filepath.Clean("./data//pokemons/")  // "data/pokemons"

abs, _ := filepath.Abs(caminho)     // caminho absoluto
```

**Atenção:** Sempre use `filepath.Join` em vez de concatenar strings com `/`.
Ela resolve o separador correto para cada sistema operacional.

### Diretórios

```go
// criar diretório (com todos os pais)
os.MkdirAll("data/pokemons/images", 0755)

// listar arquivos
entries, _ := os.ReadDir("data/pokemons")
for _, e := range entries {
    fmt.Println(e.Name(), e.IsDir())
}

// remover
os.Remove("arquivo.txt")
os.RemoveAll("data/temp")  // remove diretório e todo conteúdo
```

### Arquivos temporários

```go
// arquivo temporário
f, _ := os.CreateTemp("", "pokemon-*.txt")
defer os.Remove(f.Name())  // limpa ao final
fmt.Println("Arquivo temp:", f.Name())
f.WriteString("dados temporários")

// diretório temporário
dir, _ := os.MkdirTemp("", "pokedex-*")
defer os.RemoveAll(dir)
fmt.Println("Diretório temp:", dir)
```

### Embed directive (Go 1.16+)

Incorpora arquivos no binário em tempo de compilação:

```go
package main

import (
    _ "embed"
    "fmt"
)

//go:embed pokemon.txt
var pokemons string

//go:embed config.json
var config []byte

func main() {
    fmt.Println(pokemons)
    fmt.Println(string(config))
}
```

Útil para templates, arquivos de configuração, migrations SQL e assets estáticos.

---

## `panic` e `recover`

`panic` interrompe a execução normal da goroutine. É para **erros irrecuperáveis**,
não para controle de fluxo.

```go
package main

import "fmt"

func main() {
    fmt.Println("início")
    panic("algo deu muito errado")
    fmt.Println("fim")  // nunca executa
}
```

### `recover` — capturando panics

```go
package main

import "fmt"

func processaComRecovery() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recuperado de panic:", r)
        }
    }()

    fmt.Println("Processando...")
    panic("falha inesperada")
    fmt.Println("Nunca executa")
}

func main() {
    processaComRecovery()
    fmt.Println("Continua normalmente")  // executa!
}
```

**Regras:**
- `recover()` só funciona dentro de `defer`.
- Panic deve ser reservado para situações verdadeiramente irrecuperáveis
  (ex: falha ao ler config essencial, corrupção de memória).
- Erros de negócio NUNCA devem ser panics. Use `error` como valor de retorno.

### Swift vs Go — panic/recover

```swift
// Swift — fatalError (não recuperável)
fatalError("algo deu errado")
// não existe equivalente a recover em Swift
```

```go
// Go — panic é recuperável via recover dentro de defer
panic("erro irrecuperável")  // pode ser capturado com recover
```

---

## Command-line flags

```go
package main

import (
    "flag"
    "fmt"
)

func main() {
    // define flags
    porta := flag.Int("porta", 8080, "porta do servidor")
    host := flag.String("host", "localhost", "host do servidor")
    debug := flag.Bool("debug", false, "modo debug")

    // parse — deve ser chamado depois de todas as definições
    flag.Parse()

    fmt.Printf("Servidor em %s:%d\n", *host, *porta)
    fmt.Printf("Debug: %t\n", *debug)
    fmt.Println("Argumentos restantes:", flag.Args())
}
```

Execução:

```bash
go run main.go --porta=9090 --debug
go run main.go --host 0.0.0.0 --porta 3000
go run main.go --help  # ajuda automática
```

### Subcommands (básico)

```go
if len(os.Args) < 2 {
    fmt.Println("esperado um subcomando: serve ou migrate")
    return
}

switch os.Args[1] {
case "serve":
    cmd := flag.NewFlagSet("serve", flag.ExitOnError)
    porta := cmd.Int("porta", 8080, "porta")
    cmd.Parse(os.Args[2:])
    fmt.Println("Rodando servidor na porta", *porta)

case "migrate":
    cmd := flag.NewFlagSet("migrate", flag.ExitOnError)
    direcao := cmd.String("dir", "up", "up ou down")
    cmd.Parse(os.Args[2:])
    fmt.Println("Migrando:", *direcao)

default:
    fmt.Println("subcomando inválido")
}
```

---

## Variáveis de ambiente

```go
package main

import (
    "fmt"
    "os"
    "strconv"
)

func main() {
    // leitura simples
    dbHost := os.Getenv("DB_HOST")
    if dbHost == "" {
        dbHost = "localhost"  // fallback
    }

    // lookup — verifica se a variável existe
    dbPort, existe := os.LookupEnv("DB_PORT")
    if !existe {
        dbPort = "5432"
    }

    // conversão
    portInt, err := strconv.Atoi(dbPort)
    if err != nil {
        fmt.Println("DB_PORT inválido:", err)
        return
    }

    fmt.Printf("Banco: %s:%d\n", dbHost, portInt)

    // listar todas
    for _, env := range os.Environ() {
        fmt.Println(env)
    }

    // definir (para o processo atual)
    os.Setenv("APP_ENV", "production")
}
```

---

## Logging — `slog`

Go 1.21+ introduziu `log/slog` para logging estruturado:

```go
package main

import (
    "log/slog"
    "os"
)

func main() {
    // logger JSON (recomendado para produção)
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))

    // logger texto (recomendado para desenvolvimento)
    textLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))

    logger.Info("servidor iniciado",
        "porta", 8080,
        "host", "localhost",
    )

    logger.Warn("conexão lenta",
        "duracao_ms", 3200,
        "endpoint", "/api/pokemon",
    )

    logger.Error("falha ao buscar pokemon",
        "id", 25,
        "erro", "connection refused",
    )

    // com contexto
    logger.Info("pokemon encontrado",
        slog.String("nome", "Pikachu"),
        slog.Int("level", 25),
        slog.Group("requisicao",
            slog.String("method", "GET"),
            slog.Int64("duracao_us", 1500),
        ),
    )

    textLogger.Debug("detalhes internos", "cache_hit", true)
}
```

Saída (JSON):

```json
{"time":"2024-01-15T10:30:00.123Z","level":"INFO","msg":"servidor iniciado","porta":8080,"host":"localhost"}
{"time":"2024-01-15T10:30:00.456Z","level":"ERROR","msg":"falha ao buscar pokemon","id":25,"erro":"connection refused"}
```

### Swift vs Go — logging

```swift
// Swift — os_log / Logger (iOS 14+)
import os
let logger = Logger(subsystem: "com.app.pokedex", category: "network")
logger.info("Servidor iniciado na porta \(porta)")
logger.error("Falha ao buscar: \(error.localizedDescription)")
```

```go
// Go — slog (Go 1.21+)
logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
logger.Info("servidor iniciado", "porta", 8080)
logger.Error("falha ao buscar", "erro", err)
```

**Atenção:** `log` e `slog` coexistem. Para código novo, prefira `slog` — é
estruturado, suporta níveis e tem handlers JSON.

---

## Ordenação — `sort`

```go
package main

import (
    "fmt"
    "sort"
)

type Pokemon struct {
    Nome  string
    Level int
}

func main() {
    // tipos básicos
    numeros := []int{5, 2, 8, 1, 9}
    sort.Ints(numeros)
    fmt.Println(numeros)  // [1 2 5 8 9]

    strings := []string{"Charizard", "Pikachu", "Blastoise"}
    sort.Strings(strings)
    fmt.Println(strings)  // [Blastoise Charizard Pikachu]

    // ordenação customizada com sort.Slice
    pokemons := []Pokemon{
        {"Charizard", 36},
        {"Pikachu", 25},
        {"Blastoise", 36},
    }

    // por nível (crescente)
    sort.Slice(pokemons, func(i, j int) bool {
        return pokemons[i].Level < pokemons[j].Level
    })
    fmt.Println(pokemons)

    // por nível (decrescente), desempate por nome
    sort.Slice(pokemons, func(i, j int) bool {
        if pokemons[i].Level != pokemons[j].Level {
            return pokemons[i].Level > pokemons[j].Level
        }
        return pokemons[i].Nome < pokemons[j].Nome
    })
    fmt.Println(pokemons)

    // verificar se está ordenado
    fmt.Println(sort.IntsAreSorted(numeros))  // true
}
```

Se você implementar a interface `sort.Interface`, pode usar `sort.Sort`:

```go
type ByLevel []Pokemon

func (p ByLevel) Len() int           { return len(p) }
func (p ByLevel) Less(i, j int) bool { return p[i].Level < p[j].Level }
func (p ByLevel) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

var pokemons ByLevel = []Pokemon{{"Charizard", 36}, {"Pikachu", 25}}
sort.Sort(pokemons)
```

Na prática, `sort.Slice` cobre a maioria dos casos.

---

## URL Parsing

```go
package main

import (
    "fmt"
    "net/url"
)

func main() {
    raw := "https://pokeapi.co/api/v2/pokemon?limit=20&offset=0"

    parsed, err := url.Parse(raw)
    if err != nil {
        fmt.Println("URL inválida:", err)
        return
    }

    fmt.Println("Scheme:", parsed.Scheme)     // https
    fmt.Println("Host:", parsed.Host)         // pokeapi.co
    fmt.Println("Path:", parsed.Path)         // /api/v2/pokemon

    // query params
    limit := parsed.Query().Get("limit")      // "20"
    fmt.Println("Limit:", limit)

    // construir URL
    values := url.Values{}
    values.Add("nome", "Pikachu")
    values.Add("level", "25")
    qs := values.Encode()                     // nome=Pikachu&level=25

    u := &url.URL{
        Scheme:   "https",
        Host:     "api.exemplo.com",
        Path:     "/pokemons",
        RawQuery: qs,
    }
    fmt.Println(u.String())  // https://api.exemplo.com/pokemons?nome=Pikachu&level=25
}
```

---

## Base64 Encoding

```go
package main

import (
    "encoding/base64"
    "fmt"
)

func main() {
    dados := []byte("Pikachu usa Choque do Trovão!")

    // codificação padrão
    encoded := base64.StdEncoding.EncodeToString(dados)
    fmt.Println("Encode:", encoded)

    // decodificação
    decoded, err := base64.StdEncoding.DecodeString(encoded)
    if err != nil {
        fmt.Println("Erro:", err)
        return
    }
    fmt.Println("Decode:", string(decoded))

    // URL-safe (sem +/ no resultado)
    urlEncoded := base64.URLEncoding.EncodeToString(dados)
    fmt.Println("URL-safe:", urlEncoded)
}
```

---

## Text Templates

```go
package main

import (
    "os"
    "text/template"
)

type Pokemon struct {
    Nome  string
    Level int
    Tipos []string
}

func main() {
    tmpl := `Pokémon: {{.Nome}}
Nível: {{.Level}}
Tipos: {{range .Tipos}}{{.}} {{end}}
{{if gt .Level 50}}🔥 FORTE{{else}}📊 Normal{{end}}
`

    t := template.Must(template.New("pokemon").Parse(tmpl))

    p := Pokemon{
        Nome:  "Charizard",
        Level: 36,
        Tipos: []string{"Fogo", "Voador"},
    }

    t.Execute(os.Stdout, p)
}
```

Saída:

```
Pokémon: Charizard
Nível: 36
Tipos: Fogo Voador
📊 Normal
```

Funções úteis em templates: `{{if .Campo}}`, `{{range}}`, `{{with}}`, `eq`, `gt`,
`and`, `or`, `not`. Use `html/template` em vez de `text/template` quando o output
for HTML (protege contra XSS automaticamente).

---

## Rate Limiting básico

Implementação simples com `ticker`:

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 5 requisições por segundo
    rate := time.Second / 5
    limiter := time.NewTicker(rate)
    defer limiter.Stop()

    for i := 0; i < 10; i++ {
        <-limiter.C  // espera o próximo slot
        fmt.Printf("Requisição %d em %v\n", i, time.Now().Format("15:04:05.000"))
    }
}
```

Saída:

```
Requisição 0 em 15:04:05.001
Requisição 1 em 15:04:05.201
Requisição 2 em 15:04:05.401
...
```

Padrão burst (permite rajadas curtas, mas mantém a média):

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // bucket de 3 tokens, recarrega 1 token por segundo
    burstyLimiter := make(chan time.Time, 3)

    // enche o bucket inicial
    for i := 0; i < 3; i++ {
        burstyLimiter <- time.Now()
    }

    // recarrega 1 token por segundo
    go func() {
        for t := range time.NewTicker(time.Second).C {
            select {
            case burstyLimiter <- t:
            default:
            }
        }
    }()

    for i := 0; i < 10; i++ {
        <-burstyLimiter
        fmt.Printf("Requisição %d em %v\n", i, time.Now().Format("15:04:05.000"))
    }
}
```

---

## `os.Exit`

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    fmt.Println("Processamento iniciado...")

    if err := processa(); err != nil {
        fmt.Println("Erro fatal:", err)
        os.Exit(1)   // sai com código de erro
    }

    os.Exit(0)       // sai com sucesso
    // defer NÃO executa após os.Exit
}

func processa() error {
    return fmt.Errorf("falha catastrófica")
}
```

**Atenção:** `defer` **não** executa após `os.Exit`. Se precisar de limpeza antes
de sair, faça-a antes de chamar `os.Exit`.

---

## Exercícios da Fase 06

### 1. Log estruturado

Crie uma função que usa `slog` para logar os detalhes de uma requisição HTTP
(método, path, status, duração). Use atributos estruturados.

<details>
<summary>Gabarito</summary>

```go
package main

import (
    "log/slog"
    "os"
    "time"
)

func logRequisicao(method, path string, status int, duracao time.Duration) {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
    logger.Info("requisicao",
        slog.String("method", method),
        slog.String("path", path),
        slog.Int("status", status),
        slog.Int64("duracao_us", duracao.Microseconds()),
    )
}

func main() {
    inicio := time.Now()
    time.Sleep(15 * time.Millisecond) // simula processamento
    logRequisicao("GET", "/api/pokemon", 200, time.Since(inicio))
}
```
</details>

### 2. Ordenação customizada

Crie um slice de structs `Movimento{Nome string, Poder int, Precisao int}` e
ordene por poder decrescente, desempatando por precisão decrescente.

<details>
<summary>Gabarito</summary>

```go
package main

import (
    "fmt"
    "sort"
)

type Movimento struct {
    Nome     string
    Poder    int
    Precisao int
}

func main() {
    movimentos := []Movimento{
        {"Choque do Trovão", 90, 100},
        {"Lança-chamas", 90, 100},
        {"Tackle", 40, 100},
        {"Hidro Bomba", 110, 80},
        {"Trovão", 110, 70},
    }

    sort.Slice(movimentos, func(i, j int) bool {
        if movimentos[i].Poder != movimentos[j].Poder {
            return movimentos[i].Poder > movimentos[j].Poder
        }
        if movimentos[i].Precisao != movimentos[j].Precisao {
            return movimentos[i].Precisao > movimentos[j].Precisao
        }
        return movimentos[i].Nome < movimentos[j].Nome
    })

    for _, m := range movimentos {
        fmt.Printf("%-20s Poder: %3d Precisão: %3d%%\n", m.Nome, m.Poder, m.Precisao)
    }
}
```
</details>

### 3. Template de e-mail

Crie um template `text/template` para um e-mail de confirmação de captura:

```
Assunto: Pokémon capturado!
Treinador Ash capturou Pikachu (nível 25) usando uma Pokébola.
```

Use uma struct `Captura` com campos `Treinador`, `Pokemon`, `Nivel`, `Item`.

<details>
<summary>Gabarito</summary>

```go
package main

import (
    "os"
    "text/template"
)

type Captura struct {
    Treinador string
    Pokemon   string
    Nivel     int
    Item      string
}

func main() {
    tmpl := `Assunto: Pokémon capturado!
Treinador {{.Treinador}} capturou {{.Pokemon}} (nível {{.Nivel}}) usando {{.Item}}.
`

    t := template.Must(template.New("captura").Parse(tmpl))
    c := Captura{Treinador: "Ash", Pokemon: "Pikachu", Nivel: 25, Item: "Pokébola"}
    t.Execute(os.Stdout, c)
}
```
</details>

### 4. Embed de arquivo

Use `//go:embed` para incorporar um arquivo de configuração JSON no binário e
desserializá-lo para uma struct.

<details>
<summary>Gabarito</summary>

```go
package main

import (
    _ "embed"
    "encoding/json"
    "fmt"
)

//go:embed config.json
var configRaw []byte

type Config struct {
    Porta int    `json:"porta"`
    Host  string `json:"host"`
    Debug bool   `json:"debug"`
}

func main() {
    var cfg Config
    if err := json.Unmarshal(configRaw, &cfg); err != nil {
        panic(err)
    }
    fmt.Printf("Servidor: %s:%d (debug: %t)\n", cfg.Host, cfg.Porta, cfg.Debug)
}
```

Crie `config.json`:

```json
{"porta": 8080, "host": "localhost", "debug": true}
```
</details>

---

**Voltar ao início:** [ROTEIRO](ROTEIRO.md)

[Voltar à fase anterior: Fase 05 — Testes, Web & Ferramentas](fase-05-testes-web.md)
