package gb

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

/**
必须换行，Indent设置换行无效
<?xml version="1.0" encoding="GB2312"?>
<Query>
<CmdType>Catalog</CmdType>
<SN>1745</SN>
<DeviceID>34020000001310000001</DeviceID>
</Query>
 */
func Marshal(v interface{}) []byte {
	data := []byte{}
	buf := bytes.NewBuffer(data)
	buf.Write([]byte(fmt.Sprintf(Header, EncodeGB2312)))

	enc := xml.NewEncoder(buf)
	enc.Indent("", "")
	if err := enc.Encode(v); err != nil {
		fmt.Printf("【xml Marshal】error: %v\n", err)
	}
	return bytes.ReplaceAll(buf.Bytes(), []byte("><"), []byte(">\n<"))
}

func Unmarshal(payload []byte, v interface{})error{
	payload = bytes.ReplaceAll(payload, []byte(EncodeGB2312), []byte(""))
	return xml.Unmarshal(payload, v)
}
