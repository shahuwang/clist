// Copyright 2016 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"github.com/gizak/termui"
	"os"
)

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}

	strs := []string{
		"[0] [github.com/gizak/termui]",
		"[1] [你好，世界]",
		"[2] [こんにちは世界]",
		"[3] [color output](fg-white,bg-green)",
		"[4] [output.go]",
		"[5] [random_out.go]",
		"[6] [dashboard.go]",
		"[7] [nsf/termbox-go]"}

	ls := termui.NewList()
	ls.Items = strs
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "List"
	ls.Height = termui.TermHeight()
	ls.Width = termui.TermWidth()
	ls.Y = 0
	current := 3
	termui.Render(ls)
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
		fmt.Print("hello world")
	})
	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		if current < len(ls.Items)-1 {
			curStr := ls.Items[current]
			length := len(curStr)
			ls.Items[current] = curStr[0 : length-19]
			nextStr := ls.Items[current+1]
			ls.Items[current+1] = nextStr + "(fg-white,bg-green)"
			termui.Render(ls)
			current = current + 1
		}
	})
	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		if current > 0 {
			curStr := ls.Items[current]
			length := len(curStr)
			ls.Items[current] = curStr[0 : length-19]
			nextStr := ls.Items[current-1]
			ls.Items[current-1] = nextStr + "(fg-white,bg-green)"
			termui.Render(ls)
			current = current - 1
		}
	})
	termui.Handle("/sys/wnd/resize", func(termui.Event) {
		ls.Height = termui.TermHeight()
		ls.Width = termui.TermWidth()
		termui.Render(ls)
	})
	termui.Handle("/sys/kbd/<enter>", func(termui.Event) {
		fmt.Print("hellodddd")
		termui.StopLoop()
	})
	termui.Loop()
	termui.Close()
	fmt.Println("xxxxxxxxx")
	fmt.Printf("hedddxxx")
	os.Stdout.WriteString("\radsfsadf")
	os.Stdout.Sync()
	err = os.Setenv("MYLIST", "HELLO")
	if err != nil {
		fmt.Println(err.Error())
	}
}
