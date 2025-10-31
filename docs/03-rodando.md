## 🚀 Desenvolvimento Local

### Pré-requisitos

```bash
# Verificar versões
go version    # 1.21+
docker --version  # 20+ (opcional)
```

### Setup Inicial

```bash
# 1. Configurar .env
cp env.example .env
# Edite .env com suas configurações

# 2. Instalar dependências
go mod tidy

# 3. Rodar migrations
make migrate

# 4. Popul

ar banco (opcional)
make seed
```

### Iniciar Servidor

```bash
# Modo desenvolvimento (hot-reload com Air)
make dev

# Ou build e run
make build
make run

# Ou diretamente
go run cmd/server/main.go
```

Servidor rodando em **http://localhost:80**

### Logs Esperados

```
INFO[2025-01-15T10:30:00Z] Server starting on port 80
INFO[2025-01-15T10:30:00Z] Database connected successfully
INFO[2025-01-15T10:30:00Z] Migrations completed
```

## 🐳 Com Docker Compose

### Subir Tudo

```bash
# Subir banco + API
docker compose up -d

# Ver logs
docker compose logs -f api

# Status
docker compose ps
```

### Migrations e Seed (Docker)

```bash
# Aplicar migrations
docker compose exec api make migrate

# Rodar seed
docker compose exec api make seed

# Ver logs
docker compose logs -f api
```

### Parar

```bash
# Parar containers
docker compose down

# Parar e remover volumes
docker compose down -v
```

## 🧪 Testando Endpoints

### Health Check

```bash
curl http://localhost:80/health

# Resposta esperada:
# {
#   "status": "ok",
#   "message": "API is healthy",
#   "version": "1.0.0",
#   "timestamp": {"time": "now"}
# }
```

### Contatos

**Listar contatos:**
```bash
curl http://localhost:80/api/v1/contacts

# Com paginação
curl "http://localhost:80/api/v1/contacts?page=1&limit=5"
```

**Criar contato:**
```bash
curl -X POST http://localhost:80/api/v1/contacts \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@example.com",
    "phone": "+55 11 99999-9999",
    "company": "Tech Corp"
  }'

# Resposta (201 Created):
# {
#   "id": 1,
#   "name": "João Silva",
#   "email": "joao@example.com",
#   "phone": "+55 11 99999-9999",
#   "company": "Tech Corp",
#   "created_at": "2025-01-15T10:30:00Z",
#   "updated_at": "2025-01-15T10:30:00Z"
# }
```

**Buscar por ID:**
```bash
curl http://localhost:80/api/v1/contacts/1
```

**Atualizar contato:**
```bash
curl -X PUT http://localhost:80/api/v1/contacts/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva Atualizado",
    "company": "Nova Empresa"
  }'
```

**Deletar contato:**
```bash
curl -X DELETE http://localhost:80/api/v1/contacts/1

# Resposta: 204 No Content
```

**Buscar contatos:**
```bash
# Busca por nome, email ou empresa
curl "http://localhost:80/api/v1/contacts/search?q=João"

# Com paginação
curl "http://localhost:80/api/v1/contacts/search?q=Tech&page=1&limit=5"
```

## 🧪 Testes Automatizados

### Rodar Todos os Testes

```bash
# Testes básicos
make test

# Com coverage
make test-coverage

# Verbose
go test -v ./...

# Testes específicos
go test ./internal/handlers/...
```

### Testes Disponíveis

- ✅ Health check endpoint
- ✅ CRUD de contatos
- ✅ Validação de dados
- ✅ Busca e paginação
- ✅ Soft delete
- ✅ Error handling

### Output Esperado

```
=== RUN   TestHealthCheck
--- PASS: TestHealthCheck (0.00s)
=== RUN   TestCreateContact
--- PASS: TestCreateContact (0.01s)
=== RUN   TestGetContacts
--- PASS: TestGetContacts (0.01s)
=== RUN   TestGetContact
--- PASS: TestGetContact (0.01s)
=== RUN   TestUpdateContact
--- PASS: TestUpdateContact (0.01s)
=== RUN   TestDeleteContact
--- PASS: TestDeleteContact (0.01s)
=== RUN   TestSearchContacts
--- PASS: TestSearchContacts (0.01s)
PASS
ok      api-contacts-go/tests    0.087s
```

## 🔍 Debug e Troubleshooting

### Ver Queries SQL

Editar `.env`:

```env
ENVIRONMENT=development
```

Ou via código (`internal/database/database.go`):

```go
gormLogger = logger.Default.LogMode(logger.Info) // Mostra queries SQL
```

### Verificar Conexão do Banco

```bash
# Testar conexão
go run -c "
package main
import (
    \"api-contacts-go/internal/database\"
    \"fmt\"
)
func main() {
    db, err := database.Initialize(\"$DATABASE\")
    if err != nil {
        fmt.Println(\"Error:\", err)
    } else {
        fmt.Println(\"Database connected successfully\")
    }
}
"
```

### Verificar Migrations

```bash
# Status atual
migrate -path migrations -database "$DATABASE" version

# Histórico
migrate -path migrations -database "$DATABASE" history
```

**⚠️ Migrations "Dirty" (Auto-fix):**

Se você ver este warning nos logs:
```
Failed to run migrations: Dirty database version X
```

**Não se preocupe!** O sistema detecta e corrige automaticamente:
1. ✅ Detecta estado "dirty" da migration
2. ✅ Força limpeza da versão atual
3. ✅ Continua com a inicialização normalmente

Isso pode acontecer se:
- Uma migration anterior falhou no meio
- O container foi reiniciado durante uma migration
- Houve erro de conexão durante migration

O servidor **continua funcionando normalmente** após o auto-fix.

### Logs Detalhados

```bash
# Desenvolvimento
ENVIRONMENT=development go run cmd/server/main.go

# Docker
docker compose logs -f api
```

### Verificar Porta em Uso

```bash
# Linux/Mac
lsof -i :80

# Windows
netstat -ano | findstr :80
```

## 📈 Performance Testing

### Simples (cURL)

```bash
# Medir latência
time curl http://localhost:80/api/v1/contacts
```

### com Apache Bench

```bash
# 1000 requests, 10 concurrent
ab -n 1000 -c 10 http://localhost:80/api/v1/contacts
```

### com Artillery

```bash
npm install -g artillery

# Criar config.yml
artillery quick --count 100 --num 10 http://localhost:80/health

# Resultados esperados:
# - p95 latency: < 10ms
# - Requests/sec: > 2000
```

## 🔄 Hot Reload

Para desenvolvimento com hot-reload, instale Air:

```bash
go install github.com/cosmtrek/air@latest
```

Crie `.air.toml`:

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

Execute:

```bash
air
```

## 🎯 Checklist de Validação

Antes de considerar pronto:

- [ ] `go mod tidy` sem erros
- [ ] `make migrate` aplica migrations
- [ ] `make seed` popula dados
- [ ] `make dev` inicia servidor
- [ ] `curl /health` retorna status ok
- [ ] `curl /api/v1/contacts` retorna lista
- [ ] `make test` passa todos os testes
- [ ] Docker Compose sobe corretamente
- [ ] Logs estruturados aparecem
- [ ] Soft delete funciona
- [ ] Busca funciona

## 🚀 Próximos passos

Continue para [Deploy](./04-deploy.md) para colocar em produção.
