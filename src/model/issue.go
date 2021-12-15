package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type issueStatus int

const (
	IssueOpen issueStatus = iota
	IssueClose
)

//	现在issue的序列号?这东西还是绑定到user上好了
type UserIssueInfo struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	UserID       primitive.ObjectID `bson:"user_id" json:"user_id"`
	IssueNumbers int                `bson:"issue_numbers" json:"issue_numbers"`
}

type IssueModel struct {
	ID primitive.ObjectID `bson:"_id" json:"_id"`

	// Schrodinger
	ReposityID primitive.ObjectID `bson:"reposity_id" json:"reposity_id"`

	AutherID primitive.ObjectID `bson:"auther_id" json:"auther_id"`

	Number int         `bson:"number" json:"number"`
	Status issueStatus `bson:"status" json:"status"`

	Title string `bson:"title" json:"title"` //issue只有一个title

	CommentNums int            `bson:"comment_nums" json:"comment_nums"`
	Comments    []IssueComment `bson:"comments" json:"comments"`
}

type IssueComment struct {
	OwnerID primitive.ObjectID `bson:"owner_id" json:"owner_id"`
	Number  int                `bson:"number" json:"number"`
	Body    string             `bson:"body" json:"body"`
}
