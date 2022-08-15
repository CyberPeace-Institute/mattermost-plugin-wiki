package sqlstore

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/jmoiron/sqlx"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/pkg/errors"
)

type Migration struct {
	fromVersion   semver.Version
	toVersion     semver.Version
	migrationFunc func(sqlx.Ext, *SQLStore) error
}

const MySQLCharset = "DEFAULT CHARACTER SET utf8mb4"

var migrations = []Migration{
	{
		fromVersion: semver.MustParse("0.0.0"),
		toVersion:   semver.MustParse("0.1.0"),
		migrationFunc: func(e sqlx.Ext, sqlStore *SQLStore) error {
			if e.DriverName() == model.DatabaseDriverMysql {
				if _, err := e.Exec(`
					CREATE TABLE IF NOT EXISTS CPI_Wiki_System (
						SKey VARCHAR(64) PRIMARY KEY,
						SValue VARCHAR(1024) NULL
					)
				` + MySQLCharset); err != nil {
					return errors.Wrapf(err, "failed creating table CPI_Wiki_System")
				}

				if _, err := e.Exec(`
					CREATE TABLE IF NOT EXISTS CPI_WikiDocs (
						ID VARCHAR(26) PRIMARY KEY,
						Name VARCHAR(1024) NOT NULL,
						Content VARCHAR(32765) NOT NULL,
						Description VARCHAR(4096) NOT NULL,
						Status VARCHAR(26) NOT NULL,
						OwnerUserID VARCHAR(26) NOT NULL,
						TeamID VARCHAR(26) NOT NULL,
						ChannelID VARCHAR(26) NOT NULL,
						CreateAt BIGINT NOT NULL,
						UpdateAt BIGINT NOT NULL DEFAULT 0,
						DeleteAt BIGINT NOT NULL DEFAULT 0,
						INDEX CPI_WikiDocs_TeamID (TeamID)
						INDEX CPI_WikiDocs_ChannelID (ChannelID)
					)
				` + MySQLCharset); err != nil {
					return errors.Wrapf(err, "failed creating table CPI_WikiDocs")
				}

			} else {
				if _, err := e.Exec(`
					CREATE TABLE IF NOT EXISTS CPI_Wiki_System (
						SKey VARCHAR(64) PRIMARY KEY,
						SValue VARCHAR(1024) NULL
					);
				`); err != nil {
					return errors.Wrapf(err, "failed creating table CPI_Wiki_System")
				}

				if _, err := e.Exec(`
					CREATE TABLE IF NOT EXISTS CPI_WikiDocs (
						ID TEXT PRIMARY KEY,
						Name TEXT NOT NULL,
						Content TEXT NOT NULL,
						Description TEXT NOT NULL,
						Status TEXT NOT NULL,
						OwnerUserID TEXT NOT NULL,
						TeamID TEXT NOT NULL,
						ChannelID TEXT NOT NULL,
						CreateAt BIGINT NOT NULL,
						UpdateAt BIGINT NOT NULL DEFAULT 0,
						DeleteAt BIGINT NOT NULL DEFAULT 0
					);
				`); err != nil {
					return errors.Wrapf(err, "failed creating table CPI_WikiDocs")
				}

				if _, err := e.Exec(createPGIndex("CPI_WikiDocs_TeamID", "CPI_WikiDocs", "TeamID")); err != nil {
					return errors.Wrapf(err, "failed creating index CPI_WikiDocs_TeamID")
				}

				if _, err := e.Exec(createPGIndex("CPI_WikiDocs_ChannelID", "CPI_WikiDocs", "ChannelID")); err != nil {
					return errors.Wrapf(err, "failed creating index CPI_WikiDocs_ChannelID")
				}
			}

			return nil
		},
	},
}

// 'IF NOT EXISTS' syntax is not supported in Postgres 9.4, so we need
// this workaround to make the migration idempotent
var createPGIndex = func(indexName, tableName, columns string) string {
	return fmt.Sprintf(`
		DO
		$$
		BEGIN
			IF to_regclass('%s') IS NULL THEN
				CREATE INDEX %s ON %s (%s);
			END IF;
		END
		$$;
	`, indexName, indexName, tableName, columns)
}
