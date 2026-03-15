## SL Forms

Бэкенд сервиса форм на Go (Gin, PostgreSQL, Redis).

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
