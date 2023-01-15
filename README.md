# Bored Dao Bot


- ✨Magic ✨ bot for magic guys✨


## Installation

1. Clone repo
2. run `go mod tidy`
3. create .env and write env
4. Do not forget to migrate db `db/migrations`
5. run `go build github.com/arandich/telegram-dao/cmd/main`

Configure environments in .env

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

## User struct:

```sh
        Id        int
	Username  string
	RoleId    int
	Karma     int
	Tokens    int
	CreatedAt time.Time
```

### Event struct
```sh
        Id     int
	Name   string
	Date   time.Time
	UserId int
	Status string
```

## Implemented commands
```sh
        "/start" - Приветственное слово + вывод информации о пользователе
	"/инфо" - Вывод информации о пользователе
	"/участники" - Доступ[Член клуба] Вывод список участников и их роли
	"/мои_активности" - Вывод информации о активностях пользователя
```

