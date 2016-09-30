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
	"sort"
	"strings"
	"time"
)

// 排序接口实现
type ByFreq []Slist

func (a ByFreq) Len() int {
	return len(a)
}

func (a ByFreq) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByFreq) Less(i, j int) bool {
	return a[i].Freq > a[j].Freq
}

func workspace() string {
	home, err := user.Current()
	if err != nil {
		fmt.Errorf("%s", err.Error())
		panic(err)
	}
	workspace := fmt.Sprintf("%s/.clist/", home.HomeDir)
	if _, err := os.Stat(workspace); os.IsNotExist(err) {
		os.Mkdir(workspace, 0764)
	}
	return workspace
}

func read_raw_input() []string {
	//在用户的home目录下创建文件夹 .clist
	// 在里面存放文件 clist， 用于给使用者添加命令
	// 此处读取命令，然后交给其他函数整理
	wk := workspace()
	clist := path.Join(wk, "clist")
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
	wk := workspace()
	clist := path.Join(wk, "clist")
	if _, err := os.Stat(clist); os.IsNotExist(err) {
		os.Create(clist)
	}
	slist := path.Join(wk, "slist")
	if _, err := os.Stat(slist); os.IsNotExist(err) {
		os.Create(slist)
		// date, _ := time.Parse("2006-01-02 15:04:05", "1970-01-01 00:01:05")
		initial := fmt.Sprintf(`{"time": "%s", "commands": []}`, "2006-01-02 15:04:05")
		ioutil.WriteFile(slist, []byte(initial), 0764)
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
	restruct(&meta, slist, clist)
	return &meta
}

func restruct(meta *Meta, slist, clist string) *Meta {
	// 根据 clist 修改的时间和 slist保存的时间，判断是否需要重新根据 clist 生成 slist
	info, err := os.Stat(clist)
	if err != nil {
		log.Fatal(err)
	}
	modtime := info.ModTime().UTC()
	metatime, _ := time.Parse("2006-01-02 15:04:05", meta.Time)
	sort.Sort(ByFreq(meta.Commands))
	if metatime.Before(modtime) {
		lines := read_raw_input()
		cmds := make([]Slist, 0)
		for _, line := range lines {
			strs := strings.Split(line, "|")
			flag := true
			for _, cl := range meta.Commands {
				if cl.Cmd == strs[0] {
					flag = false
					if len(strs) > 1 {
						cl.Desc = strs[1]
					}
					cmds = append(cmds, cl)
					break
				}
			}
			if flag {
				var sl Slist
				sl.Cmd = strs[0]
				if len(strs) > 1 {
					sl.Desc = strs[1]
				}
				cmds = append(cmds, sl)
			}
		}
		sort.Sort(ByFreq(cmds))
		meta.Commands = cmds
	}
	meta.Time = time.Now().UTC().Format("2006-01-02 15:04:05")
	data, err := json.Marshal(meta)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile(slist, data, 0764)
	return meta
}

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	slists := read_struct_input()
	strs := make([]string, 0)
	for k, v := range slists.Commands {
		// 封装成 [1] [cmd  desc]
		var line string
		if k == 0 {
			// line = line + "(fg-white,bg-green)"
			line = fmt.Sprintf("[%d] [%s  %s](fg-white,bg-green)", k+1, v.Cmd, v.Desc)
		} else {
			line = fmt.Sprintf("[%d] %s  %s", k+1, v.Cmd, v.Desc)
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
		wk := workspace()
		ioutil.WriteFile(path.Join(wk, "cmd"), []byte("echo Bye bye!"), 0764)
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/<down>", func(termui.Event) {
		if current < len(ls.Items)-1 {
			clist := slists.Commands[current]
			ls.Items[current] = fmt.Sprintf("[%d] %s %s", current+1, clist.Cmd, clist.Desc)
			next := slists.Commands[current+1]
			ls.Items[current+1] = fmt.Sprintf("[%d] [%s %s](fg-white,bg-green)", current+2, next.Cmd, next.Desc)
			termui.Render(ls)
			current = current + 1
		}
	})
	termui.Handle("/sys/kbd/<up>", func(termui.Event) {
		if current > 0 {
			clist := slists.Commands[current]
			ls.Items[current] = fmt.Sprintf("[%d] %s %s", current+1, clist.Cmd, clist.Desc)
			next := slists.Commands[current-1]
			ls.Items[current-1] = fmt.Sprintf("[%d] [%s %s](fg-white,bg-green)", current, next.Cmd, next.Desc)
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
		if len(slists.Commands) > 0 {
			cmd := slists.Commands[current]
			cmd.Freq = cmd.Freq + 1
			slists.Commands[current] = cmd
			wk := workspace()
			restruct(slists, path.Join(wk, "slist"), path.Join(wk, "clist"))
			ioutil.WriteFile(path.Join(wk, "cmd"), []byte(cmd.Cmd), 0764)
		}
		termui.StopLoop()
	})
	termui.Loop()
	termui.Close()
}
