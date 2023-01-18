# Bored Dao Bot


- ✨Magic ✨ bot for magic guys✨


## Installation without docker

1. Clone repo
2. run `go mod tidy`
3. create .env and write env
4. Do not forget to migrate db `db/migrations`
5. run `go build github.com/arandich/telegram-dao/cmd/main`

Configure environments in .env for local setup without docker

```sh
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=
POSTGRES_HOST=
POSTGRES_PORT=
TOKEN=
```

## Development

Handle functions for commands in `internal/commands`

Write functions for db in `internal/database/query`

To create/update db struct go to `internal/database/entity`

## Docker setup:

Write environment in docker-compose

```sh
      #db service
      POSTGRES_PASSWORD=
      POSTGRES_USER=
      POSTGRES_DB=
      
      #bot service
      TOKEN: telegram api token
      POSTGRES_USER: Example posgtres
      POSTGRES_PASSWORD: Example admin 
      POSTGRES_DB: Example database
      POSTGRES_HOST: *db service name
      POSTGRES_PORT: Example 5432
      SSLMODE: Example disable
      
      #migrate service
      command:
      [ "-path", "/database", "-database",  "postgres://pgusername:pgpassword@dbservicename:5432/pgdbname?sslmode=disable", "up" ]
      
      #pgadmin service
      PGADMIN_DEFAULT_EMAIL: Example admin@gmail.com
      PGADMIN_DEFAULT_PASSWORD: Example admin
```
Run in terminal `docker composer up -d`

## Implemented commands
```sh
        "/start" - Приветственное слово + вывод информации о пользователе
	"инфо" - Вывод информации о пользователе
	"голосования" - Вывод информации о текущих голосованиях
	"мои голосования" - Вывод информации о голосованиях в которых участвовал пользователь
	"активности" - Вывод списка активностей где пользователь еще не участвовал
	"/send user amount" - Перевод токенов пользователю
	"/addEvent" - [Администратор] Создать активность с наградой ARG: date YYYY-MM-DD, reward
	"/events" - [Администратор] Вывод журнала всех активностей и кнопки принять/отклонить
	"/transactions" - [Администратор] Вывод журнала всех транзакций токенов, первые 10
	"/createVoting" - [Администратор] *Название* *Ссылка* *Дата типа YYYY-MM-DD* *Вариант1* *Вариант2* *Вариант3*
	"участники" - Доступ[Член клуба] Вывод список участников и их роли
	"мои активности" - Вывод информации о активностях пользователя
```

