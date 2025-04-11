package api

import (
	"fmt"
	"github.com/RA341/receipt-processor-challenge/service"
	"log/slog"
	"net/http"
	"os"
)

func StartServer(addr, port string) {
	mux := http.NewServeMux()
	registerEndpoints(mux)

	finalAddr := fmt.Sprintf("%s:%s", addr, port)
	slog.Info("Server listening on ", slog.String("addr", finalAddr))
	if err := http.ListenAndServe(finalAddr, nil); err != nil {
		slog.Error("Unable to start server ", err.Error())
		os.Exit(1)
	}
}

func registerEndpoints(mux *http.ServeMux) {
	receiptSrv, err := InitServices()
	if err != nil {
		slog.Error("Unable to initialize services: ", err.Error())
		os.Exit(1)
	}

	baseRoute, rHandler := NewReceiptHandler(receiptSrv)
	mux.Handle(baseRoute, rHandler)
}

func InitServices() (*service.ReceiptService, error) {
	db, err := service.NewDB()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	srv := service.NewReceiptService(db)
	return srv, nil
}
