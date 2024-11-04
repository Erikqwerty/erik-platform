# Платформеная библиотека

## db

Установка

```shell
go get github.com/erik-platform/clients/db
```

### pg - Клиент для работы с PostgreSQL в Go

Пакет pg предоставляет удобный интерфейс для работы с базой данных PostgreSQL, используя пул соединений pgxpool. Он включает в себя инструменты для выполнения запросов, сканирования данных, управления транзакциями и настройки уровня изоляции транзакций. Пакет ориентирован на гибкость и позволяет пользователям легко использовать транзакции в их приложениях.

### Основные компоненты

#### pgClient

 pgClient - это клиент для работы с базой данных PostgreSQL, который использует пул соединений pgxpool.Pool.

- Конструктор: ```New(ctx context.Context, dsn string) (db.Client, error)```

 Методы:

- ```DB() db.DB```: Возвращает объект базы данных для выполнения запросов.
- ```Close() error```: Закрывает соединение с базой данных.

#### pg

pg реализует интерфейс db.DB для выполнения SQL-запросов через пул соединений pgxpool.

- Конструктор: ```NewDB(dbc *pgxpool.Pool) db.DB```

 Методы:

- ```ScanOneContext(ctx, dest, q, args...)```: Выполняет SQL-запрос и сканирует одну запись в dest.
- ```ScanAllContext(ctx, dest, q, args...)```: Выполняет SQL-запрос и сканирует все записи в dest.
- ```ExecContext(ctx, q, args...)```: Выполняет SQL-запрос (например, INSERT, UPDATE, DELETE).
- ```QueryContext(ctx, q, args...)```: Выполняет SQL-запрос и возвращает pgx.Rows с результатами.
- ```QueryRowContext(ctx, q, args...)```: Выполняет SQL-запрос и возвращает одну строку результата.
- ```BeginTx(ctx, txOptions)```: Начинает транзакцию с указанными параметрами.
- ```Ping(ctx)```: Проверяет соединение с базой данных.
- ```Close()```: Закрывает пул соединений.

#### Транзакционный менеджер

Транзакционный менеджер предоставляет удобный способ управления транзакциями с уровнями изоляции и автоматическим откатом при ошибках.

- Конструктор: ```NewTransactionManager(db db.Transactor) db.TxManager```

  Методы:

- ```ReadCommitted(ctx, handler)```: Выполняет транзакцию с уровнем изоляции Read Committed.

##### Примеры использования

Создание клиента базы данных

```go
package main

import (
    "context"
    "log"

    "github.com/jackc/pgx/v4/pgxpool"
    "github.com/yourusername/pg"
)

func main() {
    // Создаем пул соединений
    dsn := "postgresql://user:password@localhost:5432/dbname"
    ctx := context.Background()
    pool, err := pgxpool.Connect(ctx, dsn)
    if err != nil {
        log.Fatal("Не удалось подключиться к базе данных:", err)
    }
    defer pool.Close()

    // Создаем клиент базы данных
    db := pg.NewDB(pool)

    // Выполнение запроса
    var result YourStruct
    query := db.Query{Name: "GetYourData", QueryRaw: "SELECT * FROM your_table WHERE id = $1"}
    if err := db.ScanOneContext(ctx, &result, query, 1); err != nil {
        log.Fatal("Ошибка при сканировании данных:", err)
    }
}
```

##### Использование транзакций

```go
package main

import (
    "context"
    "log"

    "github.com/jackc/pgx/v4"
    "github.com/yourusername/pg"
    "github.com/yourusername/pg/db"
)

func main() {
    ctx := context.Background()
    dsn := "postgresql://user:password@localhost:5432/dbname"
    pool, err := pgxpool.Connect(ctx, dsn)
    if err != nil {
        log.Fatal("Не удалось подключиться к базе данных:", err)
    }
    defer pool.Close()

    db := pg.NewDB(pool)
    txManager := pg.NewTransactionManager(db)

    // Выполнение транзакции
    err = txManager.ReadCommitted(ctx, func(ctx context.Context) error {
        // Ваша логика внутри транзакции
        query := db.Query{Name: "InsertData", QueryRaw: "INSERT INTO your_table (name) VALUES ($1)"}
        _, err := db.ExecContext(ctx, query, "example")
        return err
    })

    if err != nil {
        log.Fatal("Ошибка при выполнении транзакции:", err)
    }
}
```

##### Пример использования функции Pretty

Функция Pretty из пакета prettier позволяет форматировать SQL-запросы с параметрами для удобного логирования.

```go
package main

import (
    "fmt"

    "github.com/yourusername/pg/prettier"
)

func main() {
    query := "SELECT * FROM users WHERE id = $1"
    formattedQuery := prettier.Pretty(query, prettier.PlaceholderDollar, 123)
    fmt.Println("Отформатированный запрос:", formattedQuery)
}
```

##### Логирование запросов

