package section

import "time"

type (
	Repository struct {
		Postgres RepositoryPostgres `env:"POSTGRES"`
	}

	RepositoryPostgres struct {
		Address      string        `env:"ADDRESS" required:"true"`
		Username     string        `env:"USERNAME"`
		Password     string        `env:"Password" required:"true"`
		Name         string        `env:"Name" default:"catalog"`
		ReadTimeout  time.Duration `env:"READ_TIMEOUT" default:"5s" split_words:"true"`
		WriteTimeout time.Duration `env:"WRITE_TIMEOUT" default:"5s" split_words:"true"`
	}
)
