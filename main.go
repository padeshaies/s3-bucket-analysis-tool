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
	"github.com/padeshaies/s3-bucket-analysis-tool/types"
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
	buckets := make([]types.Bucket, 0)
	output, err := client.ListBuckets(ctx, &s3.ListBucketsInput{
		Prefix: &filterSettings.BucketName,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, awsBucket := range output.Buckets {
		bucket := types.Bucket{
			Name:                   *awsBucket.Name,
			Region:                 *awsBucket.BucketRegion,
			CreationDate:           *awsBucket.CreationDate,
			ObjectNumber:           0,
			TotalSize:              0,
			MostRecentModifiedDate: time.Time{},
			Cost:                   0.0,
		}

		var output *s3.ListObjectsV2Output
		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket.Name),
		}

		objectPaginator := s3.NewListObjectsV2Paginator(client, input)
		for objectPaginator.HasMorePages() {
			output, err = objectPaginator.NextPage(ctx)
			if err != nil {
				log.Fatal(err)
			}

			for _, object := range output.Contents {
				bucket.ObjectNumber++
				bucket.TotalSize += int(*object.Size)
				if object.LastModified.After(bucket.MostRecentModifiedDate) {
					bucket.MostRecentModifiedDate = *object.LastModified
				}
			}
		}

		bucket.Cost = helpers.CalculateBucketCost(bucket.TotalSize)
		buckets = append(buckets, bucket)
	}

	for _, bucket := range buckets {
		bucket.Println(displaySettings)
	}
}

func buildDisplaySettings() (types.DisplaySettings, error) {
	result := types.DisplaySettings{
		FileSize: helpers.B,
		GroupBy:  "",
		Timezone: time.Local,
	}

	flags := os.Args[1:]

	if index := slices.Index(flags, "--file-size"); index != -1 {
		if len(flags) < index+2 {
			return result, fmt.Errorf("please provide a file size unit")
		}

		result.FileSize, _ = helpers.GetUnit(flags[index+1])
	}

	if index := slices.Index(flags, "--group-by"); index != -1 {
		if len(flags) < index+2 {
			return result, fmt.Errorf("please provide a group by option")
		}

		groupBy := flags[index+1]
		if groupBy != "region" && groupBy != "bucket" {
			return result, fmt.Errorf("invalid group by option. ilease use 'region' or 'bucket'")
		}
		result.GroupBy = groupBy
	}

	if index := slices.Index(flags, "--timezone"); index != -1 {
		if len(flags) < index+2 {
			return result, fmt.Errorf("please provide a timezone")
		}

		loc, err := time.LoadLocation(flags[index+1])
		if err != nil {
			return result, fmt.Errorf("invalid timezone")
		}
		result.Timezone = loc
	}

	return result, nil
}

func buildFilterSettings() (types.SearchFilters, error) {
	result := types.SearchFilters{
		BucketName:  "",
		StorageType: "",
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
			result.BucketName = filters[index+1]
		}

		if index := slices.Index(filters, "storage-type"); index != -1 {
			result.StorageType = filters[index+1]
		}
	}

	return result, nil
}
