package api

import (
	"context"
	"net/http"
	"time"
)

func (a *APIMetric) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := a.handler.PingDB(ctx); err != nil {
		a.log.LogErr("database ping", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
