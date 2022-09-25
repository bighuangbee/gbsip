package uas

import (
	"demo/gb"
	"demo/tools"
	"fmt"
	"github.com/jart/gosip/sdp"
	"github.com/jart/gosip/sip"
)


func (this *UdpServer)Play(uacMsg *UacMsg, req *gb.PlayReq)(streamId string, err error){
	ssrc := this.GenSSRC(1)

	playSdp := sdp.New(uacMsg.uac)
	playSdp.Origin = sdp.Origin{
		User:    this.sysConf.GB28181.SipId,
		Addr: this.sysConf.Server.HttpAddr,
		ID: "0",
		Version: "0",
	}
	playSdp.Addr = this.sysConf.Media.Addr
	playSdp.Audio = nil
	playSdp.Video = &sdp.Media{
		Proto: "TCP/RTP/AVP",
		Port:   this.sysConf.Media.StreamRecvPort,
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
		User:   req.DeviceId,
		Host:   this.sysConf.Server.HttpAddr,
	}

	sipPlay.Subject = fmt.Sprintf("%s:%s,%s:%s", req.DeviceId, ssrc, this.sysConf.GB28181.SipId, ssrc)
	//sipPlay.Via = uacMsg.msg.Via //branch事务ID
	sipPlay.Via = &sip.Via{
		Protocol: "SIP",
		Version: "2.0",
		Transport: "UDP",
		Host:      req.Addr,
		Port:      req.Port,
		Param:     &sip.Param{
			Name:  "branch",
			Value: gb.GenBranch(),// IETF RFC3261 这个branch参数的值必须用”z9hG4bK”打头
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
			User:   req.DeviceId,
			Host:   req.Addr,
			Port:   req.Port,
		},
	}
	sipPlay.From = &sip.Addr{
		Uri:     &sip.URI{
			Scheme: "sip",
			User:   this.sysConf.GB28181.SipId,
			Host:	this.sysConf.GB28181.SipDomain,
		},
		Param:  &sip.Param{
			Name:  "tag",
			Value: tools.Rand(32),
		},
	}
	sipPlay.Contact = sipPlay.From

	return gb.SsrcTostreamId(ssrc), this.WriteToUac(&UacMsg{
		uac: uacMsg.uac,
		msg: sipPlay,
	})

}
