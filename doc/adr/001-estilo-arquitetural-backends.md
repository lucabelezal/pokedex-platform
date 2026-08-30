# ADR-001: Estilo Arquitetural dos Backends

**Status:** Aceito
**Data:** 2026-08-29
**Escopo:** `core/bff/mobile-bff`, `core/app/auth-service`, `core/app/pokemon-catalog-service`
**Relacionado:** [DECISIONS.md](../DECISIONS.md) (Decisão 4), [architecture/hexagonal.md](../architecture/hexagonal.md)

---

## Contexto

A plataforma tem três serviços Go:

| Serviço | Responsabilidade | Estilo atual |
|---------|------------------|--------------|
| `mobile-bff` | Orquestração voltada ao cliente, autenticação, favoritos | Hexagonal (Ports & Adapters) |
| `auth-service` | Ciclo de vida de tokens JWT | Camadas simples (handler → service → repository) |
| `pokemon-catalog-service` | Fonte canônica do catálogo | Camadas finas (handler → repository + domain) |

A base de código já reflete, portanto, uma arquitetura **híbrida**: hexagonal onde a complexidade de integração justifica, e camadas simples onde o serviço é essencialmente CRUD ou API fina.

Este ADR registra **por que** essa distribuição é intencional e **quando** evoluir um serviço para um estilo com mais boundaries. Ele também posiciona a discussão do ecossistema Spring (MVC tradicional, Spring Modulith) como referência conceitual para decisões futuras, já que os serviços são Go hoje mas os princípios de modularidade são independentes de linguagem.

## Drivers de Decisão

1. **Testabilidade do núcleo** — regras de negócio devem ser testáveis sem banco, HTTP ou rede.
2. **Velocidade de entrega** — cerimônia arquitetural tem custo real por feature; só paga quando compra algo concreto.
3. **Onboarding** — cada camada/interface adiciona carga cognitiva a quem chega.
4. **Domínio de baixa complexidade** — catálogo, favoritos e autenticação têm regras simples; não exigem 4-5 camadas.
5. **Múltiplas saídas** — o BFF é o único serviço que conversa com banco + múltiplos clients HTTP + cache; isolamento de infraestrutura tem valor real nele.
6. **Evidência empírica** — não existe prova de que mais camadas aumentam produtividade ou reduzem bugs; existe evidência de que criam overhead.

## Decisão

**Manter modelo evolutivo e híbrido, com regra de abstração explícita:**

1. **`mobile-bff`**: manter hexagonal completo (`domain` → `ports` → `service` → `adapters`). Justificado por ser orquestrador com várias saídas concretas (PostgreSQL, clients HTTP, cache) e por precisar de doubles para testes de integração. *(Registrado na Decisão 4 do DECISIONS.md.)*
2. **`auth-service` e `pokemon-catalog-service`**: manter camadas simples (handler → service → repository). Sem ports/adapters adicionais enquanto não houver segunda implementação concreta ou uma boundary que proteja o núcleo.
3. **Regra de abstração** (aplicada em todos os serviços):
   - Crie interface apenas quando houver **segunda implementação concreta à vista** (ex: mock em memória + PostgreSQL real) **ou** quando a interface reside em anel interno protegendo o núcleo da infraestrutura (DIP real).
   - Não crie interface para "poder trocar no futuro". "Talvez um dia" não é razão — é `Speculative Generality`.
   - Não separe DTO/persistence/domain models por camada só porque a arquitetura "manda". Propague o DTO de API até a aplicação; separe apenas quando o mesmo caso de uso for exposto por múltiplos canais (REST + gRPC + Kafka).
4. **Direção de evolução**: se o BFF crescer além de um limiar (ver "Critérios de Evolução"), evoluir para **modular monolith por feature/domínio com fronteiras internas**, não para mais camadas horizontais.

## Alternativas Consideradas

### A. Clean Architecture clássica (4 círculos: entities → use cases → interface adapters → frameworks)

