package helpers

import "testing"

func TestCalculateBucketCost(t *testing.T) {
	cases := []struct {
		input    int
		expected float64
	}{
		{
			input:    1,
			expected: 0.00,
		},
		{
			input:    46917849267, // before the rounds, it should be ~1.005, so 1.01 after the round
			expected: 1.01,
		},
	}

	for _, c := range cases {
		got := CalculateBucketCost(c.input)
		if got != c.expected {
			t.Errorf("CalculateBucketCost(%d) == %f, want %f", c.input, got, c.expected)
		}
	}
}
