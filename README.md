# Go-Redis

## Обзор
Этот сервис, предназначен для хранения простой статистики в формате
ключ-значение с использованием REST интерфейса. Сервис реализует подмножество
команд Redis для работы с данными, где ключ — произвольная строка, а значение —
целое число.

## Функционал
Сервис поддерживает следующие команды:
- Создать запись ключ-значение
- Получить значение указанного ключа
- Удалить указанный ключ
- Увеличить/уменьшить значение указанного ключа на 1
- Увеличить/уменьшить значение указанного ключа на указанное число

API документация доступна в файле: [api.yml](api.yml)

## Технологии
- **Язык программирования**: Go
- **Хранение данных**: Redis 
- **Сервер**: Docker и Docker-compose для развертывания и управления контейнерами

## Установка и начало работы
### Предварительные требования
Установленный Docker

### Установка проекта
```bash
git clone https://github.com/skraio/go-redis-practice
cd go-redis-practice
```

### Запуск
```bash
docker-compose up --build
```

## Примеры использования API
#### Создание записи
```bash
curl -X POST -d '{"key":"abcd", "value":16}' http://localhost:8080/record
```
Тело ответа:
```json
{
    "message": "record successfully created"
}
```

#### Получение записи
```bash
curl -X GET http://localhost:8080/record/abcd
```
Тело ответа:
```json
{
    "record": {
        "key": "abcd",
        "value": 16
    }
}
```

#### Увеличение значения по ключу
```bash
curl -X PUT http://localhost:8080/record/abcd/increment
```
Тело ответа:
```json
{
    "message": "record incremented by 1"
}
```

#### Уменьшение значения по ключу на заданное число
```bash
curl -X PUT -d '{"term": 5}' http://localhost:8080/record/abcd/decrement-by
```
Тело ответа:
```json
{
    "message": "record decremented by 5"
}
```
Получение записи:
```json
{
    "record": {
        "key": "abcd",
        "value": 12
    }
}
```

#### Удаление записи
```bash
curl -X DELETE http://localhost:8080/record/abcd
```
Тело ответа:
```json
{
    "message": "record successfully deleted"
}
```
