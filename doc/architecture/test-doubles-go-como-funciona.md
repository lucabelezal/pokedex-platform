# Test Doubles em Go — com Código Real deste Projeto

> A fase-07 (`doc/aprendizado/fase-07-testes-avancados-design.md`) explica a **teoria** dos test doubles.
> Este guia responde duas perguntas práticas com código REAL do `pokedex-platform`:
> 1. Como usar **stub, fake, spy, mock** em Go (interfaces + structs manuais)?
> 2. Onde colocar os testes: pasta `tests/` ou co-localizado junto do código?

---

## 1. Primeiro, a pergunta de fundo: por que interfaces tornam o mock fácil em Go?

Em Go, **uma struct implementa uma interface automaticamente** (satisfação implícita). Não há `implements` nem herança.

Isso significa que, para testar, você cria uma struct falsa com os **mesmos métodos** da interface e a injeta onde a interface é esperada — o código de produção nem percebe.

```go
// Produção: a interface que o service consome
type PokemonRepository interface {
    GetByID(ctx, id) (*domain.Pokemon, error)
    GetAll(ctx, page, size) (*domain.PokemonPage, error)
    ...
}

// Teste: struct falsa com os MESMOS métodos → é automaticamente um PokemonRepository
type stubPokemonRepo struct { ... }
func (s *stubPokemonRepo) GetByID(...) { ... }   // basta ter o método
```

O teste injeta o falso no construtor do service. O service não sabe (nem quer saber) se está falando com Postgres, HTTP ou um stub.

---

## 2. Os 5 Test Doubles — com código real

### 2.1 Stub — retorna valores fixos (o mais comum)

**Objetivo:** substituir uma dependência para o código sob teste rodar com dados controlados. **Não verifica nada.**

**Real no projeto** — `core/app/auth-service/internal/http/handlers_test.go:19`:

```go
// Stub: campos prontos para retornar, o teste decide o que colocar
type stubAuthService struct {
    signupResult *service.AuthResult
    signupErr    error
    loginResult  *service.AuthResult
    loginErr     error
    refreshFn    func(ctx context.Context, token string) (*service.AuthResult, error)
    logoutFn     func(ctx context.Context, token string) error
}

func (s *stubAuthService) Signup(ctx, email, password) (*service.AuthResult, error) {
    return s.signupResult, s.signupErr
}
func (s *stubAuthService) Login(ctx, email, password) (*service.AuthResult, error) {
    return s.loginResult, s.loginErr
}
```

Uso no teste — define o comportamento antes da chamada:
```go
svc := &stubAuthService{
    signupResult: &service.AuthResult{UserID: "user-1", ...},
    signupErr:    nil,
}
mux := NewMux(svc)   // handler recebe o stub como AuthService
```

**Padrão Go:** struct + campos de retorno. Nada de framework.

---

### 2.2 Fake — implementação funcional simplificada (tem lógica)

**Objetivo:** um duplo com **comportamento real** (memória, lista, etc.), não apenas valores fixos. Pode ser reutilizado em produção também.

**Real no projeto** — `internal/adapters/outbound/memory/mock_repositories.go`:

```go
// Fake: armazena em memória, tem lógica real (add, busca, dedup)
type FavoriteRepository struct {
	mu        sync.RWMutex
	favorites map[string]map[string]bool
}

func NewFavoriteRepository() *FavoriteRepository {
	return &FavoriteRepository{
		favorites: make(map[string]map[string]bool),
	}
}

func (m *FavoriteRepository) AddFavorite(_ context.Context, userID, pokemonID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.favorites[userID]; !exists {
		m.favorites[userID] = make(map[string]bool)
	}
	m.favorites[userID][pokemonID] = true
	return nil
}
```

**Detalhe importante:** este `memory` package **não é só de teste** — o `main.go` usa como **fallback** quando o Postgres está fora:
```go
// cmd/server/main.go:90
favoriteRepo = memory.NewFavoriteRepository()   // fallback em produção
```
Fake serve teste **e** degradação. Um bom teste é usar exatamente o mesmo fake.

---

### 2.3 Spy — registra chamadas (verifica "foi chamado com X?")

**Objetivo:** além de responder, **grava** as chamadas para o teste verificar interação.

**Real no projeto** — o `stubPokemonRepo` em `catalog-service/internal/http/handlers_test.go` tem os dois papéis, mas o padrão spy aparece quando você usa function fields para capturar:

```go
// Espião de facto: function field que registra e delega
getAllFn: func(ctx context.Context, page, pageSize int) (*domain.PokemonPage, error) {
    chamadas = append(chamadas, page)   // registra!
    return pokemonPage(...), nil
},
```

