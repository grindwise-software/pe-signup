package usrmgmt

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-kit/kit/endpoint"
)

type StringService interface {
	uppercase(string) (string, error)
	count(string) int
}

type stringService struct{}

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

type countRequest struct {
	S string `json:"s"`
}

type countResponse struct {
	V int `json:"v"`
}

func (stringService) uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}

	return strings.ToUpper(s), nil
}

func (stringService) count(s string) int {
	return len(s)
}

var ErrEmpty = errors.New("empty string")

func makeUppercaseEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}

		return uppercaseResponse{v, ""}, nil
	}
}

func makeCountEndpoint(svc StringService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(countRequest)
		v := svc.count(req.S)
		return countResponse{v}, nil
	}
}

// func main() {
// 	svc := stringService{}

// 	uppercaseHandler := httptransport.NewServer(
// 		makeUppercaseEndpoint(svc),
// 		decodeUppercaseRequest,
// 		encodeResponse,
// 	)

// 	countHandler := httptransport.NewServer(
// 		makeCountEndpoint(svc),
// 		decodeCountRequest,
// 		encodeResponse,
// 	)

// 	http.Handle("/uppercase", uppercaseHandler)
// 	http.Handle("/count", countHandler)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
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

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
