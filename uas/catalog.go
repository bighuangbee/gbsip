package uas

import (
	"context"
	"errors"
	"github.com/jart/gosip/sip"
	"gosip/data/domain"
	"gosip/gb"
)

//向UAC发送catalog请求
func (this *UdpServer)Catalog(uacMsg *UacMsg, catalog *gb.Query)error{

	queryCatalog := uacMsg.msg.Copy()
	queryCatalog.Method = sip.MethodMessage
	queryCatalog.CSeqMethod = sip.MethodMessage
	queryCatalog.Via.Port = queryCatalog.From.Uri.Port
	queryCatalog.Status = 0
	queryCatalog.From.Uri.User = this.SysConf.GB28181.SipId
	queryCatalog.From.Uri.Host = this.SysConf.GB28181.SipDomain
	queryCatalog.From.Uri.Port = 0
	queryCatalog.To = uacMsg.msg.From
	queryCatalog.To.Param = nil
	queryCatalog.Payload = &sip.MiscPayload{
		T: gb.MANSCDP,
		D: gb.Marshal(catalog),
	}

	if err := this.WriteToUac(&UacMsg{
		uacConn: uacMsg.uacConn,
		msg:     queryCatalog,
	}); err != nil{
		return errors.New("Catalog " + err.Error())
	}
	return nil
}


func (this *UdpServer)CatalogRespone(uacMsg *UacMsg)error{
	payload := uacMsg.msg.Payload.Data()
	catalogRespone := &gb.CatalogResponse{}
	gb.Unmarshal(payload, catalogRespone)


	if len(catalogRespone.DeviceList.Channels) > 0{
		c := catalogRespone.DeviceList.Channels[0]
		device := domain.Device{
			DeviceId:      catalogRespone.DeviceID,
			Name:          c.Name,
			Manufacturer:  c.Manufacturer,
			Model:         c.Model,
			//Firmware:      "",
			//Transport:     "",
			Status:        gb.StatusMap(c.Status),
			HostAddress:   c.IPAddress,
			Ip:            c.IPAddress,
			Port:          uint16(uacMsg.uacConn.Port),
			//Expires:       0,
			//Charset:       "",
		}
		this.Repo.Device.Save(context.Background(), &device)

		for _, channle := range catalogRespone.DeviceList.Channels {

			c := domain.Channel{
				DeviceId:      catalogRespone.DeviceID,
				ChannelId:     c.DeviceID,
				Name:          c.Name,
				Manufacturer:  c.Manufacturer,
				Model:         c.Model,
				//Firmware:      "",
				//Transport:     "",
				Status:        gb.StatusMap(c.Status),
				HostAddress:   c.IPAddress,
				Ip:            c.IPAddress,
				Port:          uint16(uacMsg.uacConn.Port),
				//Expires:       0,
				//Charset:       "",
			}
			this.Repo.Channel.Save(context.Background(), &c)

			this.UacManager.Set(&Uac{
				DeviceId:  catalogRespone.DeviceID,
				ChannelId: channle.DeviceID,
			}, uacMsg.uacConn)
		}
	}


	//回复200
	msg := new(sip.Msg)
	msg.Via = uacMsg.msg.Via
	msg.Status = sip.StatusOK
	msg.CSeq = uacMsg.msg.CSeq
	msg.CSeqMethod = sip.MethodMessage
	msg.CallID = uacMsg.msg.CallID
	msg.From = uacMsg.msg.From
	msg.To = uacMsg.msg.To
	return this.WriteToUac(&UacMsg{
		uacConn: uacMsg.uacConn,
		msg:     msg,
	})
}
