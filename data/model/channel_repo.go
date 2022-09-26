package model

import (
	"context"
	"gorm.io/gorm/clause"
	"gosip/data"
	"gosip/data/domain"
)

type ChannelRepo struct {
	Data *data.Data
}

//列表
func (this *ChannelRepo) List(ctx context.Context, deviceId string)(list []domain.Channel, total int64, err error){
	db := this.Data.DB(ctx).Model(&domain.Channel{})
	if deviceId != ""{
		db = db.Where("device_id = ?", deviceId)
	}

	if err = db.Count(&total).Error; err != nil || total == 0{
		return nil, 0, err
	}

	err = db.Find(&list).Error
	return nil, 0, err
}

//插入或更新
func (this *ChannelRepo) Save(ctx context.Context, data *domain.Channel) (err error){
	return this.Data.DB(ctx).Clauses(clause.OnConflict{
		//Columns:   []clause.Column{{Name: "deviceId"}}, //  唯一索引
		DoUpdates: clause.AssignmentColumns([]string{"name", "manufacturer", "model", "status", "host_address", "ip", "port", "expires", "charset", "updated_at", "deleted_at"}), // 更新哪些字段
	}).Create(data).Error
}

//新增
func (this *ChannelRepo) Create(ctx context.Context, data *domain.Channel) (err error){
	return this.Data.DB(ctx).Create(&data).Error
}

//获取详情
func (this *ChannelRepo) GetByDeviceId(ctx context.Context, deviceId string, channelId string) (data *domain.Channel, err error){
	raw := &domain.Channel{}
	err = this.Data.DB(ctx).Where("device_id = ?", deviceId).Where("channel_id = ?", channelId).Find(raw).Error
	return raw, err
}
