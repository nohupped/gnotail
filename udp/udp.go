package udp

import (
	"net"
	"fmt"
	log "github.com/nohupped/glog"
)

// NewUDPClient returns a *net.UDPConn struct.
func NewUDPClient(ip string, port int) *net.UDPConn{

	ServerAddr,err := net.ResolveUDPAddr("udp",fmt.Sprintf("%s:%d",ip, port))
	if err != nil {
		panic(err)
	}
	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	ConnUDP, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	if err != nil {
		panic(err)
	}
	log.Infof("UDP Client is sending data on %s:%d\n", ip, port)

	return ConnUDP

}
