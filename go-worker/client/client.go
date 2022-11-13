package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-worker/log"
	"go-worker/model"
	"go-worker/tracer"
	"io/ioutil"
	"net/http"
)

const (
	_viaCepApiBaseURL = "https://viacep.com.br/ws"
	_nodeAPIBaseURL   = "http://localhost:8081"
)

func SearchViaCepZipCode(zipCode string, ctx context.Context) (model.ViaCepAddress, error) {

	span := tracer.Span(ctx, "go-worder/client", "search-via-cep-zip-code")
	defer span.End()

	method := "GET"
	url := fmt.Sprintf("%s/%s/json/", _viaCepApiBaseURL, zipCode)

	tracer.HttpRequestSpanAttributes(span, url, method, "")
	response, err := http.Get(url)
	if err != nil {
		tracer.HttErrorSpanAttributes(span, err)
		return model.ViaCepAddress{}, err
	}

	defer response.Body.Close()

	responseBody, _ := ioutil.ReadAll(response.Body)
	logger := log.HttpResponseLogger(url, method, string(responseBody), response.StatusCode)
	tracer.HttResponseSpanAttributes(span, string(responseBody), response.StatusCode)

	if response.StatusCode != http.StatusOK {
		err := errors.New(string(responseBody))
		logger.Error(err)
		return model.ViaCepAddress{}, err
	}

	logger.Info()

	var address model.ViaCepAddress
	if err := json.Unmarshal(responseBody, &address); err != nil {
		tracer.HttErrorSpanAttributes(span, err)
		return model.ViaCepAddress{}, err
	}

	return address, nil
}

func CreateNodeApiAddress(address model.Address, ctx context.Context) error {

	span := tracer.Span(ctx, "go-worder/client", "create-node-api-address")
	defer span.End()

	method := "POST"
	url := fmt.Sprintf("%s/address/", _nodeAPIBaseURL)
	body, err := json.Marshal(address)

	if err != nil {
		tracer.HttErrorSpanAttributes(span, err)
		return err
	}

	tracer.HttpRequestSpanAttributes(span, url, method, string(body))
	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		tracer.HttErrorSpanAttributes(span, err)
		return err
	}

	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)

	logger := log.HttpResponseLogger(url, method, string(responseBody), response.StatusCode)
	tracer.HttResponseSpanAttributes(span, string(responseBody), response.StatusCode)

	if response.StatusCode != http.StatusCreated {
		err := errors.New(string(responseBody))
		logger.Error(err)
		return err
	}

	logger.Info()
	return nil
}