E no handler test (auth-service) você verifica **o que foi passado** — o assert mora no function field que o teste programa, não no método do double:
```go
// handlers_test.go:99 — o teste "espia" o argumento dentro do fn field
isAccessTokenRevokedFn: func(ctx context.Context, jti string) (bool, error) {
    if jti != "jti-ativo" {
        t.Fatalf("jti inesperado: %s", jti)   // ASSERT dentro do spy
    }
    return false, nil
},
```
O spy de verdade teria `Calls []string` e `assert.Equal(t, []string{"25"}, spy.FindByIDCalls)`.

---

### 2.4 Mock — stub + spy + verificação (o mais "completo")

**Objetivo:** pré-programa comportamento E verifica que foi chamado exatamente como esperado. **Cuidado:** é o mais rígido — muda o código e o teste quebra (por design).

**Real no projeto** — `auth-service/internal/service/auth_service_test.go:17`:

```go
type mockAuthRepository struct {
    createUserFn    func(ctx, email, passwordHash string) (*repository.User, error)
    getByEmailFn    func(ctx, email string) (*repository.User, error)
    rotateRefreshTokenFn func(ctx, current, new, userID string, expiresAt time.Time) error
    isAccessTokenRevokedFn func(ctx, jti string) (bool, error)
    ...
}
```

E cada teste "programa" só o que precisa:
```go
repo := &mockAuthRepository{
    getByEmailFn: func(ctx, email) (*repository.User, error) {
        return &repository.User{ID: "user-1", PasswordHash: string(hash)}, nil
    },
    storeRefreshTokenFn: func(...) error { return nil },
}
```

**Go não usa framework de mock pesado** (Mockito/Cuckoo). A struct com function fields É o mock idiomático. Se quiser mocks gerados, há `mockery` (gera structs de teste), mas o projeto escolheu manual por explicitude.

---

### 2.5 Dummy — só preenche assinatura (nunca é usado)

**Objetivo:** um valor que o código pede mas que não influencia o teste. Em Go, o `nil` costuma bastar:

```go
// service_test.go:141 — passa nil como catalog porque o teste não usa favoritos
svc := service.NewFavoriteService(favoriteRepo, pokemonRepo, nil)
```

---

## 3. Resumo — tabela decisória

| Double | Responde? | Verifica? | Tem lógica? | Exemplo real |
|--------|-----------|-----------|-------------|--------------|
| **Dummy** | nunca usado | não | não | `nil` no `NewFavoriteService` |
| **Stub** | valores fixos | não | não | `stubAuthService` (handlers_test) |
| **Fake** | comportamento real | não | **sim** | `memory.NewFavoriteRepository()` |
| **Spy** | valores | **sim** (registra) | às vezes | `getAllFn` capturando args |
| **Mock** | pré-programado | **sim** (rigoroso) | não | `mockAuthRepository` (service_test) |

**Regra de ouro do projeto (e da comunidade Go):**
- Prefira **stub/fake** (menos acoplado) — teste o COMPORTAMENTO, não a implementação.
- Use **spy** quando precisar confirmar interação.
- Use **mock rigoroso** com parcimônia — ele acopla o teste à implementação interna.

---

## 4. `httptest` — o fake de servidor HTTP

Para testar **adapters outbound** (clients HTTP), o Go dá `httptest.NewServer`: um servidor HTTP real de verdade que você programa a resposta.

**Real no projeto** — `internal/adapters/outbound/http/favorite_catalog_client_test.go:18`:

```go
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    assert.Equal(t, "/v1/pokemons/favorites", r.URL.Path)
    assert.Equal(t, "1,25", r.URL.Query().Get("ids"))
    _ = json.NewEncoder(w).Encode([]domain.Pokemon{...})
}))
defer srv.Close()

client := httpclient.NewFavoriteCatalogClient(srv.URL)   // aponta pro fake
result, err := client.GetFavoriteDetails(ctx, []string{"1", "25"})
```

**Isso é um Fake** (servidor HTTP real) com **Spy** (assert no request). O client de produção fala com `srv.URL` como se fosse o catalog-service real.

Para **adapters inbound** (handlers), o inverso: `httptest.NewRecorder()` captura a resposta:
```go
req := httptest.NewRequest(http.MethodGet, "/v1/pokemons", nil)
w := httptest.NewRecorder()
mux.ServeHTTP(w, req)          // roda o handler de verdade
if w.Code != http.StatusOK { ... }
```

