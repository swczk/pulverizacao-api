# API GraphQL - Sistema de Pulverização

API GraphQL desenvolvida em Go com MongoDB para gerenciamento de aplicações de pulverização agrícola.

## Estrutura do Projeto

```
pulverizacao-api/
├── main.go
├── go.mod 
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── config/
│   └── config.go
├── database/
│   └── connection.go
├── models/
│   └── aplicacao.go
├── graphql/
│   └── schema.go
└── scripts/
    └── init-mongo.js
```

## 🚀 Quick Start com Docker

### Opção 1: Docker Compose (Recomendado para desenvolvimento)

```bash
# Clone o repositório
git clone <repo-url>
cd pulverizacao-api

# Configure as variáveis de ambiente
cp .env.example .env

# Inicie todos os serviços
docker-compose up -d

# Verifique o status
docker-compose ps
```

Serviços disponíveis:
- **API GraphQL**: http://localhost:8080/graphql
- **MongoDB**: localhost:27017
- **Mongo Express**: http://localhost:8081 (admin/pass)

### Opção 2: Docker Hub

```bash
# Baixe e execute a imagem
docker run -p 8080:8080 \
  -e MONGO_URI="your-mongodb-uri" \
  -e DATABASE_NAME="pulverizacao" \
  your-dockerhub-username/pulverizacao-api:latest
```

## 📦 Instalação Local

1. Configure MongoDB Atlas:
   - Crie cluster no MongoDB Atlas
   - Configure usuário e senha
   - Adicione IP à whitelist
   - Obtenha connection string
2. Instale dependências:
   ```bash
   go mod tidy
   ```
3. Configure variáveis no `.env`:
   ```
   MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority
   DATABASE_NAME=pulverizacao
   PORT=8080
   ```
4. Execute:
   ```bash
   go run main.go
   ```

## 🛠️ Comandos Make

```bash
# Desenvolvimento
make dev          # Executar localmente
make test         # Executar testes
make build        # Build da aplicação

# Docker
make docker-build # Build da imagem Docker
make docker-run   # Executar container
make docker-push  # Push para registry

# Docker Compose
make up           # Iniciar serviços
make down         # Parar serviços
make logs         # Ver logs
make rebuild      # Rebuild completo

# Utilitários
make status       # Status dos serviços
make health       # Verificar saúde
make cleanup      # Limpeza
```

## ## 🚀 CI/CD e Deploy

### GitHub Actions

O projeto inclui pipeline automático de CI/CD:

1. **Configurar secrets no GitHub:**
   - `DOCKER_USERNAME`: seu username do Docker Hub
   - `DOCKER_PASSWORD`: token de acesso do Docker Hub

2. **Pipeline automático:**
   - Testes automatizados
   - Build multi-arquitetura (amd64/arm64)
   - Push para Docker Hub em push para `main`
   - Tags automáticas baseadas em versões

3. **Comandos para release:**
   ```bash
   # Tag de versão
   git tag v1.0.0
   git push origin v1.0.0
   
   # Será criada automaticamente:
   # - your-username/pulverizacao-api:v1.0.0
   # - your-username/pulverizacao-api:1.0
   # - your-username/pulverizacao-api:latest
   ```

### Deploy em Produção

**Opção 1: Docker Compose**
```bash
# Servidor de produção
docker-compose -f docker-compose.yml up -d
```

**Opção 2: Kubernetes**
```bash
# Usando imagem do Docker Hub
kubectl create deployment pulverizacao-api \
  --image=your-username/pulverizacao-api:latest
kubectl expose deployment pulverizacao-api \
  --port=8080 --target-port=8080
```

**Opção 3: Cloud Services**
- AWS ECS/Fargate
- Google Cloud Run
- Azure Container Instances

Uso

Acesse GraphiQL em: `http://localhost:8080/graphql`

### Queries

**Buscar todas as aplicações:**
```graphql
query {
  aplicacoes(limit: 10, offset: 0) {
    id
    operador
    dosagem
    finalizada
    talhao {
      nome
      cultura
    }
    equipamento {
      nome
      modelo
    }
  }
}
```

**Buscar aplicação por ID:**
```graphql
query {
  aplicacao(id: "507f1f77bcf86cd799439011") {
    id
    operador
    dosagem
    dataInicio
    condicaoClimatica
    observacoes
    finalizada
  }
}
```

### Mutations

**Criar aplicação:**
```graphql
mutation {
  createAplicacao(input: {
    talhaoId: "507f1f77bcf86cd799439011"
    equipamentoId: "507f1f77bcf86cd799439012"
    tipoAplicacaoId: "507f1f77bcf86cd799439013"
    dataInicio: "2025-05-28T10:00:00Z"
    dosagem: 2.5
    operador: "João Silva"
    condicaoClimatica: "Ensolarado"
    observacoes: "Aplicação normal"
  }) {
    id
    operador
    dosagem
    finalizada
  }
}
```

**Atualizar aplicação:**
```graphql
mutation {
  updateAplicacao(
    id: "507f1f77bcf86cd799439011"
    input: {
      finalizada: true
      observacoes: "Aplicação finalizada"
    }
  ) {
    id
    finalizada
    dataFim
    observacoes
  }
}
```

**Deletar aplicação:**
```graphql
mutation {
  deleteAplicacao(id: "507f1f77bcf86cd799439011")
}
```

## Modelos

### Aplicacao
- `id`: ID único
- `talhaoId`: Referência ao talhão
- `equipamentoId`: Referência ao equipamento  
- `tipoAplicacaoId`: Referência ao tipo de aplicação
- `dataInicio`: Data/hora de início
- `dataFim`: Data/hora de fim (opcional)
- `dosagem`: Dosagem aplicada
- `volumeAplicado`: Volume aplicado (opcional)
- `operador`: Nome do operador
- `condicaoClimatica`: Condições climáticas
- `observacoes`: Observações gerais
- `finalizada`: Status da aplicação

### Relacionamentos
- `talhao`: Dados do talhão (via lookup)
- `equipamento`: Dados do equipamento (via lookup)
- `tipoAplicacao`: Dados do tipo de aplicação (via lookup)

## Configuração

### MongoDB Atlas
1. Crie um cluster no MongoDB Atlas
2. Configure usuário de banco de dados em Security > Database Access
3. Adicione IP à whitelist em Security > Network Access
4. Obtenha connection string em Connect > Connect your application

### Variáveis de Ambiente (.env)
```bash
# MongoDB Atlas connection string
MONGO_URI=mongodb+srv://username:password@cluster.mongodb.net/?retryWrites=true&w=majority

# Database name
DATABASE_NAME=pulverizacao

# Server port
PORT=8080
```

**Exemplo real:**
```bash
MONGO_URI=mongodb+srv://myuser:mypass123@cluster0.abc12.mongodb.net/?retryWrites=true&w=majority
```

### Collections MongoDB
- `aplicacoes`: Aplicações de pulverização
- `talhoes`: Talhões/terrenos
- `equipamentos`: Equipamentos de pulverização
- `tipos_aplicacao`: Tipos de aplicação