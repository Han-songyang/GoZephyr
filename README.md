# GoZephyr

## 📖 Introduction

GoZephyr is an HTTP framework written in Golang!

Used to help users quickly build their own HTTP service, and provides an introduction to the general api。

## 👋 Getting Started

```bash
go get github.com/Han-songyang/GoZephyr
```

## ✈️ Quick Start

```go
package main

import "github.com/Han-songyang/GoZephyr"

func main() {
	s := GoZephyr.NewCoreServer()
	s.Get("/", func(ctx *GoZephyr.Context) {
		ctx.Resp.Write([]byte("hello, GoZephyr"))
	})
	s.Start(":8081")
}

```

⌛️ Only a small number of features have been implemented so far, and a large number of features are under development.