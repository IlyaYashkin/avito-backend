# Тестовое задание. Яшкин Илья

Запуск проекта

```sh
docker-compose up
```

Запуск тестов

```sh
docker exec avito-backend-app-1 go test ./test/ -v
```

---

## Добавление сегмента ```POST /create-segment```
#### Тело запроса
```json
{
    "name": "AVITO_VOICE_40",
    "user_percentage": 30
}
```
```user_percentage``` — процент пользователей, попадающих в сегмент. Необязательный аргумент.
#### Пример ответа
```json
{
    "data": {
        "message": "Segment created",
        "name": "AVITO_VOICE_40"
    },
    "status": "success"
}
```
---

## Удаление сегмента ```DELETE /delete-segment```
#### Тело запроса
```json
{
    "name": "AVITO_VOICE_40",
}
```
#### Пример ответа
```json
{
    "data": {
        "message": "Segment deleted",
        "name": "AVITO_VOICE_40"
    },
    "status": "success"
}
```
---

## Обновление сегментов пользователя ```POST /update-user-segments```
#### Тело запроса
```json
{
    "user_id": 1000,
    "add_segments": [
        "AVITO_VOICE_30",
        "AVITO_VOICE_40",
        {
            "segment": "AVITO_VOICE_50",
            "ttl": "2023-08-29T16:18:17Z"
        }
    ],
    "delete_segments": [
        "AVITO_VOICE_10",
        "AVITO_VOICE_20"
    ]
}
```
Массив ```add_segments``` должен состоять из названий сегментов (строка), либо из списков ключ-значение, где ```segment``` — название сегмента, ```ttl``` — время автоматического удаления пользователя из сегмента.

Формат ```ttl``` — ```2023-08-29T16:18:17Z```

Массив ```delete_segments``` должен состоять из названий сегментов.

#### Пример ответа
```json
{
    "data": {
        "added_percentage_segments": [
            "AVITO_VOICE_70",
            "AVITO_VOICE_90",
            "AVITO_VOICE_100",
        ],
        "added_segments": [
            "AVITO_VOICE_30",
            "AVITO_VOICE_40"
        ],
        "added_ttl_segments": [
            "AVITO_VOICE_50"
        ],
        "deleted_segments": [
            "AVITO_VOICE_10",
            "AVITO_VOICE_20"
        ],
        "message": "User segments updated",
        "user_id": 1000
    },
    "status": "success"
}
```

Массив ```added_percentage_segments``` заполняется теми сегментами, в которые пользователь попал автоматически. Автоматически в сегмент может попасть только новый пользователь.

---

## Получение сегментов пользователя
```GET /get-user-segments/:user_id```

URL-параметр ```user_id``` — идентификатор пользователя

#### Пример ответа
```json
{
    "data": {
        "segments": [
            "AVITO_VOICE_30",
            "AVITO_VOICE_40"
        ],
        "user": 1000
    },
    "status": "success"
}
```
---

## Получение истории добавления/удаления сегментов у пользователя
```GET /get-user-segment-log?user_id=$1&date=$2```

Query параметры
* ```user_id``` — идентификатор пользователя
* ```date``` — дата в виде Год-Месяц (2023-09)

Если передать ```date```, то данные выбираются от выбранной даты включительно. День берется последний в выбранном месяце

#### Ответ
Происходит скачивание csv файла с данными
