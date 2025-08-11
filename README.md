# 🐊 Gator

A blog aggregator in Golang. Guided by [Boot.dev](https://www.boot.dev/lessons/dca1352a-7600-4d1d-bfdf-f9d741282e55).

# Prerequisites

- A working installation of PostgreSQL v15 or later. Refer to the [documentation](https://www.postgresql.org/docs/).
- A working installation of Golang. See https://go.dev/doc/install.

# How To Use

Clone the repository and head into the root directory of the project.

```bash
git clone https://github.com/uncomfyhalomacro/gator
cd gator
```

## Setup PostgreSQL

Start the Postgres server in the background
  - Mac: brew services start postgresql@15
  - Linux: sudo service postgresql start

Connect to the server. I recommend simply using the psql client. It's the "default" client for Postgres, and it's a great way to interact with the database. While it's not as user-friendly as a GUI like PGAdmin, it's a great tool to be able to do at least basic operations with.

Enter the psql shell:
  - Mac: psql postgres
  - Linux: sudo -u postgres psql

You should see a new prompt that looks like this:

```
postgres=#
```

Create a new database.

```
CREATE DATABASE gator;
```

Connect to the new database:

```
\c gator
```

You should see a new prompt that looks like this:

```
gator=#
```

Set the user password (Linux only)

```
ALTER USER postgres PASSWORD 'postgres';
```

You can type `exit` to leave the `psql` shell.

Get your connection string. A connection string is just a URL with all of
the information needed to connect to a database. The format is:

```
protocol://username:password@host:port/database
```

## Setup the migrations

Install goose

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

While inside your local copy of the repository, head to `sql/schema/` directory
and run the command

```bash
goose postgres "postgres://postgres:@localhost:5432/gator" up
```

This should set up all the tables for the `gator` database.

## Setup a config file

Create a config file called `.gatorconfig.json` at your `$HOME` directory.

The JSON file should have this structure (when prettified):

```json
{
	"db_url": "postgres://postgres:@localhost:5432/gator?sslmode=disable"
}
```

The `db_url` is your connection string, appended with `?sslmode=disable`.

## Build the gator command

```bash
go build
./gator register yourUsername  # You are automagically logged in
./gator addfeed "Hacker News RSS" "https://hnrss.org/newest"
```

## Run the aggregator in the background

```bash
./gator agg
```

In another terminal, you can check and browse your followed feeds.

```bash
./gator browse
./gator browse 10
```

Checkout the other subcommands below.

## Subcommands

* `following` checks the list of followed RSS feeds for the currently logged in user
* `login` logins a registered user
* `reset` resets the users database. WARNING -> destructive operation
* `agg` collects posts for the registered user
* `addfeed` adds a feed for the user. the args should be in the order `'title' 'url'`. follows the feed for the user afterwards
* `follow` follows a feed URL for the currently logged in user if it exists in the database
* `unfollow` unfollows a feed URL for the currently logged in user if it exists in the database
* `register` registers a new user. fails if user exists
* `users` lists the registered users
* `browse` browse recent posts for the currently logged in user. recieves one additional parameter, a number to limit the number of posts shown
* `feeds` list of feed URLs added by the users

