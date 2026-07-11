# Contributing to Hairdress Arzamas 52 / Вклад в проект

First off, thank you for considering contributing! We welcome all kinds of contributions — bug reports, feature requests, documentation improvements, and code changes.

Прежде всего, спасибо, что хотите внести вклад! Мы принимаем любые contributions: баг-репорты, предложения новых фич, улучшения документации и изменения кода.

---

## 📋 Table of Contents / Содержание

- [Code of Conduct / Кодекс поведения](#code-of-conduct--кодекс-поведения)
- [Issues / Сообщения о проблемах](#issues--сообщения-о-проблемах)
- [Development Workflow / Процесс разработки](#development-workflow--процесс-разработки)
- [Commit Convention / Соглашение о коммитах](#commit-convention--соглашение-о-коммитах)
- [Code Style / Стиль кода](#code-style--стиль-кода)
- [Project Structure / Структура проекта](#project-structure--структура-проекта)

---

## Code of Conduct / Кодекс поведения

### EN
By participating, you agree to maintain a respectful, inclusive, and harassment-free environment. Be constructive, be kind.

### RU
Участвуя в проекте, вы соглашаетесь поддерживать уважительную, инклюзивную атмосферу без домогательств. Будьте конструктивны и доброжелательны.

---

## Issues / Сообщения о проблемах

### EN
- **Bug reports**: use the "Bug report" template. Include steps to reproduce, expected vs actual behavior, environment details (OS, Go version, etc.).
- **Feature requests**: use the "Feature request" template. Describe the problem you're solving, not just the solution.
- **Questions**: use GitHub Discussions.

Before creating an issue, search existing issues to avoid duplicates.

### RU
- **Баг-репорты**: используйте шаблон "Bug report". Опишите шаги воспроизведения, ожидаемое и фактическое поведение, окружение (ОС, версия Go и т.д.).
- **Предложения фич**: используйте шаблон "Feature request". Опишите проблему, которую решаете, а не только решение.
- **Вопросы**: используйте GitHub Discussions.

Перед созданием issue поищите среди существующих, чтобы избежать дубликатов.

---

## Development Workflow / Процесс разработки

### EN

1. **Fork** the repository.
2. **Clone** your fork:
   ```bash
   git clone https://github.com/<your-username>/hairdress_arz_52_backend.git
   ```
3. **Create a feature branch**:
   ```bash
   git checkout -b feature/my-feature
   ```
   Branch naming:
   - `feature/` — new features
   - `fix/` — bug fixes
   - `refactor/` — code refactoring
   - `docs/` — documentation changes
   - `chore/` — build, CI, dependencies
4. **Make your changes**. Follow the [Code Style](#code-style--стиль-кода).
5. **Run checks locally**:
   ```bash
   make lint    # golangci-lint
   make test    # go test -race + coverage
   make vet     # go vet ./...
   ```
6. **Commit** following the [commit convention](#commit-convention--соглашение-о-коммитах).
7. **Push** to your fork:
   ```bash
   git push origin feature/my-feature
   ```
8. Open a **Pull Request** against the `main` branch.
   - PR title should mirror the commit convention.
   - Link related issues in the description using `Closes #123`.
   - Keep PRs focused — one feature/fix per PR.
9. Wait for **CI checks** to pass. Address any review feedback.

> **Important**: After changing proto contracts in `hairdress_arz_52_contracts`, push those changes to GitHub first.

### RU

1. **Форкните** репозиторий.
2. **Клонируйте** свой форк:
   ```bash
   git clone https://github.com/<ваш-username>/hairdress_arz_52_backend.git
   ```
3. **Создайте ветку для фичи**:
   ```bash
   git checkout -b feature/my-feature
   ```
   Именование веток:
   - `feature/` — новые фичи
   - `fix/` — исправления багов
   - `refactor/` — рефакторинг
   - `docs/` — документация
   - `chore/` — сборка, CI, зависимости
4. **Внесите изменения**. Следуйте [стилю кода](#code-style--стиль-кода).
5. **Запустите проверки локально**:
   ```bash
   make lint    # golangci-lint
   make test    # go test -race + coverage
   make vet     # go vet ./...
   ```
6. **Закоммитьте** согласно [соглашению о коммитах](#commit-convention--соглашение-о-коммитах).
7. **Запушьте** в свой форк:
   ```bash
   git push origin feature/my-feature
   ```
8. Откройте **Pull Request** в ветку `main`.
   - Заголовок PR должен соответствовать соглашению о коммитах.
   - Ссылайтесь на связанные issues через `Closes #123`.
   - Держите PR сфокусированными — одна фича/исправление на один PR.
9. Дождитесь прохождения **CI-проверок**. Учтите замечания ревью.

> **Важно**: После изменения proto-контрактов в `hairdress_arz_52_contracts` не забудьте запушить изменения на GitHub.

---

## Commit Convention / Соглашение о коммитах

### EN
We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short description>

[optional body]

[optional footer]
```

**Types**:
- `feat` — new feature
- `fix` — bug fix
- `refactor` — code change without feature or fix
- `docs` — documentation only
- `test` — adding or fixing tests
- `chore` — build, CI, dependencies, etc.
- `perf` — performance improvements
- `style` — formatting, missing semicolons, etc. (not CSS)

**Scope** (optional): e.g. `auth`, `bookings`, `api`, `db`, `docker`

**Examples**:
```
feat(auth): add phone-based OTP verification
fix(bookings): prevent double-booking on overlapping timeslots
docs: update API endpoint list in README
```

### RU
Мы используем [Conventional Commits](https://www.conventionalcommits.org/):

```
<тип>(<область>): <краткое описание>

[необязательное тело]

[необязательный footer]
```

**Типы**:
- `feat` — новая фича
- `fix` — исправление бага
- `refactor` — изменение кода без новой фичи или исправления
- `docs` — только документация
- `test` — добавление или исправление тестов
- `chore` — сборка, CI, зависимости и т.д.
- `perf` — улучшение производительности
- `style` — форматирование, пропущенные точки с запятой и т.п. (не CSS)

**Область** (опционально): например, `auth`, `bookings`, `api`, `db`, `docker`

**Примеры**:
```
feat(auth): добавить OTP-верификацию по телефону
fix(bookings): предотвратить двойную запись на пересекающиеся слоты
docs: обновить список эндпоинтов в README
```

---

## Code Style / Стиль кода

### Go

### EN
- Format code with `gofumpt` (stricter `gofmt`).
- Run `golangci-lint` before committing (`make lint`).
- Follow standard Go conventions: [Effective Go](https://go.dev/doc/effective_go), [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments).
- Error handling: always check errors; use `fmt.Errorf("context: %w", err)` for wrapping.
- Naming: camelCase for variables, PascalCase for exports, acronyms all-caps (`HTTP`, `URL`, `DB`).
- Imports: standard library → third-party → internal, separated by blank lines.
- Use `any` instead of `interface{}`.
- Use `make` for slice/map initialization when length is known.
- Avoid global variables; use dependency injection via `internal/app/app.go`.

### RU
- Форматируйте код через `gofumpt` (более строгая версия `gofmt`).
- Запускайте `golangci-lint` перед коммитом (`make lint`).
- Следуйте стандартным Go-конвенциям: [Effective Go](https://go.dev/doc/effective_go), [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments).
- Обработка ошибок: всегда проверяйте ошибки; используйте `fmt.Errorf("context: %w", err)` для обёртки.
- Именование: camelCase для переменных, PascalCase для экспортируемого, аббревиатуры — в верхнем регистре (`HTTP`, `URL`, `DB`).
- Импорты: стандартная библиотека → сторонние → внутренние, разделены пустыми строками.
- Используйте `any` вместо `interface{}`.
- Используйте `make` для инициализации slice/map, если длина известна.
- Избегайте глобальных переменных; используйте dependency injection через `internal/app/app.go`.

### Python (Admin Panel)

### EN
- Format code with `ruff` (our linter and formatter).
- Follow [PEP 8](https://peps.python.org/pep-0008/) and [PEP 484](https://peps.python.org/pep-0484/) (type hints).
- Use async endpoints (FastAPI) with `async def`.
- Models use SQLAlchemy declarative style; schemas use Pydantic v2.
- Naming: snake_case for variables/functions, PascalCase for classes, UPPER_CASE for constants.

### RU
- Форматируйте код через `ruff` (наш линтер и форматтер).
- Следуйте [PEP 8](https://peps.python.org/pep-0008/) и [PEP 484](https://peps.python.org/pep-0484/) (type hints).
- Используйте асинхронные эндпоинты (FastAPI) с `async def`.
- Модели — SQLAlchemy declarative style; схемы — Pydantic v2.
- Именование: snake_case для переменных/функций, PascalCase для классов, UPPER_CASE для констант.

### General / Общее

### EN
- Write meaningful comments for non-trivial logic. Keep comments up to date.
- Do not leave commented-out code — delete it.
- For SQL: use lowercase keywords (`select`, `from`, `where`) to match the sqlc-generated style.
- For migrations: always provide both `up` and `down` scripts.
- Keep functions small and focused (single responsibility).

### RU
- Пишите осмысленные комментарии для нетривиальной логики. Поддерживайте комментарии в актуальном состоянии.
- Не оставляйте закомментированный код — удаляйте его.
- Для SQL: используйте lowercase для ключевых слов (`select`, `from`, `where`) в стиле sqlc.
- Для миграций: всегда предоставляйте `up` и `down` скрипты.
- Держите функции небольшими и сфокусированными (single responsibility).

---

## Project Structure / Структура проекта

### EN
Refer to the [README.md](../README.md) for the full project structure. Key directories:

| Directory | Description |
|-----------|-------------|
| `cmd/` | Application entry points |
| `internal/` | Private Go packages (domain, service, repository, transport) |
| `pkg/` | Public libraries (verify, mailer) |
| `migrations/` | Database migrations (Go-migrate) |
| `app/` | Python FastAPI admin panel |
| `db/query/` | SQL queries for sqlc code generation |
| `config/` | YAML configuration files |

### RU
Полная структура проекта описана в [README.md](../README.md). Ключевые директории:

| Директория | Описание |
|-----------|----------|
| `cmd/` | Точки входа приложения |
| `internal/` | Приватные Go-пакеты (domain, service, repository, transport) |
| `pkg/` | Публичные библиотеки (verify, mailer) |
| `migrations/` | Миграции БД (Go-migrate) |
| `app/` | Python FastAPI админ-панель |
| `db/query/` | SQL-запросы для генерации sqlc |
| `config/` | YAML-конфигурации |

---

## Questions? / Вопросы?

### EN
If you have questions, feel free to open a [Discussion](https://github.com/nhassl3/hairdress_arz_52_backend/discussions) or ask in the relevant issue.

### RU
Если у вас есть вопросы, открывайте [Discussion](https://github.com/nhassl3/hairdress_arz_52_backend/discussions) или спрашивайте в соответствующем issue.
