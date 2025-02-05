package types

import "time"

type DisplaySettings struct {
	FileSize int
	GroupBy  string
	Timezone *time.Location
}
