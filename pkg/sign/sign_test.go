package sign

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSign(t *testing.T) {
	expected := map[string]string{
		"nars": "bef26b7dec708ba104e2e31d183442a7",
		"sesi": "j5nZZoD50c8c1559bb2bd2a5e0cff487f3a8b78",
	}

	got, err := NewSign("../../sign.js").Sign(map[string]string{
		"cookie": "ede2c413e49d1a65566e12c27c819",
		"uid":    "647723be79400013ab",
	})
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}
