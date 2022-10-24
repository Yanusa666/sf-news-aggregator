package news

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog"
	"sf-news-aggregator/internal/config"
	"time"
)

type Model struct {
	cfg      *config.Config
	lgr      zerolog.Logger
	connPool *pgxpool.Pool
}

func NewModel(cfg *config.Config, lgr zerolog.Logger) *Model {
	lgr = lgr.With().Str("model", "news").Logger()

	pgConf, err := pgxpool.ParseConfig(cfg.Postgres.URI)
	if err != nil {
		lgr.Fatal().Err(err).Msg("failed parse PostgreSQL config")
	}

	pgPool, err := pgxpool.ConnectConfig(context.Background(), pgConf)
	if err != nil {
		lgr.Fatal().Err(err).Msg("failed connect to PostgreSQL")
	}

	err = pgPool.Ping(context.Background())
	if err != nil {
		lgr.Fatal().Err(err).Msg("unsuccessful ping attempt")
	}

	return &Model{
		cfg:      cfg,
		lgr:      lgr,
		connPool: pgPool,
	}
}

func (m *Model) Add(item *Item) error {
	lgr := m.lgr.With().Interface("item", item).Logger()

	pubDate, err := time.Parse(time.RFC1123, item.PubDate)
	if err != nil {
		lgr.Warn().Err(err).Msg("incorrect item pubDate")
		pubDate = time.Now()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = m.connPool.Exec(ctx,
		`INSERT INTO news.news(title, description, pub_date, link) 
			 VALUES ($1,$2,$3,$4)
			 ON CONFLICT ON CONSTRAINT news_link_constraint
			 DO NOTHING`, item.Title, item.Desc, pubDate, item.Link)
	if err != nil {
		lgr.Error().Err(err).Msg("upsert RSS to db failed")
		return err
	}

	//lgr.Debug().Msg("upsert RSS to db") // для отладки обогощения бд RSS записями

	return nil
}

func (m *Model) GetLast(count uint64) ([]ItemDb, error) {
	lgr := m.lgr.With().Uint64("count", count).Logger()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.connPool.Query(ctx,
		`SELECT id, title, description, pub_date, link 
			 FROM news.news
			 ORDER BY pub_date DESC LIMIT $1`, count)
	if err != nil {
		lgr.Error().Err(err).Msg("select RSS from db failed")
		return nil, err
	}

	items := make([]ItemDb, 0, 10)
	for rows.Next() {
		item := ItemDb{}
		err = rows.Scan(&(item.Id), &(item.Title), &(item.Desc), &(item.PubDate), &(item.Link))
		if err != nil {
			lgr.Error().Err(err).Msg("scan from db failed")
			return nil, err
		}
		items = append(items, item)
	}

	lgr.Debug().Msg("get last RSS from db")

	return items, nil
}

func (m *Model) Get(newId uint64) (*ItemDb, error) {
	lgr := m.lgr.With().Uint64("new_id", newId).Logger()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	item := new(ItemDb)
	err := m.connPool.QueryRow(ctx,
		`SELECT id, title, description, pub_date, link
			FROM news.news WHERE id = $1`,
		newId).Scan(&(item.Id), &(item.Title), &(item.Desc), &(item.PubDate), &(item.Link))
	if err != nil {
		lgr.Error().Err(err).Msg("select RSS from db failed")
		return nil, err
	}

	lgr.Debug().Msg("get single RSS from db")

	return item, nil
}
