package main

import (
	"demo/config"
	"demo/uas"
)


func main(){
	sysConf := config.SysConf{
		Server:  &config.Server{
			UpdAddr: ":5050",
		},
		GB28181: &config.GB28181{
			SipId:     "37070000082008000001",
			SipDomain: "3707000008",
		},
	}

	gbSip := uas.NewUdpServer(&sysConf)
	gbSip.Run()
}
