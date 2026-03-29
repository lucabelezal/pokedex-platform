# 🔷 Bruno - Pokedex API Test Collection

Coleção completa de testes da API Pokedex em [Bruno](https://www.usebruno.com/) (alternativa leve ao Postman).

## UI/UX de Referência

As próximas APIs orientadas à experiência do app seguem como referência visual:

- [Pokémon App By Junior Saraiva](https://www.figma.com/pt-br/comunidade/file/1202971127473077147/pokedex-pokemon-app)

Tabs principais consideradas na coleção:

- `pokedex` ou home
- detalhe do Pokémon
- regiões
- favoritos
- perfil

---

## 🎯 O que é Bruno?

**Bruno** é um cliente HTTP **open-source** e **leve** para testar APIs, similar ao Postman, mas com algumas vantagens:

| Aspecto | Bruno | Postman |
|--------|-------|---------|
| **Instalação** | Standalone desktop | Cloud + desktop |
| **Versionamento** | YAML simples (git-friendly) | Binário (difícil de versionar) |
| **Tamanho** | ~100MB | ~500MB |
| **Offline** | ✅ Completo suporte | ⚠️ Limitado |
| **Preço** | Gratuito (MIT) | Freemium |

**Vantagens para este projeto:**
- ✅ Colecção versionável (`bruno/` em Git)
- ✅ Equipes podem colaborar editando YAML
- ✅ Integração com CI/CD (`bru run`)
- ✅ Leve e rápido
- ✅ Sem account necessária

---

## 📥 Como Baixar e Instalar Bruno

### Opção 1: Download Oficial (Recomendado)

1. Acesse https://www.usebruno.com/downloads
2. Selecione seu sistema operacional:
   - **macOS**: intel ou apple silicon
   - **Windows**: exe installer
   - **Linux**: AppImage ou deb/rpm
3. Instale normalmente (drag-and-drop no macOS, duplo-clique no Windows, etc)

### Opção 2: Package Manager

**macOS (Homebrew):**
```bash
brew install bruno
```

**Windows (Chocolatey):**
```bash
choco install bruno
```

**Linux (apt):**
```bash
sudo apt install bruno
```

### Opção 3: Docker (se preferir)
```bash
docker run -it -v ~/Documents/bruno:/home/user/.config/bruno \
  -p 6969:6969 \
  usebruno/bruno
```

---

## 📂 Como Importar Esta Colecção

### Método 1: Abrir com `bruno` Command (Terminal)

```bash
cd /Users/lucasnascimento/Dev/GitHub/BFF/pokedex-platform
bruno .
```

Isso abre a colecção `./bruno/` automaticamente no Bruno Desktop.

### Método 2: Importar Manualmente (UI)

1. **Abra o Bruno Desktop**
2. Clique em **File** (canto superior esquerdo)
3. Selecione **"Open Collection"** ou **"Import"**
4. Navegue até `/Users/lucasnascimento/Dev/GitHub/BFF/pokedex-platform/bruno`
5. Clique em **Open** ou **Import**

### Método 3: Git Clone (Primeira Vez)

Se ainda não tem o repositório:

```bash
git clone https://github.com/[seu-user]/BFF/pokedex-platform.git
cd pokedex-platform
bruno .
```

---

## 📋 Estrutura

```
bruno/
├── opencollection.yml         # Root da coleção OpenCollection
├── health.yml                  # GET /health
├── environments/
│   ├── local.yml              # Variáveis de ambiente local
│   └── production.yml         # Variáveis de ambiente production
├── auth/
│   ├── signup.yml             # POST /auth/signup - Criar conta
│   ├── login.yml              # POST /auth/login - Autenticar
│   ├── refresh.yml            # POST /auth/refresh - Renovar token
│   └── logout.yml             # POST /auth/logout - Encerrar
├── app/
│   ├── home.yml               # GET /home - Home da pokedex
│   ├── pokemon-details.yml    # GET /pokemons/{id}/details - Tela de detalhe
│   ├── regions.yml            # GET /regions - Lista de regiões
│   ├── favorites.yml          # GET /me/favorites - Aba de favoritos
│   └── profile.yml            # GET /me - Aba de perfil
└── pokemons/
    ├── list.yml               # GET /pokemons - Listar (paginado)
    ├── search.yml             # GET /pokemons/search - Buscar por nome
    ├── details.yml            # GET /pokemons/{id} - Detalhes de um
    ├── add-favorite.yml       # POST /pokemons/{id}/favorite - Adicionar favorito
    └── remove-favorite.yml    # DELETE /pokemons/{id}/favorite - Remover favorito
```

## 🚀 Começando Rápido

### 1️⃣ Instale Bruno
https://www.usebruno.com/downloads (2 minutos)

### 2️⃣ Abra esta Colecção
```bash
cd /Users/lucasnascimento/Dev/GitHub/BFF/pokedex-platform
bruno .
```

### 3️⃣ Selecione o Ambiente
Clique em **"local"** (engrenagem no canto superior direito) para usar `http://localhost:8000`

### 4️⃣ Teste um Endpoint
Abra `health.yml` e clique em **Send** (ou Cmd+Enter)

---

## 🔧 Como Usar em Detalhes

### 1. Configurar Ambiente

**Na Interface do Bruno:**
1. Clique no ícone de **engrenagem** (⚙️) no canto superior direito
2. Selecione **"local"** (desenvolvimento local) ou **"production"** (servidor remoto)

**Variáveis dos arquivos de ambiente (`environments/local.yml` e `environments/production.yml`):**
```yaml
name: local
variables:
  - name: baseUrl
    value: http://localhost:8000
  - name: authToken
    value: cole-seu-token-aqui  # Substitua após login
```

Os arquivos de ambiente já vêm com placeholders para facilitar visualização no Bruno.

Todas as requisições usam `{{baseUrl}}` automaticamente, então não precisa mudar URLs manualmente.

### 2. Fluxo de Teste Completo

#### **A. Health Check (teste rápido)** ✅
1. Abra `health.yml`
2. Clique em **Send** (Cmd+Enter no macOS, Ctrl+Enter no Windows/Linux)
3. Espere resposta: `{"status":"ok"}`

#### **B. Criar Conta (Signup)**
1. Abra `auth > signup.yml`
2. Edite o body com um **email e senha reais** (use um email de teste)
3. Clique em **Send**
4. Resposta esperada: `{"access_token": "eyJ...", "user_id": "uuid-...", ...}` (status 201)
5. **Copie o `access_token`** da resposta

#### **C. Salvar Token para Próximas Requisições**
1. Abra **Scripts > Globals** (ícone de chave inglesa)
2. Cole no campo `authToken`:
   ```
   authToken = eyJhbGciOiJIUzI1NiIs...
   ```
3. Feche a janela (Cmd+W ou Esc)

#### **D. Fazer Login**
1. Abra `auth > login.yml`
2. Use as mesmas credenciais do signup
3. Clique em **Send**
4. Resposta: novo token + dados do usuário (status 200)

#### **E. Listar Pokémon (Public)**
1. Abra `pokemons > list.yml`
2. Clique em **Send**
3. Veja os primeiros 20 Pokémon (não requer autenticação)

#### **E.1. Home da Pokedex**
1. Abra `app > home.yml`
2. Clique em **Send**
3. A resposta já vem moldada para a tela da Pokédex, com busca, filtros e a coleção curada de cards do app

#### **F. Buscar Pokémon**
1. Abra `pokemons > search.yml`
2. Modifique a query na URL: `?q=charizard` (ao invés de pikachu)
3. Clique em **Send**
4. Veja resultados filtrados

#### **G. Ver Detalhes de um Pokémon**
1. Abra `pokemons > details.yml`
2. Mude o ID na URL (ex: `/25` para Pikachu, `/1` para Bulbasaur)
3. Clique em **Send**
4. Veja estrutura completa com stats, tipos, etc

#### **G.1. Tela de Detalhe do App**
1. Abra `app > pokemon-details.yml`
2. Ajuste o ID conforme necessário
3. Valide o payload orientado à UI (about, weaknesses, evolutions e isFavorite)

#### **G.2. Regiões**
1. Abra `app > regions.yml`
2. Clique em **Send**
3. Valide o contrato da aba de regiões (`title` + `regions[]`)

#### **G.3. Favoritos**
1. Abra `app > favorites.yml`
2. Clique em **Send**
3. Valide os estados de tela: `unauthenticated`, `empty` e `has_data`

#### **G.4. Perfil**
1. Abra `app > profile.yml`
2. Clique em **Send**
3. Valide os estados da tela de conta (`authenticated: false/true`)

#### **H. Adicionar Favorito (Autenticado)** 🔒
1. Certifique-se que `authToken` está configurado (passo C)
2. Abra `pokemons > add-favorite.yml`
3. Modifique o ID se quiser (ex: 6 para Charizard)
4. Clique em **Send**
5. Resposta: `{"status": "ok"}` (status 200)

#### **I. Remover Favorito (Autenticado)** 🔒
1. Abra `pokemons > remove-favorite.yml`
2. Clique em **Send**
3. Resposta: `{"message": "favorito removido"}` (status 200)

#### **J. Renovar Token (Refresh)**
1. Abra `auth > refresh.yml`
2. Clique em **Send**
3. Novo token é gerado (útil quando token expirar)
4. Copie o novo `access_token` e atualize em Globals

#### **K. Fazer Logout**
1. Abra `auth > logout.yml`
2. Clique em **Send**
3. Resposta: `{"message": "sessao encerrada"}` (status 200)
4. Token é invalidado no servidor

## 🔐 Autenticação

### Opção 1: Header Authorization
```
Authorization: Bearer {access_token}
```

### Opção 2: Cookie HTTP-only
Configurado automaticamente após login/signup

**Prioridade:** Header `Authorization` sobrescreve o cookie

## 📝 Variáveis Disponíveis

| Variável | Exemplo | Uso |
|----------|---------|-----|
| `{{baseUrl}}` | `http://localhost:8000` | URL base da API |
| `{{authToken}}` | `eyJhbGc...` | JWT extraído após login |
| `{{userId}}` | `550e8400-e29b...` | ID do usuário logado |
| `{{userEmail}}` | `user@example.com` | Email do usuário |

## 🧪 Testes Automatizados (opcional)

Cada requisição pode ter um `Tests` script que valida respostas:

**Exemplo para `signup.yml`:**
```javascript
pm.test("Status é 201", function() {
  pm.response.to.have.status(201);
});

pm.test("Response tem access_token", function() {
  var jsonData = pm.response.json();
  pm.expect(jsonData).to.have.property("access_token");
  pm.globals.set("authToken", jsonData.access_token);
});
```

## 🔗 Endpoint Reference

### Public (sem autenticação)
- `GET {{baseUrl}}/v1/health`
- `POST {{baseUrl}}/v1/auth/signup`
- `POST {{baseUrl}}/v1/auth/login`
- `GET {{baseUrl}}/v1/pokemons?page=0&size=20`
- `GET {{baseUrl}}/v1/pokemons/search?q=termo`
- `GET {{baseUrl}}/v1/pokemons/{id}`

### Protected (requerem Bearer token)
- `POST {{baseUrl}}/v1/auth/refresh` (header: Authorization)
- `POST {{baseUrl}}/v1/auth/logout` (header: Authorization)
- `POST {{baseUrl}}/v1/pokemons/{id}/favorite` (header: Authorization)
- `DELETE {{baseUrl}}/v1/pokemons/{id}/favorite` (header: Authorization)

## 💡 Dicas

1. **Local Development:** Use ambiente `local` com `baseUrl: http://localhost:8000`
2. **Staging/Prod:** Crie novos ambientes conforme necessário
3. **Token Expirado:** Abra `auth > refresh.yml` para renovar
4. **Deletar dados:** Limpe variáveis em `Scripts > Globals` ou `auth > logout.yml`
5. **CI/CD:** Use `bru run` para executar testes automaticamente

```bash
# Executar toda coleção
bru run .

# Executar pasta específica
bru run . --folder auth

# Com variáveis
bru run . --env local
```

## 📚 Referências

- [Bruno Docs](https://docs.usebruno.com/)
- [API Pokédex](../README.md)
- [Mobile BFF Docs](../bff/mobile-bff/README.md)

## ✅ Checklist de Teste

- [ ] Health check retorna 200
- [ ] Signup cria conta com sucesso (201)
- [ ] Login autentica e retorna token (200)
- [ ] Listar Pokémon retorna 20 primeiros (200)
- [ ] Busca por nome funciona (200)
- [ ] Detalhe de Pokémon retorna estrutura completa (200)
- [ ] Adicionar favorito funciona autenticado (200)
- [ ] Remover favorito funciona autenticado (200)
- [ ] Refresh renova token (200)
- [ ] Logout encerra sessão (200)
