# Guia de Estilo

[Visão Geral](README.md) | [Guia](guide.md) | [Decisões](decisions.md) | [Melhores Práticas](best-practices.md) | [Regras do Projeto](regras-projeto.md)

**Status:** `[Normativo] [Canônico]` — regras fundamentais que todo código deve seguir.

## Princípios de estilo

Existem alguns princípios abrangentes que resumem como pensar sobre a escrita de
código Go legível. Os seguintes são atributos de código legível, em ordem de
importância:

1. **Clareza**: O propósito e a lógica do código são claros para o leitor.
1. **Simplicidade**: O código atinge seu objetivo da maneira mais simples possível.
1. **Concisão**: O código tem uma alta relação sinal-ruído.
1. **Manutenibilidade**: O código é escrito de forma que possa ser facilmente mantido.
1. **Consistência**: O código é consistente com o restante da base de código do projeto.

### Clareza

O objetivo central da legibilidade é produzir código que seja claro para o leitor.

A clareza é alcançada principalmente com nomes eficazes, comentários úteis e
organização eficiente do código. A clareza deve ser vista pela lente do leitor,
não do autor. É mais importante que o código seja fácil de ler do que fácil de
escrever.

### Simplicidade

O código Go deve ser simples para quem o usa, lê e mantém.

Código Go deve ser escrito da maneira mais simples que atinja seus objetivos,
tanto em termos de comportamento quanto de desempenho. Dentro da base de código
do projeto, código simples:

- É fácil de ler de cima a baixo
- Não assume que você já sabe o que está fazendo
- Não tem níveis desnecessários de abstração
- Tem comentários que explicam o porquê, não o quê
- Tem documentação que se sustenta sozinha
- Tem erros úteis e falhas de teste úteis

### Menor mecanismo

Onde houver várias maneiras de expressar a mesma ideia, prefira a que usa as
ferramentas mais padrão:

1. Use um construto da linguagem (channel, slice, map, loop, struct) quando
   suficiente para o seu caso de uso.
2. Se não houver, procure uma ferramenta dentro da biblioteca padrão.
3. Por fim, considere se há uma dependência no projeto que é suficiente antes
   de introduzir uma nova.

### Concisão

Código Go conciso tem uma alta relação sinal-ruído. É fácil discernir os
detalhes relevantes, e a nomenclatura e estrutura guiam o leitor.

### Manutenibilidade

Código é editado muito mais vezes do que é escrito. Código manutenível:

- É fácil para um programador futuro modificar corretamente
- Tem APIs estruturadas para que possam crescer com elegância
- É claro sobre as suposições que faz
- Evita acoplamento desnecessário
- Tem uma suíte de testes abrangente

### Consistência

Código consistente é aquele que parece, soa e se comporta como código semelhante
em toda a base de código. Preocupações de consistência não se sobrepõem a nenhum
dos princípios acima, mas se um empate precisar ser desfeito, geralmente é
benéfico desempatar em favor da consistência.

---

## Formatação

### `gofmt`

Todo arquivo fonte Go deve estar em conformidade com a saída da ferramenta `gofmt`.

Execute `goimports` ao salvar para gerenciar imports automaticamente.
Execute `go vet` para verificar erros.

### MixedCaps

Código Go usa `MixedCaps` ou `mixedCaps` (camel case) em vez de underscores
(snake case) ao escrever nomes com múltiplas palavras. Isso se aplica mesmo
quando quebra convenções em outras linguagens.

| Ruim | Bom |
|------|-----|
| `MAX_LENGTH` | `MaxLength` |
| `max_length` | `maxLength` |
| `user_count` | `userCount` |

Exceções para underscores em nomes Go:

1. Pacotes importados apenas por código gerado podem conter underscores.
1. Funções `Test`, `Benchmark` e `Example` em arquivos `*_test.go` podem incluir underscores.
1. Bibliotecas de baixo nível que interoperam com o sistema operacional ou cgo.

**Nomes de arquivo** de código fonte não são identificadores Go e não precisam
seguir estas convenções. Podem conter underscores.

### Tamanho de linha

Não há limite fixo de tamanho de linha para código Go. Se uma linha parecer
longa demais, prefira refatorar em vez de dividi-la. Se já estiver tão curta
quanto é prático, a linha deve permanecer longa.

Não divida uma linha:

- Antes de uma mudança de indentação (ex: declaração de função, condicional)
- Para fazer uma string longa (ex: URL) caber em múltiplas linhas curtas

---

## Nomenclatura

### Nomes de pacotes

Ao nomear pacotes, escolha um nome que seja:

- Todo em minúsculas. Sem maiúsculas ou underscores.
- Não precise ser renomeado usando alias de import na maioria dos locais de chamada.
- Curto e sucinto. Lembre-se de que o nome é identificado por extenso em cada
  local de chamada.
- **Não plural**. Por exemplo, `net/url`, não `net/urls`.
- **Não** `common`, `util`, `shared`, `helper`, `lib` ou `model`. São nomes ruins e
  pouco informativos.

| Ruim | Bom |
|------|-----|
| `package utils` | `package stringutil` |
| `package common` | `package creditcard` |
| `package models` | `package domain` |

Evite selecionar nomes de pacotes que provavelmente serão sombreados por nomes
de variáveis locais comumente usados. Por exemplo, `usercount` é um nome de
pacote melhor que `count`, já que `count` é um nome de variável comum.

