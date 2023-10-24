package GoZephyr

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	s := NewCoreServer()
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Post("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})
	go s.Start(":8081")
	time.Sleep(time.Second)
	clint := http.Client{}
	resp, err := clint.Post("http://127.0.0.1:8081/user", "text/plain", nil)
	if err != nil {
		t.Fatal(err)
	}
	// 取出body中的内容，并转换成string类型
	body := resp.Body
	defer body.Close()
	buf := make([]byte, 1024)
	n, err := body.Read(buf)
	// 断言
	assert.Equal(t, "hello, user", string(buf[:n]))

	resp, err = clint.Get("http://127.0.0.1:8081/")
	if err != nil {
		t.Fatal(err)
	}
	// 取出body中的内容，并转换成string类型
	body = resp.Body
	defer body.Close()
	buf = make([]byte, 1024)
	n, err = body.Read(buf)
	// 断言
	assert.Equal(t, "hello, world", string(buf[:n]))

	resp, err = clint.Get("http://127.0.0.1:8081/aaa")
	if err != nil {
		t.Fatal(err)
	}
	// 取出body中的内容，并转换成string类型
	body = resp.Body
	defer body.Close()
	buf = make([]byte, 1024)
	n, err = body.Read(buf)
	// 断言
	assert.Equal(t, "Not Found", string(buf[:n]))
}
