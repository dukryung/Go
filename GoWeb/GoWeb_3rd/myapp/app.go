package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//User struct
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var userMap map[int]*User
var lastID int

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")

}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprint(w, "Get UserInfo by /users/{id}")

	if len(userMap) == 0 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No Users")
		return
	}

	users := []*User{}
	for _, u := range userMap {
		users = append(users, u)
	}
	data, _ := json.Marshal(users)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user, ok := userMap[id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id:", id)
		return

	}
	w.Header().Add("Content-type", "application/json")
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))

}
func CreateUsersHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	lastID++
	user.ID = lastID
	user.CreatedAt = time.Now()
	userMap[user.ID] = user
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))

}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	_, ok := userMap[id]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID:", id)
		return
	}

	delete(userMap, id)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted User ID:", id)

}

func UpdateUsersHandler(w http.ResponseWriter, r *http.Request) {
	updateUser := new(User)
	err := json.NewDecoder(r.Body).Decode(updateUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	user, ok := userMap[updateUser.ID]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID:", updateUser.ID)
		return
	}

	if updateUser.FirstName != "" {
		user.FirstName = updateUser.FirstName
	}
	if updateUser.LastName != "" {
		user.LastName = updateUser.LastName
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))
}

func NewHandler() http.Handler {
	userMap = make(map[int]*User)
	lastID = 0
	mux := mux.NewRouter()

	mux.HandleFunc("/", IndexHandler)
	mux.HandleFunc("/users", UsersHandler).Methods("GET")
	mux.HandleFunc("/users", CreateUsersHandler).Methods("POST")
	mux.HandleFunc("/users", UpdateUsersHandler).Methods("PUT")
	mux.HandleFunc("/users/{id:[0-9]+}", GetUserInfoHandler).Methods("GET")
	mux.HandleFunc("/users/{id:[0-9]+}", deleteUserHandler).Methods("DELETE")
	return mux
}
