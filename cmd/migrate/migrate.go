package migrate

import (
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	xiangqin_backend "xiangqin-backend"
	"xiangqin-backend/pkg/candidate"
)

var (
	StartCmd = &cobra.Command{
		Use:          "migrate",
		Short:        "start migrate",
		Example:      "start migrate",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func run() error {
	cfg, err := xiangqin_backend.GetConfig("config.yaml")
	if err != nil {
		return err
	}
	db, err := gorm.Open(
		postgres.Open(cfg.Postgres.DSN),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		return err
	}
	err = db.Migrator().AutoMigrate(
		&candidate.PersonalInfo{},
		&candidate.Candidate{},
	)

	if err != nil {
		return err
	}
	return nil
}
