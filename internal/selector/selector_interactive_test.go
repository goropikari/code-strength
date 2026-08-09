//nolint:testpackage // these tests cover the unexported terminal predicate.
package selector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShouldSelectInteractive(t *testing.T) {
	t.Run("requires file input", func(t *testing.T) {
		// Arrange

		// Act
		got := shouldSelectInteractive(false, true, true, true)

		// Assert
		assert.False(t, got)
	})

	t.Run("requires file output", func(t *testing.T) {
		// Arrange

		// Act
		got := shouldSelectInteractive(true, false, true, true)

		// Assert
		assert.False(t, got)
	})

	t.Run("requires terminal input", func(t *testing.T) {
		// Arrange

		// Act
		got := shouldSelectInteractive(true, true, false, true)

		// Assert
		assert.False(t, got)
	})

	t.Run("requires terminal output", func(t *testing.T) {
		// Arrange

		// Act
		got := shouldSelectInteractive(true, true, true, false)

		// Assert
		assert.False(t, got)
	})

	t.Run("returns true for two terminal files", func(t *testing.T) {
		// Arrange

		// Act
		got := shouldSelectInteractive(true, true, true, true)

		// Assert
		assert.True(t, got)
	})
}