**Rejeitada como padrão obrigatório.** A Dependency Rule (Uncle Bob) é um princípio sólido e já é respeitada no projeto (dependências apontam para dentro). Mas a implementação canônica — use case por feature, DTO por camada, interface para cada dependência — adiciona indireção e arquivos sem comprar nada para domínios simples. O próprio Uncle Bob descreve os círculos como esquemáticos ("Only Four Circles? No, the circles are schematic"); o que importa é a regra de dependência, não o número de camadas.

### B. Spring MVC tradicional (controller → service → repository)

**Referência, não aplicável ao Go.** É a linha de base pragmática: pouca cerimônia, alta velocidade inicial. A limitação conhecida é que, sem boundaries explícitas, a lógica de negócio tende a vazar para o handler e o repository ao longo dos anos (o problema que Cockburn descreve no artigo original: a camada "nova" acaba cheia de lógica). No Go, os serviços de app já mitigam isso com `domain/` isolado e `service/` orquestrando.

### C. Spring Modulith

**Referência conceitual para evolução.** Drotbohm e o projeto Spring Modulith formalizam o meio-termo: módulos por domínio, API pública por módulo, dependências verificáveis entre módulos, sem exigir Clean Architecture completa. A ideia central — **boundaries por módulo de negócio, não por camada horizontal** — é exatamente a direção que este ADR adota para o Go (ver "Direção de Evolução"). O ecossistema Spring hoje oferece isso nativamente; no Go, o equivalente é estrutura de pastas por feature + `internal/` para ocultar implementação.

### D. Modular Monolith / Vertical Slice

**Adotado como direção de evolução.** Organizar por feature (`orders/create`, `orders/cancel`) em vez de por tipo técnico (`controllers/`, `services/`, `repositories/`) reduz o custo cognitivo de navegação e mantém boundaries entre domínios. Simon Brown defende exatamente este ponto: dá para ter boundaries arquiteturais sem cerimônia de Clean Architecture.

### E. Microservices por serviço de domínio

**Rejeitado.** Fowler (MonolithFirst e Microservice Premium) é explícito: a maioria dos sistemas deve começar monolítica; o premium de distribuição (deploy, observabilidade, consistência eventual, falhas de rede) só se paga quando o monólito fica grande demais para evoluir. A plataforma atual é pequena; extrair mais serviços agora adicionaria custo sem comprar valor.

## Comparação das Abordagens

| Abordagem | Complexidade | Testabilidade | Acoplamento | Velocidade inicial | Onboarding | Evolução |
|-----------|-------------|---------------|-------------|--------------------|------------|----------|
| MVC simples | Baixa | Boa | Maior | Alta | Fácil | Média (lógica vaza) |
| MVC + módulos por feature | Baixa-média | Boa | Média-baixa | Alta | Fácil | Boa |
| Hexagonal | Média-alta | Muito boa | Baixa | Média | Médio | Boa |
| Clean clássica | Alta | Muito boa | Muito baixa | Baixa | Difícil | Boa, mas cara |
| Modular Monolith | Média | Muito boa | Baixa | Média-alta | Médio | Muito boa |

Conclusão prática da matriz: o ponto ótimo para a maioria dos backends está entre **"MVC + módulos"** e **"Hexagonal"** — não em Clean clássica.

## Evidência

### Evidência empírica

- **Estudo 2026 — "A Survey-Based Empirical Evaluation of Clean Architecture"** (Zenodo, DOI 10.5281/zenodo.20676263): survey com 106 participantes. Correlação positiva entre adoção de Clean Architecture e qualidade percebida (`r = 0.428`) e forte correlação de gerenciamento de dependências com qualidade (`r = 0.707`) e produtividade (`r = 0.691`). **Ressalva crítica**: é survey, não experimento controlado — mostra associação/percepção, não causalidade. Os próprios autores relatam curva de aprendizado íngreme e complexidade adicional, principalmente em onboarding.
- **Conclusão da evidência**: o que tem suporte é "desenvolvedores que usam Clean relatam melhor qualidade/percepção de produtividade". NÃO tem suporte "Clean aumenta produtividade em X% ou reduz bugs em Y%".

