package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shharn/todo/handler"

)

var indexFilePath = "public/index.html"

func main() {
	listenSignal()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("public"))))
	http.HandleFunc("/api/todos", handler.TodosHandler)
	http.HandleFunc("/", handler.WithFinder(
		handler.IndexHandler,
		handler.NewLocalFileSystemFinder(indexFilePath),
	))
	port := "9000"
	http.ListenAndServe(":" + port, nil)
}

func listenSignal() {
	go func() {
		sleepTime := 1
		signalChannel := make(chan os.Signal)
		signal.Notify(signalChannel, syscall.SIGTERM)
		signal.Notify(signalChannel, syscall.SIGINT)
		sig := <-signalChannel
		log.Printf("Receive signal: %+v\n", sig)
		log.Printf(fmt.Sprintf("Wait for %v second to finish remaining task", sleepTime))
		time.Sleep(time.Duration(sleepTime) * time.Second)
		os.Exit(0)
	}()
}
