package model

import "errors"

type User struct {
	ID       uint   `json:"_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	IsAdmin  bool   `json:"is_admin"`
	Verified bool   `json:"verified"`
}

func initUserModel() {

}

func GetUserWithUsername(userName string) (User, bool, error) {
	if u, ok := DB.users[userName]; ok {
		return u, true, nil
	} else {
		return u, false, errors.New("no such user")
	}
}

func AddUser(user User) (uint, error) {
	user.ID = DB.maxInt
	DB.maxInt++
	DB.users[user.Username] = user
	return user.ID, nil
}

func UpdateUser(user User) error {
	if _, ok := DB.users[user.Username]; ok {
		DB.users[user.Username] = user
		return nil
	}
	return errors.New("No such user")
}

func DeleteUser(userName string) error {
	delete(DB.users, userName)
	return nil
}
