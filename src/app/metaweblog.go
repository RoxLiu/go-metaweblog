package app

import (
	"net/http"
	"../xmlrpc"
	"log"
	"reflect"
	"io/ioutil"
	"bytes"
	"fmt"
	"time"
	"os"
	"encoding/hex"
	"crypto/md5"
	"strings"
	"runtime"
	"os/exec"
)

type Base64 struct {
}

type BlogInfo struct {
	BlogId 			string 			`xml:"blogid"`
	URL				string			`xml:"url"`
	BlogName		string			`xml:"blogName"`
}

type Enclosure struct {
	Length 			int 			`xml:"length"`
	Type			string			`xml:"type"`
	URL				string			`xml:"url"`
}

type Source struct {
	Name			string			`xml:"name"`
	URL				string			`xml:"url"`
}

type Post struct {
	DateCreated 	time.Time		`xml:"dateCreated"`
	Description		string			`xml:"description"`
	Title			string			`xml:"title"`
	Categories		[]string		`xml:"categories"`
	Enclosure		Enclosure		`xml:"enclosure"`
	Link			string			`xml:"link"`
	Permalink		string			`xml:"permalink"`
	PostId			string			`xml:"postid"`
}

type CategoryInfo struct {
	Description		string			`xml:"description"`
	HtmlURL			string			`xml:"htmlUrl"`
	RssURL			string			`xml:"rssUrl"`
	Title			string			`xml:"title"`
	CategoryId		string			`xml:"categoryid"`
}

type MediaObject struct {
	Name			string			`xml:"name"`
	Type			string			`xml:"type"`
	Bits			Base64			`xml:"bits"`
}

type MediaObjectUrl struct {
	URL				string			`xml:"url"`
}

type MethodResponse struct {
	Content			interface{}		`xml:"params"`
}

type Params struct {
	Value			interface{}		`xml:"param"`
}

type MetaweblogHandler struct {
}

var routine = map[string]interface{} {
	"blogger.getUsersBlogs": getUsersBlogs,
	"metaWeblog.getPost": getPost,
	"metaWeblog.getRecentPosts": getRecentPosts,
	"metaWeblog.newPost": newPost,
	"metaWeblog.editPost": editPost,
	"metaWeblog.deletePost": deletePost,
	"metaWeblog.newMediaObject": newMediaObject,
	"metaWeblog.getCategories": getCategories,
	"metaWeblog.getTemplate": getTemplate,
	"metaWeblog.setTemplate": setTemplate,
}

//应答请求
func (p *MetaweblogHandler)Process(ctx *Context, w http.ResponseWriter, r *http.Request) (error) {
	//打印收到的消息
	s := fmt.Sprintf("Time: %s, Request: %s, %s", time.Now().String(), r.Method, r.URL)
	log.Println(s)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read the request body: " +err.Error())
		return err
	}
	log.Println("Request Body:\n" + string(body))

	reader := bytes.NewReader(body)
	v, err := xmlrpc.DecodeRequest(reader)
	if err != nil {
		log.Println("Failed to parse request: " + err.Error())
		return err
	}

	//find the handler to process.
	f := routine[v.Name]

	if f != nil {
		//Call by reflect
		fn := reflect.ValueOf(f)

		in := make([]reflect.Value, 3)
		in[0] = reflect.ValueOf(ctx)
		in[1] = reflect.ValueOf(w)
		in[2] = reflect.ValueOf(v)

		fn.Call(in)
	} else {
		log.Println("Handler not found for message: " + v.Name)
	}

	return nil
}

func getUsersBlogs(ctx *Context, w http.ResponseWriter, r *xmlrpc.MethodRequest)  {
	var blogs []BlogInfo
	blogs = append(blogs, BlogInfo{BlogId: ctx.Conf.GitPage.Name, URL: ctx.Conf.GitPage.URL, BlogName: ctx.Conf.GitPage.Name})

	p := xmlrpc.EncodeResponse(blogs)
	log.Println("Sent The Response:\n" + string(p.Bytes()))

	w.Write(p.Bytes())
}

func calcMd5(s string) string {
	signByte := []byte(s)
	hash := md5.New()
	hash.Write(signByte)
	return hex.EncodeToString(hash.Sum(nil))
}

func listPostFiles(dir string) ([]os.FileInfo) {
	ls, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("Failed to list the post file, dir=%s, error=%s", dir, err)
		return nil
	}

	return ls
}

func getPost(ctx *Context, w http.ResponseWriter, r *xmlrpc.MethodRequest)  {
	//blogId := r.GetParameter(0).(string)
	//username := r.GetParameter(1)
	//password := r.GetParameter(2)

	obj := Post{DateCreated: time.Now().Local(), Description: "", Title: ""}
	p := xmlrpc.EncodeResponse(obj)
	log.Println("Sent The Response:\n" + string(p.Bytes()))

	w.Write(p.Bytes())
}

func getRecentPosts(ctx *Context, w http.ResponseWriter, r *xmlrpc.MethodRequest)  {
	dir := ctx.Conf.GitPage.Root + "/_posts"
	files := listPostFiles(dir)

	var posts []Post
	for _, f := range files {
		if !f.IsDir() {
			array := strings.Split(f.Name(), "-")

			post := Post{DateCreated: time.Now().Local(), Description: "", Title: strings.Join(array[3:], "-")}
			posts = append(posts, post)
		}
	}

	p := xmlrpc.EncodeResponse(posts)
	log.Println("Sent The Response:\n" + string(p.Bytes()))

	w.Write(p.Bytes())
}

