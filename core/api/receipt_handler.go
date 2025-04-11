package api

import (
	"encoding/json"
	"fmt"
	"github.com/RA341/receipt-processor-challenge/models"
	"github.com/RA341/receipt-processor-challenge/service"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	retailerPattern = regexp.MustCompile(`^[\w\s\-&]+$`)
	totalPattern    = regexp.MustCompile(`^\d+\.\d{2}$`)
	idRegex         = regexp.MustCompile(`^\S+$`)
)

type ReceiptHandler struct {
	srv *service.ReceiptService
}

// todo consider returning non pointer

func NewReceiptHandler(srv *service.ReceiptService) (string, *ReceiptHandler) {
	return "/receipts/", &ReceiptHandler{srv: srv}
}

// ReceiptsHandler is the main handler for the /receipts path.
func (rh *ReceiptHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		rh.PostProcessReceipt(w, r)
	case http.MethodGet:
		rh.GetReceiptPoints(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		slog.Warn(fmt.Sprintf("Method %s not supported", r.Method))
	}
}

func (rh *ReceiptHandler) GetReceiptPoints(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.URL.Path, "/")
	// The path should look like: "", "receipts", "{id}", "points"
	if !(len(pathSegments) == 4 && pathSegments[3] == "points") {
		http.Error(w, NotFoundErr, http.StatusBadRequest)
		return
	}

	pathId := pathSegments[2]
	if !idRegex.MatchString(pathId) {
		http.Error(w, NotFoundErr, http.StatusBadRequest)
		return
	}

	points, err := rh.srv.GetPointsById(pathId)
	if err != nil {
		http.Error(w, NotFoundErr, http.StatusBadRequest)
		return
	}

	response := models.PointsResponse{Points: points}
	marshal, err := json.Marshal(response)
	if err != nil {
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		slog.Warn("Unable to write response to client", err)
		return
	}
}

func (rh *ReceiptHandler) PostProcessReceipt(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, BadRequestErr, http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Warn("error occurred while closing request body", err)
		}
	}(r.Body)

	// Parse the JSON data into our Receipt struct
	var receipt models.Receipt
	if err := json.Unmarshal(body, &receipt); err != nil {
		http.Error(w, BadRequestErr, http.StatusBadRequest)
		return
	}

	// Validate the required fields
	if err := validateReceipt(receipt); err != nil {
		http.Error(w, BadRequestErr, http.StatusBadRequest)
		return
	}

	receiptId, err := rh.srv.NewReceipt(receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := models.IdResponse{Id: receiptId}
	marshal, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		slog.Warn("Unable to write response to client", err)
	}
}

func validateReceipt(receipt models.Receipt) error {
	if !retailerPattern.MatchString(receipt.Retailer) {
		return fmt.Errorf("invalid retailer name: must contain only alphanumeric characters, spaces, hyphens, and ampersands")
	}

	if _, err := time.Parse("2006-01-02", receipt.PurchaseDate); err != nil {
		return fmt.Errorf("invalid purchaseDate format: must be YYYY-MM-DD")
	}

	if _, err := time.Parse("15:04", receipt.PurchaseTime); err != nil {
		return fmt.Errorf("invalid purchaseTime format: must be in 24-hour format (HH:MM)")
	}

	// Validate items
	if len(receipt.Items) < 1 {
		return fmt.Errorf("at least one item is required")
	}

	// Validate each item
	for i, item := range receipt.Items {
		if item.ShortDescription == "" {
			return fmt.Errorf("item %d is missing a short description", i+1)
		}
		if !totalPattern.MatchString(item.Price) {
			return fmt.Errorf("item %d has an invalid price format: must be in format 0.00", i+1)
		}
	}

	// Validate total (pattern: "^\\d+\\.\\d{2}$")
	if !totalPattern.MatchString(receipt.Total) {
		return fmt.Errorf("invalid total format: must be in format 0.00")
	}

	return nil
}
