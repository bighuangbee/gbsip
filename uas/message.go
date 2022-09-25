package uas

import (
	"demo/gb"
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
	uac *net.UDPAddr	//IPC连接地址
	msg *sip.Msg
}


func (this *UacRequest) ToUacMsg()(*UacMsg, error){
	sipMsg, err := sip.ParseMsg(this.message)
	if err != nil{
		return nil, err
	}
	return &UacMsg{uac: this.uac, msg: sipMsg}, nil
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
		return errors.New("WriteToUDP " + err.Error())
	}
	return nil
}

//向UAC发送catalog请求
func (this *UdpServer)Catalog(uacMsg *UacMsg, catalog *gb.Query)error{

	queryCatalog := uacMsg.msg.Copy()
	queryCatalog.Method = sip.MethodMessage
	queryCatalog.CSeqMethod = sip.MethodMessage
	queryCatalog.Via.Port = queryCatalog.From.Uri.Port
	queryCatalog.Status = 0
	queryCatalog.From.Uri.User = this.sysConf.GB28181.SipId
	queryCatalog.From.Uri.Host = this.sysConf.GB28181.SipDomain
	queryCatalog.From.Uri.Port = 0
	queryCatalog.To = uacMsg.msg.From
	queryCatalog.To.Param = nil
	queryCatalog.Payload = &sip.MiscPayload{
		T: gb.MANSCDP,
		D: gb.Marshal(catalog),
	}

	if err := this.WriteToUac(&UacMsg{
		uac: uacMsg.uac,
		msg: queryCatalog,
	}); err != nil{
		return errors.New("QueueCatalog " + err.Error())
	}
	return nil
}