Для логирования запросов используйте функцию logQuery, которая форматирует и выводит SQL-запрос с параметрами. Логирование можно настроить для интеграции с вашим логгером.

Основные интерфейсы

Интерфейс ```DB```

- ```DB()``` - возвращает объект для работы с базой данных.
- ```Close()``` - закрывает соединение с базой данных.

Интерфейс ```TxManager```

- ```ReadCommitted(ctx, handler)``` - выполняет транзакцию с уровнем изоляции Read Committed.

Интерфейс ```SQLExecer```

- Интерфейсы NamedExecer и QueryExecer для выполнения запросов и сканирования данных.

## closer

Пакет `closer` предназначен для управления функциями освобождения ресурсов, таких как закрытие соединений с базой данных, закрытие файлов и других важных ресурсов в Go-приложениях. Он гарантирует, что все зарегистрированные функции очистки будут вызваны один раз, даже если метод `CloseAll` вызывается несколько раз, а также поддерживает автоматическую очистку при получении сигналов от ОС.

### Особенности

- **Добавление нескольких функций очистки** к одному менеджеру.
- **Автоматическая очистка при получении сигналов** от операционной системы.
- **Гарантированное однократное выполнение** каждой функции очистки, даже если `CloseAll` вызван несколько раз.
- **Параллельное выполнение** функций очистки с обработкой ошибок.

### Установка

```shell
go get github.com/erik-platform/closer
```

Пример использования
Базовая настройка

```go
package main

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "github.com/yourusername/closer"
)

func main() {
    // Добавляем функции очистки
    closer.Add(
        func() error {
            fmt.Println("Закрытие ресурса 1")
            return nil
        },
        func() error {
            fmt.Println("Закрытие ресурса 2")
            return nil
        },
    )

    // Настраиваем обработку сигналов ОС (например, interrupt, terminate)
    globalCloser := closer.New(syscall.SIGINT, syscall.SIGTERM)

    // Ожидание завершения всех операций очистки перед выходом
    defer closer.Wait()

    // Основная логика приложения
    fmt.Println("Приложение запущено...")

    // Если необходимо, можно вызвать очистку вручную
    closer.CloseAll()
}
```

### API

#### func Add(f ...func() error)

Добавляет одну или несколько функций очистки в глобальный Closer. Каждая функция должна возвращать ошибку, если операция очистки завершилась неудачно.

```go
closer.Add(func() error {
    // Освобождение ресурсов
    return nil
})
```

#### func Wait()

Блокирует выполнение до завершения всех функций очистки. Обычно вызывается в конце программы для гарантии, что ресурсы освобождены перед завершением.

```go
closer.Wait()
```

#### func CloseAll()

Инициирует выполнение всех зарегистрированных функций очистки одновременно. Гарантирует, что каждая функция будет вызвана один раз, даже если CloseAll вызван несколько раз.

```go
closer.CloseAll()
```

#### Тип Closer

#### func New(sig ...os.Signal) *Closer

Создает новый экземпляр Closer. Если передан срез сигналов ОС, Closer автоматически вызовет CloseAll при получении любого из этих сигналов.

```go
globalCloser := closer.New(syscall.SIGINT, syscall.SIGTERM)
```

#### func (c *Closer) Add(f ...func() error)

Добавляет функции очистки к конкретному экземпляру Closer.

```go
c := closer.New()
c.Add(func() error {
    // Освобождение конкретного ресурса
    return nil
})
```

#### func (c *Closer) Wait()

Блокирует выполнение до завершения всех зарегистрированных функций очистки.

#### func (c *Closer) CloseAll()

Вызывает все зарегистрированные функции очистки. Выполнение происходит параллельно, ошибки выводятся в лог.

## validator - Валидатор email-адресов

Пакет `validator` предоставляет функции для проверки валидности email-адресов. Он позволяет проверять список email-адресов и возвращать ошибку, если найдены невалидные.

### Установка validator

```shell
go get github.com/erik-platform/utils/validator
```

### ValidEmails

Проверяет список email-адресов. Возвращает ошибку с указанием невалидных адресов, если таковые имеются.

```func ValidEmails(emails []string) error```

- Параметры: ```emails []string``` - список email-адресов для проверки.
- Возвращаемое значение: ```error``` - ошибка с невалидными адресами, либо nil, если все адреса валидны.

Пример:

```go
emails := []string{"<valid@example.com>", "invalid-email"}
if err := validator.ValidEmails(emails); err != nil {
    fmt.Println("Найдены невалидные email-адреса:", err)
}
```

### IsValidEmail

Проверяет валидность одного email-адреса. Возвращает true, если адрес валиден.

```func IsValidEmail(email string) bool```

Пример:

```go
email := "<example@domain.com>"
if validator.IsValidEmail(email) {
    fmt.Println("Email валиден")
}
```

Регулярное выражение

Функция IsValidEmail использует регулярное выражение:

```go
const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
```
