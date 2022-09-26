package domain

import (
	"context"
	"time"
)

type IChannelRepo interface {
	//列表
	List(ctx context.Context, deviceId string)(list []Channel, total int64, err error)
	//插入或更新GetByDeviceId
	Save(ctx context.Context, data *Channel) (err error)
	//新增
	Create(ctx context.Context, data *Channel) (err error)
	//获取详情
	GetByDeviceId(ctx context.Context, deviceId string, channelId string) (data *Channel, err error)
}

type ChannelQuery struct {
	Id        int64  `json:"id,omitempty"`
	DeviceId  string `json:"deviceId,omitempty"`
	ChannelId string `json:"channelId,omitempty"`
}

type Channel struct {
	Id int64 `json:"id,omitempty"`
	DeviceId string `json:"deviceId,omitempty"` // 国标设备id
	ChannelId string `json:"channelId,omitempty"` // 国标通道id
	Name string `json:"name,omitempty"` // 国标设备名称
	Manufacturer string `json:"manufacturer,omitempty"` // 厂商
	Model string `json:"model,omitempty"` // 型号
	Firmware string `json:"firmware,omitempty"` // 固件版本
	PtzType int8 `json:"ptzType,omitempty"` // 云台类型, 0-未知 1-球机 2-半球机 3-固定枪机 4-遥控枪机
	PtzTypeText string `json:"ptzType_text,omitempty"` // 云台类型文本
	Addr string `json:"addr,omitempty"` // 位置
	Status uint8 `json:"status,omitempty"` // 在线状态 1=在线 2=离线
	//Lat string `json:"lat,omitempty"` // 纬度 精确到小数点后8位
	//Lon string `json:"lon,omitempty"` // 经度 精确到小数点后8位
	//Alt string `json:"alt,omitempty"` // 海拔 单位米 精确到小数点后2位
	Owner string `json:"owner,omitempty"`
	ParentId string `json:"parentId,omitempty"`
	RegisterWay int8 `json:"registerWay,omitempty"`
	Secrecy int32 `json:"secrecy,omitempty"`
	StreamNum int32 `json:"streamNum,omitempty"`
	HostAddress string `json:"hostAddress,omitempty"` // IPC地址
	Ip string `json:"ip,omitempty"` // IPC IP地址
	Port uint16 `json:"port,omitempty"` // IPC端口
	Expires int32 `json:"expires,omitempty"` // 通道更新周期
	Charset string `json:"charset,omitempty"` // 字符集 GB2312/UTF-8
	CreatedAt time.Time `json:"createdAt,omitempty"` // 创建时间,注册时间
	UpdatedAt *time.Time `json:"updatedAt,omitempty"` // 修改时间
	DeletedAt *time.Time `json:"deletedAt,omitempty"` // 删除时间
}
