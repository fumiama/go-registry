package registry

import (
	"testing"
)

func TestReg(t *testing.T) {
	r := NewRegedit("127.0.0.1:8888", "storbin.test", "testpwd", "testsps")
	err := r.Connect()
	if err != nil {
		t.Fatal(err)
	}
	v, err := r.Get("test")
	if err != nil && err != ErrNoSuchKey {
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
	t.Log(v)
	if err != ErrNoSuchKey {
		t.Fatal(err)
	}
	err = r.Set("test", "测试")
	if err != nil {
		t.Fatal(err)
	}
	s, err := r.Cat()
	if err != nil {
		t.Fatal(err)
	}
	v, err = s.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	if v != "测试" {
		t.Fatal("invalid test key value in store")
	}
	_, err = r.IsMd5Equal(s.Md5)
	if err != nil {
		t.Fatal(err)
	}
	err = r.Del("test")
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.IsMd5Equal(s.Md5)
	if err == nil {
		t.Fatal("md5 has not been renewed")
	}
	err = r.Close()
	if err != nil {
		t.Fatal(err)
	}
	s, err = r.Load()
	if err != nil {
		t.Fatal(err)
	}
	v, err = s.Get("test")
	if err != nil {
		t.Fatal(err)
	}
	if v != "测试" {
		t.Fatal("invalid test key value in store")
	}
}

func TestPush(t *testing.T) {
	r := NewRegedit("reilia.fumiama.top:32664", "storbin.test", "fumiama", "--")
	err := r.Connect()
	if err != nil {
		t.Fatal(err)
	}
	err = r.Set("ZeroBot-Plugin/kanban", "\t   QQ群:1048452984, TG群:t.me/zerobotplugin\n\t    开发群:752669987, 进阶开发群:705749886\n\n\t\t\t禁止用于违法用途")
	if err != nil {
		t.Fatal(err)
	}
	err = r.Close()
	if err != nil {
		t.Fatal(err)
	}
}
