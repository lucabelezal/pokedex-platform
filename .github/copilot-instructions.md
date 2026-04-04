# Instrucoes do Copilot para este repositorio

## Idioma
- Escreva respostas e explicacoes em portugues do Brasil.
- Gere mensagens de commit em portugues do Brasil.

## Estrutura do repositorio
- A raiz do repositorio concentra materiais transversais, como `doc/`, `.github/` e arquivos de colaboracao.
- A implementacao executavel da plataforma fica em `core/`.
- Dentro de `core/`, use a organizacao:
  - `core/app/` para servicos internos
  - `core/bff/` para BFFs
  - `core/gateway/` para configuracao do gateway
  - `core/infra/` para ativos de infraestrutura e dados
  - `core/docker-compose.yml` como compose principal da plataforma

## Nomes oficiais dos servicos
- Use os nomes oficiais atuais da plataforma:
  - `mobile-bff`
  - `pokemon-catalog-service`
  - `auth-service`
- Nao introduza novamente o nome legado `pokedex-service`.
- Ao falar da plataforma como um todo, `Pokedex Platform` ou `Plataforma Pokedex` continua fazendo sentido.

## Responsabilidades arquiteturais
- `mobile-bff` e o servico voltado para experiencia do cliente.
- `pokemon-catalog-service` e a fonte canonica do catalogo de Pokemon.
- `auth-service` concentra autenticacao e ciclo de vida de token.
- Evite deslocar regras canonicas de catalogo para o BFF.

## Hexagonal no mobile-bff
- Preserve a direcao das dependencias:
  - adapters inbound (HTTP) dependem de `ports/inbound` (use cases)
  - use cases dependem de `ports/outbound` (repositorios e clientes externos)
  - adapters outbound implementam `ports/outbound`
  - estrutura: `adapters/http` → `ports/inbound` ← `service` → `ports/outbound` ← `adapters/repository`
- Nao faca handlers HTTP dependerem diretamente de clients concretos de infraestrutura.
- Nao use `tests/` como dependencia de runtime.
- Ao introduzir integracoes externas, normalize erros no adapter externo ou na camada de aplicacao, nao no handler HTTP.
- Novas entidades de dominio vivem em `domain/`, nao em `ports/`.

## Convencao de commits (obrigatorio)
- Use sempre Conventional Commits.
- Formato: `tipo(escopo-opcional): descricao curta em portugues`.
- Tipos permitidos: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`, `ci`, `build`, `perf`, `revert`.
- Exemplos validos:
  - `feat(favoritos): bloquear remocao sem autenticacao`
  - `fix(bff): corrigir validacao de token JWT`
  - `docs(readme): atualizar fluxo de execucao local`

## Qualidade das mensagens
- Evite mensagens genericas como `update`, `ajustes`, `wip`.
- Descreva o impacto funcional de forma objetiva.
- Quando houver mudanca relevante, inclua corpo explicando contexto e motivacao.

## Documentacao
- Atualize `doc/` quando houver mudanca relevante de arquitetura, nomenclatura ou responsabilidade entre servicos.
- Se a mudanca afetar a estrutura do runtime, revise tambem `README.md` e `core/README.md`.

## Colaboracao com CI
- Se o workflow de validacao de commits falhar, ajuste a mensagem para o padrao Conventional Commits.
- Em revisoes, priorize padrao de commit, clareza da descricao e impacto da mudanca.

## Convencoes Go (obrigatorio)
- **Receiver**: 1-2 letras, abreviacao do tipo. Use `c` para `Client`, `s` para `Service`. Nunca `this` ou `self`.
- **Sem prefixo Get**: `Count()` nao `GetCount()`; `Name()` nao `GetName()`. Excecao: `GetXxx` aceito em interface `net/http.ResponseWriter`.
- **Initialisms uppercase**: `ID`, `DB`, `URL`, `HTTP`, `JSON`. Para `gRPC`: exported = `GRPC`, unexported = `gRPC`. Nunca `Id`, `Db`, `Url`, `Grpc`.
- **Error strings**: minusculas, sem ponto final. `"token invalido"` nao `"Token invalido."`.
- **Indent error flow**: retorne o erro imediatamente; o caminho feliz fica sem aninhamento.
- **Wrapping de erros**: use `%w` para erros que o chamador possa inspecionar com `errors.Is`/`errors.As`; use `%v` apenas para anotacao sem inspecao.
- **Declaracao de variavel**: `var x T` para zero value explicito; `x := value` quando inicializa com valor nao-zero.
- **Goroutines**: documente quando a goroutine termina. Use `context.Context` para controlar o ciclo de vida. Nunca lance goroutines sem uma estrategia de finalizacao.
- **Slices via JSON**: use `make([]T, 0)` em repositorios e use cases que retornam slices via API/JSON. Evita serializar `null` em vez de `[]`.
- **Interface compliance**: declare `var _ Interface = (*Struct)(nil)` ao final de cada arquivo de adaptador e servico. Torna erros de implementacao visiveis em tempo de compilacao.
- **Constantes de tempo**: use `time.Duration` diretamente. Ex.: `const timeout = 30 * time.Second`. Nunca `int` com cast posterior.
- **Testes**: prefira table-driven tests (`t.Run`) quando multiplos casos compartilham a mesma logica de verificacao. Nunca declare campos na struct de teste que nao sao usados no corpo do subtest.

## Guia de estilo canonico
Para qualquer revisao ou escrita de codigo Go, carregue a skill `go-style-combined` (`.github/skills/go-style-combined/SKILL.md`).
Ela e a sintese definitiva do Uber Go Style Guide + Google Go Style Guide com as decisoes especificas deste projeto.
