# 🧱 rate-limiter

Rate limiter desenvolvido em Go com suporte a controle por IP e token, utilizando Redis para armazenamento.

## Funcionalidades

- Limitação de requisições por IP e por token.
- Armazenamento de contadores e bloqueios com Redis.
- Configuração via variáveis de ambiente.
- Suporte a benchmark de carga com goroutines.
- Pronto para execução via Docker e Docker Compose.

## Configuração

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```env
PORT=8080

# Limite por IP
RATE_LIMIT_IP=10
BLOCK_DURATION_IP=300

# Limite por token
RATE_LIMIT_TOKEN_DEFAULT=3
BLOCK_DURATION_TOKEN=300

# Redis
REDIS_ADDR=redis:6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Execução

```
docker compose up -d --build
```

### Url de exemplo

```
curl http://localhost:8080/go
```
### Retorno quando não bloqueado

```
In Go We Trust!
```
### Retorno quando bloqueado

```
you have reached the maximum number of requests or actions allowed within a certain time frame
```


## Benchmark

```
docker compose run --rm benchmark
```