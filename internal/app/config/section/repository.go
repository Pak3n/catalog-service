package section

import "time"

type (
	Repository struct {
		Postgres RepositoryPostgres
	}

	RepositoryPostgres struct {
		Address        string        `required:"true"`
		Username       string        `required:"true"`
		Password       string        `required:"true"`
		Name           string        `default:"catalog" required:"true"`
		ReadTimeout    time.Duration `default:"30s" split_words:"true"`
		WriteTimeout   time.Duration `default:"30s" split_words:"true"`
		MigrationTable string        `split_words:"true" default:"schema_migrations"`
	}
)
