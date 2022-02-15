package registry

import (
	"errors"
	"testing"
)

func TestRegedit(t *testing.T) {
	reg := NewRegedit("127.0.0.1:8888", "testpwd", "testsps")
	err := reg.Connect()
	if err != nil {
		t.Fatal(err)
	}
	ret, err := reg.Get("test")
	if err != nil && !errors.Is(err, ErrNoSuchKey) {
		t.Fatal(err)
	}
	t.Log(ret)
	if ret != "" {
		err = reg.Del("test")
		if err != nil {
			t.Fatal(err)
		}
	}
	err = reg.Set("test", "测试")
	if err != nil {
		t.Fatal(err)
	}
	ret, err = reg.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	if ret != "测试" {
		t.Fail()
	}
	err = reg.Close()
	if err != nil {
		t.Fatal(err)
	}
}
