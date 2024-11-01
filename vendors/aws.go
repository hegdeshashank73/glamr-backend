package vendors

import (
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/hegdeshashank73/glamr-backend/utils"
	"github.com/spf13/viper"
)

var S3Client *s3.S3
var SESClient *ses.SES
var SQSClient *sqs.SQS
var DDBClient *dynamodb.DynamoDB
var SecretManager *secretsmanager.SecretsManager

func initAWS() {
	st := time.Now()
	defer utils.LogTimeTaken("init.initAWS", st)

	var awsSess *session.Session
	var awsConfig = aws.Config{
		Region: aws.String(viper.GetString("AWS_REGION")),
	}
	if _, ok := os.LookupEnv("AWS_EXECUTION_ENV"); ok {
		awsSess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config:            awsConfig,
		}))
	} else if user, ok := os.LookupEnv("USER"); ok && (user == "ubuntu" || user == "root") {
		awsSess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config:            awsConfig,
		}))
	} else {
		awsSess = session.Must(session.NewSessionWithOptions(session.Options{
			Profile: "duggup",
			Config:  awsConfig,
		}))
		awsConfig.Credentials = stscreds.NewCredentials(awsSess, viper.GetString("AWS_MONO_ROLE"))
	}

	S3Client = s3.New(awsSess, &awsConfig)
	SESClient = ses.New(awsSess, &awsConfig)
	SQSClient = sqs.New(awsSess, &awsConfig)
	SecretManager = secretsmanager.New(awsSess, &awsConfig)
	DDBClient = dynamodb.New(awsSess, &awsConfig)
}
