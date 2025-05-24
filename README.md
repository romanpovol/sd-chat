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
- docker-compose up --build
- docker build -t chat-cli -f Dockerfile.cli .
- docker build -t chat-gui -f Dockerfile.gui .
