# Cinema Platform - Auth Service

Мікросервіс на Go, який відповідає за безпеку, аутентифікацію, авторизацію та управління обліковими записами користувачів. Виступає як центральний вузол перевірки доступу (Identity Provider) для всієї платформи.

## 🚀 Технологічний стек

- **Мова:** Go 1.24
- **База даних:** PostgreSQL
- **Безпека:** JWT (JSON Web Tokens), `bcrypt` для хешування паролів
- **Роутинг:** Chi
- **Логування:** `log/slog`

## ✨ Основні функції

- Реєстрація та Логін: Безпечна автентифікація користувачів із видачею пари токенів (Access/Refresh).
- Керування сесіями (Refresh Token Rotation): Можливість оновлювати короткострокові токени доступу без повторного введення пароля, а також віддалено "розлогінювати" користувачів.
- RBAC (Role-Based Access Control): Розділення прав доступу. Видача ролей `user` та `admin`. Наприклад, лише `admin` може звертатися до `media-catalog` для додавання нових фільмів.
- Верифікація токенів: Надання внутрішнього API для інших мікросервісів, щоб вони могли швидко перевірити валідність токена користувача.

## 🛠 Локальний запуск

### Вимоги

- Go 1.24
- Піднята інфраструктура через репозиторій `cinema-infra` (Postgres + MinIO + RabbitMQ).

### Встановлення

1. Клонуйте репозиторій:

```bash
 git clone [https://github.com/***/cinema-auth.git](https://github.com/***/cinema-auth.git)
 cd cinema-auth
```

2. Завантажте Go-модулі:

```bash
go mod download
```

3. Створіть файл .env, або використовуйте docker-compose для встановлення env змінних:

```
PORT=8080
TIMEOUT=5
API_GATEWAY_URL=localhost:8080
DB_HOST=localhost
DB_NAME=cinema_db
DB_USER=user
DB_PASSWORD=password
SECRET_JWT_KEY=secret_jwt
```

4. Запустіть сервіс:

```bash
go run cmd/app/main.go
```

## 📁 Структура проєкту

Проєкт побудований за принципами чистої архітектури (Clean Architecture / Deep Coding) та DDD:

- cmd/ — точки входу в додаток.
- internal/app — збір усіх залежностей, та запуск.
- internal/delivery/rest/ — http-контролери запитів.
- internal/service/ — бізнес-логіка (робота з користувачами, токенами та іншим).
- internal/repository/ — робота з PostgreSQL та виконання складних SQL-запитів.
- internal/server — створення роутера, та запуск http-сервера.
- internal/domain — структури для роботи з базою даних, запитів, відповідей, помилками та іншим.
- internal/config — пакет для завантаження конфігів.
- pkg/storage/ — пакети для роботи з базою даних та інші.
