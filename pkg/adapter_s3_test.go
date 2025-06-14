package stm

import (
	"bytes"
	"compress/gzip"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type MockRoundTripper struct {
	Requests  []*http.Request
	Responses []*http.Response
}

func (m *MockRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	m.Requests = append(m.Requests, r)

	shiftResponse := m.Responses[0]
	m.Responses = m.Responses[1:]
	return shiftResponse, nil
}

func TestS3Adapter(t *testing.T) {
	rt := &MockRoundTripper{}
	rt.Responses = append(rt.Responses, &http.Response{StatusCode: http.StatusOK})

	l := NewLocation(NewConfig().SetCompress(false))

	same := "test"
	a := NewS3Adapter(same, same, same, credentials.NewStaticCredentials(same, same, same))

	a.HttpClient = &http.Client{Transport: rt}

	a.Write(l, []byte("Hello world"))

	if len(rt.Requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(rt.Requests))
	}

	requestBody, _ := io.ReadAll(rt.Requests[0].Body)

	if string(requestBody) != "Hello world" {
		t.Fatalf("expected request body to be 'Hello world', got '%s'", string(requestBody))
	}
}

func TestS3AdapterGzip(t *testing.T) {
	rt := &MockRoundTripper{}
	rt.Responses = append(rt.Responses, &http.Response{StatusCode: http.StatusOK})

	l := NewLocation(NewConfig())

	same := "test"
	a := NewS3Adapter(same, same, same, credentials.NewStaticCredentials(same, same, same))

	a.HttpClient = &http.Client{Transport: rt}

	a.Write(l, []byte("Hello world"))

	if len(rt.Requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(rt.Requests))
	}

	requestBody, _ := io.ReadAll(rt.Requests[0].Body)

	g, _ := gzip.NewReader(bytes.NewReader(requestBody))
	outFileContentDecoded, _ := io.ReadAll(g)

	if !reflect.DeepEqual(outFileContentDecoded, []byte("Hello world")) {
		t.Fatalf("unexcepted file content, got %b", outFileContentDecoded)
	}
	if a.Bytes() != nil {
		t.Fatalf("expected nil bytes, got %b", a.Bytes())
	}
}
