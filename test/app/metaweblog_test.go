package app

import (
	"os/exec"
	"log"
	"testing"
	"fmt"
	"io/ioutil"
)

func commitInWindows(s string)  {
	cmd := exec.Command("cmd", "/C", s)

	stdout, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(opBytes))
}


func TestCommitInWindows(t *testing.T) {
	s := fmt.Sprintf(`cd /d %s & git pull & git add -A & git commit -m "%s"`, " E:/rox/roxliu.github.io", "publish post")
	log.Println("Run the command in windows: " + s)
	commitInWindows(s)
}

/*
func TestEncodeResponse(t *testing.T) {
	var posts []app.Post
	post := app.Post{DateCreated: time.Now().Local(), Description: "", Title: ""}
	posts = append(posts, post)

	p := xmlrpc.EncodeResponse(posts)
	log.Println("EncodeResponse:\n" + string(p.Bytes()))
}*/