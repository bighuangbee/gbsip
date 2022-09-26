package domain

import "context"

type IStreamsRepo interface {
	//列表
	List(ctx context.Context)(list []Stream, total int64, err error)
	//插入或更新
	//Save(ctx context.Context, data *Streams) (err error)
	//插入或更新
	Update(ctx context.Context, data *Stream) (err error)
	//新增
	Create(ctx context.Context, data *Stream) (err error)
	//获取详情
	GetByDeviceId(ctx context.Context, req *ChannelQuery) (data *Stream, err error)
}

type Stream struct {
	Id int64 `json:"id,omitempty"`
	T int32 `json:"t,omitempty"`
	DeviceId string `json:"deviceId,omitempty"`
	ChannelId string `json:"channelId,omitempty"`
	StreamType string `json:"streamType,omitempty"`
	Status int32 `json:"status,omitempty"`
	Callid string `json:"callid,omitempty"`
	Stop int8 `json:"stop,omitempty"`
	Msg string `json:"msg,omitempty"`
	Cseqno uint32 `json:"cseqno,omitempty"`
	StreamId string `json:"stream_id,omitempty"`
	Hls string `json:"hls,omitempty"`
	Rtmp string `json:"rtmp,omitempty"`
	Rtsp string `json:"rtsp,omitempty"`
	Wsflv string `json:"wsflv,omitempty"`
	Stream int8 `json:"stream,omitempty"` // 1=实时流，0=历史流
}