### Nomes de funções

Funções e métodos **não** devem usar prefixo `Get` ou `get`, a menos que o
conceito subjacente use a palavra "get" (ex: HTTP GET). Prefira começar o nome
com o substantivo diretamente.

| Ruim | Bom |
|------|-----|
| `func GetCount() int` | `func Count() int` |
| `func GetUser(id string)` | `func User(id string)` |
| `func (c *Config) GetJobName()` | `func (c *Config) JobName()` |

Se a função envolver computação complexa ou chamada remota, use palavras como
`Compute` ou `Fetch` em vez de `Get`, para deixar claro que a chamada pode
demorar e pode bloquear ou falhar.

```go
// Bom:
func (s *Service) FetchPokemonDetails(ctx context.Context, id string) (*PokemonDetail, error)
```

Funções que retornam algo recebem nomes substantivados. Funções que fazem algo
recebem nomes verbais.

```go
// Bom:
func (c *Config) JobName(key string) (value string, ok bool)  // substantivo
func (c *Config) WriteTo(w io.Writer) (int64, error)          // verbo
```

### Nomes de receivers

Nomes de variáveis receiver devem ser:

- **Curtos** (geralmente uma ou duas letras)
- **Abreviações do próprio tipo**
- Aplicados de forma **consistente** a cada receiver daquele tipo
- **Nunca** `this` ou `self`
- **Não** usar underscore; omita o nome se não for usado

| Ruim | Bom |
|------|-----|
| `func (tray Tray)` | `func (t Tray)` |
| `func (this *ReportWriter)` | `func (w *ReportWriter)` |
| `func (self *Scanner)` | `func (s *Scanner)` |
| `func (info *ResearchInfo)` | `func (ri *ResearchInfo)` |

Convenções do projeto:

| Tipo | Receiver |
|------|----------|
| Service | `s` |
| Client / Repository | `c` |
| Handler | `h` |
| ResponseBuilder | `rb` |

```go
// Bom:
func (s *PokemonService) List(ctx context.Context, params SearchParams) (*PokemonPage, error)
func (c *PokemonCatalogServiceRepository) GetByID(ctx context.Context, id string) (*Pokemon, error)
func (h *Handler) GetPokemonDetails(w http.ResponseWriter, r *http.Request)
```

### Initialisms

Palavras em nomes que são inicialismos ou acrônimos (ex: `URL` e `NATO`) devem
ter a mesma capitalização em todo o nome. `URL` deve aparecer como `URL` ou `url`
(como em `urlPony` ou `URLPony`), nunca como `Url`.

| Uso | Exportado | Não exportado |
|-----|-----------|---------------|
| XML API | `XMLAPI` | `xmlAPI` |
| ID | `ID` | `id` |
| DB | `DB` | `db` |
| URL | `URL` | `url` |
| HTTP | `HTTP` | `http` |
| JSON | `JSON` | `json` |
| JWT | `JWT` | `jwt` |
| DTO | `DTO` | `dto` |
| UUID | `UUID` | `uuid` |
| CORS | `CORS` | — |
| iOS | `IOS` | `iOS` |
| gRPC | `GRPC` | `gRPC` |

| Ruim | Bom |
|------|-----|
| `userId` | `userID` |
| `UrlPony` | `urlPony` ou `URLPony` |
| `JsonDecoder` | `jsonDecoder` ou `JSONDecoder` |
| `HttpClient` | `httpClient` ou `HTTPClient` |

### Nomes de variáveis

A regra geral é que o comprimento de um nome deve ser proporcional ao tamanho do
seu escopo e inversamente proporcional ao número de vezes que é usado dentro
desse escopo.

- Uma variável criada no escopo de arquivo pode exigir múltiplas palavras.
- Uma variável com escopo em um único bloco interno pode ser uma única palavra ou
  até um ou dois caracteres.

| Escopo | Exemplo |
|--------|---------|
| Pacote | `var DefaultTimeout = 30 * time.Second` |
| Função | `user, err := s.repo.UserByID(ctx, id)` |
| Loop | `for i, v := range items` |
| Bloco curto | `n, err := w.Write(data)` |

### Nomes de constantes

Nomes de constantes devem usar `MixedCaps` como todos os outros nomes em Go.
Constantes exportadas começam com maiúscula, não exportadas com minúscula.

Nomeie constantes com base em seu **papel**, não em seu valor.

| Ruim | Bom |
|------|-----|
| `const MAX_PACKET_SIZE = 512` | `const MaxPacketSize = 512` |
| `const kMaxBufferSize = 1024` | `const maxBufferSize = 1024` |
| `const Twelve = 12` | Use `12` diretamente ou nomeie pelo papel |
| `const UserNameColumn = "username"` | Use a string diretamente se não houver papel semântico |

### Consistência local

Onde o guia de estilo não tem nada a dizer sobre um ponto particular de estilo,
os autores são livres para escolher o estilo que preferirem, a menos que o código
em proximidade (geralmente dentro do mesmo arquivo ou pacote) tenha adotado uma
posição consistente sobre o assunto.

Exemplos de considerações de estilo local **válidas**:

- Uso de `%s` ou `%v` para impressão formatada de erros
- Uso de channels com buffer em vez de mutexes

Exemplos de considerações de estilo local **inválidas**:

- Restrições de tamanho de linha para código
- Uso de bibliotecas de teste baseadas em asserção
