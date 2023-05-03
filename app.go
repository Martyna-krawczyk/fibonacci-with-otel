package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type App struct {
	r io.Reader
	l *log.Logger
}

// name is the Tracer name used to identify this instrumentation library.
const name = "fibonacci-with-otel"

func NewApp(r io.Reader, l *log.Logger) *App {
	return &App{r: r, l: l}
}

func (a *App) Run(ctx context.Context) error {
	for {
		// Each execution of the run loop, we should get a new "root" span and context
		// starts a new OpenTelemetry span for the Write() method with the name "Poll"
		// The Start() method returns a new context and a trace.Span object, which
		// is used to end the span.
		newCtx, span := otel.Tracer(name).Start(ctx, "Run")

		n, err := a.Poll(newCtx)
		if err != nil {
			span.End()
			return err
		}
		a.Write(newCtx, n)
		span.End()
	}
}

func (a *App) Poll(ctx context.Context) (uint, error) {
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()

	a.l.Print("What fibonacci no would you like to know?")

	// reading a single integer value from the input source and returning it to
	// the calling function along with any errors encountered during the reading process
	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)

	// Store n as a string to not overflow an int64.
	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String("request.n", nStr))
	return n, err
}

func (a *App) Write(ctx context.Context, n uint) {
	var span trace.Span
	ctx, span = otel.Tracer(name).Start(ctx, "Write")
	defer span.End()

	f, err := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer(name).Start(ctx, "Fibonacci")
		defer span.End()
		return Fibonacci(n)
	}(ctx)

	if err != nil {
		a.l.Printf("fibonacci(%d) = %v\n", n, err)
	} else {
		a.l.Printf("fibonacci(%d) = %d\n", n, f)
	}
}
