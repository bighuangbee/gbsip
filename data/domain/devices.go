package domain

import (
	"context"
	"time"
)

type IDeviceRepo interface {
	//设备列表
	List(ctx context.Context)(list []Devices, total int64, err error)
	//插入或更新
	Save(ctx context.Context, data *Devices) (err error)
	//新增
	Create(ctx context.Context, data *Devices) (err error)
	//获取详情
	GetByDeviceId(ctx context.Context, deviceId string) (devices *Devices, err error)
}

type Devices struct {
	Id int64 `json:"id,omitempty"`
	DeviceId string `json:"deviceId,omitempty"` // 国标设备id
	Name string `json:"name,omitempty"` // 国标设备名称
	Manufacturer string `json:"manufacturer,omitempty"` // 厂商
	Model string `json:"model,omitempty"` // 型号
	Firmware string `json:"firmware,omitempty"` // 固件版本
	Transport string `json:"transport,omitempty"` // 信令传输 UDP
	Status uint8 `json:"status,omitempty"` // 在线状态 1=在线 2=离线
	LastAliveTime time.Time `json:"lastAliveTime,omitempty"` // 最近一次心跳时间
	HostAddress string `json:"hostAddress,omitempty"` // IPC地址
	Ip string `json:"ip,omitempty"` // IPC IP地址
	Port uint16 `json:"port,omitempty"` // IPC端口
	Expires int32 `json:"expires,omitempty"` // 通道更新周期
	Charset string `json:"charset,omitempty"` // 字符集 GB2312/UTF-8
	CreatedAt time.Time `json:"createdAt,omitempty"` // 创建时间,注册时间
	UpdatedAt *time.Time `json:"updatedAt,omitempty"` // 修改时间
	DeletedAt *time.Time `json:"deletedAt,omitempty"` // 删除时间
}
