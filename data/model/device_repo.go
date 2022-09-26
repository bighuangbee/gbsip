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
func (this *DeviceRepo) List(ctx context.Context)(list []domain.Device, total int64, err error){
	db := this.Data.DB(ctx).Model(&domain.Device{})
	if err = db.Count(&total).Error; err != nil || total == 0{
		return nil, 0, err
	}

	err = db.Find(&list).Error
	return nil, 0, err
}

//插入或更新
func (this *DeviceRepo) Save(ctx context.Context, data *domain.Device) (err error){
	return this.Data.DB(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "deviceId"}}, //  唯一索引
		DoUpdates: clause.AssignmentColumns([]string{"name", "manufacturer", "model", "firmware", "transport", "status", "last_alive_time", "host_address", "ip", "port", "expires", "charset", "updated_at", "deleted_at"}), // 更新哪些字段
	}).Create(data).Error
}

//新增
func (this *DeviceRepo) Create(ctx context.Context, data *domain.Device) (err error){
	return this.Data.DB(ctx).Create(&data).Error
}

//获取详情
func (this *DeviceRepo) GetByDeviceId(ctx context.Context, deviceId string) (data *domain.Device, err error){
	var raw domain.Device
	err = this.Data.DB(ctx).Where("device_id = ?", deviceId).Find(&raw).Error
	return &raw, err
}

