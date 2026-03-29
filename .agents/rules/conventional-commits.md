---
trigger: always_on
---

# Rule: Antigravity Standard (Conventional Commits & Idioma PT-BR)

## Contexto
Esta regra deve ser aplicada em todas as interações neste workspace, especialmente ao gerar mensagens de commit, documentação ou explicações de código.

## Ativação
- **Modo:** Always On (Sempre Ativo)
- **Descrição:** Garante que a comunicação e as mensagens de commit sigam o padrão definido em português.

## 🌐 Regra de Idioma
- **Comunicação Geral:** Todas as respostas, pensamentos e explicações do Agente devem ser obrigatoriamente em **Português (Brasil)**.
- **Mensagens de Commit:** A descrição do commit deve ser escrita em **Português**.

## 🚀 Padrão Conventional Commits (PT-BR)
As mensagens de commit devem seguir estritamente a estrutura:
`<tipo>[escopo opcional]: <descrição em português>`

### 1. Tipos Permitidos (Prefixos)
Os prefixos permanecem em inglês para compatibilidade com ferramentas de automação, mas a descrição deve ser em português:
- **feat**: Nova funcionalidade.
- **fix**: Correção de bug.
- **docs**: Alterações na documentação.
- **style**: Formatação e estilo (sem alteração de lógica).
- **refactor**: Refatoração de código.
- **perf**: Melhoria de performance.
- **test**: Adição ou modificação de testes.
- **build**: Alterações no sistema de build ou dependências.
- **ci**: Alterações em scripts de integração contínua.
- **chore**: Tarefas de manutenção geral.

### 2. Regras de Formatação
- **Descrição:** Deve começar com letra minúscula.
- **Pontuação:** Não utilizar ponto final ao final da frase.
- **Verbo:** Preferencialmente utilizar o infinitivo ou presente (ex: "adicionar", "corrige").
- **Comprimento:** O cabeçalho não deve exceder 72 caracteres.

## Exemplos de Uso
- ✅ `feat(api): adicionar autenticação de usuário`
- ✅ `fix(ui): corrigir alinhamento do botão de gravidade`
- ✅ `docs: atualizar o manual de instalação`
- ❌ `feat(api): add user auth` (Descrição em inglês proibida)
- ❌ `Adicionando nova rota` (Falta o tipo/prefixo)

## Comportamento do Agente
Ao ser solicitado um commit ("commita isso"), o Agente deve analisar as mudanças e sugerir o comando no formato:
`git commit -m "tipo(escopo): descrição em português"`