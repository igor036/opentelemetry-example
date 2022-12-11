package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-worker/internal/log"
	"go-worker/internal/metrics"
	"go-worker/internal/settings"
	"go-worker/internal/tracer"
	"go-worker/model"
	"io/ioutil"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

func SearchViaCepZipCode(zipCode string, ctx context.Context) (model.ViaCepAddress, error) {

	metricsAttrs := metrics.BuildMetrictsAttrs(attribute.Key("api").String("via-cep"))
	span := tracer.Span(ctx, "go-worder/client", "search-via-cep-zip-code")
	defer span.End()

	method := "GET"
	url := fmt.Sprintf("%s/%s/json/", settings.Env.ViaCep.BaseURL, zipCode)
	tracer.HttpRequestSpanAttributes(span, url, method, "")

	startTime := time.Now()
	response, err := http.Get(url)
	endTime := time.Now()

	metrics.ResponseTimeMetric(endTime.Sub(startTime), metricsAttrs)

	if err != nil {
		tracer.HttErrorSpanAttributes(span, err)
		metrics.ErrorRequestCountMetric(metricsAttrs)
		return model.ViaCepAddress{}, err
	}

	defer response.Body.Close()

	responseBody, _ := ioutil.ReadAll(response.Body)
	logger := log.HttpResponseLogger(url, method, string(responseBody), response.StatusCode)
	tracer.HttResponseSpanAttributes(span, string(responseBody), response.StatusCode)
	metrics.SuccessRequestCountMetric(metricsAttrs)

	if response.StatusCode != http.StatusOK {
		err := errors.New(string(responseBody))
		logger.Error(err)
		return model.ViaCepAddress{}, err
	}

	var address model.ViaCepAddress
	if err := json.Unmarshal(responseBody, &address); err != nil {
		tracer.HttErrorSpanAttributes(span, err)
		return model.ViaCepAddress{}, err
	}

	return address, nil
}

func CreateNodeApiAddress(address model.Address, ctx context.Context) error {

	metricsAttrs := metrics.BuildMetrictsAttrs(attribute.Key("api").String("node-api"))
	span := tracer.Span(ctx, "go-worder/client", "create-node-api-address")
	defer span.End()

	method := "POST"
	url := fmt.Sprintf("%s/address/", settings.Env.NodeAPI.BaseURL)
	body, err := json.Marshal(address)

	if err != nil {
		tracer.HttErrorSpanAttributes(span, err)
		return err
	}

	tracer.HttpRequestSpanAttributes(span, url, method, string(body))

	startTime := time.Now()
	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	endTime := time.Now()

	metrics.ResponseTimeMetric(endTime.Sub(startTime), metricsAttrs)

	if err != nil {
		tracer.HttErrorSpanAttributes(span, err)
		metrics.ErrorRequestCountMetric(metricsAttrs)
		return err
	}

	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)

	logger := log.HttpResponseLogger(url, method, string(responseBody), response.StatusCode)
	tracer.HttResponseSpanAttributes(span, string(responseBody), response.StatusCode)
	metrics.SuccessRequestCountMetric(metricsAttrs)

	if response.StatusCode != http.StatusCreated {
		err := errors.New(string(responseBody))
		logger.Error(err)
		return err
	}

	logger.Info()
	return nil
}