---

## 5. Onde colocar os testes? `tests/` vs co-localizado

### A pergunta real

Em JavaScript (Jest) o padrão é `__tests__/` ou `.test.ts` ao lado do código. Em Go, existem dois padrões:

| Padrão | Onde | Quem usa no projeto |
|--------|------|---------------------|
| **Co-localizado** | `internal/service/auth_service_test.go` ao lado de `auth_service.go` | ✅ todos os serviços (auth, catalog, bff) |
| **Pasta separada** | `tests/integration/` | ✅ mobile-bff (integração) |

### A resposta curta

**Co-localizado é o padrão idiomatico de Go** (e o que a comunidade recomenda):
- `go test ./...` acha sozinho (arquivo `*_test.go` no mesmo pacote).
- Fácil navegar: `foo.go` ⇄ `foo_test.go`.
- Testes de **unidade** (mesmo pacote) têm acesso a funções não exportadas (white-box).
- O editor "jump to test" funciona por padrão.

**Pasta `tests/` separada** é **exceção** — usada apenas para **integração** (requer DB real).

### A decisão que foi executada no projeto

O mobile-bff **antes** usava `tests/unit/` (package `unit`) separado do código, com `tests/mocks/` re-exportando fakes. Foi **migrado** para o padrão co-localizado:

```
antes:                                       depois:
tests/unit/handlers_test.go          →  internal/adapters/inbound/http/handlers_test.go
tests/unit/service_test.go           →  internal/service/service_test.go
tests/unit/domain_test.go            →  internal/domain/domain_test.go
tests/unit/auth_client_test.go       →  internal/adapters/outbound/http/auth_service_client_test.go
tests/mocks/mock_repositories.go     →  internal/adapters/outbound/memory/ (import direto)
tests/integration/                   →  mantido (padrão correto p/ integração)
```

**Como ficou sem perder black-box:** os testes migrados usam `package xxx_test` (mesmo diretório, pacote externo). Isso mantém o teste **black-box** (só API pública — como um cliente usaria) MAS co-localizado (navegação fácil, sem `tests/mocks/`). O `internal/` já esconde o que precisa esconder; não é preciso uma pasta separada para isso.

**O que restou em `tests/`:** apenas `tests/integration/` — testes que sobem Postgres real, precisam de env vars e não devem rodar no `go test` comum.

### Recomendação concreta (o estado atual)

| Tipo de teste | Onde colocar | Exemplo |
|---------------|--------------|---------|
| **Unit de service/handler/repo** | **co-localizado** (`internal/.../*_test.go`) | `internal/service/service_test.go` ✅ |
| **Unit de adapter HTTP client** | co-localizado (ao lado do `.go`) | `internal/adapters/outbound/http/*_test.go` ✅ |
| **Teste de domínio puro** | co-localizado | `internal/domain/domain_test.go` ✅ |
| **Integração (Postgres real)** | `tests/integration/` (fora de `internal/`) | mobile-bff ✅ |

**Por que integração separada?** Precisa de build tag (`//go:build integration`), Docker, env vars — não deve rodar no `go test` comum. Isolar em `tests/integration/` deixa claro.


---

## 6. Padrão de ouro para escrever o teste

1. **Encontre a dependência a falsificar** — procure a interface que o código sob teste recebe.
2. **Escolha o double certo** — stub (dado fixo) / fake (lógica) / spy (interação) / mock (rigoroso).
3. **Implemente os métodos** da interface na struct de teste.
4. **Injete no construtor** — `service.NewXxx(stub)`.
5. **Teste o comportamento observável**, não a implementação interna.

```go
// O teste não sabe se o repo é Postgres, HTTP ou stub:
svc := service.NewPokemonService(stubRepo, stubFavRepo)

page, err := svc.ListPokemons(ctx, 0, 20, "")

assert.NoError(t, err)
assert.Equal(t, 2, len(page.Content))   // comportamento
```

---

## Referências no projeto

- Teoria de test doubles: `doc/aprendizado/fase-07-testes-avancados-design.md`
- Stub real: `core/app/auth-service/internal/http/handlers_test.go`
- Mock real: `core/app/auth-service/internal/service/auth_service_test.go`
- Fake real (produção + teste): `core/bff/mobile-bff/internal/adapters/outbound/memory/mock_repositories.go`
- httptest real: `core/bff/mobile-bff/internal/adapters/outbound/http/favorite_catalog_client_test.go`
- Regra arquitetural: `AGENTS.md` — "O package `tests/` nunca é importado por código de produção"