func newPost(ctx *Context, w http.ResponseWriter, r *xmlrpc.MethodRequest)  {
	//blogId := r.GetParameter(0).(string)
	//username := r.GetParameter(1)
	//password := r.GetParameter(2)
	post := r.GetParameter(3).(xmlrpc.Struct)

	title := post["title"].(string)
	content := post["description"].(string)

	publish := r.GetParameter(4).(bool)

	year, month, day := time.Now().Date()

	date := fmt.Sprintf("%d-%02d-%02d", year, month, day)
	path := ctx.Conf.GitPage.Root + "/_posts/" + date + "-" + title + ".html"

	//todo: stash before write and commit

	//write the post to file under the blog folder.
	err := writePost(path, title, content)
	if err != nil {
		log.Panicf("Failed to write file: %s, error=%s", path, err)
		p := xmlrpc.EncodeFault(err.Error())
		w.Write([]byte(p.Bytes()))
		log.Println("Sent The Response:\n" + string(p.Bytes()))

		return
	}

	//invoke git to commit and push
	if publish {
		commit(ctx.Conf.GitPage.Root)
	}

	uri := fmt.Sprintf("%d/%02d/%s", year, month, title + ".html")
	p := xmlrpc.EncodeResponse(uri)
	log.Println("Sent The Response:\n" + string(p.Bytes()))

	w.Write([]byte(p.Bytes()))
}

func writePost(path string, title string, content string) (error) {
	s :=
`---
title: ` + title + "\n"
	s +=
`layout: post
---
`
    s += content

	return ioutil.WriteFile(path, []byte(s), os.ModeAppend)
}

func commit(dir string)  {
	if strings.Contains(runtime.GOOS, "windows") {
		s := fmt.Sprintf(`cd /d %s && git pull && git add -A && git commit -m "%s" && git push`, dir, "publish-post")
		log.Println("Run the command in windows: " + s)
		runCommand("cmd", "/C", s)
	} else {
		s := fmt.Sprintf(`cd %s && git pull && git add -A && git commit -m "%s" && git push`, dir, "publish-post")
		runCommand("/bin/bash", "-c", s)
	}
}

func runCommand(c string, args... string)  {
	cmd := exec.Command(c, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Println(err)
	}
	// 读取输出结果
	bs, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(bs))
}

func editPost(ctx *Context, w http.ResponseWriter, r *xmlrpc.MethodRequest)  {
	newPost(ctx, w, r)
}

func deletePost(ctx *Context, w http.ResponseWriter, r *xmlrpc.MethodRequest)  {
	log.Println("This Methos Is Not Implemented.")
}

func newMediaObject(ctx *Context, w http.ResponseWriter, r *xmlrpc.MethodRequest)  {
	//blogid := r.GetParameter(0).(string)
	//username := r.GetParameter(1)
	//password := r.GetParameter(2)
	media := r.GetParameter(3).(xmlrpc.Struct)

	name := media["name"].(string)
	//t := media["type"].(string)
	bits := media["bits"].([]byte)

	prefix := "/assets/posts/"
	//remove the "Open-Live-Writer/" from the name
	name = strings.TrimLeft(name, "Open-Live-Writer/")

	//mkdir
	rs := []rune(name)
	dir := string(rs[0: strings.LastIndex(name, "/") + 1])
	os.MkdirAll(ctx.Conf.GitPage.Root + prefix + dir, os.ModePerm)

	//write file
	path := ctx.Conf.GitPage.Root + prefix + name
	err := ioutil.WriteFile(path, bits, os.ModeAppend)
	if err != nil {
		p := xmlrpc.EncodeFault(err)
		log.Println("Sent The Response:\n" + string(p.Bytes()))

		w.Write(p.Bytes())
		return
	}

	//write the media data to file under the ${blogId} folder.
	obj := MediaObjectUrl{URL: ctx.Conf.GitPage.URL + prefix + name}
	p := xmlrpc.EncodeResponse(obj)
	log.Println("Sent The Response:\n" + string(p.Bytes()))

	w.Write(p.Bytes())
}

func getCategories(ctx *Context, w http.ResponseWriter, r *xmlrpc.MethodRequest)  {
	//category := CategoryInfo{Description: "1829", HtmlURL: ctx.Conf.GitPage.URL, RssURL: ctx.Conf.GitPage.URL, Title: "GitPage Blog", CategoryId: "category-001"}
	category := CategoryInfo{Description: "1829", HtmlURL: "http://locahost:12345/blog/categories", RssURL: "http://locahost:12345/blog/categories", Title: "GitPage Blog", CategoryId: "category-001"}
	categories := []CategoryInfo {category}

	p := xmlrpc.EncodeResponse(categories)
	log.Println("Sent The Response:\n" + string(p.Bytes()))

	w.Write(p.Bytes())
}

func getTemplate(ctx *Context, w http.ResponseWriter, r *xmlrpc.MethodRequest)  {
	log.Println("This Methos Is Not Implemented.")
}

func setTemplate(ctx *Context, w http.ResponseWriter, r *xmlrpc.MethodRequest)  {
	log.Println("This Methos Is Not Implemented.")
}
