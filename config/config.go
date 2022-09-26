package config

type SysConf struct {
	Server *Server
	GB28181 *GB28181
	Media *Media
	Database *Database
}

type Server struct {
	UpdPort  uint16
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

type Database struct {
	Address, UserName, Password, DBName, Driver string
	Timeout int
}
