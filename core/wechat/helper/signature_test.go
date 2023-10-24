package helper

import (
	"testing"
)

func TestSignature(t *testing.T) {
	// abc sig
	abc := "f3671fdccd493ec17b3880c9c84a7ccaea097b8a"
	if abc != Signature("jsapi_ticket=kgt8ON7yVITDhtdwci0qebo8WhD9eU7mj-dGbTK6GO5oQPgnJFdkl5Rw0LWpGXOjdp_uZ8BW7JCbPoYU3I5CzQ&", "noncestr=xGMN96jOCBUa070D&", "timestamp=1698145746&", "url=http://dev.scrm2.greenbirds.cn:17803/") {
		t.Error("测试失败")
	}
	t.Log("测试通过")
}

//
