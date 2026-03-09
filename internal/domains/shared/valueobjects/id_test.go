package valueobjects_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/darmayasa221/polymarket-go/internal/domains/shared/valueobjects"
)

func TestNewID(t *testing.T) {
	t.Parallel()

	id1 := valueobjects.NewID()
	id2 := valueobjects.NewID()
	assert.False(t, id1.IsEmpty())
	assert.NotEqual(t, id1, id2) // IDs must be unique
	assert.Len(t, id1.String(), 36)
}

func TestID_Equals(t *testing.T) {
	t.Parallel()

	id := valueobjects.NewID()
	same := id
	assert.True(t, id.Equals(same))
	assert.False(t, id.Equals(valueobjects.NewID()))
}
