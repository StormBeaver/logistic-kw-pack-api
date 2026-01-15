Проект выполнен в рамках курса [Route256 2021 OZON](https://rutracker.org/forum/viewtopic.php?t=6201055)

#### Репозиторий с заданиями
https://github.com/ozonmp/omp-docs

# О проекте
Проект включает в себя помимо настоящего репозиториия, ещё три:
* https://github.com/StormBeaver/logistic-pack-bot — телеграм бот;
* https://github.com/StormBeaver/logistic-pack-retranslator - ретранслятор;
* https://github.com/StormBeaver/logistic-pack-facade — фасад.

## gRPC-server
Проект представленный в данном репозитории является gRPC-server'ом, который объединяет три проекта указанных выше. 

### API
Protobuf контракт API описан в `api\logisticPack\logistic_pack_api\v1\logistic_pack_api.proto`, предоставляет CRUD-методы. Код grpc объектов, методов, валидации вынесен в отдельный модуль `pkg\logistic-package-api`.

В обработчиках запросов имеется возможность поднять уровень логирования с помощью метаданных запроса. Для этого необходимо передать по ключу `"grpc-metadata-log-level"` значение `"debug"` или другой иной уровень логирования.

### Repository
Данные хранятся в Postgres. Методы пакета `internal/repo` повторяют методы `internal/api`, однако, при этом, реализуют паттерн [transactional outbox](https://microservices.io/patterns/data/transactional-outbox.html). Помимо изменения основной таблицы, хранящей записи о сущностях домена `pack`, при выполнении CUD-методов добавляется запись в таблицу `packs_events`, описывающая произошедшие изменения.

Текст SQL-запросов собирается с помощью [squirell](https://github.com/Masterminds/squirrel).

#### Миграции
Миграции описаны в `migrations`.

Индексы созданы:
* на столбцах `id, removed` в таблице `packs`, поскольку по нему идут условия `WHERE` в запросах `Get`, `Remove`, `Update`;
* на столбце `id` в таблице `packs_events`, поскольку по нему идёт условие `WHERE` в запросах ретранслятора (см. О проекте).

## Метрики
Сервис собирает метрики с помощью [Prometheus](https://github.com/prometheus/client_golang).

### Метрики gRPC-server'a
grpc-сервер собирает две метрики:
* `totalPackNotFound` _Counter_ — общее количество NotFound событий;
* `totalPackCUDEvents` _Counter_ — общее количество CUD событий.

Метрики grpc-сервера доступны на `:9100/metrics`.

### Grafana и Prometheus

С заданным интервалом Prometheus считывает метрики сервисов. Метрики запрашиваются Grafan'ой и отображаются на графике:
`logistic_pack_api_not_found_total` и `logistic_pack_api_cud_found_total` — рассчитывая количество происходящих событий в секунду в окне в 1 минуту.

Grafana доступна на `:3000`.

## Jaeger

Обработчики запросов gRPC-server'a записывает трейсы и отправляет их в Jaeger. Трейсы имеют вложенные спаны методов репозиторя и SQL-запросов.

Jaeger доступен на `:16686`.

## Docker

Описаны докерфайлы для образов gRPS-server'a.

## Makefile

В мейкфайле описаны команды для локального запуска приложения, для генерации protobuf модулей с помощью утилиты buf, сборки докер-образов, и их запуска.