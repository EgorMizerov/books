# books

## API

```
GET /books        - Список книг
GET /books/:id    - Получить книгу
POST /books       - Добавить книгу
PUT /books/:id    - Изменить книгу
DELETE /books/:id - Удалить книгу
```

## Модели
book
```
id string
tite string
author *author
```
author
```
Firstname string
Lastname string
```
