package setiteration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractDate(t *testing.T) {
	t.Run("returns yyyy-mm-dd in title", func(t *testing.T) {
		title := "foo: 2022-12-31"

		got, err := ExtractDate(title)
		require.NoError(t, err)

		want := "2022-12-31"
		assert.Equal(t, want, got)
	})

	t.Run("error if there is no yyyy-mm-dd in title", func(t *testing.T) {
		title := "foo"

		_, err := ExtractDate(title)
		require.Error(t, err)
	})
}
