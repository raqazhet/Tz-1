# Todo List Microservice

Этот микросервис предоставляет RESTful API для управления списком задач (Todo List).

## Требования

- Go (версия >= 1.17)
- Docker
- Docker Compose
- Git
## Установка и Запуск приложения

1. Клонируйте репозиторий проекта на свой компьютер:
```shell
git clone https://github.com/raqazhet/Tz-1.git
```
2. Перейдите в директорию проекта:
```shell
cd ./Tz-1
```
3. Запустите приложение с помощью Docker Compose:
```shell
docker-compose up
```
4. Чтобы остановить приложение, выполните команду:
```shell
docker compose down
```


### API Endpoints
#### ** Формат обмена данными JSON.**
#### Swagger документация доступна по адресу http://localhost:7777/swagger/index.html
## Создание задачи

1. Метод: POST
- URL: /api/v1/users/todo-list/tasks
- Тело запроса:

```json
{
   "title": "Купить книгу",
   "activeAt": "2023-08-04"
}
```
- Создание новой задачи

## Обновление задачи

2. Метод: PUT
- URL: /api/v1/users/todo-list/tasks/:id- Тело запроса:
```json
{
   "title": "Купить книгу - Высоконагруженные приложения",
   "activeAt": "2023-08-05"
}
```
- Обновление существующей задачи.

## Удаление задачи

3. Метод: DELETE
- URL: /api/v1/users/todo-list/tasks/:id
   - Удаление задачи.


## Пометить задачу выполненной

4. Метод: PUT
- URL: /api/v1/users/todo-list/tasks/:id/done
  -  Помечает задачу как выполненную.

## Список задач

5. Метод: GET
- URL: /api/v1/users/todo-list/tasks
  -  Получает список задач.


##  Тестирование
Запуск unit тестов
```shell
go test -cover ./...
```
