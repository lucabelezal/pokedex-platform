# Gateway

## Objetivo

O Kong é o ponto de entrada da plataforma. Ele fornece uma camada pública única de acesso na frente das aplicações internas.

## Papel Atual

- Receber tráfego externo.
- Roteirizar requisições para o `mobile-bff`.
- Esconder a topologia interna dos serviços para os clientes.

## Configuração Atual

O repositório usa uma configuração declarativa do Kong em `core/gateway/kong/kong.yml`.

No momento:

- o Kong expõe `/v1`
- as requisições são encaminhadas para `http://mobile-bff:8080/api`

## Notas Arquiteturais

O gateway está intencionalmente enxuto neste momento. Isso é uma boa escolha para uma fase inicial do projeto, porque mantém o roteamento simples enquanto o domínio ainda está evoluindo.

## Oportunidades De Melhoria

- Documentar com mais clareza as rotas públicas e a responsabilidade de cada rota.
- Adicionar tracing de requisições ou correlation headers entre gateway e serviços downstream.
- Decidir se a autenticação deve continuar totalmente no BFF ou se parte dela pode migrar para o gateway no futuro.
- Adicionar rate limiting e padronização de erros caso a plataforma fique mais exposta externamente.
