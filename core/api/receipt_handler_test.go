package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/RA341/receipt-processor-challenge/models"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	// Invalid character '!'
	retailerIncorrect = `{
  "retailer": "Walgreens!", 
  "purchaseDate": "2022-01-02",
  "purchaseTime": "08:13",
  "total": "2.65",
  "items": [
    {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
    {"shortDescription": "Dasani", "price": "1.40"}
  ]
}`

	// Invalid format (MM-DD-YYYY instead of YYYY-MM-DD)
	purchaseDateIncorrect = `{
  "retailer": "Walgreens",
  "purchaseDate": "01-02-2022",
  "purchaseTime": "08:13",
  "total": "2.65",
  "items": [
    {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
    {"shortDescription": "Dasani", "price": "1.40"}
  ]
}`

	// Invalid format (includes AM/PM, not pure HH:MM)
	purchaseTimeIncorrect = `{
  "retailer": "Walgreens",
  "purchaseDate": "2022-01-02",
  "purchaseTime": "08:13 AM", 
  "total": "2.65",
  "items": [
    {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
    {"shortDescription": "Dasani", "price": "1.40"}
  ]
}`

	// Invalid format (only one decimal place)
	totalIncorrect = `{
  "retailer": "Walgreens",
  "purchaseDate": "2022-01-02",
  "purchaseTime": "08:13",
  "total": "2.6",
  "items": [
    {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
    {"shortDescription": "Dasani", "price": "1.40"}
  ]
}`
	// some field removed
	removedFields = `{
  "purchaseDate": "2022-01-02",
  "purchaseTime": "08:13",
  "items": [
    {"shortDescription": "Pepsi - 12-oz", "price": "1.25"},
    {"shortDescription": "Dasani", "price": "1.40"}
  ]
}`
)

func TestReceiptHandler_PostProcessReceipt_200_morning_receipt(t *testing.T) {
	bodyBytes, err := os.ReadFile("../../examples/morning-receipt.json")
	if err != nil {
		t.Fatalf("Failed to load request body: %v", err)
	}
	requestBody := bytes.NewReader(bodyBytes)

	var data models.Receipt
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal request body: %v", err)
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/receipts/process", requestBody)
	resp := httptest.NewRecorder()

	receiptSrv, err := initServices()
	if err != nil {
		t.Fatalf("Failed to init services: %v", err)
	}
	_, handler := NewReceiptHandler(receiptSrv)

	handler.ServeHTTP(resp, req)

	expectedStatus := http.StatusOK
	if status := resp.Code; status != expectedStatus {
		t.Logf("Response body: %s", resp.Body.String())
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	var responseBody models.IdResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Could not unmarshal response body: %v\nBody: %s", err, resp.Body.String())
	}

	// Check specific fields in the response JSON
	if !idRegex.MatchString(responseBody.Id) {
		t.Fatalf("handler returned an id that failed regex: %s, ID was %s", idRegex.String(), responseBody.Id)
	}

	expectedContentType := "application/json" // Assuming your handler sets this
	if ctype := resp.Header().Get("Content-Type"); ctype != expectedContentType {
		t.Fatalf("handler returned wrong Content-Type: got %v want %v",
			ctype, expectedContentType)
	}
}

func TestReceiptHandler_PostProcessReceipt_200_simple_receipt(t *testing.T) {
	bodyBytes, err := os.ReadFile("../../examples/simple-receipt.json")
	if err != nil {
		t.Fatalf("Failed to load request body: %v", err)
	}
	requestBody := bytes.NewReader(bodyBytes)

	var data models.Receipt
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal request body: %v", err)
		return
	}

	req := httptest.NewRequest(http.MethodPost, "/receipts/process", requestBody)
	resp := httptest.NewRecorder()

	receiptSrv, err := initServices()
	if err != nil {
		t.Fatalf("Failed to init services: %v", err)
	}
	_, handler := NewReceiptHandler(receiptSrv)

	handler.ServeHTTP(resp, req)

	expectedStatus := http.StatusOK
	if status := resp.Code; status != expectedStatus {
		t.Logf("Response body: %s", resp.Body.String())
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, expectedStatus)
	}

	var responseBody models.IdResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Could not unmarshal response body: %v\nBody: %s", err, resp.Body.String())
	}

	// Check specific fields in the response JSON
	if !idRegex.MatchString(responseBody.Id) {
		t.Fatalf("handler returned an id that failed regex: %s, ID was %s", idRegex.String(), responseBody.Id)
	}

	expectedContentType := "application/json" // Assuming your handler sets this
	if ctype := resp.Header().Get("Content-Type"); ctype != expectedContentType {
		t.Fatalf("handler returned wrong Content-Type: got %v want %v",
			ctype, expectedContentType)
	}
}

