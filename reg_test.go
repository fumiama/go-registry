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

func TestPush(t *testing.T) {
	r := NewRegedit("reilia.eastasia.azurecontainer.io:32664", "fumiama", "--")
	err := r.Connect()
	if err != nil {
		t.Fatal(err)
	}
	/*
		err = r.Set("ZeroBot-Plugin/kanban", "QQ群:1048452984\n禁止用于商业用途")
		if err != nil {
			t.Fatal(err)
		}
		m, err := hex.DecodeString("48f7021b54af01360a99e0f0e4937bc2")
		if err != nil {
			t.Fatal(err)
		}
		err = r.Set("data/BookReview/bookreview.db", helper.BytesToString(m))
		if err != nil {
			t.Fatal(err)
		}
		m, err = hex.DecodeString("ae4448ea660c10ffa025bab643566c04")
		if err != nil {
			t.Fatal(err)
		}
		err = r.Set("data/Diana/text.pb", helper.BytesToString(m))
		if err != nil {
			t.Fatal(err)
		}
		m, err = hex.DecodeString("f4d6b91f203c89eb42b0a92518a29502")
		if err != nil {
			t.Fatal(err)
		}
		err = r.Set("data/Diana/text.db", helper.BytesToString(m))
		if err != nil {
			t.Fatal(err)
		}
		m, err = hex.DecodeString("5dde28c2d5c55cc2fd2e9f8276786606")
		if err != nil {
			t.Fatal(err)
		}
		err = r.Set("data/Omikuji/kuji.db", helper.BytesToString(m))
		if err != nil {
			t.Fatal(err)
		}
		m, err = hex.DecodeString("86d9be1db96a74094f9660d47d8d3af7")
		if err != nil {
			t.Fatal(err)
		}
		err = r.Set("data/Reborn/rate.json", helper.BytesToString(m))
		if err != nil {
			t.Fatal(err)
		}
		m, err = hex.DecodeString("9ed2e3922ffe50fda24a566fbd152f7b")
		if err != nil {
			t.Fatal(err)
		}
		err = r.Set("data/SetuTime/SetuTime.db", helper.BytesToString(m))
		if err != nil {
			t.Fatal(err)
		}
		m, err = hex.DecodeString("a02630a6811e93c398a7361e04c1c865")
		if err != nil {
			t.Fatal(err)
		}
		err = r.Set("data/VtbQuotation/vtb.db", helper.BytesToString(m))
		if err != nil {
			t.Fatal(err)
		}
	*/
	err = r.Close()
	if err != nil {
		t.Fatal(err)
	}
}
