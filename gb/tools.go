package gb

import (
	"demo/tools"
	"fmt"
	"strconv"
)

// IETF RFC3261 这个branch参数的值必须用”z9hG4bK”打头
func GenBranch()string{
	return "z9hG4bK"+tools.Rand(32)
}



// zlm接收到的ssrc为16进制。发起请求的ssrc为10进制
func SsrcTostreamId(ssrc string) string {
	if ssrc[0:1] == "0" {
		ssrc = ssrc[1:]
	}
	num, _ := strconv.Atoi(ssrc)
	return fmt.Sprintf("%08X", num)
}
