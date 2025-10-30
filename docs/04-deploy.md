# Deploy na Shard Cloud

## 🚀 Deploy na Shard Cloud

A Shard Cloud oferece hospedagem moderna e confiável para seus projetos Go. Siga este guia para fazer deploy da sua API em minutos.

### 📋 Pré-requisitos

- Conta na [Shard Cloud](https://shardcloud.app)
- Projeto compilado e funcionando localmente
- Arquivo `.shardcloud` configurado
- Banco PostgreSQL (pode usar o da Shard Cloud)

## 🔧 Configuração do projeto

### 1. Criar arquivo `.shardcloud`

Crie um arquivo `.shardcloud` na raiz do projeto:

```bash
DISPLAY_NAME=Contacts API
ENTRYPOINT=bin/main
MEMORY=512
VERSION=recommended
SUBDOMAIN=contacts-api
START=go mod tidy && go build -o bin/main ./cmd/server && make migrate && make seed && ./bin/main
DESCRIPTION=API REST para gerenciamento de contatos com Fiber e PostgreSQL
```

### 2. Configurar variáveis de ambiente

Configure as variáveis no dashboard da Shard Cloud:

```env
# Database - REQUIRED
DATABASE=postgres://user:password@host:port/database?ssl=true

# Server
PORT=80
ENVIRONMENT=production

# Logging
LOG_LEVEL=info
```

## 📦 Preparação para deploy

### 1. Testar build localmente

```bash
# Instalar dependências
go mod tidy

# Compilar projeto
go build -o bin/main ./cmd/server

# Executar migrações
make migrate

# Popular banco com dados de exemplo
make seed

# Testar aplicação
./bin/main
```

### 2. Verificar funcionamento

```bash
# Testar health endpoint
curl http://localhost/health

# Testar API
curl http://localhost/api/v1/contacts
```

## 🚀 Deploy na Shard Cloud

### Método 1: Upload direto (Recomendado)

1. **Acesse o Dashboard**
   - Vá para [Shard Cloud Dashboard](https://shardcloud.app/dash)
   - Faça login na sua conta

2. **Criar nova aplicação**
   - Clique em **"New app"**
   - Selecione **"Upload"**

3. **Preparar arquivos**
   - Zip toda a pasta do projeto (incluindo `.shardcloud`)
   - Certifique-se de que o `go.mod` está incluído

4. **Upload e deploy**
   - Arraste o arquivo ZIP ou clique para selecionar
   - Aguarde o processamento (alguns minutos)
   - Sua aplicação estará disponível em `https://contacts-api.shardweb.app`

### Método 2: Deploy via Git

1. **Conectar repositório**
   - No dashboard, clique em **"New app"**
   - Selecione **"Git Repository"**
   - Conecte seu repositório GitHub/GitLab

2. **Configurar build**
   - **Build command:** `go mod tidy && go build -o bin/main ./cmd/server`
   - **Start command:** `./bin/main`
   - **Go version:** `1.23` (recomendado)

3. **Deploy automático**
   - Cada push na branch principal fará deploy automático
   - Configure webhooks se necessário

## 🗄️ Banco de dados

### Usar PostgreSQL da Shard Cloud

1. **Criar banco**
   - Vá para [Databases Dashboard](https://shardcloud.app/dash/databases)
   - Clique em **"New Database"**
   - Selecione **PostgreSQL**
   - Escolha a quantidade de RAM

2. **Configurar conexão**
   - Copie a string de conexão do dashboard
   - Configure como variável `DATABASE` na aplicação
   - Exemplo: `postgres://user:pass@host:port/db?ssl=true`

3. **Executar migrações**
   - As migrações são executadas automaticamente na inicialização
   - Verifique logs para confirmar sucesso

### Banco externo

Se preferir usar banco externo:

```env
DATABASE=postgres://user:password@external-host:5432/database?ssl=true
```

## 🌐 Configurações avançadas

### Subdomínio personalizado

No arquivo `.shardcloud`:

```bash
SUBDOMAIN=minha-api
```

Sua aplicação ficará disponível em: `https://minha-api.shardweb.app`

### Domínio personalizado

1. **Configurar DNS**
   - Adicione um registro CNAME apontando para `contacts-api.shardweb.app`
   - Ou configure A record com o IP fornecido

2. **Ativar no dashboard**
   - Vá para configurações da aplicação
   - Adicione seu domínio personalizado
   - Configure certificado SSL (automático)

### Variáveis de ambiente

Configure variáveis sensíveis no dashboard:

1. Acesse configurações da aplicação
2. Vá para **"Environment Variables"**
3. Adicione suas variáveis:
   ```
   DATABASE=postgres://user:pass@host:port/db?ssl=true
   PORT=80
   ENVIRONMENT=production
   LOG_LEVEL=info
   ```

## 🔍 Monitoramento e logs

### Logs da aplicação

- Acesse o dashboard da aplicação
- Vá para a aba **"Logs"**
- Monitore erros e performance em tempo real

### Métricas

- **Uptime:** Monitoramento automático
- **Performance:** Métricas de resposta
- **Tráfego:** Estatísticas de acesso
- **Database:** Monitoramento de conexões

### Health checks

A aplicação inclui endpoints de monitoramento:

- `GET /health` - Status geral da API
- `GET /api/v1/contacts` - Lista de contatos
- `POST /api/v1/contacts` - Criar novo contato

## 🔒 Segurança

### HTTPS automático

- Todos os deploys na Shard Cloud incluem HTTPS automático
- Certificados SSL gerenciados automaticamente
- Renovação automática

### Validação de dados

A aplicação usa validação robusta:

- Validação automática de entrada
- Middleware de validação
- Mensagens de erro padronizadas

## 🚦 CI/CD com GitHub Actions

Crie `.github/workflows/deploy.yml`:

```yaml
name: Deploy to Shard Cloud

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      
      - name: Cache Go modules
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
      
      - name: Build
        run: go build -o bin/main ./cmd/server
      
      - name: Deploy to Shard Cloud
        run: |
          # Zip project
          zip -r deploy.zip . -x "bin/*" "*.git*"
          
          # Upload to Shard Cloud (configure API token)
          curl -X POST \
            -H "Authorization: Bearer ${{ secrets.SHARD_TOKEN }}" \
            -F "file=@deploy.zip" \
            https://api.shardcloud.app/deploy
```

## 🐛 Troubleshooting

### Build falha

```bash
# Limpar cache Go
go clean -cache

# Verificar dependências
go mod verify

# Compilar com debug
go build -v -o bin/main ./cmd/server
```

### Aplicação não inicia

1. Verifique logs no dashboard
2. Confirme se `ENTRYPOINT` está correto
3. Teste localmente com `./bin/main`

### Erro de conexão com banco

1. Verifique string de conexão `DATABASE`
2. Confirme se banco está acessível
3. Teste conexão localmente

### Erro de migração

```bash
# Executar migrações manualmente
make migrate

# Verificar status das migrações
make migrate-status
```

## ✅ Checklist de deploy

- [ ] Arquivo `.shardcloud` configurado
- [ ] Projeto compila sem erros (`go build`)
- [ ] Testado localmente (`./bin/main`)
- [ ] Banco PostgreSQL configurado
- [ ] Variáveis de ambiente configuradas
- [ ] Projeto zipado ou conectado ao Git
- [ ] Deploy realizado no dashboard
- [ ] Aplicação acessível via URL
- [ ] Health endpoint funcionando (`/health`)
- [ ] API endpoints testados (`/api/v1/contacts`)
- [ ] HTTPS ativo
- [ ] Logs monitorados

## 🎉 Sucesso!

Sua API está no ar na Shard Cloud! 

### Próximos passos:

1. **Teste completo:** Verifique todos os endpoints
2. **Documentação:** Configure Swagger se necessário
3. **Monitoramento:** Configure alertas de uptime
4. **Backup:** Configure backup do banco de dados
5. **Otimização:** Monitore métricas e otimize performance

### URLs importantes:

- **Dashboard:** https://shardcloud.app/dash
- **Documentação:** https://docs.shardcloud.app/quickstart
- **Suporte:** https://shardcloud.app/support

---

**Precisa de ajuda?** Consulte a [documentação oficial da Shard Cloud](https://docs.shardcloud.app/quickstart) ou entre em contato com o suporte.