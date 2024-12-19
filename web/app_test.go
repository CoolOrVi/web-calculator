package app_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	app "github.com/coolorvi/web-calculator/web"
)

type TestCase struct {
	Name           string
	RequestBody    interface{} // Change to interface{} to allow different request types
	ExpectedStatus int
	ExpectedBody   app.CalcResponse
}

func TestCalculateHandler(t *testing.T) {
	handler := http.HandlerFunc(app.CalculateHandler)
	tests := []TestCase{
		{
			Name: "Valid expression",
			RequestBody: app.CalcRequest{
				Expression: "3 + 5 * 2",
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody: app.CalcResponse{
				Result: "13.000000",
			},
		},
		{
			Name: "Invalid character",
			RequestBody: app.CalcRequest{
				Expression: "3 + 5 * 2a",
			},
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody: app.CalcResponse{
				Error: "Expression is not valid",
			},
		},
		{
			Name: "Division by zero",
			RequestBody: app.CalcRequest{
				Expression: "10 / 0",
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedBody: app.CalcResponse{
				Error: "Internal server error",
			},
		},
		{
			Name: "Empty expression",
			RequestBody: app.CalcRequest{
				Expression: "",
			},
			ExpectedStatus: http.StatusUnprocessableEntity,
			ExpectedBody: app.CalcResponse{
				Error: "Expression is not valid",
			},
		},
		{
			Name:           "Invalid JSON",          // Updated to handle malformed JSON
			RequestBody:    `{"expression": 3 + 5}`, // Invalid JSON (no quotes around expression)
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody: app.CalcResponse{
				Error: "Invalid JSON",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			body, err := json.Marshal(test.RequestBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != test.ExpectedStatus {
				t.Errorf("Expected status %d, got %d", test.ExpectedStatus, rec.Code)
			}

			var respBody app.CalcResponse
			if err := json.NewDecoder(rec.Body).Decode(&respBody); err != nil {
				t.Fatalf("Failed to decode response body: %v", err)
			}

			if respBody != test.ExpectedBody {
				t.Errorf("Expected body %+v, got %+v", test.ExpectedBody, respBody)
			}
		})
	}
}
