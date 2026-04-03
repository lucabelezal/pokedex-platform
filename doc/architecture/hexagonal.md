# Arquitetura Hexagonal (Ports & Adapters)

> Baseado no artigo original de Alistair Cockburn (2005) e nas convenções adotadas neste projeto.

---

## 1. Conceito Central

A arquitetura hexagonal organiza um sistema em torno de uma **aplicação central** (o hexágono) que não conhece o mundo externo. Toda comunicação com o exterior ocorre exclusivamente por meio de **portas** (_ports_) e **adaptadores** (_adapters_).

O objetivo é garantir que a lógica de negócio possa ser testada e executada de forma independente de qualquer tecnologia de infraestrutura: banco de dados, framework HTTP, serviço externo, CLI, etc.

---

## 2. Elementos da Arquitetura

### Hexágono (Aplicação)

O núcleo do sistema. Contém:

- **Domínio**: entidades, regras de negócio, erros de domínio.
- **Casos de uso**: orquestração da lógica de negócio.

O hexágono **não importa** pacotes de infraestrutura. Toda dependência de recurso externo é expressa como uma interface definida pelo próprio hexágono.

### Portas (_Ports_)

Interfaces que o hexágono expõe ou exige. Existem dois tipos:

| Tipo | Também chamado | Direção do fluxo |
|------|----------------|-----------------|
| **Inbound** | Driving port / Primary port | Externo → Aplicação |
| **Outbound** | Driven port / Secondary port | Aplicação → Externo |

**Porta inbound** define o contrato que o mundo externo usa para acionar a aplicação — por exemplo, casos de uso que um handler HTTP invoca.

**Porta outbound** define o contrato que a aplicação usa para acessar recursos externos — por exemplo, repositórios de banco de dados ou clientes de serviços externos.

### Adaptadores (_Adapters_)

Implementações concretas das portas. Traduzem entre o protocolo do mundo externo e a linguagem do hexágono.

| Tipo | Também chamado | Exemplo |
|------|----------------|---------|
| **Adaptador inbound** | Driving adapter / Primary adapter | Handler HTTP, CLI, consumer de fila |
| **Adaptador outbound** | Driven adapter / Secondary adapter | Repository PostgreSQL, client HTTP externo |

---

## 3. Fluxo de Dependências

A regra fundamental é: **as dependências apontam para dentro do hexágono**.

```
[Mundo Externo]
      │
      ▼
[Adaptador Inbound]──────────────────────────────────────────────
      │ usa porta inbound                                         │
      ▼                                                           │
[Caso de Uso]                                                     │
      │ usa porta outbound                                        │
      ▼                                                           │
[Porta Outbound] ← implementada por → [Adaptador Outbound]───────
                                              │
                                              ▼
                                      [Recurso Externo]
                                   (banco, serviço, cache)
```

O caso de uso **nunca depende** de um adaptador concreto — apenas de uma interface (porta). Isso permite substituir qualquer implementação externa sem alterar a lógica de negócio.

---

## 4. Estrutura de Pastas neste Projeto

O `mobile-bff` implementa a arquitetura hexagonal da seguinte forma:

```
internal/
├── domain/                        # Hexágono: entidades e erros de domínio
│   ├── auth_session.go            # AuthSession (movida dos ports para o domínio)
│   ├── errors.go
│   └── pokemon.go
│
├── ports/
│   ├── inbound/                   # Contratos que os adaptadores inbound consomem
│   │   └── usecase.go             # PokemonUseCase, FavoriteUseCase, AuthUseCase
│   └── outbound/                  # Contratos que a aplicação exige de recursos externos
│       ├── auth.go                # AuthProvider
│       └── repository.go         # PokemonRepository, FavoriteRepository
│
├── service/                       # Hexágono: implementações dos casos de uso
│   ├── auth_service.go            # implementa inbound.AuthUseCase
│   └── pokemon_service.go         # implementa inbound.PokemonUseCase e FavoriteUseCase
│
└── adapters/
    ├── http/                      # Adaptadores inbound: recebem requisições externas
    │   ├── handlers.go            # depende de inbound.*
    │   └── enriched_response_builder.go
    └── repository/               # Adaptadores outbound: acessam recursos externos
        ├── auth_service_client.go # implementa outbound.AuthProvider
        └── mock_repository.go    # implementa outbound.PokemonRepository e FavoriteRepository
```

**Regra de dependência aplicada:**

```
adapters/http  →  ports/inbound  ←  service
                                        │
                                   ports/outbound
                                        │
                                  adapters/repository
```

Os handlers HTTP importam apenas `ports/inbound`. Os serviços (casos de uso) importam `ports/outbound`. Os adaptadores de repositório implementam `ports/outbound`. Nenhuma camada interna depende de uma camada externa.

---

## 5. Por Que Separar Inbound e Outbound?

Manter todos os contratos no mesmo pacote `ports` é funcional, mas obscurece a direção do fluxo. Separar em sub-pacotes torna explícito:

- **Quem aciona a aplicação** (`inbound`): casos de uso que os adaptadores HTTP chamam.
- **De quem a aplicação depende** (`outbound`): repositórios e clientes externos que os serviços usam.

Isso também evita que adaptadores inbound (HTTP) dependam inadvertidamente de contratos outbound (repositórios).

---

## 6. Testabilidade

A principal vantagem da arquitetura hexagonal é a facilidade de testar:

**Testes unitários de casos de uso**: substituem os adaptadores outbound por stubs simples.

```go
// Stub que implementa outbound.AuthProvider sem infraestrutura real
type stubAuthProvider struct {
    session *domain.AuthSession
    err     error
}

func (s *stubAuthProvider) Login(ctx context.Context, email, password string) (*domain.AuthSession, error) {
    return s.session, s.err
}
```

**Testes de adaptadores inbound**: substituem os casos de uso por stubs que implementam `inbound.*`.

```go
// Stub que implementa inbound.AuthUseCase para testar handlers HTTP
type stubAuthUseCase struct {
    session *domain.AuthSession
    err     error
}
```

Desta forma, é possível testar toda a lógica de negócio sem banco de dados, servidor HTTP ou rede.

---

## 7. Checklist de Conformidade

Use este checklist ao adicionar código novo:

- [ ] O domínio (`domain/`) não importa nenhum pacote externo a ele mesmo.
- [ ] Os casos de uso (`service/`) importam apenas `domain/` e `ports/outbound/`.
- [ ] Os adaptadores inbound (`adapters/http/`) importam apenas `ports/inbound/`.
- [ ] Os adaptadores outbound (`adapters/repository/`) implementam `ports/outbound/` e importam apenas `domain/`.
- [ ] Nenhum handler HTTP importa diretamente um repositório concreto.
- [ ] Novas entidades de domínio vivem em `domain/`, não em `ports/`.

---

## Referências

- Alistair Cockburn, *Hexagonal Architecture* (2005) — [https://alistair.cockburn.us/hexagonal-architecture/](https://alistair.cockburn.us/hexagonal-architecture/)
- *Pattern: Ports and Adapters* — Livro "Implementing Domain-Driven Design", Vaughn Vernon
