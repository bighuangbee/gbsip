package model

import (
	"context"
	"gorm.io/gorm/clause"
	"gosip/data"
	"gosip/data/domain"
)

type DeviceRepo struct {
	Data *data.Data
}

//设备列表
func (this *DeviceRepo) List(ctx context.Context)(list []domain.Devices, total int64, err error){
	db := this.Data.DB(ctx)
	if err = db.Count(&total).Error; err != nil || total == 0{
		return nil, 0, err
	}

	err = db.Find(&list).Error
	return nil, 0, err
}

//插入或更新
func (this *DeviceRepo) Save(ctx context.Context, data *domain.Devices) (err error){
	return this.Data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "deviceId"}}, //  唯一索引
		DoUpdates: clause.AssignmentColumns([]string{"name", "manufacturer", "model", "firmware", "transport", "status", "last_alive_time", "host_address", "ip", "port", "expires", "charset", "updated_at", "deleted_at"}), // 更新哪些字段
	}).Create(data).Error
}

//新增
func (this *DeviceRepo) Create(ctx context.Context, data *domain.Devices) (err error){
	return this.Data.DB(ctx).Create(&data).Error
}

//获取详情
func (this *DeviceRepo) GetByDeviceId(ctx context.Context, deviceId string) (data *domain.Devices, err error){
	err = this.Data.DB(ctx).Where("deviceId = ?", deviceId).Find(&data).Error
	return nil, err
}

