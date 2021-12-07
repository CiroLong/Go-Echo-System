package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const colNameUserImage = "userImageModel"

var colUserImage *mongo.Collection

func initModelUserImage() {
	colUserImage = MongoDB.Collection(colNameUserImage)
}

type Image struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	UserId   primitive.ObjectID `bson:"user_id" json:"user_id"`
	Filename string             `bson:"filename" json:"filename"`
	//Format   string             `bson:"format" json:"format"`
}

func AddImage(image Image) (string, error) {
	image.ID = primitive.NewObjectID()
	_, err := colUserImage.InsertOne(context.Background(), image)
	if err != nil {
		return "", err
	}
	return image.ID.Hex(), nil
}

func UpdateImage(idHex string, info bson.M) error {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}
	update := bson.M{
		"$set": info,
	}
	_, err = colUserImage.UpdateOne(context.Background(), bson.M{"_id": id}, update)
	return err
}

func GetImageWithIdHex(idHex string) (Image, bool, error) {
	var image Image

	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return image, false, err
	}
	err = colUserImage.FindOne(context.Background(), bson.M{"_id": id}).Decode(&image)
	if err == mongo.ErrNoDocuments {
		return image, false, nil
	}
	if err != nil {
		return image, false, err
	}
	return image, true, nil
}

func GetImageWithFilename(filename string) (Image, bool, error) {
	var image Image
	err := colUserImage.FindOne(context.Background(), bson.M{"filename": filename}).Decode(&image)
	if err == mongo.ErrNoDocuments {
		return image, false, nil
	}
	if err != nil {
		return image, false, err
	}
	return image, true, nil
}

func GetImageWithUserId(userId primitive.ObjectID) (Image, bool, error) {
	var image Image
	err := colUserImage.FindOne(context.Background(), bson.M{"user_id": userId}).Decode(&image)
	if err == mongo.ErrNoDocuments {
		return image, false, nil
	}
	if err != nil {
		return image, false, err
	}
	return image, true, nil
}

func DeleteImage(idHex string) error {
	id, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}

	_, err = colUserImage.DeleteOne(context.Background(), bson.M{"_id": id})
	return err
}
