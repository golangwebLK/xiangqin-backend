package migrate

import (
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	xiangqin_backend "xiangqin-backend"
	"xiangqin-backend/pkg/candidate"
	"xiangqin-backend/pkg/company"
	"xiangqin-backend/pkg/user"
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
	if err = db.Migrator().AutoMigrate(
		&candidate.PersonalInfo{},
		&candidate.PersonalLike{},
		&user.User{},
		&user.Content{},
		&user.Permission{},
		&company.Company{},
	); err != nil {
		return err
	}
	systemCode := uuid.NewString()
	customerCode := uuid.NewString()
	companyCode := uuid.NewString()
	userCode := uuid.NewString()
	customerDataCode := uuid.NewString()
	contents := [...]user.Content{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name:       "系统管理",
			Logo:       "icon-settings",
			Code:       systemCode,
			ParentCode: "",
			Target:     "",
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Name:       "企业管理",
			Logo:       "icon-storage",
			Code:       companyCode,
			ParentCode: systemCode,
			Target:     "/company",
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			Name:       "用户管理",
			Logo:       "icon-user",
			Code:       userCode,
			ParentCode: systemCode,
			Target:     "/user",
		},
		{
			Model: gorm.Model{
				ID: 4,
			},
			Name:       "客户管理",
			Logo:       "icon-user-group",
			Code:       customerCode,
			ParentCode: "",
			Target:     "",
		},
		{
			Model: gorm.Model{
				ID: 5,
			},
			Name:       "客户信息匹配",
			Logo:       "icon-book",
			Code:       customerDataCode,
			ParentCode: customerCode,
			Target:     "/customer",
		},
	}
	if err = db.Create(&contents).Error; err != nil {
		return err
	}

	permission := [...]user.Permission{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Role:      "SuperManager",
			ContentID: companyCode,
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Role:      "Manager",
			ContentID: userCode,
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			Role:      "Manager",
			ContentID: customerDataCode,
		},
		{
			Model: gorm.Model{
				ID: 4,
			},
			Role:      "User",
			ContentID: customerDataCode,
		},
	}
	if err = db.Create(&permission).Error; err != nil {
		return err
	}
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte("123456"), 12)
	superUser := user.User{
		Model: gorm.Model{
			ID: 1,
		},
		Name:        "金康网络科技",
		Username:    "root",
		Password:    string(passwordHash),
		IsUser:      true,
		Role:        "SuperManager",
		CompanyCode: uuid.NewString(),
	}
	if err = db.Create(&superUser).Error; err != nil {
		return err
	}
	return nil
}
