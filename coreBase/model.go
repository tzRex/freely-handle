package corebase

import (
	"github.com/tzRex/freely-handle/before"
	"gorm.io/gorm"
)

type IModel interface {
	TableName() string
	GroupName() string
}

// 因为gorm自带的model返回的时间类型格式不规范，所以这里自己重新定义一下
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt customTime     `gorm:"comment:创建时间" json:"createdAt"`
	UpdatedAt customTime     `gorm:"comment:更新时间" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 用于搜集迁移失败的数据
var ModelMigrateErr = []string{}

// 数据库表迁移
func MigrateTable(model IModel) {
	err := before.Database.AutoMigrate(model)
	if err != nil {
		name := model.TableName()
		ModelMigrateErr = append(ModelMigrateErr, name)
	}
}

// 数据库中间表迁移
func MigrateJoinTable(mainModel IModel, col string, followModel IModel) {
	err := before.Database.SetupJoinTable(mainModel, col, followModel)
	if err != nil {
		name := followModel.TableName()
		ModelMigrateErr = append(ModelMigrateErr, name)
	}
}
