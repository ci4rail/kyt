package iothubclient

import (
	"encoding/json"
	"testing"
)

func TestMakeStaticDeviceInfo(t *testing.T) {

	tt := map[*DeviceInfo]string{
		{
			"firmwareVersion": "klgeplgepleglpleplgpelp",
			"hwrev":           "1.0",
		}: `{"versions":{"firmwareVersion":"klgeplgepleglpleplgpelp","hwrev":"1.0"}}`,
		{}: `{"versions":{}}`,
	}

	for d, exp := range tt {
		s := makeStaticDeviceInfo(*d)

		b, _ := json.Marshal(s)
		t.Logf("%s", b)
		if string(b) != exp {
			t.Errorf("want %v, got %v", exp, string(b))
		}
	}
}
