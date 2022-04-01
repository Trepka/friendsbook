package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

const (
	host       = "localhost"
	port       = 5432
	user       = "postgres"
	dbpassword = "postgres"
	dbname     = "friendsbook"
)

var schema = `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name TEXT,
	age INT,
	friends INT[] NOT NULL DEFAULT '{}'::INT[]
  );
`

func ConnectDB() PostgressUsersStorage {
	var storage PostgressUsersStorage
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", user, dbname, dbpassword)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Database connection failed")
	}
	db.MustExec(schema)
	storage.db = db
	return storage
}

func (p *PostgressUsersStorage) GetAllUsers() (UsersDb, error) {
	user := User{}
	allUsers := make(UsersDb)
	rows, err := p.db.Queryx("SELECT * FROM users")
	if err != nil {
		log.Fatal("can't read rows")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Age, pq.Array(&user.Friends))
		if err != nil {
			log.Fatalf("can't scan row %v", err)
		}

		allUsers[user.ID] = user
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return allUsers, nil
}

func (p *PostgressUsersStorage) GetUser(id string) (User, error) {
	var u User
	row := p.db.QueryRowx("SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&u.ID, &u.Name, &u.Age, pq.Array(&u.Friends))
	if err != nil {
		log.Fatalf("can't scan row %v", err)
	}
	return u, nil
}

func (p PostgressUsersStorage) AddUser(u User) {
	tx := p.db.MustBegin()
	tx.MustExec("INSERT INTO users (name, age, friends) VALUES ($1, $2, $3)", u.Name, u.Age, pq.Array(u.Friends))
	tx.Commit()
}

func (p *PostgressUsersStorage) DeleteUser(id string) {
	_, err := p.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Printf("can't delete user %v", err)
	}
}

func (p PostgressUsersStorage) UpdateUserAge(id string, age int) {
	_, err := p.db.Exec(`UPDATE users SET age = $2 WHERE id = $1`, id, age)
	if err != nil {
		log.Println("can't update age")
	}
}

func (p PostgressUsersStorage) GetFriends(id int) []int64 {
	var u User
	row := p.db.QueryRowx(`SELECT friends FROM users WHERE id = $1`, id)
	err := row.Scan(pq.Array(&u.Friends))
	if err != nil {
		log.Fatal("can't read row")
	}
	return u.Friends
}

func (p PostgressUsersStorage) MakeFriends(id int, friendsList []int64) {
	_, err := p.db.Exec(`UPDATE users SET friends = $2 WHERE id = $1`, id, pq.Array(friendsList))
	if err != nil {
		log.Println("can't update friends")
	}
}
