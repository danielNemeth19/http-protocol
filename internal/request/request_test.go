package request

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data string
	numBytesPerRead int
	pos int
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := min(cr.pos + cr.numBytesPerRead, len(cr.data))
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	return n, nil
}

func TestRequestFromReader_EOF(t *testing.T) {
	reader := &chunkReader{
		data: "GE",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, r.RequestLine, RequestLine{})
}

func TestRequestFromReader_ParsesRequestLineGet(t *testing.T) {
	reader := &chunkReader{
		data: "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestRequestFromReader_ParsesRequestLineWithPath(t *testing.T) {
	reader := &chunkReader{
		data: "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 1,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestRequestFromReader_ParsesRequestLineWithPathPost(t *testing.T) {
	reader := &chunkReader{
		data: "POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 80,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "POST", r.RequestLine.Method)
	assert.Equal(t, "/coffee", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
}

func TestRequestFromReader_RaisesErrorIfPartMissing(t *testing.T) {
	reader := &chunkReader{
		data: "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 4,
	}
	_, err := RequestFromReader(reader)
	require.Error(t, err)
	require.EqualError(t, err, "Request line supposed to have three parts, got: 2\n")
}

func TestRequestFromReader_RaisesErrorIfMethodInvalid(t *testing.T) {
	reader := &chunkReader{
		data: "BLA /coffee HTTP/2.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 12,
	}
	_, err := RequestFromReader(reader)
	require.Error(t, err)
	require.EqualError(t, err, "BLA is not a valid method\n")
}

func TestRequestFromReader_RaisesErrorIfVersionUnsupported(t *testing.T) {
	reader := &chunkReader{
		data: "GET /coffee HTTP/2.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 5,
	}
	_, err := RequestFromReader(reader)
	require.Error(t, err)
	require.EqualError(t, err, "HTTP Version is unsupported: 2.1\n")
}

