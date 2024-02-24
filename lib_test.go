package dunia

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestRun(t *testing.T) {
	_, err := Init()
	assert.Equal(t, err, nil, err)
}