func TestReceiptHandler_PostProcessReceipt_400(t *testing.T) {
	cases := []string{
		purchaseTimeIncorrect,
		purchaseDateIncorrect,
		totalIncorrect,
		retailerIncorrect,
		removedFields,
	}

	receiptSrv, err := initServices()
	if err != nil {
		t.Fatalf("Failed to init services: %v", err)
	}

	_, handler := NewReceiptHandler(receiptSrv)

	for _, c := range cases {
		requestBody := bytes.NewReader([]byte(c))
		req := httptest.NewRequest(http.MethodPost, "/receipts/process", requestBody)
		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		expectedStatus := http.StatusBadRequest
		if status := resp.Code; status != expectedStatus {
			t.Logf("Response body: %s", resp.Body.String())
			fatalErr(t, "handler returned wrong status code", status, expectedStatus)
		}

		body := strings.TrimSpace(resp.Body.String())
		if body != BadRequestErr {
			fatalErr(t, "handler returned wrong body", body, BadRequestErr)
		}

		expectedContentType := "text/plain; charset=utf-8" // Assuming your handler sets this
		if ctype := resp.Header().Get("Content-Type"); ctype != expectedContentType {
			fatalErr(t, "handler returned wrong Content-Type", ctype, expectedContentType)
		}
	}
}

func TestReceiptHandler_GetReceiptPoints_404(t *testing.T) {
	receiptSrv, err := initServices()
	if err != nil {
		t.Fatalf("Failed to init services: %v", err)
	}

	pointId := 6969
	target := fmt.Sprintf("/receipts/%d/points", pointId)
	requestBody := bytes.NewReader([]byte(""))

	req := httptest.NewRequest(http.MethodGet, target, requestBody)
	resp := httptest.NewRecorder()

	_, handler := NewReceiptHandler(receiptSrv)

	handler.ServeHTTP(resp, req)

	expectedStatus := http.StatusNotFound
	if status := resp.Code; status != expectedStatus {
		t.Logf("Response body: %s", resp.Body.String())
		fatalErr(t, "handler returned wrong status code", status, expectedStatus)
	}

	if strings.TrimSpace(resp.Body.String()) != NotFoundErr {
		fatalErr(t, "handler returned an invalid status code", resp.Body.String(), NotFoundErr)
	}

	expectedContentType := "text/plain; charset=utf-8" // Assuming your handler sets this
	if ctype := resp.Header().Get("Content-Type"); ctype != expectedContentType {
		fatalErr(t, "handler returned wrong Content-Type: got %v want %v",
			ctype, expectedContentType)
	}
}

func TestReceiptHandler_GetReceiptPoints_morning_receipt(t *testing.T) {
	bodyBytes, err := os.ReadFile("../../examples/morning-receipt.json")
	if err != nil {
		t.Fatalf("Failed to load request body: %v", err)
	}
	runTestGetPoints(t, 15, string(bodyBytes))
}

func TestReceiptHandler_GetReceiptPoints_simple_receipt(t *testing.T) {
	bodyBytes, err := os.ReadFile("../../examples/simple-receipt.json")
	if err != nil {
		t.Fatalf("Failed to load request body: %v", err)
	}
	runTestGetPoints(t, 31, string(bodyBytes))
}

func runTestGetPoints(t *testing.T, expected int64, payload string) {
	receiptSrv, err := initServices()
	if err != nil {
		t.Fatalf("Failed to init services: %v", err)
	}
	var data models.Receipt
	err = json.Unmarshal([]byte(payload), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal request body: %v", err)
	}

	pointId, err := receiptSrv.NewReceipt(data)
	if err != nil {
		t.Fatalf("Failed to create receipt: %v", err)
		return
	}

	target := fmt.Sprintf("/receipts/%s/points", pointId)
	requestBody := bytes.NewReader([]byte(""))

	req := httptest.NewRequest(http.MethodGet, target, requestBody)
	resp := httptest.NewRecorder()

	_, handler := NewReceiptHandler(receiptSrv)

	handler.ServeHTTP(resp, req)

	expectedStatus := http.StatusOK
	if status := resp.Code; status != expectedStatus {
		t.Logf("Response body: %s", resp.Body.String())
		fatalErr(t, "handler returned wrong status code", status, expectedStatus)
	}

	var responseBody models.PointsResponse
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	if err != nil {
		fatalErr(t, "Could not unmarshal response body", err, resp.Body.String())
	}

	if responseBody.Points != expected {
		fatalErr(t, "handler returned an id that failed regex", responseBody.Points, expected)
	}

	expectedContentType := "application/json" // Assuming your handler sets this
	if ctype := resp.Header().Get("Content-Type"); ctype != expectedContentType {
		fatalErr(t, "handler returned wrong Content-Type: got %v want %v",
			ctype, expectedContentType)
	}
}

func fatalErr(t *testing.T, message string, got any, want any) {
	t.Fatalf("%s: \ngot: %v\nwant: %v", message, got, want)
}
