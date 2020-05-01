package app

import (
	"context"
	"github.com/barcodepro/go-service-template/service/store"
)

func Start(ctx context.Context, config *Config) error {
	sc := store.Config{
		PostgresURL: config.PostgresURL,
	}
	st, err := store.NewStore(ctx, &sc)
	if err != nil {
		return err
	}
	defer st.Close()

	srv := newServer(config, st)

	return srv.httpserver.ListenAndServe()
}
