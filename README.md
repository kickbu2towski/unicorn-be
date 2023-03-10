## Running locally
> NOTE: All of the commands below assume your working directory is set to the root folder.

You need to set the env variable `UNICORN_DB_DSN`. This variable contains connection parameters which is needed to establish connection to your PostgreSQL database.

It is usually of this form `password://username:password@host:port/dbname?param1=true&param2=false`

For sending emails using your Gmail account, you need to configure these env variables:

- `UNICORN_MAIL_SENDER` (eg: John Doe)
- `UNICORN_MAIL_USERNAME` (eg: jonhdoe@gmail.com)
- `UNICORN_MAIL_PASSWORD` 

The env variable `UNICORN_MAIL_PASSWORD` is app password for your Gmail account. If you don't know how to get this, watch [this video](https://youtu.be/L9TbZxpykLQ?t=564).

Once you set these variables, run the SQL migrations. For this you need to install the [migrate tool](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) in your system.

Once installed, up all the migrations:

```bash
migrate -source file://path-to-migrations -database $UNICORN_DB_DSN up
```

Once the migrations ran successfully, you can start the app:

```bash
go run ./cmd/api
```

You can configure other runtime variables when starting the app. Use the help command to know about these variables:

```bash
go run ./cmd/api --help
```