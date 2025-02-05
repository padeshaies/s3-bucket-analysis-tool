package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/padeshaies/s3-bucket-analysis-tool/helpers"
)

func main() {

	displaySettings, err := buildDisplaySettings()
	if err != nil {
		log.Fatal(err)
	}

	filterSettings, err := buildFilterSettings()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client := s3.NewFromConfig(cfg)

	output, err := client.ListBuckets(ctx, &s3.ListBucketsInput{
		Prefix: &filterSettings.bucketName,
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

			bucketInfo.cost = helpers.CalculateBucketCost(bucketInfo.totalSize)

		}

		fmt.Printf("  - Number of files: %v\n", bucketInfo.fileNumber)
		fmt.Printf("  - Total size: %v\n", helpers.FormatFileSize(bucketInfo.totalSize, displaySettings.fileSize))
		fmt.Printf("  - Most recent modified date: %v\n", bucketInfo.mostRecentModifiedDate.In(loc))
		fmt.Printf("  - Cost: $%v\n", bucketInfo.cost)
	}
}

type DisplaySettings struct {
	fileSize int
	groupBy  string
}

func buildDisplaySettings() (DisplaySettings, error) {
	result := DisplaySettings{
		fileSize: helpers.B,
		groupBy:  "",
	}

	flags := os.Args[1:]

	if index := slices.Index(flags, "--file-size"); index != -1 {
		if len(flags) < index+2 {
			return result, fmt.Errorf("please provide a file size unit")
		}

		result.fileSize, _ = helpers.GetUnit(flags[index+1])
	}

	if index := slices.Index(flags, "--group-by"); index != -1 {
		if len(flags) < index+2 {
			return result, fmt.Errorf("please provide a group by option")
		}

		groupBy := flags[index+1]
		if groupBy != "region" && groupBy != "bucket" {
			return result, fmt.Errorf("invalid group by option. ilease use 'region' or 'bucket'")
		}
		result.groupBy = groupBy
	}

	return result, nil
}

type FilterSettings struct {
	bucketName  string
	storageType string
}

func buildFilterSettings() (FilterSettings, error) {
	result := FilterSettings{
		bucketName:  "",
		storageType: "",
	}

	flags := os.Args[1:]

	if index := slices.Index(flags, "--filters"); index != -1 {
		if len(flags) < index+2 {
			log.Fatal("Please provide a filter option")
		}

		filtersArguments := strings.Split(flags[index+1], ":")
		if len(filtersArguments) != 2 {
			log.Fatal("Invalid filter option. Please use 'bucket' or 'storage-type' and the value separated by a colon")
		}
		filters := strings.Split(filtersArguments[1], ";")

		if index := slices.Index(filters, "bucket"); index != -1 {
			result.bucketName = filters[index+1]
		}

		if index := slices.Index(filters, "storage-type"); index != -1 {
			result.storageType = filters[index+1]
		}
	}

	return result, nil
}
