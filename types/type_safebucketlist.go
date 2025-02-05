package types

import "sync"

type SafeBucketList struct {
	Buckets *[]*Bucket
	Lock    sync.Mutex
}
