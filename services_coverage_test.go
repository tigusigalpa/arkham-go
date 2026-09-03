package arkham

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var contextInterface = reflect.TypeOf((*context.Context)(nil)).Elem()

// TestServiceEndpointsIssueRequests exercises every REST service method against
// a local server. The API response is JSON null, which can be decoded into any
// response type while keeping this contract test independent of response models.
func TestServiceEndpointsIssueRequests(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if got := r.Header.Get("API-Key"); got != testAPIKey {
			t.Errorf("API-Key header = %q, want %q", got, testAPIKey)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("null"))
	}))
	defer server.Close()

	client, err := NewClient(testAPIKey, WithBaseURL(server.URL), WithMaxRetries(0))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}

	clientValue := reflect.ValueOf(client).Elem()
	for i := 0; i < clientValue.NumField(); i++ {
		service := clientValue.Field(i)
		if !isService(service) {
			continue
		}

		for j := 0; j < service.NumMethod(); j++ {
			method := service.Type().Method(j)
			if clientValue.Type().Field(i).Name == "Streams" && method.Name == "Connect" {
				continue
			}

			boundMethod := service.Method(j)
			results := boundMethod.Call(endpointArguments(boundMethod.Type()))
			if err := results[len(results)-1].Interface(); err != nil {
				t.Errorf("%s.%s returned an unexpected error: %v", clientValue.Type().Field(i).Name, method.Name, err)
			}
		}
	}

	if requests == 0 {
		t.Fatal("service methods did not issue any HTTP requests")
	}
}

func isService(value reflect.Value) bool {
	return value.Kind() == reflect.Pointer && !value.IsNil() && strings.HasSuffix(value.Type().Elem().Name(), "Service")
}

func endpointArguments(methodType reflect.Type) []reflect.Value {
	arguments := make([]reflect.Value, methodType.NumIn())
	for i := range arguments {
		arguments[i] = endpointArgument(methodType.In(i))
	}
	return arguments
}

func endpointArgument(argumentType reflect.Type) reflect.Value {
	if argumentType == contextInterface {
		return reflect.ValueOf(context.Background())
	}

	switch argumentType.Kind() {
	case reflect.String:
		return reflect.ValueOf("test").Convert(argumentType)
	case reflect.Bool:
		return reflect.ValueOf(true).Convert(argumentType)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.ValueOf(int64(1)).Convert(argumentType)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return reflect.ValueOf(uint64(1)).Convert(argumentType)
	case reflect.Slice:
		value := reflect.MakeSlice(argumentType, 1, 1)
		value.Index(0).Set(endpointArgument(argumentType.Elem()))
		return value
	case reflect.Pointer:
		value := reflect.New(argumentType.Elem())
		populateEndpointValue(value.Elem())
		return value
	case reflect.Struct:
		value := reflect.New(argumentType).Elem()
		populateEndpointValue(value)
		return value
	default:
		return reflect.Zero(argumentType)
	}
}

func populateEndpointValue(value reflect.Value) {
	if value.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if !field.CanSet() || field.Kind() == reflect.Pointer {
			continue
		}
		if field.Kind() == reflect.Struct {
			populateEndpointValue(field)
			continue
		}
		field.Set(endpointArgument(field.Type()))
	}
}
