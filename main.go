package main

import (
	"gosip/config"
	"gosip/uas"
)


func main(){
	sysConf := config.SysConf{
		Server:  &config.Server{
			UpdPort:  5050,
			HttpAddr: "192.168.80.107",	//
		},
		GB28181: &config.GB28181{
			SipId:     "37070000082008000001",
			SipDomain: "3707000008",
		},
		Media: &config.Media{
			Addr: "192.168.80.107",
			Port: 8090,
			StreamRecvPort: 10000,
		},
		Database: &config.Database{
			Address:  "localhost",
			UserName: "root",
			Password: "Hiscene2022",
			DBName:   "gbsip",
			Driver:   "mysql",
			Timeout:  10,
		},
	}

	gbSip := uas.NewUdpServer(&sysConf)
	gbSip.Run()
}
