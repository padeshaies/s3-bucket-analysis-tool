package helpers

import (
	"fmt"
	"testing"
)

func TestFormatFileSize(t *testing.T) {
	cases := []struct {
		input struct {
			fileSize int
			unit     int
		}
		expected string
	}{
		{ // Format into bytes
			input: struct {
				fileSize int
				unit     int
			}{
				fileSize: 1,
				unit:     B,
			},
			expected: "1 bytes",
		}, // Format into kilobytes
		{
			input: struct {
				fileSize int
				unit     int
			}{
				fileSize: 1024,
				unit:     KB,
			},
			expected: "1.00 KB",
		},
		{ // Format into megabytes
			input: struct {
				fileSize int
				unit     int
			}{
				fileSize: 1572864,
				unit:     MB,
			},
			expected: "1.50 MB",
		},
		{ // Format into gigabytes
			input: struct {
				fileSize int
				unit     int
			}{
				fileSize: 5647881994,
				unit:     GB,
			},
			expected: "5.26 GB",
		}, // Format into terabytes
		{
			input: struct {
				fileSize int
				unit     int
			}{
				fileSize: 123456789012387,
				unit:     TB,
			},
			expected: "112.28 TB",
		}, // Format into unknown unit should default to bytes
		{
			input: struct {
				fileSize int
				unit     int
			}{
				fileSize: 95,
				unit:     10,
			},
			expected: "95 bytes",
		},
	}

	for _, c := range cases {
		got := FormatFileSize(c.input.fileSize, c.input.unit)
		if got != c.expected {
			t.Errorf("FormatFileSize(%d, %d) == %s, want %s", c.input.fileSize, c.input.unit, got, c.expected)
		}
	}
}

func TestGetUnit(t *testing.T) {
	cases := []struct {
		input    string
		expected struct {
			unit int
			err  error
		}
	}{
		{
			input: "B",
			expected: struct {
				unit int
				err  error
			}{
				unit: B,
				err:  nil,
			},
		},
		{
			input: "KB",
			expected: struct {
				unit int
				err  error
			}{
				unit: KB,
				err:  nil,
			},
		},
		{
			input: "MB",
			expected: struct {
				unit int
				err  error
			}{
				unit: MB,
				err:  nil,
			},
		},
		{
			input: "GB",
			expected: struct {
				unit int
				err  error
			}{
				unit: GB,
				err:  nil,
			},
		},
		{
			input: "TB",
			expected: struct {
				unit int
				err  error
			}{
				unit: TB,
				err:  nil,
			},
		},
		{
			input: "allo maman!",
			expected: struct {
				unit int
				err  error
			}{
				unit: B,
				err:  fmt.Errorf("invalid unit ALLO MAMAN!"),
			},
		},
	}

	for _, c := range cases {
		got, err := GetUnit(c.input)
		if err != nil && err.Error() != c.expected.err.Error() {
			t.Errorf("GetUnit(%s) returned an error: %s", c.input, err)
		}
		if got != c.expected.unit {
			t.Errorf("GetUnit(%s) == %d, want %d", c.input, got, c.expected)
		}
	}
}
