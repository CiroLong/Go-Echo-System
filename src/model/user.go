package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

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

func (u *User) HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost) //加密存储
	return string(h), err
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plain))
	return err == nil
}

func GetUserWithUsername(userName string) (User, bool, error) {
	if u, ok := DB.users[userName]; ok {
		return u, true, nil
	} else {
		return u, false, nil
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
