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

func StartApp(port string) {
	usersDb := database.ConnectDB()
	storage := Storage{}
	storage.UsersRepository = usersDb

	router := chi.NewRouter()

	SetHandlers(storage, router)
	http.ListenAndServe(":"+port, router)
}

func SetHandlers(storage Storage, router *chi.Mux) {
	router.Get("/user", storage.GetAllUsers)
	router.Get("/user/{id}", storage.GetUser)
	router.Post("/user", storage.AddUser)
	router.Get("/friends/{id}", storage.GetFriends)
	router.Post("/make_friends", storage.MakeFriends)
	router.Delete("/user/{id}", storage.DeleteUser)
	router.Put("/user/{id}", storage.UpdateUserAge)
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
	w.Write([]byte(fmt.Sprintf("???????????????????????? %s ????????????", deletedUserName.Name)))
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

	w.Write([]byte("?????????????? ???????????????????????? ?????????????? ????????????????"))
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
	msg := fmt.Sprintf("%s ?? %s ???????????? ????????????", user.Name, friend.Name)
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
