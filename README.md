# RabbitMQ Chat Application

Реализация сетевого чата на RabbitMQ с CLI и GUI клиентами

## Features
- Поддержка двух типов клиентов: консольный (CLI) и графический (GUI)
- Автоматическое создание каналов при переключении
- Без сохранения истории сообщений
- Работа через RabbitMQ Exchange (topic)
- Кроссплатформенная GUI версия на Fyne
- Docker-сборка для всех компонентов

## Требования
- Go 1.21+
- Docker и docker-compose (для сборки контейнеров)
- X11 сервер (для GUI версии в Linux)

## Быстрый старт
Описан в Docker

## Структура проекта
├── cmd/
│   ├── cli-client/     # Консольный клиент
│   └── gui-client/     # Графический клиент
├── internal/
│   └── rabbitmq/       # Общая логика работы с RabbitMQ
├── docker-compose.yml  # Конфигурация окружения
├── Dockerfile.cli      # Сборка CLI клиента
├── Dockerfile.gui      # Сборка GUI клиента
├── go.sum              # Зависимости Go
└── go.mod              # Зависимости Go


