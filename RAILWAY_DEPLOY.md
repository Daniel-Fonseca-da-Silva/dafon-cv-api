# Deploy no Railway

Este documento contém as instruções para fazer deploy da API no Railway.

## Pré-requisitos

1. Conta no Railway (https://railway.app)
2. Projeto conectado ao GitHub
3. Banco de dados MySQL configurado (pode usar o MySQL do Railway ou externo)

## Passos para Deploy

### 1. Conectar o Repositório

1. Acesse o Railway Dashboard
2. Clique em "New Project"
3. Selecione "Deploy from GitHub repo"
4. Escolha este repositório

### 2. Configurar Variáveis de Ambiente

No Railway Dashboard, vá em "Variables" e configure as seguintes variáveis:

#### Obrigatórias:
```bash
# Porta (Railway define automaticamente)
PORT=8080

# Modo do Gin
GIN_MODE=release

# Configuração do Banco de Dados
DB_HOST=seu-host-mysql
DB_PORT=3306
DB_USER=seu-usuario-mysql
DB_PASSWORD=sua-senha-mysql
DB_NAME=nome-do-banco
DB_SSL_MODE=require

# Configuração de Email (Resend)
RESEND_API_KEY=sua-chave-resend
MAIL_FROM=seu-email@dominio.com

# URL da aplicação (será definida automaticamente pelo Railway)
APP_URL=https://seu-projeto.railway.app
```

#### Opcionais:
```bash
# Configuração de Worker Pool
WORKER_POOL_NUM_WORKERS=5
WORKER_POOL_QUEUE_SIZE=100
```

### 3. Configurar Banco de Dados

#### Opção 1: MySQL do Railway
1. No Railway Dashboard, clique em "New Service"
2. Selecione "Database" → "MySQL"
3. Railway criará automaticamente as variáveis de ambiente do banco
4. Use essas variáveis no seu serviço principal

#### Opção 2: Banco Externo
Configure as variáveis de ambiente manualmente com os dados do seu banco externo.

### 4. Deploy

1. O Railway detectará automaticamente o `Dockerfile`
2. O build será executado automaticamente
3. A aplicação estará disponível na URL fornecida pelo Railway

## Estrutura de Arquivos para Railway

- `railway.json`: Configuração específica do Railway
- `Dockerfile`: Imagem Docker otimizada
- `.railwayignore`: Arquivos ignorados no deploy
- `RAILWAY_DEPLOY.md`: Esta documentação

## Verificação do Deploy

Após o deploy, você pode verificar se a aplicação está funcionando:

1. Acesse a URL fornecida pelo Railway
2. Teste o endpoint de health: `https://seu-projeto.railway.app/health`
3. Verifique os logs no Railway Dashboard

## Troubleshooting

### Problemas Comuns:

1. **Erro de conexão com banco**: Verifique se as variáveis de ambiente do banco estão corretas
2. **CORS errors**: Verifique se `APP_URL` está configurada corretamente
3. **Build falha**: Verifique se todas as dependências estão no `go.mod`

### Logs:
- Acesse o Railway Dashboard → Seu Projeto → Deployments → Logs
- Os logs da aplicação Go aparecerão aqui

## Comandos Úteis

```bash
# Instalar Railway CLI (opcional)
npm install -g @railway/cli

# Login no Railway
railway login

# Deploy local (para testes)
railway up
```

## Monitoramento

O Railway fornece métricas básicas:
- CPU e Memory usage
- Request logs
- Deploy history
- Environment variables

Para monitoramento mais avançado, considere integrar com ferramentas como:
- Sentry (error tracking)
- DataDog (APM)
- New Relic (APM)
