package main

import (
	"go-worker/handler"
	"go-worker/internal/metrics"
)

func main() {
	metrics.StartupMetrics()
	handler.ImportZipCodeHandler()
}