### Autores e fontes autoritativas

| Autor/Fonte | Por que é referência | Ponto central para este ADR |
|-------------|----------------------|------------------------------|
| **Alistair Cockburn** (criador da Hexagonal) | Fonte primária, 2005 | O objetivo original é isolar o núcleo de UI/banco para teste e reuso headless — **não** criar muitas camadas. "O hexágono não é um hexágono porque o número seis importa" — são quantos ports você precisar. |
| **Martin Fowler** | Trade-offs e evolução | Presentation Domain Data Layering: as 3 camadas são uma boa modularização padrão. MonolithFirst e Microservice Premium: não comece distribuído. |
| **Robert C. Martin (Uncle Bob)** | Autor da Clean Architecture | A Dependency Rule (dependências apontam para dentro) é o princípio que vale; os círculos são esquemáticos, não um número fixo de camadas. |
| **Eric Evans / Vaughn Vernon** | DDD | Boundaries (Bounded Context) justificam-se pela complexidade do domínio. DDD ≠ Clean/Hexagonal; DDD dá a razão para separar, os estilos dão o mecanismo. |
| **Simon Brown** | Modular Monolith, C4 | É possível ter boundaries arquiteturais sem cerimônia; arquitetura está no código, não no número de pastas. |
| **Mark Richards** | "Software Architecture: The Hard Parts" | Arquitetura é trade-off, não religião; toda decisão tem custo. |
| **Oliver Drotbohm** | Spring Modulith | O meio-termo institucionalizado: módulos por domínio com dependências verificáveis, evoluindo da simplicidade para a sofisticação conforme a necessidade. |
| **AWS Prescriptive Guidance** | Prática em produção | Advertência explícita: o adapter só se justifica se houver múltiplas entradas/saídas ou troca provável de infraestrutura; caso contrário, vira "outra camada para manter" + possível latência extra. |

### Críticas de overengineering (para calibrar o lado contrário)

- **Victor Rentea — "Overengineering in Onion/Hexagonal Architectures"**: interface útil existe se (a) tem ≥2 implementações, (b) protege um anel interno (DIP), ou (c) está em biblioteca de cliente. Separe DTOs apenas com múltiplos canais. Não separe persistence model do domain model por dogma.
- **Reza Enayati — "What 'clean architecture' actually costs you"**: custo por feature (9 arquivos para uma rota que seria 30 linhas), navegação mais lenta (go-to-definition cai na interface), onboarding mais caro. Regra prática: se você consegue dizer em duas frases o que a arquitetura compra no próximo mês, vale; se não, é decoração.
- **r/golang — "Why Clean Architecture ... Don't Belong in GoLang"**: interfaces implícitas do Go tornam a cerimônia de interfaces explícitas redundante; copiar estrutura de Java/Spring para Go cria abstrações sem comportamento.
- **Reddit r/softwarearchitecture — "Is hexagonal architecture good design or just extra layers?"**: consenso prático — hexagonal paga quando há domínio complexo ou múltiplas integrações; é cerimônia em CRUD simples.

## Consequências

### Positivas

- Cada serviço tem exatamente a estrutura que seu problema exige — nada de interfaces que nunca terão segunda implementação.
- O domínio dos serviços de app fica testável sem infraestrutura (o `domain/` isolado já garante isso).
- O BFF mantém testabilidade de integração com doubles (adapters em memória), que é onde ele mais precisa.
- Decisão documentada → novas features seguem o padrão sem rediscutir arquitetura toda vez.

### Negativas / Trade-offs aceitos

- **Inconsistência estrutural entre serviços**: BFF hexagonal vs app services em camadas. Um dev novo precisa entender os dois estilos. Mitigação: este ADR + regras de abstração explícitas.
- **Acoplamento maior nos services de app**: handler ↔ service ↔ repository compartilham DTOs. Aceito porque não há segunda implementação concreta nem múltiplos canais.
- **Risco de vazamento de regra de negócio** para handler/repository nos services de app ao longo do tempo. Mitigação: revisão de código deve repelir regra de negócio fora de `domain/`/`service/`.

