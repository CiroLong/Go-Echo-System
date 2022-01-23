package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// 链接数据库
const colNameUserIssue = "userIssueModel"

var colUserIssue *mongo.Collection

func initModelUserIssue() {
	colUserIssue = MongoDB.Collection(colNameUserIssue)
}

// 链接数据库

type issueStatus int

const (
	IssueOpen issueStatus = iota
	IssueClose
)

type UserIssueInfo struct { // 这个信息也放入该数据库中
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	IssueNumbers int                `bson:"issue_numbers" json:"issue_numbers"`
}

type IssueModel struct { // 直接存入数据库
	ID primitive.ObjectID `bson:"_id" json:"_id"`

	// Schrodinger
	// ReposityID primitive.ObjectID `bson:"reposity_id" json:"reposity_id"`

	AutherID primitive.ObjectID `bson:"auther_id" json:"auther_id"` // 发起issue的

	Title  string      `bson:"title" json:"title"`   // issue只有一个title
	Number int         `bson:"number" json:"number"` // issue编号
	Status issueStatus `bson:"status" json:"status"` // 状态: open or close

	CommentNums int            `bson:"comment_nums" json:"comment_nums"` // 评论数
	Comments    []IssueComment `bson:"comments" json:"comments"`         // 评论内容
}

type IssueComment struct {
	OwnerID primitive.ObjectID `bson:"owner_id" json:"owner_id"` // 评论者
	Number  int                `bson:"number" json:"number"`     // 编号
	Body    string             `bson:"body" json:"body"`
}
