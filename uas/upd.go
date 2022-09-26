package uas

import (
	"context"
	"errors"
	"fmt"
	"github.com/jart/gosip/sdp"
	"github.com/jart/gosip/sip"
	"gosip/config"
	"gosip/data"
	"gosip/data/model"
	"gosip/gb"
	"gosip/tools"
	"gosip/tools/log"
	"net"
	"strconv"
	"sync"
)

const bufferSize uint16 = 65535 - 20 - 8 // IPv4 max size - IPv4 Header size - UDP Header size
const messagePoolSize = 1000

type UdpServer struct {
	SysConf     *config.SysConf
	messagePool chan *UacRequest //消息池
	UDPConn     *net.UDPConn         //UDP服务-连接
	UDPAddr     *net.UDPAddr		//UDP服务-地址
	ssrc        int

	UacManager UacManager //IPC连接集合

	data *data.Data

	Repo *model.Repo
}

func NewUdpServer(sysConf *config.SysConf) *UdpServer {
	server :=  &UdpServer{SysConf: sysConf,
		messagePool: make(chan *UacRequest, messagePoolSize),
		UacManager: &UacConn{
			m:   sync.RWMutex{},
			Uac: make(map[string]*net.UDPAddr),
		},
	}
	server.ssrc=10
	go server.requestHandler()

	logger := log.New("./logs")
	db, err := data.New(&data.Options{
		Address:  sysConf.Database.Address,
		UserName: sysConf.Database.UserName,
		Password: sysConf.Database.Password,
		DBName:   sysConf.Database.DBName,
		Logger:   logger,
	})
	if err != nil{
		fmt.Println(err)
		return nil
	}
	server.data = &data.Data{
		Db: db,
	}
	server.Repo = &model.Repo{
		Device:  &model.DeviceRepo{Data: server.data},
		Channel: &model.ChannelRepo{Data: server.data},
		Stream: &model.StreamRepo{Data: server.data},
	}


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
	udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", this.SysConf.Server.UpdPort))
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


func (this *UdpServer)distribute(uacMsg *UacMsg)(err error){

	if uacMsg.msg.Method == sip.MethodRegister{
		err = this.Register(uacMsg)

		//time.AfterFunc(time.Second, func() {
		//	streamId, _ := this.Play(uacMsg.uacConn, &gb.PlayReq{
		//		ChannelId:  uacMsg.msg.From.Uri.User,
		//		Addr:      uacMsg.msg.Request.Host,
		//		Port:      uacMsg.msg.Request.Port,
		//	})
		//	fmt.Println("------------------streamId:",streamId)
		//})
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


					deviceId := uacMsg.msg.From.Uri.User
					_, total, err := this.Repo.Channel.List(context.Background(), deviceId)
					fmt.Println("-------------- ", err, total, deviceId)
					if err != nil{
						return err
					}
					if total == 0{
						err = this.Catalog(uacMsg, &gb.Query{
							Payload: gb.Payload{
								CmdType:  gb.CmdTypeCatalog,
								SN:       base.SN,
								DeviceID: deviceId,
							},
						})
					}


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
