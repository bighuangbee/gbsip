package gb

import (
	"encoding/xml"
)

const EncodeUTF8 = "UTF-8"
const EncodeGB2312 = "GB2312"

const MANSCDP = "Application/MANSCDP+xml"
const SDP = "application/sdp"

const Header = `<?xml version="1.0" encoding="%s"?>` + "\n"

type Payload struct {
	CmdType string `xml:"CmdType"`
	SN string `xml:"SN"`
	DeviceID string `xml:"DeviceID"`
}

type Keepalive struct {
	Payload
	Status string `xml:"Status"`
}

type Query struct {
	XMLName xml.Name `xml:"Query"`
	Payload
}

//下级SIP发送的设备信息
type CatalogResponse struct {
	XMLName xml.Name `xml:"Response"`
	Payload
	SumNum	string
	DeviceList struct {
		Channels []Channel  `xml:"Item"`
	}  `xml:"DeviceList"`
}


type Channel struct {
	DeviceID string
	Name string
	Manufacturer string
	Model string
	Owner       string
	CivilCode   string
	Address     string
	Parental    string
	ParentID    string
	RegisterWay string
	Secrecy     int
	StreamNum   int
	IPAddress   string
	Status      string
	Info        Channelnfo
}
type Channelnfo struct {
	PTZType int8
	DownloadSpeed string
}

//请求实时播放
type PlayReq struct {
	DeviceId string
	ChannelId string
	Addr	string
	Port	uint16
}
