package api

import (
	"fmt"
	"github.com/RA341/receipt-processor-challenge/service"
	u "github.com/RA341/receipt-processor-challenge/utils"
	"log/slog"
	"net/http"
	"os"
)

func StartServer(addr string) {
	mux := http.NewServeMux()
	registerEndpoints(mux)

	slog.Info("Server listening on", slog.String("addr", addr))
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("Unable to start server ", u.ErrLog(err))
		os.Exit(1)
	}
}

func registerEndpoints(mux *http.ServeMux) {
	receiptSrv, err := initServices()
	if err != nil {
		slog.Error("Unable to initialize services:", u.ErrLog(err))
		os.Exit(1)
	}

	baseRoute, rHandler := NewReceiptHandler(receiptSrv)
	mux.Handle(baseRoute, rHandler)
}

func initServices() (*service.ReceiptService, error) {
	db, err := service.NewDB()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to db: %v", err)
	}

	srv := service.NewReceiptService(db)
	return srv, nil
}
