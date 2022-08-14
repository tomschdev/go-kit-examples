package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// StringService provides operations on strings.
type StringService interface { // go-kit models a service as an interface
	Uppercase(string) (string, error)
	Count(string) int
}

// stringService is a concrete implementation of StringService
type stringService struct{}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

// Implementation methods of StringService via stringService struct
func (stringService) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func (stringService) Count(s string) int {
	return len(s)
}

// For each method of the service interface, we define request and response structs
type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

// Endpoints are a primary abstraction in go-kit. An endpoint represents a single RPC (method in our service interface)
// Each method in our service interface is provisioned an endpoint
// Input: StringService interface
// Output: endpoint.Endpoint

// GoDoc Endpoint:
// type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

// Endpoint is the fundamental building block of servers and clients.
// It represents a single RPC method.
func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	// Essentially, the endpoint is a function that takes context and any request type
	// and returns any response type and an error
	// Recall: the empty interface can hold any type because all types implement at least no methods

	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest) // For the Uppercase endpoint, we know we require an uppercaseRequest from our request interface{}
		// The above operation is a type assertion, it will extract the uppercaseRequest from the empty interface
		v, err := svc.Uppercase(req.S) // Insert the S field of the request into the Uppercase implementation
		// Return the response of the service's method, depending on err
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil // Notice the error is contained in the uppercaseResponse struct
		}
		return uppercaseResponse{v, ""}, nil
	}
}

// Again, general procedure of endpoint:
// 1) return a func taking context and request, and returning any response and error
// 2) extract specific request struct
// 3) invoke service's method corresponding to the endpoint, with request as struct as input
// 4) return response struct with error encoding in response struct
// Note use of adapter pattern
func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.Count(req.S)
		return countResponse{v}, nil
	}
}

// Transports expose the service to the network. In this first example we utilize JSON over HTTP.
// Stitch it all together to create handlers and map them to routes, then serve.
func main() {
	// Create instance of stringService struct
	svc := stringService{}

	// Create handler by linking endpoint to transport, providing:
	// endpoint builder with concrete service struct inserted (builder accepts interface - so need methods to be implemented)
	// request decode function
	// response encode function
	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	// Link route to handler
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)

	// Serve
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Decode functions are just helpers to be used by the handlers to convert json request to a struct
// Therefore, they receive http.Request and return interface{} (any type of struct)
func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil { // Decoder takes request body and Decode(<pointer to empty request>)
		return nil, err
	}
	return request, nil
}

func decodeCountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request countRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

// Encode functions are used by handler to convert struct to json
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response) // Encoder takes response writer and then Encode(<empty response>)
}
