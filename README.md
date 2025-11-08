# Интеллектуальный обучающий веб-модуль-чат- бот «Виртуальный преподаватель по охране труда»

Проект состоит из 2 сервисов:
- ai-server
- tg-bot

## Запуск ai-server

Создать и активировать виртуальное окружение:

```shell
python3 -m venv venv
source venv/bin/activate
```

Установить зависимости:

```shell
pip install -r  ./cmd/ai-server/requirements.txt
```

Запустить сервер:

```shell
python3 ./cmd/ai-server/main.py
```

## Запуск tg-bot

```shell
go run cmd/tg-bot/main.go
```