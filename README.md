# Congratulations

## Описание
Это тестовое задание от компании Rutube. Цель - сделать удобный сервис для уведомлений о приближении дня рождения.
Приложение реализовано в кастичной микросервисной аръитектуре. Каждый сервис может быть запущен абсолютно автономно, но
для этого требуется внести небольшие изменения. Всего в приложении 4 микросервиса - "Authorization", "Employees", "Notifications"
и "Subdcribe". Как главную sql базу данных я использовал PostgreSQL. Так же сервис "Employees" использует Redis. 
Для синхронизации между Redis и Postgres я использовал Kafka. Все приложение спроектировано и реализовано при помощи Docker. 
Главный вход в приложение -Rest API. Для описания api я использовал swagger. Все API работают на протоколе HTTPS симметричный ключ лежит в каждом микросервисе в папке /keys в файле symmetric-key.bin. Для билда изображений в докер для всех микросервисов я использую кеш, поэтому долго будет только в первый раз, но есть вам надо будет что то подправить, то дальше будет все только быстрее и быстрее.

## Инструменты

1. Docker version 27.1.1
2. go version go1.22.5
3. GNU Make 4.3 (optional)

## Устновка и запуск приложения

1. Клонируйте репозиторий:  
    ```bash
    git clone https://github.com/l1qwie/Congratulations.git
    
2. Перейдите в главную директория приложения:
    ```bash
    cd Congratulations

3. Создайте docker-network для приложения:
    ```bash
    make net

    или

    ```bash
    docker network create congratulations

## WARNING
Прежде чем начать, определитесь, что вам нужно:

1. Вы хотите запустить все Rest API
2. Вы хотите запустить сервисы по отдельности

## Запуск всего приложения (всех Rest APIs)
1. С самого начала требуется запулить все нужные images из docker:
    ```bash
    make launch-all

    или

    docker compose -f docker-compose.yml up -d

2. После того, как Postgres, Kafka и Redis запустилиь, теперь нужно сделать image приложения
    ```bash
    make build

    или

    docker build . -t congratulations

3. Вы почти уже все сделали! Теперь осталось все запустить:
    ```bash
    make up

    или

    docker run --rm --name notif --network congratulations congratulations /app/bin

Поздравляю! Приложение работает!

## Запуск сервисов по отдельности с тестами
1. С самого начала требуется запулить все нужные images из docker. Это запустит бд без дополнительных данных (для тестов):
    ```bash
    make launch-test

    или

    docker compose -f docker-compose-test.yml up -d

2. Выберите один из сервисов "Authorization", "Employees", "Notifications", "Subdcribe" и перейдите в директорию с таким же названием.
    ```bash
    cd {Название микросервиса}

3. Все микросервисы запускаются одинаково: 
- для создания image
    ```bash
    make build
- для запуска:
    ```bash
    make up

## P.S. 
Если по какой-то неведомой причине вам требуется остановить или перезаписать Postgres, Kafka или Redis то в главной директории (Congratulations)
есть команда:
    ```bash
    make delete
