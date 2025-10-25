# WB Order Data Service

Микросервис для обработки заказов с веб-интерфейсом. Получает заказы из Kafka, сохраняет в PostgreSQL и кэширует в памяти для быстрого доступа.

## Функциональность

- **Прием заказов** из Kafka в реальном времени
- **Сохранение в PostgreSQL** 
- **HTTP REST API**
- **Веб-интерфейс**
- **Восстановление кэша**

## Технологии

- **Backend:** Go 1.21+
- **Database:** PostgreSQL 15+
- **Message Broker:** Apache Kafka
- **HTTP Router:** Gorilla Mux
- **Frontend:** JavaScript + HTML

## Команды 
- make migrate      
- make run           
- make producer     
- make test          
- make migrate       

