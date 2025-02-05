package helpers

import (
	"fmt"
	"strings"
)

const (
	B = iota
	KB
	MB
	GB
	TB
)

func FormatFileSize(fileSize, unit int) string {
	switch unit {
	case B:
		return fmt.Sprintf("%d bytes", fileSize)
	case KB:
		return fmt.Sprintf("%d KB", fileSize/1024)
	case MB:
		return fmt.Sprintf("%d MB", fileSize/1024/1024)
	case GB:
		return fmt.Sprintf("%d GB", fileSize/1024/1024/1024)
	case TB:
		return fmt.Sprintf("%d TB", fileSize/1024/1024/1024/1024)
	default:
		return fmt.Sprintf("%d bytes", fileSize)
	}
}

func GetUnit(unit string) (int, error) {
	unit = strings.ToUpper(unit)

	switch unit {
	case "B":
		return B, nil
	case "KB":
		return KB, nil
	case "MB":
		return MB, nil
	case "GB":
		return GB, nil
	case "TB":
		return TB, nil
	default:
		return B, fmt.Errorf("invalid unit %s", unit)
	}
}
