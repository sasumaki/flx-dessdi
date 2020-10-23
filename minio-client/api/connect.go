package api

import (
	"github.com/minio/minio-go"
)

// GetClient returns minio client
func GetClient(endpoint string, accessKeyID string, secretKey string) (*minio.Client, error) {

	useSSL := false
	minioClient, err := minio.New(endpoint, accessKeyID, secretKey, useSSL)
	if err != nil {
		return minioClient, err
	}

	return minioClient, nil
}
func Main() {

}

// TestConnection test connection to bugget
func TestConnection(endpoint string, accessKeyID string, secretKey string) (bool, error) {
	client, err := GetClient(endpoint, accessKeyID, secretKey)

	if err != nil {
		return false, err
	}

	_, err = client.ListBuckets()
	if err != nil {
		return false, err
	}
	return true, nil
}
