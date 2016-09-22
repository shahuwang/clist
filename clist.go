// Copyright 2016 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gizak/termui"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
)

func read_raw_input() []string {
	//在用户的home目录下创建文件夹 .clist
	// 在里面存放文件 clist， 用于给使用者添加命令
	// 此处读取命令，然后交给其他函数整理
	home, err := user.Current()
	if err != nil {
		fmt.Errorf("%s", err.Error())
		panic(err)
	}
	workspace := fmt.Sprintf("%s/.clist/", home.HomeDir)
	if _, err := os.Stat(workspace); os.IsNotExist(err) {
		os.Mkdir(workspace, 0764)
	}
	clist := path.Join(workspace, "clist")
	if _, err = os.Stat(clist); os.IsNotExist(err) {
		os.Create(clist)
	}
	file, err := os.Open(clist)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	lines := make([]string, 0)
	// 每行文本都以 | 号做数据的分割
	// 目前的需求是：  命令|说明
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

type Slist struct {
	Cmd  string
	Desc string
	Freq int32
}

type Meta struct {
	Time     string
	Commands []Slist
}

func read_struct_input() *Meta {
	home, err := user.Current()
	if err != nil {
		fmt.Errorf("%s", err.Error())
		panic(err)
	}
	workspace := fmt.Sprintf("%s/.clist/", home.HomeDir)
	if _, err := os.Stat(workspace); os.IsNotExist(err) {
		os.Mkdir(workspace, 0764)
	}
	slist := path.Join(workspace, "slist")
	if _, err = os.Stat(slist); os.IsNotExist(err) {
		os.Create(slist)
	}
	file, err := os.Open(slist)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	var meta Meta
	err = json.Unmarshal(data, &meta)
	if err != nil {
		log.Fatal(err)
	}
	return &meta
}

// func encode(lines []string)[]string {
// 	// 将 命令|说明
// }

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	slists := read_struct_input()
	strs := make([]string, 0)
	for k, v := range slists {
		// 封装成 [1] [cmd  desc]
		line := fmt.Sprintf("[%d] [%s  %s]", k+1, v.Cmd, v.Desc)
		if k == 0 {
			line = line + "(fg-white,bg-green)"
		}
		strs = append(strs, line)
	}

	ls := termui.NewList()
	ls.Items = strs
	ls.ItemFgColor = termui.ColorYellow
	ls.BorderLabel = "List"
	ls.Height = termui.TermHeight()
	ls.Width = termui.TermWidth()
	ls.Y = 0
	current := 0
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
	fmt.Println(strs[0])
	fmt.Println(strs[1])
}
