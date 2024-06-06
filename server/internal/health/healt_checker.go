package health

import (
	"database/sql"
	"net/http"
)

type HealthChecker struct {
	DB *sql.DB
}

func NewHealthChecker(db *sql.DB) *HealthChecker {
	return &HealthChecker{DB: db}
}

func (h *HealthChecker) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.DB.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error pinging the database"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
