# VOVAN_AI 🤖

**VOVAN_AI** — это AI-ассистент и агент для кодинга, работающий через OpenAI-совместимый API. Поддерживает два режима работы: обычный чат и режим агента с возможностью выполнения bash-команд.

### Настройка

Создайте файл `.env` на основе `.env.example`:

```bash
cp .env.example .env
```

Отредактируйте `.env`:

```env
BASE_URL = "https://api.openai.com/v1"      # URL API (можно заменить на любой OpenAI-совместимый)
API_KEY = "ваш-api-ключ"                    # API ключ
MODEL_NAME = "gpt-4o-mini"                  # Название модели
```

### Запуск

**Режим чата (по умолчанию):**

```bash
go run main.go
# или явно:
go run main.go -mode chat
```

**Режим агента (с выполнением команд):**

```bash
go run main.go -mode agent
```