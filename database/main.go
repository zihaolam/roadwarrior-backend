package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
	"github.com/zihaolam/roadwarrior-backend/config"
)

var DB = NewDBConnection()
var AWSSession, _ = NewAWSSession()

func NewAWSSession() (*session.Session, error) {
	if config.ENV == "dev" {
		return session.NewSessionWithOptions(session.Options{Profile: "localzone-project"})
	}
	conf := aws.NewConfig().WithRegion("ap-southeast-1")
	return session.Must(session.NewSession(conf)), nil
}

func NewDBConnection() *dynamo.DB {
	if config.ENV == "dev" {
		return dynamo.New(AWSSession, &aws.Config{Region: aws.String("ap-southeast-1"), Endpoint: aws.String("http://localhost:5392")})
	}
	return dynamo.New(AWSSession, &aws.Config{Region: aws.String("ap-southeast-1")})
}

func GetTable(tableName string) dynamo.Table {
	return DB.Table(tableName)
}

func CreateTable(tableName string, from interface{}) error {
	return DB.CreateTable(tableName, from).Run()
}
