package uas

import (
	"demo/config"
	"demo/gb"
	"demo/tools"
	"errors"
	"fmt"
	"github.com/jart/gosip/sdp"
	"github.com/jart/gosip/sip"
	"net"
	"strconv"
	"sync"
	"time"
)

const bufferSize uint16 = 65535 - 20 - 8 // IPv4 max size - IPv4 Header size - UDP Header size
const messagePoolSize = 1000

type UdpServer struct {
	sysConf *config.SysConf
	messagePool chan *UacRequest //消息池
	UDPConn *net.UDPConn         //UDP服务-连接
	UDPAddr *net.UDPAddr		//UDP服务-地址
	ssrc int

	UacConns *UacConn	//IPC连接集合
}

func NewUdpServer(sysConf *config.SysConf) *UdpServer {
	server :=  &UdpServer{sysConf: sysConf,
		messagePool: make(chan *UacRequest, messagePoolSize),
		UacConns: &UacConn{
			m:   sync.RWMutex{},
			Uac: make(map[string]*net.UDPAddr),
		},
	}
	server.ssrc=10
	go server.requestHandler()

	fmt.Println("【UdpServer】Run Success, Listening Port:", sysConf.Server.UpdPort)
	return server
}

func (this *UdpServer) WriteToUac(msg *UacMsg)error{
	fmt.Println(fmt.Sprintf("-------------WriteToUac:\n%s\n------WriteToUac end\n", msg.msg.String()))
	_, err := this.UDPConn.WriteToUDP([]byte(msg.msg.String()), msg.uacConn)
	return err
}

//isRTP 0=实时流，1=历史流
func (this *UdpServer) GenSSRC(isRTP int)string{
	this.ssrc += 1
	return strconv.Itoa(isRTP)+ tools.RepairSuff(this.ssrc, 9)

}

func (this *UdpServer)Run(){
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", this.sysConf.Server.UpdPort))
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
			fmt.Println("read udp failed, err", err)
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
					case message := <- this.messagePool:{
						var sipMsg *UacMsg
						var err error
						if sipMsg, err = message.ToUacMsg(); err != nil{
							fmt.Println("ToUacMsg err:", err)
							continue
						}
						if err := this.distribute(sipMsg); err != nil{
							fmt.Println(err)
						}
					}
				}
			}
		}()

	}
}


var f = true
func (this *UdpServer)distribute(uacMsg *UacMsg)(err error){

	if uacMsg.msg.Method == sip.MethodRegister{
		err = this.Register(uacMsg)

		time.AfterFunc(time.Second, func() {
			streamId, _ := this.Play(uacMsg, &gb.PlayReq{
				ChannelId:  "34020000001310000001",
				Addr:      "192.168.80.2",
				Port:      5060,
			})
			fmt.Println("------------------streamId:",streamId)
		})
	}else {
		if uacMsg.msg.Payload != nil {
			payload := uacMsg.msg.Payload.Data()

			if uacMsg.msg.Status == sip.StatusOK {
				var sdpMsg *sdp.SDP
				if sdpMsg, err = sdp.Parse(string(payload)); err != nil{
					return err
				}

				if sdpMsg.Session == "Play" {
					err = this.PlayRespone(uacMsg)
				}
			}else{
				base := gb.Payload{}
				if err = gb.Unmarshal(payload, &base); err != nil {
					return errors.New("CmdType parse " + err.Error())
				}

				switch base.CmdType {
				case gb.CmdTypeKeepalive:{
					err = this.Keepalive(uacMsg)
					if err != nil{
						return err
					}


					if !f {
						return
					}
					f = true
					err = this.Catalog(uacMsg, &gb.Query{
						Payload: gb.Payload{
							CmdType:  "Catalog",
							SN:       base.SN,
							DeviceID: uacMsg.msg.From.Uri.User,
						},
					})

					//todo 认证、存储
				}
				case gb.CmdTypeCatalog:
					err = this.CatalogRespone(uacMsg)

				}
			}

		}
	}

	return
}
