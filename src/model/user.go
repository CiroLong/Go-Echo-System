package model

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
)

const colNameUser = "userModel"
const secretCost = bcrypt.DefaultCost // 	bcrypt包加密用

var colUser *mongo.Collection

func initModelUser() {
	colUser = MongoDB.Collection(colNameUser)

	//?? 这里插入成功了
	user := User{
		ID: primitive.NewObjectID(),
		Username: "asdasd",
	}
	_, err  := colUser.InsertOne(context.Background(), user)
	if err != nil {log.Println(err.Error())}
}

//	姓名、个人说明、邮箱、网站、头像
// 	使用*string的原因
// 	If you want to split null and "", you should use *string instead of string.
type User struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	Username     string             `bson:"username" json:"username"`
	PasswordHash string             `bson:"password" json:"password"`
	Statement    string             `bson:"statement" json:"statement"`
	Email        string             `bson:"email" json:"email"`
	Image        *string            `bson:"image" json:"image"`

	//Phone        string             `bson:"phone" json:"phone"`
	IsAdmin  bool `bson:"is_admin" json:"is_admin"` // 管理
	Verified bool `bson:"verified" json:"verified"` // 已验证
}

func (u *User) HashPassword(plain string) (string, error) {
	if len(plain) == 0 {
		return "", errors.New("password should not be empty")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(plain), secretCost) //加密存储
	return string(h), err
}

func (u *User) CheckPassword(plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plain))
	return err == nil
}

// 	以下为直接使用的ks 的代码
//	看懂了=
// 	有需要在添加其他操作函数

func GetUserWithID(idHex string) (User, bool, error) {
	var user User

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return user, false, err
	}

	err = colUser.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return user, false, nil
	}
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

func GetUserWithUsername(username string) (User, bool, error) {
	var user User
	err := colUser.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return user, false, nil
	}
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

func GetAllUsers() ([]User, error) {
	var users []User
	result, err := colUser.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	err = result.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func IsUserAdmin(idHex string) (bool, error) {
	user, found, err := GetUserWithID(idHex)
	if !found {
		return false, errors.New("user with _id " + idHex + " not found")
	}
	if err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}

func AddUser(user User) (string, error) {
	user.ID = primitive.NewObjectID()
	_, err := colUser.InsertOne(context.Background(), user)
	if err != nil {
		return "", err
	}
	return user.ID.Hex(), nil
}

func UpdateUser(idHex string, info bson.M) error {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": info,
	}

	_, err = colUser.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	//第二个参数指定更新对象，第三个参数指定更新内容
	return err
}

func DeleteUser(idHex string) error {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}

	_, err = colUser.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
