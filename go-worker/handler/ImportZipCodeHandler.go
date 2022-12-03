package handler

import (
	"context"
	"encoding/json"
	"go-worker/client"
	"go-worker/internal/settings"
	"go-worker/internal/tracer"
	"go-worker/listener"
	"go-worker/model"
	"log"
	"os"
	"os/signal"

	"github.com/aws/aws-sdk-go/service/sqs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func ImportZipCodeHandler() {
	listener.DispatchMessage(30, settings.Env.SQS.ImportZipCodeQueue, process)
}

func process(message *sqs.Message) {
	sqsMessage := listener.UnmarshalSQSMessage(message)
	importZipCode := unmarshalSQSImportZipcodeAddress(sqsMessage.Message)
	importZipCodeAddress(importZipCode.ZipCode)
}

func importZipCodeAddress(zipCode string) {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	shutDown, err := tracer.InitTrace()
	if err != nil {
		log.Panic(err)
	}

	tr := otel.GetTracerProvider().Tracer("go-worker/main")
	ctx, span := tr.Start(ctx, "main", trace.WithSpanKind(trace.SpanKindServer))
	tracer.ZipCodeSpanAttribute(span, zipCode)

	defer tracer.RunShutdown(shutDown, ctx)
	defer span.End()

	viaCepAddress, err := client.SearchViaCepZipCode(zipCode, ctx)
	if err != nil {
		return
	}

	viaCepAddress.NormalizdZipCode()
	err = client.CreateNodeApiAddress(model.Address(viaCepAddress), ctx)
	if err != nil {
		return
	}
}

func unmarshalSQSImportZipcodeAddress(message string) model.SQSImportZipcodeAddress {

	var sqsMessage model.SQSImportZipcodeAddress

	if err := json.Unmarshal([]byte(message), &sqsMessage); err != nil {
		panic(err)
	}

	return sqsMessage
}
