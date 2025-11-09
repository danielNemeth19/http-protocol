# HTTP/1.1 RFCs

| RFC      | Status/Notes                                      | Link                                               |
|----------|---------------------------------------------------|----------------------------------------------------|
| RFC 2616 | Deprecated by RFC 7231                            | [RFC 2616](https://datatracker.ietf.org/doc/html/rfc2616) |
| RFC 7231 | Active, widely referenced, verbose                | [RFC 7231](https://datatracker.ietf.org/doc/html/rfc7231) |
| RFC 9110 | Covers HTTP "semantics"                           | [RFC 9110](https://datatracker.ietf.org/doc/html/rfc9110) |
| RFC 9112 | Easier to read, relies on RFC 9110 understanding  | [RFC 9112](https://datatracker.ietf.org/doc/html/rfc9112) |

**Notes:**  
- [RFC 9110](https://datatracker.ietf.org/doc/html/rfc9110) and [RFC 9112](https://datatracker.ietf.org/doc/html/rfc9112) have better separation of information.  
- [RFC 7231](https://datatracker.ietf.org/doc/html/rfc7231) can be verbose; [RFC 9112](https://datatracker.ietf.org/doc/html/rfc9112) is concise but assumes familiarity with [RFC 9110](https://datatracker.ietf.org/doc/html/rfc9110).

---
# TCP Chapter

## Run TCP Listener and Redirect Output

```sh
go run ./cmd/tcplistener | tee /tmp/tcp.txt
```

## Send Message via Netcat

In another shell:

```sh
printf "Do you have what it takes to be an engineer at TheStartup™?\r\n" | nc -w 1 127.0.0.1 42069
```

Netcat will transmit the message with 1 second timeout time

---

# TCP Chapter: UDP Sender

## Run UDP Sender

```sh
go run ./cmd/udpsender
```

## Start UDP Listener

In a separate terminal (-l option starts up an upd listener):

```sh
nc -u -l 42069
```

Messages sent in the UDP sender terminal should appear in the listener terminal.

---

# Requests Chapter: TCP to HTTP

## Run TCP Listener and Redirect Output

```sh
go run ./cmd/tcplistener | tee /tmp/rawget.http
```

## Send HTTP GET Request

From another shell:

```sh
curl http://localhost:42069/coffee
```

*Note: The request will hang since the TCP listener only listens; the request will come in but not be processed.*

---
## Requests Chapter: TCP to HTTP
- run the tcplistener and redirect its output:
go run ./cmd/tcplistener | tee /tmp/rawget.http

- from another shell send a GET request to it:
curl http://localhost:42069/coffee

- the request will hang since TCP listener only listens but request will come in

---

# Explanations

- **Reading from a network** is conceptually similar to reading from a file:
    - From a file, we *pull* data (control how much we read).
    - From a network, data is *pushed* to us; we must handle incoming bytes.
- The interface for both is the same: a stream of ordered bytes.

---

# `RequestFromReader` Key Concepts
- Buffer (`buf`) is initialized to hold incoming data.
    - in tests, `numBytesPerRead` in `chunkReader` can mimic the network, sending arbitrary amount of bytes
- we need to track the amount of bytes that we read from the stream (`readToIndex`)
- if this goes above the capacity of the buffer:
    - we create new buffer with capacity `2*cap(buf)`
    - and copy bytes already in buffer (i.e. read but not parsed yet) to the new buffer
- on every read, we always read from the position of `readToIndex`:
    - this is to make sure that we don't override bytes that we already read but haven't parsed yet
    - (see `remainderBytes` about this later)
- after we read:
    - we increase `readToIndex` with the number of bytes read
    - we call `parse` passing in the data, which is whatever we have in the buffer up to `readIndex` (can be some garbage bytes after the index)
- `parse` will return `parsedBytes` or an error:
    - it only returns non-null `parsedBytes`, if it was able to parse anything from the last chunk of bytes we read
    - this could be the requestLine or header(s) or body
    - if there's anything parsed:
        - we need to check if there's any part of the read bytes, that is not parsed yet
        - i.e is `readToIndex - parsedBytes` greater than 0?
        - e.g.: we read 8 bytes in the last chunk:
            - out of which the first 3 completes the request line -> this is parsed
            - but the remaining 5 bytes will need to be carried over for subsequent reading and parsing
        - so whatever is still left to be parsed:
            - gets copied to the start of the buffer
            - bytes not impacted by the copy will stay in the buffer, in this example after buf[5:]
        - this means:
            - we don't flush the buffer entirely
            - we just move unparsed data to the start of it
            - finally we re-set `readToIndex` to be at `remainderBytes`, so that we continue to read from the stream into the buffer from this position

---

# Debug Notes

To run a test with `dlv`:

```sh
dlv test ./internal/request/ -- -test.run TestRequestFromReader_EOF
```

Set a breakpoint:

```sh
break ./internal/request/request.go:96
```
