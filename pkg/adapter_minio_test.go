package stm

import (
	"bytes"
	"compress/gzip"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestMinioAdapter(t *testing.T) {
	rt := &MockRoundTripper{}
	rt.Responses = append(
		rt.Responses,
		&http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(
				bytes.NewReader(
					[]byte(
						`<?xml version="1.0" encoding="UTF-8"?><main>1</main>`,
					),
				),
			),
		},
		&http.Response{StatusCode: http.StatusOK},
	)

	l := NewLocation(NewConfig().SetCompress(false))

	same := "test"
	a := NewMinioAdapter(same, same, credentials.NewStaticV4(same, same, same))

	a.Transport = rt

	a.Write(l, []byte("Hello world"))

	if len(rt.Requests) != 2 {
		t.Fatalf("expected 2 request, got %d", len(rt.Requests))
	}

	requestBody, _ := io.ReadAll(rt.Requests[1].Body)
	requestBodyRows := strings.Split(string(requestBody), "\n")

	if strings.TrimSpace(requestBodyRows[1]) != "Hello world" {
		t.Fatalf("expected request body to be 'Hello world', got '%s'", string(requestBody))
	}
	if a.Bytes() != nil {
		t.Fatalf("expected nil bytes, got %b", a.Bytes())
	}
}

func TestMinioAdapterGzip(t *testing.T) {
	rt := &MockRoundTripper{}
	rt.Responses = append(
		rt.Responses,
		&http.Response{
			StatusCode: http.StatusOK,
			Body: io.NopCloser(
				bytes.NewReader(
					[]byte(
						`<?xml version="1.0" encoding="UTF-8"?><main>1</main>`,
					),
				),
			),
		},
		&http.Response{StatusCode: http.StatusOK},
	)

	l := NewLocation(NewConfig())

	same := "test"
	a := NewMinioAdapter(same, same, credentials.NewStaticV4(same, same, same))

	a.Transport = rt

	a.Write(l, []byte("Hello world"))

	if len(rt.Requests) != 2 {
		t.Fatalf("expected 2 request, got %d", len(rt.Requests))
	}

	requestBody, _ := io.ReadAll(rt.Requests[1].Body)
	requestBodyRows := strings.Split(string(requestBody), "\n")

	g, _ := gzip.NewReader(bytes.NewReader([]byte(requestBodyRows[1])))
	outFileContentDecoded, _ := io.ReadAll(g)

	if string(outFileContentDecoded) != "Hello world" {
		t.Fatalf("expected request body to be 'Hello world', got '%s'", string(requestBody))
	}
	if a.Bytes() != nil {
		t.Fatalf("expected nil bytes, got %b", a.Bytes())
	}
}
