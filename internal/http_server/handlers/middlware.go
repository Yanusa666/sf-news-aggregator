package handlers

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sf-news-aggregator/internal/constants"
)

func (h *Handler) Middlware(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		requestId := r.Header.Get(constants.RequestIdKey)
		ctx = context.WithValue(ctx, constants.RequestIdKey, requestId)
		r = r.WithContext(ctx)

		w.Header().Add(constants.RequestIdKey, requestId)
		w.Header().Add("Content-Type", "application/json")

		handle(w, r, ps)
	}
}
