package main

import (
	"context"
	"go-worker/client"
	"go-worker/internal/tracer"
	"go-worker/model"
	"log"
	"os"
	"os/signal"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutDown, err := tracer.InitTrace()
	if err != nil {
		log.Panic(err)
	}

	tr := otel.GetTracerProvider().Tracer("go-worker/main")
	ctx, span := tr.Start(ctx, "main", trace.WithSpanKind(trace.SpanKindServer))

	defer tracer.RunShutdown(shutDown, ctx)
	defer span.End()

	viaCepAddress, err := client.SearchViaCepZipCode("01001000", ctx)
	if err != nil {
		return
	}

	viaCepAddress.NormalizdZipCode()
	err = client.CreateNodeApiAddress(model.Address(viaCepAddress), ctx)
	if err != nil {
		return
	}
}
