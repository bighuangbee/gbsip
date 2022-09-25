package config

type SysConf struct {
	Server *Server
	GB28181 *GB28181
	Media *Media
}

type Server struct {
	UpdAddr string
	HttpAddr string
}



type GB28181 struct {
	SipId string
	SipDomain string
}

type Media struct{
	Addr string
	Port uint16
	StreamRecvPort uint16
}
