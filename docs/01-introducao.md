## 📖 O que é este template?

API REST completa para gerenciamento de contatos construída com **Fiber**, **GORM** e **PostgreSQL**. Inclui CRUD completo, validação, migrations, busca e paginação.

## 🎯 Casos de uso

- **Agenda de contatos:** Sistema de gerenciamento de contatos pessoais
- **CRM simples:** Base para sistema de relacionamento com clientes
- **Backend para apps:** API para aplicativos de contatos (mobile/web)
- **Aprendizado:** Exemplo de API REST com Go
- **Base para projetos:** Ponto de partida para sistemas maiores
- **Microserviço:** Componente de sistema de gestão

## ✨ Características principais

### API REST Completa

- ✅ CRUD completo de contatos (Create, Read, Update, Delete)
- ✅ Busca por nome, email e empresa
- ✅ Paginação de resultados
- ✅ Validação robusta de dados
- ✅ Soft delete (exclusão lógica)
- ✅ Filtros e ordenação

### Performance e Escalabilidade

- ⚡ Fiber (framework mais rápido do Go)
- 🗄️ GORM ORM (type-safe e otimizado)
- 📊 Queries otimizadas com índices
- 🔄 Connection pooling
- 📦 Build otimizado com Docker

### Qualidade de Código

- ✅ Validação com go-playground/validator
- ✅ Tratamento robusto de erros
- ✅ Logs estruturados (Logrus)
- ✅ Type safety com Go
- ✅ Testes automatizados

### DevOps

- 🐳 Docker e Docker Compose
- 🔄 Migrations com golang-migrate
- 🌱 Seeds para desenvolvimento
- 🏥 Health check endpoint
- 📊 Logs estruturados JSON

## 🏗️ Arquitetura

```
api-contacts-go/
├── cmd/
│   └── server/
│       └── main.go        # Entry point
├── internal/
│   ├── config/            # Configurações
│   ├── database/          # Conexão DB e migrations
│   ├── models/            # GORM models
│   ├── handlers/          # HTTP handlers
│   ├── services/          # Lógica de negócio
│   └── middleware/        # Middlewares
├── migrations/            # SQL migrations
├── seed/                  # Dados iniciais
├── tests/                 # Testes
└── docs/                  # Documentação
```

### Stack Tecnológica

- **Runtime:** Go 1.21+
- **Framework:** Fiber v2
- **ORM:** GORM
- **Migrations:** golang-migrate
- **Database:** PostgreSQL
- **Validation:** go-playground/validator
- **Logs:** Logrus
- **Tests:** testing + testify
- **Container:** Docker + Docker Compose

## 📊 Modelo de Dados

### Contact (Contato)
```go
type Contact struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    Name      string         `json:"name" gorm:"not null"`
    Email     string         `json:"email" gorm:"uniqueIndex;not null"`
    Phone     string         `json:"phone" gorm:"size:20"`
    Company   string         `json:"company" gorm:"size:100"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
```

### Campos

- **id:** ID único auto-incremento
- **name:** Nome do contato (obrigatório, 2-100 chars)
- **email:** Email único (obrigatório, formato válido)
- **phone:** Telefone opcional (10-20 chars)
- **company:** Empresa opcional (max 100 chars)
- **created_at:** Data de criação automática
- **updated_at:** Data de última atualização automática
- **deleted_at:** Soft delete (exclusão lógica)

## 🔗 Endpoints Disponíveis

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/health` | Health check + status do banco |
| GET | `/api/v1/contacts` | Listar contatos (paginado) |
| GET | `/api/v1/contacts/:id` | Buscar contato por ID |
| POST | `/api/v1/contacts` | Criar novo contato |
| PUT | `/api/v1/contacts/:id` | Atualizar contato |
| DELETE | `/api/v1/contacts/:id` | Deletar contato (soft delete) |
| GET | `/api/v1/contacts/search` | Buscar contatos por texto |

## 🚀 Quick Start

```bash
# Clone e acesse
cd api-contacts-go

# Suba com Docker Compose
docker compose up -d

# Rode migrations
docker compose exec api make migrate

# Popule com dados
docker compose exec api make seed

# Teste
curl http://localhost:80/health
curl http://localhost:80/api/v1/contacts
```

## 📈 Performance Esperada

- **Latência:** < 5ms para queries simples
- **Throughput:** > 20k requests/segundo
- **Memória:** ~20MB em idle
- **Startup:** < 1 segundo

## 🔄 Próximos passos

Continue para [Configuração](./02-configuracao.md) para setup detalhado.
