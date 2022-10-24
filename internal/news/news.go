package news

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"net/url"
	"sf-news-aggregator/internal/config"
	"time"
)

type News struct {
	cfg       *config.Config
	lgr       zerolog.Logger
	addresses []*url.URL
	ctx       context.Context
	Model     *Model
}

func NewNews(ctx context.Context, cfg *config.Config, lgr zerolog.Logger) *News {
	addresses := make([]*url.URL, 0, len(cfg.RSS))
	for _, rawURL := range cfg.RSS {
		addr, err := url.Parse(rawURL)
		if err != nil {
			lgr.Fatal().Str("URL", rawURL).Msg("incorrect RSS URL")
		}

		addresses = append(addresses, addr)
	}

	model := NewModel(cfg, lgr)

	n := &News{
		ctx:       ctx,
		cfg:       cfg,
		lgr:       lgr,
		addresses: addresses,
		Model:     model,
	}

	go n.initEnrichment()

	return n
}

func (n *News) Shutdown() error {
	n.Model.connPool.Close()
	return nil
}

func (n *News) collectFeeds(addr *url.URL) (*Rss, error) {
	lgr := n.lgr.With().Str("addr", addr.String()).Logger()

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			DisableCompression:  true,
			MaxIdleConnsPerHost: -1,
		},
	}

	req, err := http.NewRequest("GET", addr.String(), nil)
	if err != nil {
		lgr.Error().Err(err).Msg("http.NewRequest failed")
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		lgr.Error().Err(err).Msg("client.Do failed")
		return nil, err
	}
	defer resp.Body.Close()

	respStatus := resp.StatusCode
	if respStatus != http.StatusOK {
		err = fmt.Errorf("incorrect status %d", respStatus)
		lgr.Error().Err(err).Msg("incorrect response")
		return nil, err
	}

	rss := new(Rss)
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(rss)
	if err != nil {
		lgr.Error().Err(err).Msg("decode xml failed")
		return nil, err
	}

	//lgr.Debug().Interface("rss", rss).Msg("receive RSS") // для отладки обогощения бд RSS записями

	return rss, nil
}

func (n *News) enrich() {
	for _, addr := range n.addresses {
		go func(addr *url.URL) {
			rss, err := n.collectFeeds(addr)
			if err != nil {
				return
			}

			for _, item := range rss.Channel.Items {
				_ = n.Model.Add(&item)
			}

		}(addr)
	}
}

func (n *News) initEnrichment() {
	ticker := time.NewTicker(time.Duration(n.cfg.RequestPeriod) * time.Second)

	for {
		select {
		case <-n.ctx.Done():
			n.lgr.Debug().Msg("enrichment end")
			return

		case <-ticker.C:
			n.enrich()
		}
	}
}
