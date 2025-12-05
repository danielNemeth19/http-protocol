package main

import (
	"log"
	"os"
	"os/signal"
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
