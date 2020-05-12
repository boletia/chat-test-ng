package dynamo

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	guuid "github.com/google/uuid"
)

var (
	awsRegion   = "us-east-1"
	dynamoTable = "streaming-chat-cognito-test"
)

// DyDB dynamo db handler
type DyDB struct {
	table string
	DB    *dynamodb.DynamoDB
}

// New dydb worker
func New() DyDB {
	sess := session.Must(
		session.NewSession(&aws.Config{
			Region: &awsRegion},
		))
	svc := dynamodb.New(sess)

	return DyDB{
		DB:    svc,
		table: dynamoTable,
	}
}

// Write send messages to dynamo
func (dy DyDB) Write(msg []byte) error {
	var newItem map[string]*dynamodb.AttributeValue
	var objmap map[string]json.RawMessage
	var err error

	if err := json.Unmarshal(msg, &objmap); err != nil {
		return err
	}

	message := struct {
		NickName       string `json:"nickname"`
		Message        string `json:"message"`
		EventSubdomain string `json:"event_subdomain"`
		Avatar         string `json:"avatar"`
		ID             string `json:"id,omitempty"`
	}{}

	if err := json.Unmarshal(objmap["data"], &message); err != nil {
		return err
	}

	id := guuid.New()
	message.ID = id.String()

	if newItem, err = dynamodbattribute.MarshalMap(message); err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      newItem,
		TableName: aws.String(dy.table),
	}

	_, err = dy.DB.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}
