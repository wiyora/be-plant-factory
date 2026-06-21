package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/tern/v2/migrate"
	"github.com/rizalarfiyan/be-plant-factory/migrations"
	"github.com/rs/zerolog"
)

func AutoMigrate(dsn string, log zerolog.Logger) error {
	log.Info().Msg("Auto-migrating database schema")

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to database for migration")
		return err
	}

	defer func() {
		if err := conn.Close(ctx); err != nil {
			log.Error().Err(err).Msg("Failed to close database connection after migration")
		}
	}()

	if err := conn.Ping(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to ping database for migration")
		return err
	}

	mig, err := migrate.NewMigrator(ctx, conn, "schema_version")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create migrator")
		return err
	}

	mig.OnStart = func(sequence int32, name, direction, sql string) {
		log.Info().Int32("sequence", sequence).Str("name", name).Str("direction", direction).Msg("Starting migration")
	}

	currentVersion, err := mig.GetCurrentVersion(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current migration version")
		return err
	}

	if err := mig.LoadMigrations(migrations.FS); err != nil {
		log.Error().Err(err).Msg("Failed to load migrations")
		return err
	}

	if err := mig.Migrate(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to apply migrations")
		return err
	}

	newVersion, err := mig.GetCurrentVersion(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current migration version")
		return err
	}

	totalUpMigrations := newVersion - currentVersion
	if totalUpMigrations > 0 {
		log.Info().Msgf("Applied %d migrations", totalUpMigrations)
	} else {
		log.Info().Msg("No new migrations to apply")
	}

	return nil
}
