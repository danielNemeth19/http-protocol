package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidSingleHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)
}

func TestValidSingleHeaderWithExtraSpace(t *testing.T) {
	headers := NewHeaders()
	data := []byte("  Host:localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 24, n)
	assert.False(t, done)
}

func TestValidTwoHeadersWithExistingHeaders(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Host:localhost:42069\r\n")
	headers.Parse(data)
	data = []byte("Content-Type: application/json\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, "application/json", headers["Content-Type"])
	assert.Equal(t, 32, n)
	assert.False(t, done)
}

func TestNotEnoughData(t *testing.T) {
	headers := NewHeaders()
	data := []byte("Content-Type: application")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.False(t, done)
	assert.Equal(t, 0, n)
}

func TestInvalidFieldLine(t *testing.T) {
	headers := NewHeaders()
	data := []byte("InvalidFieldLine\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	require.EqualError(t, err, "Field line supposed to have a ':' separator")
	require.NotNil(t, headers)
	assert.False(t, done)
	assert.Equal(t, 0, n)
}

func TestValidDone(t *testing.T) {
	headers := NewHeaders()
	data := []byte("\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.True(t, done)
	assert.Equal(t, 2, n)
}

func TestInvalidSpacingHeader(t *testing.T) {
	headers := NewHeaders()
	data := []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	require.EqualError(t, err, "Invalid field name")
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
