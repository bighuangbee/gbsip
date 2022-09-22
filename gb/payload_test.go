package gb

import (
	"fmt"
	"testing"
)

var notifyPayload =
`<?xml version="1.0" encoding="GB2312" standalone="yes" ?>
<Keepalive>
<CmdType>Keepalive</CmdType>
<SN>1247</SN>
<DeviceID>34020000001320000001</DeviceID>
<Status>OK</Status>
</Keepalive>`



var queueCatalog =
`<?xml version="1.0" encoding="UTF-8"?>
<Query>
<CmdType>Catalog</CmdType>
<SN>49013560</SN>
<DeviceID>34020000001110000001</DeviceID>
</Query>`

func Test_xml(t *testing.T) {
	s := Keepalive{}
	if err := Unmarshal([]byte(notifyPayload), &s); err != nil{
		t.Error(err)
		return
	}
	fmt.Println("notifyPayload struct:", s)

	qc := Query{}
	if err := Unmarshal([]byte(queueCatalog), &qc); err != nil{
		t.Error(err)
		return
	}
	fmt.Println("NewQueryCatalog struct:", qc)

	queryCatalog := Query{
		Payload:   Payload{
			CmdType:  "Catalog",
			SN:       "112233",
			DeviceID: "5544332211",
		},
	}
	fmt.Println("queryCatalog string:", string(Marshal(&queryCatalog)))
}


var deviceList = `
<?xml version="1.0" encoding="GB2312" standalone="yes" ?>
<Response>
<CmdType>Catalog</CmdType>
<SN>1803</SN>
<DeviceID>34020000001320000001</DeviceID>
<SumNum>1</SumNum>
<DeviceList Num="1">
	<Item>
		<DeviceID>34020000001310000001</DeviceID>
		<Name>IPC80.172</Name>
		<Manufacturer>Dahua</Manufacturer>
		<Model>DH-P20A2-4G-PV</Model>
		<Owner>0</Owner>
		<CivilCode>340200</CivilCode>
		<Address>axy</Address>
		<Parental>0</Parental>
		<ParentID>34020000001320000001</ParentID>
		<RegisterWay>1</RegisterWay>
		<Secrecy>0</Secrecy>
		<StreamNum>2</StreamNum>
		<IPAddress>192.168.80.2</IPAddress>
		<Status>ON</Status>
		<Info>
			<PTZType>3</PTZType>
			<DownloadSpeed>1/2/4/8</DownloadSpeed>
		</Info>
	</Item>
</DeviceList>
</Response>`

func Test_CatalogResponse(t *testing.T) {
	s := CatalogResponse{}
	if err := Unmarshal([]byte(deviceList), &s); err != nil {
		t.Error(err)
		return
	}
	fmt.Println("deviceList struct:", s)
}
