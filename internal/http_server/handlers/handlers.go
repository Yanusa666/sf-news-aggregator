package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
	"sf-news-aggregator/internal/config"
	"sf-news-aggregator/internal/constants"
	"sf-news-aggregator/internal/news"
	"strconv"
	"time"
)

type Handler struct {
	cfg  *config.Config
	lgr  zerolog.Logger
	news *news.News
}

func NewHandler(cfg *config.Config, lgr zerolog.Logger, news *news.News) *Handler {
	return &Handler{
		cfg:  cfg,
		lgr:  lgr,
		news: news,
	}
}

func (h *Handler) GetNews(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	countStr := ps.ByName("count")
	lgr := h.lgr.With().
		Str("handler", "GetNews").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Str("count", countStr)).
		Logger()

	count, err := strconv.ParseUint(countStr, 10, 64)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect count: %s", countStr)})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	items, err := h.news.Model.GetLast(count)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	news := make([]NewEntity, 0, 10)
	for _, i := range items {
		news = append(news, NewEntity{
			Id:      i.Id,
			Title:   i.Title,
			Link:    i.Link,
			Desc:    i.Desc,
			PubDate: i.PubDate.Format(time.RFC3339),
		})
	}

	lgr.Debug().Msg("executed")

	resp, _ := json.Marshal(GetNewsResp{News: news})
	fmt.Fprintf(w, string(resp))
}

func (h *Handler) GetNew(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	newIdStr := ps.ByName("new_id")
	lgr := h.lgr.With().
		Str("handler", "GetNew").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Str("new_id", newIdStr)).
		Logger()

	newId, err := strconv.ParseUint(newIdStr, 10, 64)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect new_id: %s", newIdStr)})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	item, err := h.news.Model.Get(newId)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	resp, _ := json.Marshal(&GetNewResp{
		Id:      item.Id,
		Title:   item.Title,
		Desc:    item.Desc,
		PubDate: item.PubDate.Format(time.RFC3339),
		Link:    item.Link,
	})

	lgr.Debug().Msg("executed")

	fmt.Fprintf(w, string(resp))
}
