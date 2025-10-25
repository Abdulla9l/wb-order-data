# WB Order Data Service

Микросервис для обработки заказов с веб-интерфейсом. Получает заказы из Kafka, сохраняет в PostgreSQL и кэширует в памяти для быстрого доступа.

## Функциональность

- **Прием заказов** из Kafka в реальном времени
- **Сохранение в PostgreSQL** с нормализованной структурой
- **In-memory кэширование** последних заказов (LRU)
- **HTTP REST API** для поиска заказов по ID
- **Веб-интерфейс** для удобного просмотра данных
- **Восстановление кэша** из БД при перезапуске

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