### Riscos e Mitigação

| Risco | Mitigação |
|-------|-----------|
| BFF vira god service com cerimônia excessiva | Aplicar regra de abstração: só port quando houver segunda implementação ou DIP real |
| Regra de negócio vaza para camadas externas nos services de app | Checklist de conformidade (como o de `hexagonal.md`) + review |
| Alguém "completa" a arquitetura por dogma (adiciona ports sem necessidade) | Este ADR é normativo; citar regra de abstração em review |
| Domínio cresce e as 3 camadas já não bastam | Acionar "Critérios de Evolução" abaixo |

## Critérios de Evolução

**Adicionar ports/adapters a um serviço de app quando** (qualquer um):
- aparecer uma segunda implementação concreta para uma dependência (ex: repository alternativo, provider de pagamento, fila);
- o mesmo caso de uso precisar ser exposto por múltiplos canais (REST + gRPC + consumer de Kafka);
- um teste exigir double estável de um client externo (adapter em memória) e a interface de mock for menos custosa que o acoplamento direto.

**Evoluir o `mobile-bff` para modular monolith por feature quando**:
- o número de features/aggregates dificultar navegação (a organização por tipo técnico — `ports/`, `service/`, `adapters/` — começa a esconder o que o sistema faz);
- surgirem dois bounded contexts claramente distintos no BFF (ex: "favoritos" e "perfil") que devam evoluir independentemente;
- equipe crescer além de ~4 devs trabalhando em paralelo no mesmo módulo.

Nesse ponto, a estrutura recomendada é:

```
internal/
  favorites/          ← módulo por domínio
    application/
    domain/
    infrastructure/
    web/
  auth/
    application/
    domain/
    infrastructure/
    web/
```

com ports/adapters apenas nos módulos que realmente têm múltiplas saídas — não em todos.

**Não evoluir** quando a motivação for só "padronizar igual ao BFF" ou "ficar mais clean". Sem driver concreto, a mudança é decoração.

## O teste definitivo para cada abstração

Antes de criar interface, DTO, adapter ou camada, responder:

1. Que dependência estou impedindo?
2. Que mudança estou tornando barata?
3. Que teste fica melhor?
4. Quanta complexidade estou adicionando?
5. Essa mudança é plausível no horizonte do sistema?

Se as respostas forem "não sei", "talvez" ou "nenhuma", a abstração provavelmente não deveria existir.

## Verificação arquitetural (fitness functions)

Arquitetura documentada é insuficiente. Uma regra como "`domain/` não pode importar infraestrutura" só vale de verdade se algo falha quando violada.

No Go:

- convenções de pacote + `internal/` para ocultar implementação;
- testes de dependência que percorrem os imports (fitness functions);
- regras de CI que falham o build em violação;
- `go list` / análise estática para detectar ciclos ou imports ilegais entre módulos/features.

Transformar cada regra estrutural importante em teste automatizado (ex: "feature A não acessa `internal/` de feature B", "módulos não formam ciclos").

## Métricas para avaliar se a decisão está funcionando

Este ADR não deve ser avaliado só por opinião:

- **Change amplification** — arquivos alterados por feature (alvo: menor que clean clássica);
- **Navegação** — saltos até a implementação concreta (go-to-definition não deve cair numa interface sem uso);
- **Velocidade** — tempo médio por feature;
- **Testes** — tempo de execução da suíte;
- **Acoplamento** — dependências cruzadas entre módulos/features;
- **Onboarding** — tempo para novo dev fazer uma mudança pequena.

Não é preciso formalizar todos como dashboard. São preferíveis a argumentos puramente estéticos. Exemplo: `POST /orders` = 3 arquivos (handler/service/repository) vs 10+ na clean clássica. A pergunta é se os 10 compraram uma mudança mais barata; se não, é custo sem retorno.

