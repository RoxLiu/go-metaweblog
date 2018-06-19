package app

import (
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v2"
)

type Context struct {
	Conf *Conf
}

type Conf struct {
	Server Server  		`yaml:"server"`
	GitPage GitPage
}

type Server struct {
	//Host string
	Port string
}

type GitPage struct {
	Name string
	URL string
	Root string
}

func initConf() *Conf {
	//dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(dir)

	var c Conf
	f, err := ioutil.ReadFile("./conf.yml")
	if err != nil {
		log.Fatalf("Failed to read conf.file: %v ", err)
		return &c
	}

	err = yaml.Unmarshal(f, &c)
	if err != nil {
		log.Fatalf("Failed to Unmarshal to yml: %v ", err)
	}

	return &c
}

func InitContext() (*Context, error) {
	var ctx  Context
	ctx.Conf = initConf()

	return &ctx, nil
}