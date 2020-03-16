package aws

import (
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
)

func (m *Manager) GetCodeURI(funcName string) (string, error) {
	return "", nil
}

func (m *Manager) GetImage(funcName string) (string, error) {
	return "", nil
}

func (m *Manager) SaveCode(funcName string, file string) (error) {
	body, err := os.Open(file)
	if err != nil {
		return err
	}
	_, err = m.s3Client.PutObject(&s3.PutObjectInput{
		Bucket: &m.userCodeBucket,
		Key: &funcName,
		Body: body,
	})
	if err != nil {
		return err
	}
	return nil
}