## Checklist de code review arquitetural

**Interfaces**
- [ ] Existe uso real da interface?
- [ ] Há mais de uma implementação relevante (ou DIP real)?
- [ ] A interface pertence ao consumidor?

**Dependências**
- [ ] `domain/` conhece framework ou SQL?
- [ ] Infraestrutura invade application?
- [ ] Existem dependências circulares?

**DTOs**
- [ ] Os contratos evoluem independentemente?
- [ ] Existe múltiplo canal justificando a separação?

**Estrutura**
- [ ] A feature está fácil de encontrar?
- [ ] Código que muda junto está próximo?
- [ ] Existem arquivos que apenas delegam?

## Nota sobre o ecossistema Spring (referência para futuros projetos)

Para um futuro backend Kotlin/Spring, a mesma lógica aplica: não partir de Clean Architecture como dogma. A via recomendada seria Spring MVC + modularização por domínio (equivalente ao Spring Modulith) + hexagonal seletiva onde houver múltiplas integrações. O próprio ecossistema Spring reforça isso: o Modulith existe justamente para oferecer boundaries por módulo sem impor a cerimônia completa.

## Referências

**Fontes primárias (autores):**
- Alistair Cockburn, *The Hexagonal (Ports & Adapters) Architecture* (2005) — https://alistair.cockburn.us/hexagonal-architecture/
- Robert C. Martin, *The Clean Architecture* (2012) — https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- Martin Fowler, *Presentation Domain Data Layering* — https://martinfowler.com/bliki/PresentationDomainDataLayering.html
- Martin Fowler, *Monolith First* — https://martinfowler.com/bliki/MonolithFirst.html
- Martin Fowler, *Microservice Premium* — https://martinfowler.com/bliki/MicroservicePremium.html
- Eric Evans, *Domain-Driven Design* (2004) — https://www.domainlanguage.com/
- Vaughn Vernon, *Implementing Domain-Driven Design* (2013)
- Mark Richards & Neal Ford, *Software Architecture: The Hard Parts* (2021)

**Meio-termo e modularidade:**
- Simon Brown, *Modular Monoliths* — https://simonbrown.je/ / palestras C4
- Oliver Drotbohm, *Spring Modulith* — https://docs.spring.io/spring-modulith/
- Netflix, *Ready for changes with Hexagonal Architecture* (2020) — https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749

**Evidência empírica:**
- Senadheera, Somaweera & Sandaruwan, *A Survey-Based Empirical Evaluation of Clean Architecture* (2026), IMPETUS, DOI 10.5281/zenodo.20676263 — https://doi.org/10.5281/zenodo.20676263

**Guia prático de produção:**
- AWS Prescriptive Guidance, *Hexagonal architecture pattern* — https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/hexagonal-architecture.html

**Críticas / calibragem de overengineering:**
- Victor Rentea, *Overengineering in Onion/Hexagonal Architectures* — https://victorrentea.ro/blog/overengineering-in-onion-hexagonal-architectures/
- Reza Enayati, *What 'clean architecture' actually costs you* — https://rezaenayati.me/backend/clean-architecture-cost
- r/golang, *Why Clean Architecture and Over-Engineered Layering Don't Belong in GoLang* — https://www.reddit.com/r/golang/
- r/softwarearchitecture, *Is hexagonal architecture good design or just extra layers?* — https://www.reddit.com/r/softwarearchitecture/

**Vídeos:**
- Alistair Cockburn, *Entre puertos y adaptadores* (Ágiles 2020, legendas ES/EN) — https://www.youtube.com/watch?v=Sc_B6dQ6di0
- Alistair Cockburn, *Alistair in the Hexagone* — https://www.youtube.com/watch?v=th4AgBcrEHA
- Three Dots Labs, *Is Clean Architecture Overengineering?* (No Silver Bullet)
- Oliver Gierke, *Why Hexagonal et. al. are answers to the wrong question* — https://youtu.be/co3acmgP2Ng
