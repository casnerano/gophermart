# Gophermart - система лояльности

## Запуск
```bash
make init && make start
```

## Чек-лиск на доработку

### Приложение
- [x] ~~Graceful shutdown~~
- [x] ~~Облегчить main-функцию, вынести часть инциализации в сервисы~~
- [x] ~~Ошибки уровня ниже Warning только на dev-окружении~~
- [ ] Декомпозировать часть перегруеных сервисов
- [ ] Добавить больше бизнесовы логов в хендлеры
- [ ] Покрыть тестами хендлеры

### Инфраструктура
- [x] ~~Запуск dev-окружения в докере~~
- [x] ~~Добавить healthcheck для сервиса RabbitMQ~~
- [x] ~~Автоматический запуск миграций~~
- [ ] Настроить автодеплой на сервер
- [ ] DSN подключений к PostgreSQL и RabbitMQ формировать на уровне приложения
- [ ] Доступы к сервисам хранить в виде секретов
