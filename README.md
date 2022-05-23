# Bookslib

RestAPI app simulate minimum functionality of social net. Supports get, post, delete queries.

## Features

* **/GET/user** - return all users instances from storage
* **/GET/user/{id}** - return user instances from storage by ID
* **/POST/user** - add new user to the storage. ID of new user autoincrement by Postgres. Example of json to make new user see below

```json
{
    "user_name":"Boris",
    "user_age":22,
    "user_friends": [1,2]
}
```

* **/GET/friends/{id}** - return all users friends.
* **/POST/make_friends** - add friend to user list if friends(array of ID). Source_id - ID user who make request to friend, target_id - new user friend. Example of json to add new friend

```json
{
    "source_id": 5,
    "target_id": 6
}
```

* **/DELETE/user/{id}** - delete user from the storage by ID.
* **/PUT/user/{id}** - change user age. Example of json to change user age

```json
{"new_user_age":40}
```

## Used in project

* router based on [go-chi](https://github.com/go-chi/chi) library.
* [sqlx](https://github.com/jmoiron/sqlx) library for connect to database.
* database - PostgreSQL.
