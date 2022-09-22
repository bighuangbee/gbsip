package config

type SysConf struct {
	Server *Server
	GB28181 *GB28181
}

type Server struct {
	UpdAddr string
	HttpAddr string
}



type GB28181 struct {
	SipId string
	SipDomain string
}
