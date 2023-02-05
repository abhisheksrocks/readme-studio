package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type RequestHeader struct {
	key   string
	value string
}

func (e *RequestHeader) GetKey() string {
	return e.key
}

func (e *RequestHeader) GetValue() string {
	return e.value
}

type RequestHeaderError string

const (
	RequestHeaderErrorEmptyKey   RequestHeaderError = "empty key"
	RequestHeaderErrorEmptyValue RequestHeaderError = "empty value"
)

func makeContentTypeHeader(value string) RequestHeader {
	return RequestHeader{
		key:   "Content-Type",
		value: value,
	}
}
func makeDefaultContentTypeHeader() RequestHeader {
	return makeContentTypeHeader("application/json")
}

func makeAuthorizationHeader(token string) RequestHeader {
	return RequestHeader{
		key:   "Authorization",
		value: "bearer " + token,
	}
}

// type GraphQlRequest struct {
// 	EndpointURL string
// 	Query       string
// 	Headers     []RequestHeader
// }

// func (r *GraphQlRequest) call() {}

type GraphQlRequestErrorMessage string

const (
	GraphQlRequestErrorUnknown          GraphQlRequestErrorMessage = "unknown error"
	GraphQlRequestErrorInvalidEndpoint  GraphQlRequestErrorMessage = "couldn't parse endpoint"
	GraphQlRequestErrorInvalidQueryData GraphQlRequestErrorMessage = "invalid query data"
	GraphQlRequestErrorInvalidResponse  GraphQlRequestErrorMessage = "couldn't parse response"
)

type ErrorDataSource string

const (
	ErrorDataSourceUnknown ErrorDataSource = "unknown"
	ErrorDataSourceGithub  ErrorDataSource = "github"
	ErrorDataSourceUs      ErrorDataSource = "us"
)

type ErrorData struct {
	Source  ErrorDataSource
	Message string
	URL     string
}

type GithubError struct {
	Message          string `json:"message"`
	DocumentationURL string `json:"documentation_url"`
}

func (src GithubError) toErrorData() ErrorData {
	return ErrorData{
		Source:  ErrorDataSourceGithub,
		Message: (src.Message),
		URL:     src.DocumentationURL,
	}
}

func makeRequest(endpointURL string, query string, headers []RequestHeader, client *http.Client, result interface{}) *ErrorData {
	queryMap := map[string]string{
		"query": query,
	}
	queryRaw, err := json.Marshal(queryMap)
	if err != nil {
		return &ErrorData{
			Source:  ErrorDataSourceUs,
			Message: string(GraphQlRequestErrorInvalidQueryData),
		}
	}
	_, err = url.Parse(endpointURL)
	if err != nil {
		return &ErrorData{
			Source:  ErrorDataSourceUs,
			Message: string(GraphQlRequestErrorInvalidEndpoint),
		}
	}

	requestData, err := http.NewRequest(http.MethodPost, endpointURL, bytes.NewBuffer(queryRaw))
	if err != nil {
		return &ErrorData{
			Source:  ErrorDataSourceUs,
			Message: (err.Error()),
		}
	}
	for _, head := range headers {
		if empty(head.key) {
			return &ErrorData{
				Source:  ErrorDataSourceUs,
				Message: string(RequestHeaderErrorEmptyKey),
			}
		}
		if empty(head.value) {
			return &ErrorData{
				Source:  ErrorDataSourceUs,
				Message: string(RequestHeaderErrorEmptyValue),
			}
		}
		requestData.Header.Add(head.key, head.value)
	}
	response, err := client.Do(requestData)
	if err != nil {
		// return err
		return &ErrorData{
			Source:  ErrorDataSourceUnknown,
			Message: (err.Error()),
		}
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		// return errors.New(GraphQlRequestErrorInvalidResponse)
		return &ErrorData{
			Source:  ErrorDataSourceUs,
			Message: string(GraphQlRequestErrorInvalidResponse),
		}
	}
	if response.StatusCode >= http.StatusOK && response.StatusCode < 400 {
		err = json.Unmarshal(data, &result)
		if err != nil {
			return &ErrorData{
				Source:  ErrorDataSourceUs,
				Message: string(GraphQlRequestErrorInvalidResponse),
			}
		}
		return nil
	}
	var possibleGithubError GithubError
	err = json.Unmarshal(data, &possibleGithubError)
	if err != nil {
		return &ErrorData{
			Source:  ErrorDataSourceUs,
			Message: string(GraphQlRequestErrorInvalidResponse),
		}
	}
	toReturn := possibleGithubError.toErrorData()
	if notEmpty(string(toReturn.Message)) {
		return &ErrorData{
			Source:  ErrorDataSourceGithub,
			Message: toReturn.Message,
		}
	}
	return &toReturn
	// return errors.New(string(data))
}
