package uas

import (
	"demo/config"
	"demo/gb"
	"demo/tools"
	"fmt"
	"github.com/jart/gosip/sdp"
	"github.com/jart/gosip/sip"
	"net"
	"strconv"
	"time"
)

const bufferSize uint16 = 65535 - 20 - 8 // IPv4 max size - IPv4 Header size - UDP Header size
const messagePoolSize = 1000

type UdpServer struct {
	sysConf *config.SysConf
	messagePool chan *UacRequest //消息池
	UDPConn *net.UDPConn         //UDP服务器
	UDPAddr *net.UDPAddr
	ssrc int
}

func NewUdpServer(sysConf *config.SysConf) *UdpServer {
	server :=  &UdpServer{sysConf: sysConf,
		messagePool: make(chan *UacRequest, messagePoolSize),
	}
	server.ssrc=10
	go server.requestHandler()

	fmt.Println("【UdpServer】Run Success, Listening Port:", sysConf.Server.UpdAddr)
	return server
}

func (this *UdpServer) WriteToUac(msg *UacMsg)error{
	fmt.Println(fmt.Sprintf("-------------WriteToUac:\n%s\n------WriteToUac end\n", msg.msg.String()))
	_, err := this.UDPConn.WriteToUDP([]byte(msg.msg.String()), msg.uac)
	return err
}

//isRTP 0=实时流，1=历史流
func (this *UdpServer) GenSSRC(isRTP int)string{
	this.ssrc += 1
	return strconv.Itoa(isRTP)+ tools.RepairSuff(this.ssrc, 9)

}

func (this *UdpServer)Run(){
	udpAddr, err := net.ResolveUDPAddr("udp", this.sysConf.Server.UpdAddr)
	if err != nil {
		panic(err)
	}

	updServerConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		panic(err)
	}

	this.UDPConn = updServerConn
	this.UDPAddr = udpAddr

	for {
		var buf = make([]byte, bufferSize)
		n, updClient, err := updServerConn.ReadFromUDP(buf[0:])
		if err != nil {
			fmt.Println("read udp failed,err", err)
			continue
		}

		fmt.Println(fmt.Sprintf("【Reqeust】form: %s, data:\n%s\n", updClient.String(), string(buf[:n])))

		this.messagePool <- &UacRequest{
			uac:     updClient,
			message: buf[:n],
		}
	}
}



func (this *UdpServer) requestHandler(){
	for i := 0; i < 1; i++ {
		go func() {
			for {
				select {
				case message := <- this.messagePool:
					sipMsg, err := message.ToUacMsg()
					if err != nil{
						fmt.Println("ToUacMsg err:", err)
						continue
					}
					if err := this.distribute(sipMsg); err != nil{
						fmt.Println(err)
					}
				}
			}
		}()

	}
}



var f = true
func (this *UdpServer)distribute(sipMsg *UacMsg)(err error){

	//fmt.Println("parseMessage:",  sip.MethodRegister ,*msgReceive)

	if sipMsg.msg.Payload != nil{
		payload := sipMsg.msg.Payload.Data()

		keepalive := &gb.Keepalive{}
		if err = gb.Unmarshal(payload, keepalive); err == nil{
			if keepalive.CmdType == gb.CmdTypeKeepalive{
				err = this.Keepalive(sipMsg)

				if !f{
					return
				}
				f=true
				return this.Catalog(sipMsg, &gb.Query{
					Payload: gb.Payload{
						CmdType:  "Catalog",
						SN:       keepalive.SN,
						DeviceID: sipMsg.msg.From.Uri.User,
					},
				})

				//todo 认证、存储
			}
		}


		if sipMsg.msg.CSeqMethod == "INVITE"{

			sdpMsg, err := sdp.Parse(string(payload))
			if err != nil{
				return err
			}
			if sdpMsg.Session == "Play"{
				m := sipMsg.msg.Copy()
				m.Request = m.From.Uri
				m.Status = 0
				m.Method = "ACK"
				m.CSeqMethod = "ACK"
				m.Payload = nil
				m.From.Uri.User = this.sysConf.GB28181.SipId
				m.Via.Port=5050

				//ddd := strings.ReplaceAll(msg.msg.String(), ("3707000008 SIP"), ("192.168.80.2:5060 SIP"))
				fmt.Println("------------===========Play\n", m.String())
				return this.WriteToUac(&UacMsg{
					uac: sipMsg.uac,
					msg: m,
				})

			}
		}

		catalogRespone := &gb.CatalogResponse{}
		if err = gb.Unmarshal(payload, catalogRespone); err == nil{
			//fmt.Println("------catalogRespone:\n", catalogRespone)

			//回复200

			msg := new(sip.Msg)
			msg.Via = sipMsg.msg.Via
			msg.Status = sip.StatusOK
			msg.CSeq = sipMsg.msg.CSeq
			msg.CSeqMethod = sip.MethodMessage
			msg.CallID = sipMsg.msg.CallID
			msg.From = sipMsg.msg.From
			msg.To = sipMsg.msg.To
			return this.WriteToUac(&UacMsg{
				uac: sipMsg.uac,
				msg: msg,
			})

		}
	}

	if sipMsg.msg.Method == sip.MethodRegister{
		err = this.Register(sipMsg)

		time.AfterFunc(time.Second, func() {
			streamId, _ := this.Play(sipMsg, &gb.PlayReq{
				DeviceId:  "34020000001320000001",
				Addr:      "192.168.80.2",
				Port:      5060,
			})
			fmt.Println("------------------streamId:",streamId)
		})

	}

	return
}
