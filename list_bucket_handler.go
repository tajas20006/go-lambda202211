// mockgen用のおまじない
//go:generate mockgen -source=$GOFILE -package=$GOPACKAGE -destination=mock_$GOFILE

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	S3BUCKET = "go-lambda202211"
)

// テストがしやすくなるように使用する関数を持ったinterfaceを定義する
// 今回はmockgenを使ってモックを生成する
type s3API interface {
	ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

type Object struct {
	Key   string
	Bytes int64
}

// listObjects は指定されたbucketのオブジェクトをリストアップし、そのkeyとsize(bytes)を返す
func listObjects(ctx context.Context, svc s3API, bucket string) (objects []Object, err error) {
	objects = []Object{}

	result, err := svc.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return
	}

	for _, object := range result.Contents {
		objects = append(objects, Object{
			Key:   *object.Key,
			Bytes: object.Size,
		})
	}
	return
}

// Handlerが実際の処理開始部分
func Handler(ctx context.Context) {
	// s3のサービスを取得
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-west-2"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	svc := s3.NewFromConfig(cfg)

	// s3バケットの中身を見てみる(このLambdaのソースを置いてある)
	objects, err := listObjects(ctx, svc, S3BUCKET)
	if err != nil {
		log.Fatalf("failed to list objects, %v", err)
	}

	fmt.Println("Objects:")
	for _, object := range objects {
		fmt.Println(object.Key, ": ", object.Bytes, "bytes")
	}
}

func main() {
	// lambdaを使うときのおまじない
	lambda.Start(Handler)
}
