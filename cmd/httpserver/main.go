package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/danielNemeth19/http-protocol/internal/request"
	"github.com/danielNemeth19/http-protocol/internal/response"
	"github.com/danielNemeth19/http-protocol/internal/server"
)

const port = 42069

func myHandler(w io.Writer, req *request.Request) *server.HandlerError {
	target := req.RequestLine.RequestTarget 
	if target == "/yourproblem" {
		resp := server.HandlerError{
			Code: response.StatusBadRequest,
			Message: "Your problem is not my problem\n",
		}
		return &resp
	}
	if target == "/myproblem" {
		resp := server.HandlerError{
			Code: response.StatusInternalServerError,
			Message: "Woopsie, my bad\n",
		}
		return &resp
	}
	w.Write([]byte("All good, frfr\n"))
	return nil
}

func main()  {
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
