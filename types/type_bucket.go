package types

import (
	"fmt"
	"sync"
	"time"

	"github.com/padeshaies/s3-bucket-analysis-tool/helpers"
)

type Bucket struct {
	Name                   string
	Region                 string
	StorageTypes           []string
	CreationDate           time.Time
	MostRecentModifiedDate time.Time
	ObjectsNumber          map[string]int
	ObjectsSize            map[string]int

	Lock sync.Mutex
}

func (b *Bucket) TotalSize() int {
	totalSize := 0
	for _, size := range b.ObjectsSize {
		totalSize += size
	}
	return totalSize
}

func (b *Bucket) TotalObjectNumber() int {
	totalObjectNumber := 0
	for _, number := range b.ObjectsNumber {
		totalObjectNumber += number
	}
	return totalObjectNumber
}

func (b *Bucket) TotalCost() (float64, error) {
	totalCost := 0.0
	for _, storageType := range b.StorageTypes {
		cost, err := helpers.CalculateObjectsCostByStorageType(storageType, b.Region, b.TotalSize(), b.TotalObjectNumber())
		if err != nil {
			return 0.0, err
		}
		totalCost += cost
	}
	return totalCost, nil
}

func (b *Bucket) Println(displaySettings DisplaySettings) {
	fmt.Printf("Name: %v\n", b.Name)
	fmt.Printf("  - Region: %v\n", b.Region)
	fmt.Printf("  - CreationDate: %v\n", b.CreationDate.In(displaySettings.Timezone))
	fmt.Printf("  - Number of files: %v\n", b.TotalObjectNumber())
	fmt.Printf("  - Total size: %v\n", helpers.FormatFileSize(b.TotalSize(), displaySettings.FileSize))
	fmt.Printf("  - Most recent modified date: %v\n", b.MostRecentModifiedDate.In(displaySettings.Timezone))
	fmt.Printf("  - Storage types: %v\n", b.StorageTypes)

	totalCost, err := b.TotalCost()
	if err != nil {
		panic(err)
	}
	fmt.Printf("  - Cost: $%v per month (only for storage)\n", totalCost)
}
