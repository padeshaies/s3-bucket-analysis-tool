package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	prefix := ""
	output, err := client.ListBuckets(ctx, &s3.ListBucketsInput{
		Prefix: &prefix,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, bucket := range output.Buckets {
		fmt.Printf("Name: %v\n", *bucket.Name)

		// bug: invalid memory address or nil pointer dereference
		fmt.Printf("  - Region: %v\n", *bucket.BucketRegion)

		loc, _ := time.LoadLocation("Local")
		fmt.Printf("  - CreationDate: %v\n", bucket.CreationDate.In(loc))

		var output *s3.ListObjectsV2Output
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(*bucket.Name),
		}

		bucketInfo := struct {
			fileNumber             int
			totalSize              int
			mostRecentModifiedDate time.Time
			cost                   float64
		}{}

		objectPaginator := s3.NewListObjectsV2Paginator(client, input)
		for objectPaginator.HasMorePages() {
			output, err = objectPaginator.NextPage(ctx)
			if err != nil {
				log.Fatal(err)
			}

			for _, object := range output.Contents {
				bucketInfo.fileNumber++
				bucketInfo.totalSize += int(*object.Size)
				if object.LastModified.After(bucketInfo.mostRecentModifiedDate) {
					bucketInfo.mostRecentModifiedDate = *object.LastModified
				}
			}

			bucketInfo.cost = CalculateBucketCost(bucketInfo.totalSize)

		}

		fmt.Printf("  - Number of files: %v\n", bucketInfo.fileNumber)
		fmt.Printf("  - Total size: %v bytes\n", bucketInfo.totalSize)
		fmt.Printf("  - Most recent modified date: %v\n", bucketInfo.mostRecentModifiedDate.In(loc))
		fmt.Printf("  - Cost: $%v\n", bucketInfo.cost)
	}
}
