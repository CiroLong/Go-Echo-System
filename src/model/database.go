package model

import (
	"errors"
	"fmt"
)

type Database struct {
	users  map[uint]User
	maxInt uint //自增ID, 不用填写
}

var DB Database

func initDb() {
	DB.users = make(map[uint]User)
	DB.maxInt = 1
}

func (d *Database) Insert(u *User) error {
	if _, ok := d.users[u.ID]; ok {
		return errors.New("the user exits, please use Update to update")
	}
	d.users[d.maxInt] = *u
	u.ID = d.maxInt
	d.maxInt++
	return nil
}

// 不太好写， 那就先查找到对应的User才能改
func (d *Database) Update(u *User) error {
	_, ok := d.users[u.ID]
	if !ok {
		return errors.New("no such user")
	}
	d.users[u.ID] = *u
	return nil
}

func (d *Database) Delete(user *User) error {
	delete(d.users, user.ID)
	return nil
}

func (d *Database) Find(id uint) (User, error) {
	if user, ok := d.users[id]; ok {
		return user, nil
	} else {
		return user, errors.New(fmt.Sprint("can't find user with id ", id))
	}
}
