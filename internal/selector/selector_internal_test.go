//nolint:testpackage // these tests cover the unexported output formatter.
package selector

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintMatches(t *testing.T) {
	t.Run("each visible directory has a one-based index", func(t *testing.T) {
		// Arrange
		out := &strings.Builder{}

		// Act
		printMatches(out, []string{"docs"}, map[string]bool{})

		// Assert
		assert.Equal(t, "1: [ ] docs\n", out.String())
	})
}
