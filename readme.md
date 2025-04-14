# 🧱 rate-limiter

Rate limiter desenvolvido em Go com controle por IP e token, utilizando Redis para armazenamento.

### Como o Rate Limiter Funciona

O rate limiter controla o número de requisições feitas por um identificador (IP ou token) dentro de um intervalo de tempo. Ele usa o Redis para:

- **Contar** quantas requisições foram feitas.
- **Bloquear** temporariamente quem ultrapassar o limite.

#### Lógica Interna

1. Ao receber uma requisição:
   - O identificador (IP ou token) é extraído.
   - O Redis incrementa o contador associado a ele.
   - Se for a primeira requisição, uma expiração curta (ex: 1 segundo) é aplicada.
2. Se o contador ultrapassar o limite:
   - O identificador é **bloqueado** no Redis por um tempo determinado.
   - Requisições seguintes são negadas até o tempo expirar.

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
