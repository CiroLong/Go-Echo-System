package model

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//type Database struct {
//	users  map[string]User //算了，我拿username做键好了
//	maxInt uint            //自增ID, 不用填写
//}

var MongoDB *mongo.Database

func initDB() {
	// Set client options
	// url来自手动设置，后期考虑使用config包从外部文件导入
	var clientOptions = options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	// context 包是go并发编程中使用的用于同步 goroutine 的包，之后阅读《Go 语言设计与实现》学习
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Println("package model: error when init mongo client with uri mongodb://localhost:27017")
		log.Panic(err)
	}
	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Println("package model: unsuccessful ping to mongodb with uri mongodb://localhost:27017")
		log.Panic(err)
	}
	log.Println("Connected to MongoDB!")

	MongoDB = client.Database("echo-for-github")
}
