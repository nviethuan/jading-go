package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadToS3(key string) {
	bucketName := "minedb"

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-southeast-1"),
	)

	if err != nil {
		log.Fatalf("failed to load SDK configuration, %v", err)
	}
	client := s3.NewFromConfig(cfg)

	home, err := os.UserHomeDir()
	if err != nil {
		panic("failed to get user home directory")
	}

	dbPathAccount := home + "/app/data/" + key

	accountDB, err := os.ReadFile(dbPathAccount)
	if err != nil {
		log.Fatalf("failed to read accounts.db, %v", err)
	}

	params := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   strings.NewReader(string(accountDB)),
	}

	_, err = client.PutObject(context.TODO(), params)
	if err != nil {
		log.Fatalf("failed to put object, %v", err)
	}
}

func main() {
	key := ""
	flag.StringVar(&key, "key", "", "key")
	flag.Parse()
	UploadToS3(key)
}
