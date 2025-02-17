# Notifier

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/Rabbitmq-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=Prometheus&logoColor=white)
![Grafana](https://img.shields.io/badge/grafana-%23F46800.svg?style=for-the-badge&logo=grafana&logoColor=white)

O Notifier é um sistema de notificações que fornece informações sobre o clima e o tempo,
obtidas do CPTEC (Centro de Previsão de Tempo e Estudos Climáticos).
O sistema permite a criação de usuários, opt-out (para parar de receber notificações)
e o envio de notificações agendadas ou "o mais breve possível".
As notificações incluem a previsão do tempo para uma cidade específica nos próximos 4 dias e,
no caso de cidades litorâneas, a previsão de ondas para o dia atual.
Atualmente, o sistema suporta apenas notificações via webhook, Mas está arquitetado
para receber suporte para outras formas de notificação.

![img.png](arquitetura_sistema.png)

## Funcionalidades

### Criação de usuários

Cadastro de usuários com informações como nome, e-mail,
número de telefone e endpoint de webhook.

### Opt-out

Usuários podem optar por não receber mais notificações.

### Notificações agendadas

Envio de notificações em horários específicos.

### Notificações imediatas

Envio de notificações "o mais breve possível".

### Integração com CPTEC

Consulta de previsão do tempo e de ondas para cidades brasileiras.

### Sistema de filas

Utilização de RabbitMQ para gerenciamento de mensagens e notificações.

### Observabilidade

Dashboard com métricas de disponibilidade e desempenho da API,
utilizando Prometheus e Grafana.

## Arquitetura da Solução

### Componentes

- **API Usuários de Notificação**\
Aplicação em Go responsável por expor rotas http para criação de usuários,
opt-out de usuário, e solicitar notificação.

- **Worker Producer**\
Aplicação em Go que periodicamente busca notificações pendentes no banco de dados
e as envia para o serviço de mensageria (RabbitMQ) com referências do usuário,
cidade e tipo de notificação.

- **Worker Consumer**\
Aplicação em Go que consome mensagens da fila do RabbitMQ,
consulta a API do CPTEC nas rotas para obter a previsão do tempo e envia a notificação
ao usuário com base no tipo de notificação configurada (webhook, e-mail, SMS, etc.).

- **Banco de dados**\
Banco de dados PostgreSQL para persistência dos dados de usuários e notificações.

- **Serviço de Mensageria**\
Brooker RabbitMQ com uma exchange chamada notifications e filas para cada tipo
de notificação, já pensando na escalabilidade da aplicação.

      - webhook.notifications
      - email.notifications
      - sms.notifications
      - push.notifications

- **Prometheus e Grafana**\
Coleta de métricas de desempenho, emissão de alertas em caso de
indisponibilidade da api e dashboards com os dados de uso de recurso por parte
da aplicação.

## 🛠️ Configuração e Uso

### Pré-requisitos

- Docker
- Docker Compose

### Como Iniciar o Sistema

- Primeira inicialização:

```bash
docker compose up --build
```

- Inicializações subsequentes:

```bash
docker compose up
```

### Endpoints

- **Criar um usuário**:\
`POST localhost:8081/api/v1/users`

```json
{
    "name": "Michael Scott",
    "email": "m.s@dundermifflin.com",
    "phone_number": "21982438803",
    "webhook": "www.google.com"
}
```

- **Desativar recebimento de notificação de um usuário**:\
`PUT localhost:8081/api/v1/users/opt-out/{id}`

- **Solicitar uma Notificação**:\
`PUT localhost:8081/api/v1/notification`

```json
{
   "date": "2025-02-16T04:32:00Z",
   "city_name": "rio de janeiro",
   "user_id": 1,
   "notification_type": "webhook"
}
```
O campo "date" pode ser omitido para enviar uma notificação imediata.

### Como testar o sistema

1. Criar usuário e guardar id retornado.

2. Solicitar um envio de notificação imediata (sem o campo "date") e validar se o webhook foi chamado.

3. Solicitar um envio de notificação para uma data e horário da sua escolha e validar se o webhook foi chamado no momento correto.

4. Realizar o opt-out do usuário.

5. Solicitar um envio de notificação para o usuário. Verificar se não foi chamado o webhook e analisar o log do consumer com o motivo do não envio.
```bash
ERR Failed to process message queue=webhook.notifications error="user is not accepting notifications"
```

6. Caso não tenha um webhook para validar, o log do consumer indica o conteúdo da notificação (previsão do tempo) e o motivo da falha.
```bash
consumer    | Feb 17 05:41:06.153 INF Sending Webhook usuario=1 city="rio de janeiro" content="{"previsão_do_tempo":{"nome":"Rio de Janeiro","uf":"RJ","atualizacao":"2025-02-16","previsao":[{"dia":"2025-02-17","tempo":"pn","maxima":37,"minima":26,"iuv":0},{"dia":"2025-02-18","tempo":"pn","maxima":38,"minima":27,"iuv":0},{"dia":"2025-02-19","tempo":"pn","maxima":33,"minima":25,"iuv":0},{"dia":"2025-02-20","tempo":"pn","maxima":34,"minima":24,"iuv":0}]},"ondas_do_dia":{"nome":"Rio de Janeiro","uf":"RJ","atualizacao":"16-02-2025","manha":{"dia":"16-02-2025 12h Z","agitacao":"Fraco","altura":"1.4","direcao":"E","vento":"6.1","vento_dir":"ENE"},"tarde":{"dia":"16-02-2025 18h Z","agitacao":"Fraco","altura":"1.5","direcao":"ESE","vento":"8.7","vento_dir":"E"},"noite":{"dia":"16-02-2025 21h Z","agitacao":"Fraco","altura":"1.5","direcao":"ESE","vento":"8.9","vento_dir":"ENE"}}}"
consumer    | Feb 17 05:41:06.153 ERR Failed to process message queue=webhook.notifications error="could not request the webhook"
```

## 🛠️ Tecnologias Utilizadas

Go: Linguagem principal para desenvolvimento das APIs e workers.

Docker: Conteinerização da aplicação.

PostgreSQL: Banco de dados para armazenamento de usuários e notificações.

RabbitMQ: Sistema de mensageria para gerenciamento de filas.

Prometheus: Coleta de métricas de desempenho e emissão de
alertas em caso de indisponibilidade da api, garantindo a
rápida ação para restaurar o sistema.

Grafana: Visualização das métricas coletadas pelo Prometheus.

Squirrel: Biblioteca para construção de queries SQL em Go.

## Próximos Passos

- Adicionar suporte para outros tipos de notificação (e-mail, SMS, push).

- Implementar autenticação e autorização para as APIs.

## Referências

- [Api do CPTEC](http://servicos.cptec.inpe.br/XML/):
   1. **listaCidade**: Obter identificador da cidade.
   1. **previsao**: Obter previsão do tempo para 4 dias.
   1. **ondas**: Obter previsão de ondas.

- [Dashboard de Observabilidade no compose](https://grafana.com/docs/grafana/latest/administration/provisioning/)

- [Instrumentando aplicação com Prometheus](https://prometheus.io/docs/guides/go-application/)

- [Criando alertas no Prometheus](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)

- [Documentação do Framework GIN](https://gin-gonic.com/)
