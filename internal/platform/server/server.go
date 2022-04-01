package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"friendsbook/internal/platform/database"

	"github.com/go-chi/chi/v5"
)

type Storage struct {
	UsersRepository database.PostgressUsersStorage
}

func StartApp() {
	usersDb := database.ConnectDB()
	r := chi.NewRouter()
	storage := Storage{}
	storage.UsersRepository = usersDb
	r.Get("/user", storage.GetAllUsers)
	r.Get("/user/{id}", storage.GetUser)
	r.Post("/user", storage.AddUser)
	r.Get("/friends/{id}", storage.GetFriends)
	r.Post("/make_friends", storage.MakeFriends)
	r.Delete("/user/{id}", storage.DeleteUser)
	r.Put("/user/{id}", storage.UpdateUserAge)
	http.ListenAndServe(":8080", r)
}

func (s *Storage) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json charset=utf-8")
	allUsers, err := s.UsersRepository.GetAllUsers()
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewEncoder(w).Encode(allUsers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s Storage) AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json charset=utf-8")
	user := database.User{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s.UsersRepository.AddUser(user)
	w.WriteHeader(http.StatusCreated)

}

func (s Storage) GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json charset=utf-8")
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := database.User{}
	user, err = s.UsersRepository.GetUser(strconv.FormatInt(int64(id), 10))
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Storage) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	deletedUserName := database.User{}
	deletedUserName, err = s.UsersRepository.GetUser(strconv.FormatInt(int64(id), 10))

	s.UsersRepository.DeleteUser(strconv.FormatInt(int64(id), 10))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte(fmt.Sprintf("пользователь %s удален", deletedUserName.Name)))
}

func (s Storage) UpdateUserAge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json charset=utf-8")
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userAge := database.UserAge{}
	err = json.NewDecoder(r.Body).Decode(&userAge)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.UsersRepository.UpdateUserAge(strconv.FormatInt(int64(id), 10), userAge.Age)

	w.Write([]byte("возраст пользователя успешно обновлен"))
}

func (s Storage) MakeFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json charset=utf-8")
	friendsRequest := database.FriendRequest{}
	err := json.NewDecoder(r.Body).Decode(&friendsRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := s.UsersRepository.GetUser(strconv.Itoa(friendsRequest.SourceId))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	friend, _ := s.UsersRepository.GetUser(strconv.Itoa(friendsRequest.TargetId))
	if !Contains(user.Friends, friend.ID) {
		user.Friends = append(user.Friends, int64(friend.ID))
		s.UsersRepository.MakeFriends(user.ID, user.Friends)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	msg := fmt.Sprintf("%s и %s теперь друзья", user.Name, friend.Name)
	w.Write([]byte(msg))
}

func (s Storage) GetFriends(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json charset=utf-8")
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	friends := s.UsersRepository.GetFriends(int(id))
	err = json.NewEncoder(w).Encode(friends)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func Contains(friends []int64, userID int) bool {
	for _, value := range friends {
		if value == int64(userID) {
			return true
		}
	}
	return false
}
