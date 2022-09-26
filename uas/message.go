package uas

import (
	"errors"
	"github.com/jart/gosip/sip"
	"net"
)

//User Agent Client 原始数据
type UacRequest struct {
	uac *net.UDPAddr	//IPC连接地址
	message []byte
}

//User Agent Client 经过sip库解析
type UacMsg struct {
	uacConn *net.UDPAddr //IPC连接地址
	msg     *sip.Msg
}


func (this *UacRequest) ToUacMsg()(*UacMsg, error){
	sipMsg, err := sip.ParseMsg(this.message)
	if err != nil{
		return nil, err
	}
	return &UacMsg{uacConn: this.uac, msg: sipMsg}, nil
}



func (this *UdpServer)Keepalive(uacMsg *UacMsg)error{
	respone := uacMsg
	respone.msg.Status = sip.StatusOK
	respone.msg.Payload = nil

	if err := this.WriteToUac(respone); err != nil{
		return errors.New("Keepalive " + err.Error())
	}
	return nil
}

func (this *UdpServer)Register(uacMsg *UacMsg)error{
	respone := uacMsg
	respone.msg.Status = sip.StatusOK

	//回复
	if err := this.WriteToUac(respone); err != nil{
		return errors.New("Register " + err.Error())
	}
	return nil
}

