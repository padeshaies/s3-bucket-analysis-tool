package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/padeshaies/s3-bucket-analysis-tool/helpers"
	"github.com/padeshaies/s3-bucket-analysis-tool/types"
)

var cfg aws.Config

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
	bucketList := &types.SafeBucketList{
		Buckets: &[]*types.Bucket{},
		Lock:    sync.Mutex{},
	}

	bucketPaginator := s3.NewListBucketsPaginator(client, &s3.ListBucketsInput{
		Prefix: aws.String(filterSettings.BucketName),
	})

	var tasks sync.WaitGroup
	for bucketPaginator.HasMorePages() {
		output, err := bucketPaginator.NextPage(ctx)
		if err != nil {
			log.Fatal(err)
		}

		tasks.Add(1)
		go analyzeBucketPage(output, client, ctx, bucketList, &tasks, filterSettings)
	}

	tasks.Wait()

	for _, bucket := range *bucketList.Buckets {
		bucket.Println(displaySettings)
	}
}

func analyzeBucketPage(page *s3.ListBucketsOutput, client *s3.Client, ctx context.Context, bucketList *types.SafeBucketList, tasks *sync.WaitGroup, filterSettings types.SearchFilters) {
	for _, awsBucket := range page.Buckets {
		bucket := types.Bucket{
			Name:                   *awsBucket.Name,
			Region:                 *awsBucket.BucketRegion,
			CreationDate:           *awsBucket.CreationDate,
			ObjectNumber:           0,
			TotalSize:              0,
			MostRecentModifiedDate: time.Time{},
			Cost:                   0.0,
			Lock:                   sync.Mutex{},
		}

		input := &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket.Name),
		}

		// Adjust the client to the bucket region if necessary
		// TODO - Fix this
		/* var regionClient *s3.Client
		if client.Options().Region != bucket.Region {
			newCfg := cfg.Copy()
			newCfg.Region = bucket.Region
			regionClient = s3.NewFromConfig(newCfg)
		} else {
			regionClient = client
		}

		fmt.Println("Searching region " + regionClient.Options().Region + " for bucket " + bucket.Name) */

		objectPaginator := s3.NewListObjectsV2Paginator(client, input)

		var tasks sync.WaitGroup
		for objectPaginator.HasMorePages() {
			output, err := objectPaginator.NextPage(ctx)
			if err != nil {
				log.Fatal(err)
			}

			tasks.Add(1)
			go analyzeBucketObjectPage(output, &bucket, &tasks, filterSettings)
		}

		tasks.Wait()

		if filterSettings.StorageType == "" || (filterSettings.StorageType != "" && bucket.ObjectNumber > 0) {
			bucket.Cost = helpers.CalculateBucketCost(bucket.TotalSize)

			bucketList.Lock.Lock()
			*bucketList.Buckets = append(*bucketList.Buckets, &bucket)
			bucketList.Lock.Unlock()
		}
	}

	tasks.Done()
}

func analyzeBucketObjectPage(page *s3.ListObjectsV2Output, bucket *types.Bucket, tasks *sync.WaitGroup, filterSettings types.SearchFilters) {
	for _, object := range page.Contents {
		bucket.Lock.Lock()

		// Apply storage type filter
		if filterSettings.StorageType != "" && string(object.StorageClass) != filterSettings.StorageType {
			bucket.Lock.Unlock()
			continue
		}

		bucket.ObjectNumber++
		bucket.TotalSize += int(*object.Size)
		if object.LastModified.After(bucket.MostRecentModifiedDate) {
			bucket.MostRecentModifiedDate = *object.LastModified
		}

		if !slices.Contains(bucket.StorageTypes, string(object.StorageClass)) {
			bucket.StorageTypes = append(bucket.StorageTypes, string(object.StorageClass))
		}

		bucket.Lock.Unlock()
	}

	tasks.Done()
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

		filtersArgument := strings.Split(flags[index+1], ";")

		for _, filter := range filtersArgument {
			keyValue := strings.Split(filter, ":")

			if len(keyValue) != 2 {
				log.Fatal("Invalid filter option. Please use a key and a value separated by a colon")
			}

			key := keyValue[0]
			if key != "bucket" && key != "storage-type" {
				log.Fatal("Invalid filter option. Please use 'bucket' or 'storage-type'")
			}

			if key == "bucket" {
				result.BucketName = keyValue[1]
			}

			if key == "storage-type" {
				storageType := strings.ToUpper(keyValue[1])

				if (storageType != "STANDARD") &&
					(storageType != "REDUCED_REDUNDANCY") &&
					(storageType != "GLACIER") &&
					(storageType != "STANDARD_IA") &&
					(storageType != "ONEZONE_IA") &&
					(storageType != "INTELLIGENT_TIERING") &&
					(storageType != "DEEP_ARCHIVE") &&
					(storageType != "OUTPOSTS") &&
					(storageType != "GLACIER_IR") &&
					(storageType != "SNOW") &&
					(storageType != "EXPRESS_ONEZONE") {
					return result, fmt.Errorf("invalid storage type")
				}

				result.StorageType = storageType
			}
		}
	}

	return result, nil
}
