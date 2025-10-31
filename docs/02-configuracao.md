## 🔐 Variáveis de Ambiente

### Arquivo `.env`

Copie `env.example` para `.env` e configure:

```env
# Database (obrigatório)
DATABASE=postgresql://user:password@localhost:5432/contacts_db

# Server (opcional)
PORT=80

# Environment
ENVIRONMENT=development
```

### Variáveis Detalhadas

#### `DATABASE` (obrigatório)

String de conexão PostgreSQL:

```env
DATABASE=postgresql://USER:PASSWORD@HOST:PORT/DATABASE
```

Exemplos:

```env
# Local
DATABASE=postgresql://contactuser:contactpass@localhost:5432/contacts_db

# Docker Compose
DATABASE=postgresql://contactuser:contactpass@db:5432/contacts_db

# Supabase
DATABASE=postgresql://user:pass@db.xxx.supabase.co:5432/postgres

# Railway
DATABASE=postgresql://user:pass@containers-us-west-1.railway.app:5432/railway
```

#### `PORT` (opcional, padrão: 80)

Porta onde o servidor escutará:

```env
PORT=80      # Produção (padrão)
PORT=3000    # Desenvolvimento alternativo
```

#### `ENVIRONMENT` (opcional, padrão: development)

Ambiente de execução:

```env
ENVIRONMENT=development   # Logs verbosos, debug
ENVIRONMENT=production    # Logs JSON, otimizações
```

## 🗄️ Configuração do Banco de Dados

### Opção 1: Docker Compose (Recomendado)

```bash
docker compose up -d db
```

Credenciais padrão:
- **User:** contactuser
- **Password:** contactpass
- **Database:** contacts_db
- **Port:** 5432

### Opção 2: PostgreSQL Local

```bash
# Criar usuário e banco
psql -U postgres
CREATE USER contactuser WITH PASSWORD 'contactpass';
CREATE DATABASE contacts_db OWNER contactuser;
GRANT ALL PRIVILEGES ON DATABASE contacts_db TO contactuser;
```

### Opção 3: PostgreSQL em Cloud

**Supabase (Grátis):**
1. Crie projeto em https://supabase.com
2. Vá em Settings > Database
3. Copie Connection String
4. Cole no `.env`

**Railway:**
1. Crie projeto em https://railway.app
2. Adicione PostgreSQL plugin
3. Copie `DATABASE_URL`

## 🔄 Migrations

### Configurar golang-migrate

```bash
# Instalar migrate CLI
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Criar migration
migrate create -ext sql -dir migrations initial

# Aplicar migrations
migrate -path migrations -database "$DATABASE" up

# Ver status
migrate -path migrations -database "$DATABASE" version
```

### Criar Nova Migration

```bash
# 1. Criar arquivos de migration
migrate create -ext sql -dir migrations add_phone_field

# 2. Editar arquivos gerados:
# - migrations/XXXXXX_add_phone_field.up.sql
# - migrations/XXXXXX_add_phone_field.down.sql

# 3. Aplicar
migrate -path migrations -database "$DATABASE" up
```

### Rollback

```bash
# Voltar uma migration
migrate -path migrations -database "$DATABASE" down 1

# Voltar para versão específica
migrate -path migrations -database "$DATABASE" goto 1
```

## 🌱 Seeds

### Rodar Seeds

```bash
make seed
# ou
go run seed/seed.go
```

Isso criará 10 contatos de exemplo.

### Customizar Seeds

Edite `seed/seed.go`:

```go
contacts := []models.Contact{
    {
        Name:    "Seu Nome",
        Email:   "seu@email.com",
        Phone:   "+55 11 99999-9999",
        Company: "Sua Empresa",
    },
    // ...
}
```

## 🐳 Docker

### Build Customizado

```bash
# Build da imagem
docker build -t api-contacts-go .

# Run com variáveis
docker run -p 80:80 \
  -e DATABASE=postgresql://user:pass@host:5432/db \
  api-contacts-go
```

### Docker Compose Personalizado

```yaml
version: '3.8'
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${DB_USER:-contactuser}
      POSTGRES_PASSWORD: ${DB_PASSWORD:-contactpass}
      POSTGRES_DB: ${DB_NAME:-contacts_db}
    volumes:
      - postgres_data:/var/lib/postgresql/data

  api:
    build: .
    environment:
      DATABASE: postgresql://${DB_USER:-contactuser}:${DB_PASSWORD:-contactpass}@db:5432/${DB_NAME:-contacts_db}
    depends_on:
      - db
```

## 🔧 Configuração Avançada

### Logs Estruturados

Editar `cmd/server/main.go`:

```go
import "github.com/sirupsen/logrus"

// Configurar nível de log
if cfg.Environment == "development" {
    logrus.SetLevel(logrus.DebugLevel)
    logrus.SetFormatter(&logrus.TextFormatter{})
} else {
    logrus.SetLevel(logrus.InfoLevel)
    logrus.SetFormatter(&logrus.JSONFormatter{})
}
```

### Connection Pool

Editar `internal/database/database.go`:

```go
db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
    Logger: gormLogger,
    // Configurações de pool
    PrepareStmt: true,
    DisableForeignKeyConstraintWhenMigrating: true,
})

// Configurar pool SQL
sqlDB, err := db.DB()
if err != nil {
    return nil, err
}

sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
sqlDB.SetConnMaxLifetime(time.Hour)
```

### CORS

Editar `cmd/server/main.go`:

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://meusite.com,https://app.meusite.com",
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders: "Origin,Content-Type,Accept,Authorization",
    AllowCredentials: true,
}))
```

### Rate Limiting (Opcional)

```bash
go get github.com/gofiber/fiber/v2/middleware/limiter
```

```go
import "github.com/gofiber/fiber/v2/middleware/limiter"

app.Use(limiter.New(limiter.Config{
    Max:        100,
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c *fiber.Ctx) string {
        return c.IP()
    },
}))
```

## 🎯 Próximos passos

Continue para [Rodando](./03-rodando.md) para executar e testar a API.
