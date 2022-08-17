package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	tmpFrom = "./testdata/tmpfrom.txt"
	tmpTo   = "./testdata/tmpto.txt"
)

func TestCopy(t *testing.T) {
	cases := []struct {
		name     string
		text     string
		offset   int64
		limit    int64
		expected string
	}{
		{"base test", "11111222223333344444", 0, 0, "11111222223333344444"},
		{"test with offset", "11111222223333344444", 5, 0, "222223333344444"},
		{"test with limit", "11111222223333344444", 0, 5, "11111"},
		{"test with both", "11111222223333344444", 5, 5, "22222"},
	}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			err := os.WriteFile(tmpFrom, []byte(tcase.text), 0o755)
			require.NoError(t, err)

			err = Copy(
				tmpFrom,
				tmpTo,
				tcase.offset,
				tcase.limit,
			)
			require.NoError(t, err)

			b, err := os.ReadFile(tmpTo)
			require.NoError(t, err)

			require.Equal(t, tcase.expected, string(b))
		})
	}

	err := os.Remove(tmpFrom)
	require.NoError(t, err)

	err = os.Remove(tmpTo)
	require.NoError(t, err)
}
