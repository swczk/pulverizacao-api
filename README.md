# API GraphQL - Dados Geográficos

API GraphQL desenvolvida em Go com MongoDB para gerenciamento exclusivo de dados geográficos de aplicações de pulverização agrícola.

## Arquitetura do Sistema

Esta API faz parte de um sistema maior onde:
- **Spring Boot API (porta 8080)**: Gerencia CRUD de equipamentos, talhões, aplicações, etc.
- **GraphQL API (porta 8081)**: Gerencia EXCLUSIVAMENTE dados geográficos

## Estrutura do Projeto

```
graphql-api/
├── main.go
├── go.mod 
├── Dockerfile
├── Makefile
├── config/
│   └── config.go
├── database/
│   └── connection.go
├── models/
│   └── geo.go
├── graphql/
│   └── schema.go
├── scripts/
│   └── init-mongo.js
└── test/
    └── geo.http
```

## 🚀 Quick Start

### Docker Compose (no projeto principal)

A API GraphQL é executada como parte do sistema completo:

```bash
# No diretório raiz do projeto
docker compose -f compose.dev.yaml up -d
```

Serviços disponíveis:
- **API GraphQL**: http://localhost:8081/graphql
- **Spring Boot API**: http://localhost:8080
- **MongoDB**: localhost:27017

## 📊 Uso da API GraphQL

Acesse GraphiQL em: `http://localhost:8081/graphql`

### Operações Disponíveis

#### Queries

**Buscar trajetória geográfica por ID da aplicação:**
```graphql
query {
  geoTrajetoria(aplicacaoId: "507f1f77bcf86cd799439011") {
    aplicacaoId
    pontoInicial {
      latitude
      longitude
      timestamp
      altitude
      speed
      accuracy
    }
    pontoFinal {
      latitude
      longitude
      timestamp
      altitude
      speed
      accuracy
    }
    trajetoria {
      latitude
      longitude
      timestamp
      altitude
      speed
      accuracy
    }
    areaCobertura
    distanciaPercorrida
    createdAt
    updatedAt
  }
}
```

**Buscar todas as trajetórias com paginação:**
```graphql
query {
  geoTrajetorias(limit: 10, offset: 0) {
    aplicacaoId
    pontoInicial {
      latitude
      longitude
      timestamp
    }
    pontoFinal {
      latitude
      longitude
      timestamp
    }
    areaCobertura
    distanciaPercorrida
    createdAt
  }
}
```

#### Mutations

**Criar trajetória geográfica:**
```graphql
mutation {
  createGeoTrajetoria(input: {
    aplicacaoId: "507f1f77bcf86cd799439011"
    pontoInicial: {
      latitude: -25.4284
      longitude: -49.2733
      timestamp: "2025-06-16T06:00:00Z"
      accuracy: 3.5
      altitude: 945.2
      speed: 0.0
    }
    pontoFinal: {
      latitude: -25.4310
      longitude: -49.2690
      timestamp: "2025-06-16T07:45:00Z"
      accuracy: 4.1
      altitude: 952.8
      speed: 12.5
    }
    trajetoria: [
      {
        latitude: -25.4284
        longitude: -49.2733
        timestamp: "2025-06-16T06:00:00Z"
        accuracy: 3.5
        altitude: 945.2
        speed: 0.0
      }
    ]
    areaCobertura: 5.2
    distanciaPercorrida: 850.5
  }) {
    aplicacaoId
    areaCobertura
    distanciaPercorrida
  }
}
```

**Atualizar trajetória geográfica:**
```graphql
mutation {
  updateGeoTrajetoria(
    aplicacaoId: "507f1f77bcf86cd799439011"
    input: {
      pontoFinal: {
        latitude: -25.4315
        longitude: -49.2685
        timestamp: "2025-06-16T08:00:00Z"
        accuracy: 3.8
        altitude: 955.0
        speed: 8.2
      }
      novosPontos: [
        {
          latitude: -25.4312
          longitude: -49.2687
          timestamp: "2025-06-16T07:50:00Z"
          accuracy: 3.9
          altitude: 953.5
          speed: 10.1
        }
      ]
      areaCobertura: 6.1
      distanciaPercorrida: 920.8
    }
  ) {
    aplicacaoId
    areaCobertura
    distanciaPercorrida
    updatedAt
  }
}
```

**Deletar trajetória geográfica:**
```graphql
mutation {
  deleteGeoTrajetoria(aplicacaoId: "507f1f77bcf86cd799439011")
}
```

## Modelos de Dados

### GeoPoint
- `latitude`: Latitude do ponto
- `longitude`: Longitude do ponto
- `timestamp`: Timestamp do ponto
- `altitude`: Altitude (opcional)
- `speed`: Velocidade (opcional)
- `accuracy`: Precisão GPS (opcional)

### GeoTrajetoria
- `aplicacaoId`: ID da aplicação (referência para Spring Boot API)
- `pontoInicial`: Ponto inicial da trajetória
- `pontoFinal`: Ponto final da trajetória (opcional)
- `trajetoria`: Array de pontos da trajetória
- `areaCobertura`: Área coberta em hectares
- `distanciaPercorrida`: Distância total percorrida em metros
- `createdAt`: Data de criação
- `updatedAt`: Data de atualização

## Collections MongoDB

- `geo_trajetorias`: Trajetórias geográficas das aplicações

## Configuração

### Variáveis de Ambiente
```bash
MONGO_URI=mongodb://localhost:27017
DATABASE_NAME=pulverizacao
PORT=8081
```

## Testes

Use os arquivos `.http` na pasta `test/` para testar as operações:

```bash
# Teste todas as operações geográficas
curl -X POST http://localhost:8081/graphql \
  -H "Content-Type: application/json" \
  -d @test/geo.http
```