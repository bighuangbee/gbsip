package model

import (
	"context"
	"errors"
	"gosip/data"
	"gosip/data/domain"
)

type StreamRepo struct {
	Data *data.Data
}

//列表
func (this *StreamRepo) List(ctx context.Context)(list []domain.Stream, total int64, err error){
	db := this.Data.DB(ctx).Model(&domain.Stream{})

	if err = db.Count(&total).Error; err != nil || total == 0{
		return nil, 0, err
	}

	err = db.Find(&list).Error
	return nil, 0, err
}

//更新
func (this *StreamRepo) Update(ctx context.Context, data *domain.Stream) (err error){
	db := this.Data.DB(ctx)
	if data.Id > 0{
		db = db.Where("id = ?", data.Id)
	}else if data.DeviceId != "" && data.ChannelId != ""{
		db = db.Where("device_id = ? and channel_id = ?", data.DeviceId, data.ChannelId)
	}else{
		return errors.New("缺少Device channelId或id")
	}
	return db.Updates(data).Error
}

//新增
func (this *StreamRepo) Create(ctx context.Context, data *domain.Stream) (err error){
	return this.Data.DB(ctx).Create(&data).Error
}

//获取详情
func (this *StreamRepo) GetByDeviceId(ctx context.Context, req *domain.ChannelQuery) (data *domain.Stream, err error){
	db := this.Data.DB(ctx)
	if req.Id > 0{
		db = db.Where("id = ?", req.Id)
	}else if req.DeviceId != "" && req.ChannelId != ""{
		db = db.Where("device_id = ?", req.DeviceId).Where("channel_id = ?", req.ChannelId)
	}else{
		return data, errors.New("缺少Device channelId或id")
	}
	var raw domain.Stream
	err = db.Find(&raw).Error
	return &raw, err
}
