# API de Contatos (Contacts API)

API REST completa para gerenciamento de contatos com Fiber, GORM e PostgreSQL. CRUD completo, validação, migrations e health check.

## 🎯 Características

- ✅ CRUD completo de contatos
- ✅ Validação de dados robusta
- ✅ Migrations com golang-migrate
- ✅ PostgreSQL com GORM
- ✅ Fiber (framework rápido)
- ✅ Docker e Docker Compose
- ✅ Health check endpoint
- ✅ CORS configurado
- ✅ Logs estruturados
- ✅ Testes automatizados

## 📋 Requisitos

- Go 1.21+
- PostgreSQL 14+ (ou use Docker Compose)
- Docker (opcional)

## 🚀 Como rodar

### Com Docker Compose (Recomendado)

```bash
# Copiar .env (para Docker Compose)
cp env.docker.example .env

# Subir banco e aplicação
docker compose up -d

# Rodar migrations
docker compose exec api make migrate

# Rodar seed
docker compose exec api make seed

# Acesse: http://localhost:80/health
```

### Sem Docker

```bash
# Configurar DATABASE no .env
cp .env.example .env

# Instalar dependências
go mod tidy

# Rodar migrations
make migrate

# Rodar seed
make seed

# Iniciar servidor
make dev
```

## 📦 Scripts

```bash
make dev          # Servidor de desenvolvimento
make build        # Build da aplicação
make run          # Executar binário
make migrate      # Aplicar migrations
make migrate-new  # Criar nova migration
make seed         # Popular banco
make test         # Executar testes
make lint         # Linter (golangci-lint)
make fmt          # Formatar código (gofmt)
```

**💡 Dica:** Migrations são executadas automaticamente na inicialização. Se houver warnings sobre estado "dirty", o sistema corrige automaticamente.

## 🔗 Endpoints

### Health Check
```
GET /health
```

### Contatos

```
GET    /contacts           # Listar todos (com paginação)
GET    /contacts/:id       # Buscar por ID
POST   /contacts           # Criar novo
PUT    /contacts/:id       # Atualizar
DELETE /contacts/:id       # Deletar
GET    /contacts/search    # Buscar por nome/email
```

### Exemplos de uso

**Criar contato:**
```bash
curl -X POST http://localhost:80/contacts \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@example.com",
    "phone": "+55 11 99999-9999"
  }'
```

**Listar com paginação:**
```bash
# Página 1, 10 itens
GET /contacts?page=1&limit=10

# Busca por nome
GET /contacts/search?q=João
```

## 📂 Estrutura

```
api-contacts-go/
├── cmd/
│   └── server/
│       └── main.go        # Entry point
├── internal/
│   ├── config/            # Configurações
│   ├── database/          # Conexão DB
│   ├── models/            # GORM models
│   ├── handlers/          # HTTP handlers
│   ├── services/          # Lógica de negócio
│   └── middleware/        # Middlewares
├── migrations/            # SQL migrations
├── seed/                  # Dados iniciais
├── tests/                 # Testes
└── docs/                  # Documentação
```

## 🗄️ Banco de dados

### Modelo de Contact

```go
type Contact struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Name      string    `json:"name" gorm:"not null"`
    Email     string    `json:"email" gorm:"uniqueIndex;not null"`
    Phone     string    `json:"phone"`
    Company   string    `json:"company"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Conexão

A aplicação lê a string de conexão da variável `DATABASE`:

```env
# Local (sem SSL)
DATABASE=postgres://user:password@localhost:5432/contacts_db

# Produção (com SSL)
DATABASE=postgres://user:password@postgres.example.com:5432/db?sslmode=require
```

**Nota sobre SSL:** O código automaticamente converte `ssl=true` para `sslmode=require` (formato correto do PostgreSQL). Para bancos remotos sem `sslmode` especificado, `sslmode=require` é adicionado automaticamente.

## 🐳 Docker

### Build

```bash
docker build -t api-contacts-go .
```

### Run

```bash
docker run -p 80:80 \
  -e DATABASE=postgres://user:pass@host:5432/db \
  api-contacts-go
```

## 📊 Logs

A aplicação usa logs estruturados (JSON em produção):

```json
{
  "level": "info",
  "msg": "Server started on :80",
  "time": "2025-01-15T10:30:00Z"
}
```

## 🧪 Testes

```bash
# Rodar todos os testes
make test

# Com coverage
go test -cover ./...

# Testes específicos
go test ./internal/handlers/...
```

## 📄 Licença

MIT

---