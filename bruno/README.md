# 🔷 Bruno - Pokedex API Test Collection

Coleção completa de testes da API Pokedex em [Bruno](https://www.usebruno.com/) (alternativa leve ao Postman).

## 📋 Estrutura

```
bruno/
├── health.yml                  # GET /health
├── auth/
│   ├── signup.yml             # POST /auth/signup - Criar conta
│   ├── login.yml              # POST /auth/login - Autenticar
│   ├── refresh.yml            # POST /auth/refresh - Renovar token
│   └── logout.yml             # POST /auth/logout - Encerrar
└── pokemons/
    ├── list.yml               # GET /pokemons - Listar (paginado)
    ├── search.yml             # GET /pokemons/search - Buscar por nome
    ├── details.yml            # GET /pokemons/{id} - Detalhes de um
    ├── add-favorite.yml       # POST /pokemons/{id}/favorite - Adicionar favorito
    └── remove-favorite.yml    # DELETE /pokemons/{id}/favorite - Remover favorito
```

## 🚀 Como Usar

### 1. Abrir no Bruno Desktop
```bash
# Instale Bruno: https://www.usebruno.com/downloads
bruno .
```
Ou abra Bruno e importe a coleção: `File > Open Collection > ./bruno`

### 2. Configurar Ambiente

1. Clique em **environments** (engrenagem no canto superior)
2. Selecione **"local"** (localhost:8000)
3. As variáveis `{{baseUrl}}`, `{{authToken}}`, `{{userId}}` estarão disponíveis

Para production:
- Selecione **"production"** e atualize `baseUrl`: `https://api.pokedex.com`

### 3. Fluxo de Teste Completo

#### **A. Criar Conta (Signup)**
1. Abra `auth > signup.yml`
2. Edite o body com um email/senha reais
3. Clique em **Send**
4. Copie o valor de `access_token` da resposta
5. Execute `Scripts > Globals` e cole em `authToken`

#### **B. Fazer Login**
1. Abra `auth > login.yml`
2. Edite com suas credenciais
3. Clique em **Send**
4. Bruno extrai automaticamente `authToken` (se configurado o teste)

#### **C. Listar Pokémon**
1. Abra `pokemons > list.yml`
2. Clique em **Send**
3. Veja 20 Pokémon da primeira página

#### **D. Buscar Pokémon**
1. Abra `pokemons > search.yml`
2. Modifique a query: `?q=charizard` (ao invés de pikachu)
3. Clique em **Send**

#### **E. Ver Detalhes**
1. Abra `pokemons > details.yml`
2. Mude o ID na URL (ex: 1 para Bulbasaur, 25 para Pikachu)
3. Clique em **Send**

#### **F. Adicionar Favorito (requer login)**
1. Primeiro, faça login (passo B)
2. Abra `pokemons > add-favorite.yml`
3. Clique em **Send**
4. Resposta: `{"status": "ok"}`

#### **G. Remover Favorito**
1. Abra `pokemons > remove-favorite.yml`
2. Clique em **Send**

#### **H. Renovar Token**
1. Abra `auth > refresh.yml`
2. Clique em **Send**
3. Novo token é retornado e pode ser usado nas próximas requisições

#### **I. Fazer Logout**
1. Abra `auth > logout.yml`
2. Clique em **Send**
3. Sessão encerrada (token invalidado no servidor)

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
- `POST {{baseUrl}}/api/v1/pokemons/{id}/favorite` (header: Authorization)
- `DELETE {{baseUrl}}/api/v1/pokemons/{id}/favorite` (header: Authorization)

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
