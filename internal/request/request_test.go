package request

import (
	"io"
	"testing"

	"github.com/danielNemeth19/http-protocol/internal/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type chunkReader struct {
	data            string
	numBytesPerRead int
	pos             int
}

func (cr *chunkReader) Read(p []byte) (n int, err error) {
	if cr.pos >= len(cr.data) {
		return 0, io.EOF
	}
	endIndex := min(cr.pos+cr.numBytesPerRead, len(cr.data))
	n = copy(p, cr.data[cr.pos:endIndex])
	cr.pos += n
	return n, nil
}

func TestRequestFromReader_EOF(t *testing.T) {
	reader := &chunkReader{
		data:            "GE",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, r.RequestLine, RequestLine{})
}

func TestRequestFromReader_ParsesRequestLineGet(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
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
		data:            "GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
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
		data:            "POST /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
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
		data:            "/coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 4,
	}
	_, err := RequestFromReader(reader)
	require.Error(t, err)
	require.EqualError(t, err, "Error during parsing: Request line supposed to have three parts, got: 2\n")
}

func TestRequestFromReader_RaisesErrorIfMethodInvalid(t *testing.T) {
	reader := &chunkReader{
		data:            "BLA /coffee HTTP/2.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 12,
	}
	_, err := RequestFromReader(reader)
	require.Error(t, err)
	require.EqualError(t, err, "Error during parsing: BLA is not a valid method\n")
}

func TestRequestFromReader_RaisesErrorIfVersionUnsupported(t *testing.T) {
	reader := &chunkReader{
		data:            "GET /coffee HTTP/2.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 5,
	}
	_, err := RequestFromReader(reader)
	require.Error(t, err)
	require.EqualError(t, err, "Error during parsing: HTTP Version is unsupported: 2.1\n")
}

func TestRequestFromReader_ParsesHeaders(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 50,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "localhost:42069", r.Headers["host"])
	assert.Equal(t, "curl/7.81.0", r.Headers["user-agent"])
	assert.Equal(t, "*/*", r.Headers["accept"])
}

func TestRequestFromReader_ParsesMultipleHeaders(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\na: 1\r\nb: 2\r\nc: 3\r\n\r\n",
		numBytesPerRead: 50,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "1", r.Headers["a"])
	assert.Equal(t, "2", r.Headers["b"])
	assert.Equal(t, "3", r.Headers["c"])
}

func TestRequestFromReader_ParsesEmptyHeaders(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\n\r\n",
		numBytesPerRead: 50,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	h := headers.NewHeaders()
	assert.Equal(t, h, r.Headers)
}

func TestRequestFromReader_MalformedHeader(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}
	_, err := RequestFromReader(reader)
	require.Error(t, err)
	require.EqualError(t, err, "Error during parsing: Invalid header key: key contains invalid char  ")
}

func TestRequestFromReader_ParsesCaseInsensitiveHeaders(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHOST: localhost:42069\r\nCACHE-CONTROL: no-cache\r\n\r\n",
		numBytesPerRead: 7,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "localhost:42069", r.Headers["host"])
	assert.Equal(t, "no-cache", r.Headers["cache-control"])
}

func TestRequestFromReader_MissingEndOfHeaders(t *testing.T) {
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nConnection: keep-alive\r\n",
		numBytesPerRead: 3,
	}
	r, err := RequestFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "GET", r.RequestLine.Method)
	assert.Equal(t, "/", r.RequestLine.RequestTarget)
	assert.Equal(t, "1.1", r.RequestLine.HttpVersion)
	assert.Equal(t, "keep-alive", r.Headers["connection"])
}

// func TestRequestFromReader_StandardBody(t *testing.T) {
	// reader := &chunkReader{
		// data: "POST /submit HTTP/1.1\r\n" +
			// "Host: localhost:42069\r\n" +
			// "Content-Length: 13\r\n" +
			// "\r\n" +
			// "hello world!\n",
		// numBytesPerRead: 3,
	// }
	// r, err := RequestFromReader(reader)
	// require.NoError(t, err)
	// require.NotNil(t, r)
	// assert.Equal(t, "hello world!\n", string(r.Body))
// }

// func TestRequestFromReader_BodyShorterThanReportContentLength(t *testing.T) {
	// reader := &chunkReader{
		// data: "POST /submit HTTP/1.1\r\n" +
			// "Host: localhost:42069\r\n" +
			// "Content-Length: 20\r\n" +
			// "\r\n" +
			// "partial content",
		// numBytesPerRead: 3,
	// }
	// r, err := RequestFromReader(reader)
	// require.Error(t, err)
	// require.NotNil(t, r)
// }
