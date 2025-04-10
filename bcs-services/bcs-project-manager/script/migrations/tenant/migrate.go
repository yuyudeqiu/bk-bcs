package tenant

import (
	"context"
	"fmt"

	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/common/page"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/config"
	"github.com/Tencent/bk-bcs/bcs-services/bcs-project-manager/internal/store"
)

var (
	model store.ProjectModel
)

func initDB() error {
	// mongo
	store.InitMongo(&config.MongoConfig{
		Address:        config.GlobalConf.Mongo.Address,
		Replicaset:     config.GlobalConf.Mongo.Replicaset,
		ConnectTimeout: config.GlobalConf.Mongo.ConnectTimeout,
		Database:       config.GlobalConf.Mongo.Database,
		Username:       config.GlobalConf.Mongo.Username,
		Password:       config.GlobalConf.Mongo.Password,
		MaxPoolSize:    config.GlobalConf.Mongo.MaxPoolSize,
		MinPoolSize:    config.GlobalConf.Mongo.MinPoolSize,
		Encrypted:      config.GlobalConf.Mongo.Encrypted,
	})
	model = store.New(store.GetMongo())
	return nil
}

func InitProject() error {
	err := initDB()
	if err != nil {
		return err
	}
	// 使用 model，在 porject 表中，把每个记录中，如果没有tenantId的记录，加上 tenantId，值为default
	ctx := context.Background()

	// 查询所有记录
	projects, _, err := model.ListProjects(ctx, nil, &page.Pagination{All: true})
	if err != nil {
		return fmt.Errorf("failed to list projects: %v", err)
	}

	// 更新每条记录
	for _, p := range projects {
		// 如果 TenantID 字段为空字符串，说明这个字段不存在
		if p.TenantID == "" {
			p.TenantID = "default"
			err = model.UpdateProject(ctx, &p)
			if err != nil {
				return fmt.Errorf("failed to update project %s: %v", p.ProjectID, err)
			}
		}
	}

	return nil
}
