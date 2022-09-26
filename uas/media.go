package uas

import (
	"gosip/gb"
	"gosip/tools"
	"fmt"
	"github.com/jart/gosip/sdp"
	"github.com/jart/gosip/sip"
	"net"
)


func (this *UdpServer)Play(uacConn *net.UDPAddr, req *gb.PlayReq)(streamId string, err error){
	ssrc := this.GenSSRC(1)

	playSdp := sdp.New(uacConn)
	playSdp.Origin = sdp.Origin{
		User:    this.SysConf.GB28181.SipId,
		Addr: this.SysConf.Server.HttpAddr,
		ID: "0",
		Version: "0",
	}
	playSdp.Addr = this.SysConf.Media.Addr
	playSdp.Audio = nil
	playSdp.Video = &sdp.Media{
		Proto: "TCP/RTP/AVP",
		Port:   this.SysConf.Media.StreamRecvPort,
		Codecs:  []sdp.Codec{
			{PT: 96, Name: "PS", Rate: 90000},
			{PT: 98, Name: "H264", Rate: 90000},
			{PT: 97, Name: "MPEG4", Rate: 90000},
		},
	}
	playSdp.Session = "Play"
	playSdp.Time = "0 0"
	playSdp.SendOnly = false
	playSdp.RecvOnly = true
	playSdp.Attrs = [][2]string{[2]string{"setup","passive"}, [2]string{"connection","new"}}
	playSdp.Other = [][2]string{[2]string{"y", ssrc}}


	sipPlay := new(sip.Msg)
	sipPlay.CallID = tools.Rand(32)
	sipPlay.CSeq = 12
	sipPlay.Request = &sip.URI{
		User:   req.ChannelId,
		Host:   this.SysConf.Server.HttpAddr,
	}

	sipPlay.Subject = fmt.Sprintf("%s:%s,%s:%s", req.ChannelId, ssrc, this.SysConf.GB28181.SipId, ssrc)
	//sipPlay.Via = uacMsg.msg.Via //branch事务ID
	sipPlay.Via = &sip.Via{
		Protocol: "SIP",
		Version: "2.0",
		Transport: "UDP",
		Host:      req.Addr,
		Port:      req.Port,
		Param:     &sip.Param{
			Name:  "branch",
			Value: gb.GenBranch(),
			Next: &sip.Param{Name: "rport"},
		},
	}

	sipPlay.Payload = &sip.MiscPayload{
		T: gb.SDP,
		D: playSdp.Data(),
	}

	sipPlay.Method = "INVITE"
	sipPlay.CSeqMethod = "INVITE"
	sipPlay.To = &sip.Addr{
		Uri:     &sip.URI{
			User:   req.ChannelId,
			Host:   req.Addr,
			Port:   req.Port,
		},
	}
	sipPlay.From = &sip.Addr{
		Uri:     &sip.URI{
			Scheme: "sip",
			User:   this.SysConf.GB28181.SipId,
			Host:	this.SysConf.GB28181.SipDomain,
		},
		Param:  &sip.Param{
			Name:  "tag",
			Value: tools.Rand(32),
		},
	}
	sipPlay.Contact = sipPlay.From

	return gb.SsrcTostreamId(ssrc), this.WriteToUac(&UacMsg{
		uacConn: uacConn,
		msg:     sipPlay,
	})

}

func (this *UdpServer)PlayRespone(uacMsg *UacMsg)(err error){
	m := uacMsg.msg.Copy()
	m.Request = m.From.Uri
	m.Status = 0
	m.Method = "ACK"
	m.CSeqMethod = "ACK"
	m.Payload = nil
	m.From.Uri.User = this.SysConf.GB28181.SipId
	m.Via.Port = this.SysConf.Server.UpdPort

	return this.WriteToUac(&UacMsg{
		uacConn: uacMsg.uacConn,
		msg:     m,
	})
}
