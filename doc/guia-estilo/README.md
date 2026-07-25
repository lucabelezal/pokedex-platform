# Guia de Estilo Go — Pokedex Platform

[Guia](guide.md) | [Decisões](decisions.md) | [Melhores Práticas](best-practices.md) | [Regras do Projeto](regras-projeto.md)

## Sobre

O Guia de Estilo Go e os documentos auxiliares estabelecem as melhores práticas para
escrever código Go legível e idiomático na Pokedex Platform. A adesão ao guia não é
absoluta, e estes documentos nunca serão exaustivos. A intenção é minimizar as dúvidas
na escrita de Go legível e unificar a orientação de estilo dada por qualquer pessoa
que revise código Go no projeto.

Este guia é uma síntese do [Guia de Estilo Go da Uber][uber-style] (traduzido para
pt-BR) e do [Guia de Estilo Go do Google][google-style], adaptado às convenções e à
arquitetura da Pokedex Platform.

[uber-style]: https://github.com/alcir-junior-caju/uber-go-style-guide-pt-br
[google-style]: https://google.github.io/styleguide/go

## Documentos

| Documento | Público-alvo | Normativo | Canônico |
|-----------|-------------|-----------|----------|
| **[Guia](guide.md)** | Todos | Sim | Sim |
| **[Decisões](decisions.md)** | Revisores | Sim | Não |
| **[Melhores Práticas](best-practices.md)** | Interessados | Não | Não |
| **[Regras do Projeto](regras-projeto.md)** | Todos | Sim | Sim |

## Documentos

1. **O [Guia](guide.md)** estabelece os fundamentos do estilo Go no projeto.
   Este documento é definitivo e serve como base para as recomendações nos
   documentos de Decisões e Melhores Práticas.

1. **[Decisões](decisions.md)** é um documento que resume decisões sobre pontos
   específicos de estilo e discute o raciocínio por trás delas quando apropriado.

   Essas decisões podem mudar ocasionalmente com base em novos dados, novas
   funcionalidades da linguagem, novas bibliotecas ou padrões emergentes.

1. **[Melhores Práticas](best-practices.md)** documenta padrões que evoluíram ao
   longo do tempo e que resolvem problemas comuns, têm boa legibilidade e são
   robustos para necessidades de manutenção de código.

   Essas melhores práticas não são canônicas, mas o uso delas é encorajado para
   manter a base de código uniforme e consistente.

1. **[Regras do Projeto](regras-projeto.md)** documenta as regras específicas da
   Pokedex Platform: arquitetura hexagonal, estrutura de diretórios, padrões de
   commit e convenções de CI/CD.

## Objetivos

Estes documentos pretendem:

- Concordar com um conjunto de princípios para ponderar estilos alternativos
- Codificar questões estabelecidas de estilo Go
- Documentar e fornecer exemplos canônicos de expressões idiomáticas Go
- Documentar os prós e contras de várias decisões de estilo
- Ajudar a minimizar surpresas em revisões de código
- Ajudar revisores a usar terminologia e orientação consistentes

Estes documentos **não** pretendem:

- Ser uma lista exaustiva de comentários de revisão de código
- Listar todas as regras que todos devem lembrar e seguir o tempo todo
- Substituir o bom senso no uso de recursos e estilos da linguagem
- Justificar mudanças em larga escala para eliminar diferenças de estilo

## Definições

**Canônico**: Estabelece regras prescritivas e duradouras.

Dentro destes documentos, "canônico" descreve algo que é considerado um padrão
que todo código (antigo e novo) deve seguir e que não se espera que mude
substancialmente ao longo do tempo.

**Normativo**: Destinado a estabelecer consistência.

Dentro destes documentos, "normativo" descreve algo que é um elemento de estilo
acordado para uso por revisores de código Go, para que as sugestões, terminologia
e justificativas sejam consistentes. Esses elementos podem mudar ao longo do tempo.

**Idiomático**: Comum e familiar.

Dentro destes documentos, "idiomático" refere-se a algo que é prevalente em código
Go e se tornou um padrão familiar e fácil de reconhecer.

## Referências externas

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Common Mistakes](https://go.dev/wiki/CommonMistakes)
- [Go Proverbs](https://go-proverbs.github.io/)
- [Uber Go Style Guide (pt-BR)](https://github.com/alcir-junior-caju/uber-go-style-guide-pt-br)
- [Google Go Style Guide](https://google.github.io/styleguide/go)
