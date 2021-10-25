package model

type Database struct {
	users  map[string]User //算了，我拿username做键好了
	maxInt uint            //自增ID, 不用填写
}

var DB Database

func initDb() {
	DB.users = make(map[string]User)
	DB.maxInt = 1
}
