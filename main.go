package main

import (
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	// creates a new logger that writes to os.Stdout (the terminal) with no 
	// prefix and no flags
	l := log.New(os.Stdout, "", 0)

	// creates a buffered channel of os.Signal values with a buffer size of 1. 
	// This channel is used to receive interrupt signals
	sigCh := make(chan os.Signal, 1)

	// registers the sigCh channel to receive os.Interrupt signals. When the 
	// program receives an interrupt signal, it will send the signal to the 
	// sigCh channel
	
	signal.Notify(sigCh, os.Interrupt)
	// creates an unbuffered channel of error values. This channel is used to 
	// receive errors from the application's Run() method
	errCh := make(chan error)

	// creates a new instance of an application using an input reader from 
	// os.Stdin and the logger l
	app := NewApp(os.Stdin, l)

	// starts a new goroutine that runs the application's Run() method with an 
	// empty context and sends any errors returned by Run() to the errCh channel
	go func() {
		errCh <- app.Run(context.Background())
	}()

	// waits for either a signal from sigCh or an error from errCh. If a signal 
	// is received, the program logs a message and returns, terminating the 
	// program. If an error is received, the program logs the error and 
	// terminates with a fatal error.
	select {
	case <-sigCh:
		l.Println("\ngoodbye")
		return
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}
