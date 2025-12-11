package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/danielNemeth19/http-protocol/internal/request"
	"github.com/danielNemeth19/http-protocol/internal/response"
	"github.com/danielNemeth19/http-protocol/internal/server"
)

const port = 42069

func myHandler(w *response.Writer, req *request.Request) {
	target := req.RequestLine.RequestTarget
	if target == "/yourproblem" {
		resp := server.HandlerError{
			Code:    response.StatusBadRequest,
			Message: response.BadRequestHTML,
		}
		w.WriteStatusLine(resp.Code)
		headers := response.GetDefaultHeaders(len(resp.Message))
		headers = response.ReplaceHeader(map[string]string{"Content-Type": "text/html"}, headers)
		w.WriteHeaders(headers)
		w.WriteBody([]byte(resp.Message))
		return
	}
	if target == "/myproblem" {
		resp := server.HandlerError{
			Code:    response.StatusInternalServerError,
			Message: response.InternalServerErrorHTML,
		}
		w.WriteStatusLine(resp.Code)
		headers := response.GetDefaultHeaders(len(resp.Message))
		headers = response.ReplaceHeader(map[string]string{"Content-Type": "text/html"}, headers)
		w.WriteHeaders(headers)
		w.WriteBody([]byte(resp.Message))
		return
	}
	toGet, found := strings.CutPrefix(target, "/httpbin")
	if found {
		// toTarget := "https://httpbin.org" + toGet
		toTarget := "http://localhost:8080" + toGet
		resp, _ := http.Get(toTarget)
		w.WriteStatusLine(response.StatusOK)
		headers := response.GetChunkedHeaders()
		headers.Set("Trailer", "X-Content-SHA256")
		headers.Set("Trailer", "X-Content-Lenght")
		w.WriteHeaders(headers)
		buf := make([]byte, 1024)
		var content []byte
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				content = append(content, buf[:n]...)
			}
			if err == io.EOF {
				if n > 0 {
					w.WriteChunkedBody(buf[:n])
				}
				w.WriteChunkedBodyDone()
				size := sha256.Sum256(buf)
				fmt.Printf("%x\n", size)
				return
			}
			if err != nil {
				fmt.Println(err)
				return
			}
			n, err = w.WriteChunkedBody(buf[:n])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	w.WriteStatusLine(response.StatusOK)
	headers := response.GetDefaultHeaders(len(response.SuccessHTML))
	headers = response.ReplaceHeader(map[string]string{"Content-Type": "text/html"}, headers)
	w.WriteHeaders(headers)
	w.WriteBody([]byte(response.SuccessHTML))
}

func main() {
	server, err := server.Serve(port, myHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)

	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
