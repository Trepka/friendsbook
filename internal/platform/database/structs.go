package database

import "github.com/jmoiron/sqlx"

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"user_name"`
	Age     int     `json:"user_age"`
	Friends []int64 `json:"user_friends"`
}

type UserAge struct {
	Age int `json:"new_user_age"`
}

type FriendRequest struct {
	SourceId int `json:"source_id"`
	TargetId int `json:"target_id"`
}

type PostgressUsersStorage struct {
	db *sqlx.DB
}

type UsersDb map[int]User
