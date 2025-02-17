# Notifier
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) ![RabbitMQ](https://img.shields.io/badge/Rabbitmq-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white) ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)


O Notifier é um sistema de notificações com informações de clima e tempo obtidas do CPTEC (Centro de Previsão de Tempo e Estudos Climáticos). O sistema possui recursos de criação de usuário, opt-out de usuário (não permite mais recebimento de notificação) e notificações agendadas ou "o mais breve possível". A notificação é a previsão do tempo para determinada cidade para os próximos 4 dias e a previsão de ondas do dia atual caso seja uma cidade litorânea. Atualmente é suportado apenas um tipo de notificação: webhook (rota de uma aplicação web).

## Arquitetura da solução

![img.png](arquitetura_sistema.png)

### Componentes

- **API Usuários de Notificação**  
Aplicação em Go responsável por expor rotas http para criação de usuários, opt-out de usuário, e solicitar notificação.

- Worker Producer  
Aplicação em Go que periodicamente busca notificações que precisam ser enviadas no banco de dados e envia uma mensagem para o serviço de mensageria com a referencia do usuário e cidade e tipo de notificação.

- Worker Consumer  
Aplicação em Go que monitora uma fila específica no serviço de mensageria e consome mensagens. A partir do conteúdo da mensagem, busca a informação relevante do usuário com base no tipo de notificação, consulta a API Externa de previsão do tempo com o nome da cidade e envia a notificação.

- Banco de dados  
Banco de dados PostgreSQL que armazena dados de usuários e notificações.

- Serviço de Mensageria  
Brooker RabbitMQ com a Exchange notifications e filas conectadas para cada tipo de notificação: webhook.notifications, email.notifications, sms.notifications e push.notifications.

- API Externa de previsão do tempo  
API do CPTEC cujas recursos são utilizados: **listaCidade** para obter identificador da cidade, **previsao** para obter previsão do tempo para 4 dias e **ondas** para obter previsão de ondas. Referência: http://servicos.cptec.inpe.br/XML/

## 🛠️ Setup

### Pré-Requisitos

- docker compose

### Como iniciar o sistema

- Primeira inicialização
```bash
docker compose up --build
```

- Proximas inicializações
```bash
docker compose up
```
### Como interagir com o sistema

- Criação de usuário  
Todo usuário criado vem ativado(aceita notificação) por padrão. Atenção ao webhook que será chamado ao recebimento de notificação.
POST localhost:8081/api/v1/users
```json
{
    "name": "Michael Scott",
    "email": "m.s@dundermifflin.com",
    "phone_number":"21982438803",
    "webhook":"www.google.com"
}
```

- Cria solicitação de notificação 
A data precisa estar no formato ISO 8601
POST localhost:8081/api/v1/notification
```json
{
   "date": "2025-02-16T04:32:00Z",
   "city_name": "rio de janeiro",
   "user_id": 1,
   "notification_type": "webhook"
}
```