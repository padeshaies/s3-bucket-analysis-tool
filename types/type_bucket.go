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
	StorageType            string
	CreationDate           time.Time
	ObjectNumber           int
	TotalSize              int
	MostRecentModifiedDate time.Time
	Cost                   float64

	Lock sync.Mutex
}

func (b *Bucket) Println(displaySettings DisplaySettings) {
	fmt.Printf("Name: %v\n", b.Name)
	fmt.Printf("  - Region: %v\n", b.Region)
	fmt.Printf("  - CreationDate: %v\n", b.CreationDate.In(displaySettings.Timezone))
	fmt.Printf("  - Number of files: %v\n", b.ObjectNumber)
	fmt.Printf("  - Total size: %v\n", helpers.FormatFileSize(b.TotalSize, displaySettings.FileSize))
	fmt.Printf("  - Most recent modified date: %v\n", b.MostRecentModifiedDate.In(displaySettings.Timezone))
	fmt.Printf("  - Cost: $%v\n", b.Cost)
}
