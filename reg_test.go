package registry

import (
	"testing"
)

func TestReg(t *testing.T) {
	r := NewRegedit("127.0.0.1:8888", "testpwd", "testsps")
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
	err = r.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestPush(t *testing.T) {
	r := NewRegedit("reilia.fumiama.top:32664", "fumiama", "--")
	err := r.Connect()
	if err != nil {
		t.Fatal(err)
	}
	err = r.Set("ZeroBot-Plugin/kanban", "由于tx更改登录鉴权, 近期gocqzbp闪退增多,\n请及时处理以免对md5服务器产生过量请求, 导致IP被封.")
	if err != nil {
		t.Fatal(err)
	}
	err = r.Close()
	if err != nil {
		t.Fatal(err)
	}
}
