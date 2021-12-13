package registry

import "testing"

func TestReg(t *testing.T) {
	r := NewRegedit("127.0.0.1:8888", "testpwd", "testsps")
	err := r.Connect()
	if err != nil {
		t.Fatal(err)
	}
	v, err := r.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	err = r.Set("test", "测试")
	if err != nil {
		t.Fatal(err)
	}
	v, err = r.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	err = r.Del("test")
	if err != nil {
		t.Fatal(err)
	}
	v, err = r.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(v)
	err = r.Close()
	if err != nil {
		t.Fatal(err)
	}
}
