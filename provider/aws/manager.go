package aws

import (
	"path"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Manager exposes the interfaces about resource CRUD on aws
type Manager struct {
	lambdaClient *lambda.Lambda
	s3Client *s3.S3

	account string
	lambdaRole string

	userCodeBucket string
	awsCodeBucket string
	
	addonsDir string
}

// Config defines the neccessary settings about aws.Manager
type Config struct {
	RegionID string `yaml:"regionID"`
	AccessKeyID string `yaml:"accessKeyID"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Account string `yaml:"account"`
	LambdaRole string `yaml:"lambdaRole"`

	UserCodeBucket string `yaml:"userCodeBucket"`
	AwsCodeBucket string `yaml:"awsCodeBucket"`
}

// NewManagerWithConfig create a manager with given configuration
func NewManagerWithConfig(config Config) (*Manager, error){
	credentail := credentials.NewStaticCredentials(config.AccessKeyID, config.AccessKeySecret, "")
	awsConfig := aws.NewConfig().WithRegion(config.RegionID).WithEndpointDiscovery(true)
	awsConfig = awsConfig.WithCredentials(credentail)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: *awsConfig,
	}))
	lambdaClient := lambda.New(sess)
	s3Client := s3.New(sess)

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		lambdaClient: lambdaClient,
		s3Client: s3Client,
		userCodeBucket: config.UserCodeBucket,
		awsCodeBucket: config.AwsCodeBucket,
		addonsDir: path.Join(home, ".jfManager", "aws"),
		account: config.Account,
		lambdaRole: config.LambdaRole,
	}
	return manager, nil
}