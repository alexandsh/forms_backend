## SL Forms

Бэкенд сервиса форм на Go (Gin, PostgreSQL, Redis).

### Запуск без Docker

#### 1. Установить зависимости

- **Go 1.25.3**
- **PostgreSQL 18**
- **Redis**

#### 2. Создать базу данных

1. Запустить PostgreSQL.
2. Создать базу и пользователя (ниже пример):

```sql
CREATE DATABASE sl_forms;
CREATE USER sl_forms WITH ENCRYPTED PASSWORD 'sl_forms_password';
GRANT ALL PRIVILEGES ON DATABASE sl_forms TO sl_forms;
```

3. Применить SQL‑скрипт из раздела «Инициализация базы данных (создание таблиц)» к базе `sl_forms`.

#### 3. Запустить Redis

Просто запустить локальный сервер Redis (по умолчанию он слушает `localhost:6379`).

#### 4. Настроить переменные окружения

В корне проекта можно создать файл `.env` со значениями:

```dotenv
DATABASE_URL=postgres://sl_forms:sl_forms_password@localhost:5432/sl_forms?sslmode=disable
REDIS_ADDR=localhost:6379
JWT_SECRET=your_jwt_secret_here
```

Либо выставить эти переменные вручную в терминале перед запуском:

```bash
set DATABASE_URL=postgres://sl_forms:sl_forms_password@localhost:5432/sl_forms?sslmode=disable
set REDIS_ADDR=localhost:6379
set JWT_SECRET=your_jwt_secret_here
```

#### 5. Установить Go‑зависимости

В корне проекта:

```bash
go mod tidy
go mod download
```

#### 6. Запуск приложения

Из корня репозитория:

```bash
go run cmd/server/main.go
```

После этого API будет доступен на `http://localhost:8080`, а Swagger — по адресу `http://localhost:8080/swagger/index.html`.


### Работа с Postman-коллекцией

Postman-коллекция реализована .json формате.

1. Запустить Postman
2. Items -> три точки -> Import
3. Выбрать файл коллекции

### Запуск через Docker Compose

Требуется установленный Docker и Docker Compose.

1. Собрать и запустить все сервисы:

```bash
docker-compose up --build
```

2. После запуска будут подняты:
- **PostgreSQL**: хост `localhost`, порт `5432`, база `sl_forms`, пользователь `sl_forms`, пароль `sl_forms_password`.
- **Redis**: хост `localhost`, порт `6379`.
- **Приложение**: хост `localhost`, порт `8080`.

Переменные окружения для приложения внутри контейнера:
- `DATABASE_URL=postgres://sl_forms:sl_forms_password@db:5432/sl_forms?sslmode=disable`
- `REDIS_ADDR=redis:6379`
- `JWT_SECRET=supersecretjwt` (заменить на свой секрет).

### Инициализация базы данных (создание таблиц)

SQL-скрипт, который нужно выполнить в базе `sl_forms`, чтобы создать все необходимые для работы таблицы:

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

CREATE TABLE forms (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL
);

CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    form_id INT NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    position INT NOT NULL
);

CREATE TABLE options (
    id SERIAL PRIMARY KEY,
    question_id INT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    value TEXT NOT NULL,
    position INT NOT NULL
);

CREATE TABLE responses (
    id SERIAL PRIMARY KEY,
    form_id INT NOT NULL REFERENCES forms(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    response_id INT NOT NULL REFERENCES responses(id) ON DELETE CASCADE,
    question_id INT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    option_id INT REFERENCES options(id),
    text_value TEXT
);
```
