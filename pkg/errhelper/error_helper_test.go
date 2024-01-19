package errhelper

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func someFunc() error {
	return errors.WithStack(errors.New("dd"))
}

func TestWithFileLine(t *testing.T) {
	err := someFunc()
	str := WithFileLine(err).Error()
	assert.Equal(t, strings.Contains(str, "/"), false)
}
