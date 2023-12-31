# Тестовое задание. Яшкин Илья

Запуск проекта

```sh
docker-compose up
```

Запуск тестов

```sh
docker exec avito-backend-app-1 go test ./test/ -v
```

[Коллекция postman](https://github.com/IlyaYashkin/avito-backend/blob/master/Avito%20backend.postman_collection.json)

## Возникшие вопросы
#### Доп. задание 1

Получение отчета по **пользователю** за определенный период, на вход год-месяц.

Про передачу идентификатора пользователя не было сказано, но метод ```/get-user-segment-log``` дополнительно принимает ```user_id```.

#### Доп. задание 2

Хотел сделать очистку сегментов при помощи периодической горутины, но из-за

> В методе получения сегментов пользователя мы должны получить АКТУАЛЬНУЮ информацию о сегментах пользователя с задержкой не более 1 минуты после добавления сегмента.

решил, что лучше вызывать функцию очистки сегментов перед выполнением метода получения сегментов пользователя.

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

Формат ```ttl``` — ```2006-01-02T15:04:05Z```, RFC3336 без временной зоны

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
* ```date``` — дата в формате Год-Месяц (```2023-09```)

Если передать ```date```, то данные выбираются от выбранной даты включительно. День берется последний в выбранном месяце.

Query параметры можно не передавать, в таком случае скачается история сегментов всех пользователей.

#### Ответ
Происходит скачивание csv файла с данными.

Заголовок csv файла: ```user_id,segment_name,operation,operation_time```.
