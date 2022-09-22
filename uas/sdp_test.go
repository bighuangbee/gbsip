package uas

import (
	"fmt"
	"github.com/jart/gosip/sdp"
	"testing"
)

var playSipStr = `
INVITE sip:34020000001320000001@192.168.80.2:5060 SIP/2.0
Supported: 
Allow: INVITE, ACK, CANCEL, MESSAGE, REGISTER
Via: SIP/2.0/UDP 192.168.80.107:5061;branch=z9hG4bK4bJeLeSZ7fOxbSHmYeFuSY6akuec3gjk;rport
CSeq: 1 INVITE
From: "sipserver" <sip:37070000082008000001@3707000008>;tag=mP3TE7WexIsdW1Ov9ye2
To: <sip:34020000001320000001@192.168.80.2:5060>
Call-ID: 6SU0zGt38EtZ5NVJyuJtly4kG3Zz8oAl
Contact: "sipserver" <sip:37070000082008000001@3707000008>;tag=mP3TE7WexIsdW1Ov9ye2
Max-Forwards: 70
User-Agent: GoSIP
Content-Type: application/sdp
Content-Length: 263
Subject: 34020000001320000001:0700000008,37070000082008000001:0700000008
`

var playSdpStr= "v=0\r\n" +
"o=37070000082008000001 0 0 IN IP4 192.168.80.107\r\n" +
"s=Play\r\n" +
"c=IN IP4 192.168.80.107\r\n" +
"t=0 0\r\n" +
"m=video 10000 TCP/RTP/AVP 96 98 97\r\n" +
"a=recvonly\r\n" +
"a=setup:passive\r\n" +
"a=connection:new\r\n" +
"a=rtpmap:96 PS/90000\r\n" +
"a=rtpmap:98 H264/90000\r\n" +
"a=rtpmap:97 MPEG4/90000\r\n" +
" =0700000008"

func Test_Play(t *testing.T) {
	playSdp, err := sdp.Parse(playSdpStr)
	if err != nil{
		fmt.Println(err)
		return
	}



	fmt.Println("playSdp:", playSdp)

	fmt.Println("--playSdp.Video.Proto,", playSdp.Origin.Addr, playSdp.Addr, playSdp.Video.Port )
	fmt.Println("--playSdp.Video.Port,", playSdp.Origin.User, playSdp.Origin.ID, playSdp.Origin.Version)
}

