package uas

import (
	"fmt"
	"net"
	"sync"
)

type Uac struct {
	DeviceId string
	ChannelId string
}

func (this *Uac) Key()string {
	return fmt.Sprintf("%s_%s", this.DeviceId, this.ChannelId)
}

//IPC连接集合
type UacConn struct {
	m sync.RWMutex
	Uac map[string]*net.UDPAddr //IPC连接集合
}

func (this *UacConn) Get(uac *Uac) (*net.UDPAddr, bool) {
	this.m.RLock()
	uacConn, ok := this.Uac[uac.Key()]
	this.m.RUnlock()
	return uacConn, ok
}

func (this *UacConn) Set(uac *Uac, uacConn *net.UDPAddr) {
	this.m.Lock()
	this.Uac[uac.Key()] = uacConn
	this.m.Unlock()
}
