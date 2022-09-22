package uas

import (
	"demo/config"
	"demo/gb"
	"fmt"
	"github.com/jart/gosip/sip"
	"net"
)

const bufferSize uint16 = 65535 - 20 - 8 // IPv4 max size - IPv4 Header size - UDP Header size
const messagePoolSize = 1000

var passwd = "111"

type UdpServer struct {
	sysConf *config.SysConf
	messagePool chan *UacRequest //消息池
	UDPConn *net.UDPConn         //UDP服务器
}

func NewUdpServer(sysConf *config.SysConf) *UdpServer {
	server :=  &UdpServer{sysConf: sysConf,
		messagePool: make(chan *UacRequest, messagePoolSize),
	}

	go server.requestHandler()

	fmt.Println("【UdpServer】Run Success, Listening Port:", sysConf.Server.UpdAddr)
	return server
}

func (this *UdpServer) WriteUac(msg *SipMsg)error{
	_, err := this.UDPConn.WriteToUDP([]byte(msg.msg.String()), msg.uac)
	return err
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
					sipMsg, err := message.ToSipMsg()
					if err != nil{
						fmt.Println("ToSipMsg err:", err)
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


func (this *UdpServer)distribute(sipMsg *SipMsg)(err error){

	//fmt.Println("parseMessage:",  sip.MethodRegister ,*msgReceive)

	if sipMsg.msg.Payload != nil{
		payload := sipMsg.msg.Payload.Data()

		keepalive := &gb.Keepalive{}
		if err = gb.Unmarshal(payload, keepalive); err == nil{
			if keepalive.CmdType == gb.CmdTypeKeepalive{
				err = this.Keepalive(sipMsg)

				this.Catalog(sipMsg, &gb.Query{
					Payload: gb.Payload{
						CmdType:  "Catalog",
						SN:       keepalive.SN,
						DeviceID: sipMsg.msg.From.Uri.User,
					},
				})

				//todo 认证、存储
			}
		}


		catalogRespone := &gb.CatalogResponse{}
		if err = gb.Unmarshal(payload, catalogRespone); err == nil{
			fmt.Println("------catalogRespone:\n", catalogRespone)

			//回复200
		}
	}

	if sipMsg.msg.Method == sip.MethodRegister{
		err = this.Register(sipMsg)


	}

	return
}
