package inmem

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKey(t *testing.T) {
	require.Equal(t, "inmemory", Key)
}
