# Телеграм бот "Менеджер трат"

## Подготовка

Указать DSN поднятого PostgreSQL:

- ExpenseManager-CBR-Telegram/cmd/bot/main.go:31

    - `db, err := gorm.Open(postgres.Open("host=localhost port=5432 user=postgres password=postgres"))`

- ExpenseManager-CBR-Telegram/Makefile:10

    - `DSN := "host=localhost port=5432 user=postgres password=postgres sslmode=disable"`

## Запуск

`make goose-up && make run` - запустит бота локально

## Команды

### Траты

- ***/add <сумма>*** - добавляет новую трату без комментария, в качестве даты берет текущую

- ***/add <сумма>; <комментарий>*** - добавляет новую трату с комментарием, в качестве даты берет текущую

- ***/add <сумма>; <комментарий>; <dd.mm.yyyy>*** - добавляет новую трату с комментарием и выставляет соответствующую дату

- ***/add <сумма>; ; <dd.mm.yyyy>*** - добавляет новую трату без комментария и выставляет соответствующую дату

### Отчёты

- ***/spent_<week|month|year>*** - собирает отчет за указанный промежуток. Отчет отправляет в виде текста

- ***/spent*** - собирает отчет за всё время. Отчет отправляет в виде текста

### Дополнительно

- ***/change_currency*** - сменить основную валюту пользователя. После этой команды отчеты будут пересчитываться в валюту по курсу в соответсвующую дату

## Архитектура

```
.
├── bin
├── cmd
├── data                             - (секрет)
│      └── config.yaml               - конфигурационный файл 
├── internal
│      ├── config
│      ├── database                  - содержит методы к базе данных
│      ├── domain                    - содержит доменные сущности (вынесено из database, чтобы не тянуть
│      │                               зависимости от database и поспользоваться утиной типизацией)
│      ├── helpers
│      │      ├── date               - доп функция по конвертации даты
│      │      └── money              - доп функции по работе с копейками валют
│      ├── infrastructure 
│      │      ├── cbr_gateway        - клиент cbr - апи для получения курсов валют
│      │      └── tg_gateway         - клиент телеграмма
│      ├── mocks                     - вынесенные mock-объекты для тестирования 
│      ├── model
│      │      └── messages           - содержит основную бизнес-логику трат и отчётов
│      ├── services                  - оборачивает клиент cbr в необходимую бизнес-логику
│      └── worker                    - воркеры:
│                                      1 - выполняет функции контроллера и отлавливает команды,
│                                      2 - запущенный в отдельном потоке ходит в cbr и обновляет курсы валют по таймеру 
├── migrations                       - миграции для базы данных  
├── config.example.yaml              - пример конфигурационного файл
├── go.mod
├── go.sum
├── README.md
└── Makefile
```

### Особенности
- Функционал покрыт unit тестами
- Для gracefully завершения использован Context
- Код написан так, что новые валюты добавляются с минимальными изменениями (добавлением кода валюты в config.yaml)
- Для запроса курса валют используются API ЦБ и данные по курсам действительно актуальные
- Происходит минимум запросов в API для достижения цели