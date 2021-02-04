package iothubclient

import (
	"encoding/json"
	"testing"
)

func TestMakeStaticDeviceInfo(t *testing.T) {

	tt := map[*DeviceInfo]string{
		{
			"firmwareversion": "klgeplgepleglpleplgpelp",
			"hwrev":           "1.0",
		}: `{"verions":{"firmwareversion":"klgeplgepleglpleplgpelp","hwrev":"1.0"}}`,
		{}: `{"verions":{}}`,
	}

	for d, exp := range tt {
		s := makeStaticDeviceInfo(*d)

		b, _ := json.Marshal(s)
		t.Logf("%s", b)
		if string(b) != exp {
			t.Errorf("want %v, got %v", exp, b)
		}
	}
}
