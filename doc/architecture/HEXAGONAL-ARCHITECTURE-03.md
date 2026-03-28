# Arquitetura Hexagonal 03

## O Que Vamos Ver

- Se esse estilo só vale a pena para projetos grandes.
- Quais são os trade-offs reais.
- Por que a crítica ao boilerplate é parcialmente verdadeira.
- Como decidir quanto de arquitetura é suficiente.

## Hexagonal Só Vale Para Projetos Grandes?

Não necessariamente, mas normalmente ela passa a ter mais valor quando o projeto começa a ter pelo menos uma destas características:

- múltiplas integrações
- múltiplos desenvolvedores
- um BFF que orquestra mais de uma dependência
- mais de um serviço com fronteiras claras de responsabilidade
- código que se espera que viva e evolua por algum tempo

Para um CRUD muito pequeno, esse estilo pode sim ser excessivo.

Para uma plataforma com BFF, gateway, auth, catálogo e preocupações de infraestrutura, alguma disciplina arquitetural costuma compensar.

Então a pergunta melhor não é "o projeto é grande?".

A pergunta melhor é:

"Does this project already suffer, or will it soon suffer, from unclear boundaries?"
"Este projeto já sofre, ou em breve sofrerá, com fronteiras pouco claras?"

## A Crítica Ao Boilerplate

Essa crítica é real.

Se um time adota arquitetura hexagonal de forma ruim, frequentemente cria:

- DTOs repetitivos sem propósito
- interfaces usadas uma única vez sem uma fronteira relevante
- camadas demais que apenas repassam dados
- uma falsa sensação de qualidade de design

É por isso que alguns desenvolvedores chamam conversa de arquitetura de "lorota". Eles geralmente já viram a versão ruim.

A correção não é abandonar estrutura. A correção é exigir que cada camada tenha uma razão de existir.

Perguntas úteis:

- Esta interface protege uma fronteira real?
- Este adapter isola um detalhe de tecnologia?
- Este service implementa um caso de uso de verdade?
- Este modelo pertence ao domínio ou apenas ao transporte?

Se a resposta for "não", a camada provavelmente é ruído.

## Por Que Empresas Ainda Vão Nessa Direção

Mesmo com a cerimônia extra, muitas empresas continuam indo em direção a hexagonal, clean architecture ou designs influenciados por DDD porque essas abordagens otimizam mudança de longo prazo.

Os benefícios mais comuns são:

- teste mais fácil dos fluxos de negócio
- substituição mais fácil de detalhes técnicos
- onboarding mais simples por meio de fronteiras mais claras
- menos acoplamento entre código de framework e regras de negócio
- colaboração melhor porque as responsabilidades estão nomeadas e separadas

Em outras palavras, empresas muitas vezes aceitam um pouco mais de código em troca de menos confusão.

## A Posição Que Queremos Neste Projeto

O melhor caminho para este projeto não é nem:

- "everything in handlers and repositories"

nem:

- "ten abstractions before any real logic"

O melhor caminho é o do meio:

- usar arquitetura para tornar as fronteiras explícitas
- manter os nomes alinhados com a responsabilidade de negócio
- evitar camadas desnecessárias
- refatorar apenas onde o acoplamento começa a machucar

É por isso que uma versão leve e amigável para Go da arquitetura hexagonal é uma boa escolha aqui.

## Heurísticas De Decisão Para Mudanças Futuras

Ao adicionar código novo, estas perguntas ajudam:

### Isso Deve Ir Para O Domínio?

Coloque ali se isso expressar uma regra de negócio que não deve depender de transporte ou persistência.

### Isso Deve Ser Um Caso De Uso?

Coloque na camada de service de aplicação se isso coordena um fluxo de negócio entre ports.

### Isso Deve Ser Um Adapter?

Coloque ali se isso traduz entre a aplicação e uma tecnologia ou protocolo específico.

### Isso Deve Virar Um Novo Serviço?

Considere isso quando:

- a capacidade tiver posse clara
- as preocupações de escala forem diferentes
- o domínio passar a ter significado próprio
- múltiplos clientes precisarem da mesma capacidade independentemente do BFF atual

## Reflexão Final

Arquitetura não é sobre provar sofisticação. É sobre tornar mudança mais segura.

Se uma estrutura ajuda o time a responder "onde este código deve ficar?" com menos hesitação, ela já está fazendo um trabalho valioso.

Se ela cria mais confusão do que clareza, está pesada demais.

O resultado certo para este projeto é uma arquitetura que seja:

- clara o suficiente para guiar o time
- pequena o suficiente para continuar prática
- explícita o suficiente para sustentar crescimento

## Recapitulação

Arquitetura hexagonal não é bala de prata e nem exigência universal. Ela é um trade-off. Para este projeto, esse trade-off vale a pena porque o sistema já tem múltiplas responsabilidades e pontos de integração. A chave é manter a abordagem leve, intencional e sustentada por valor real de manutenção.
