package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

type Node struct {
	Id       int
	Adcode   string `xorm:"-"`
	Center   string
	Citycode string
	Level    string
	Name     string
	Children []Node `json:"districts" xorm:"-"`
	Parent   int
}

type Body struct {
	Districts []Node
}

var engine *xorm.Engine

func findChildren(key string, deep int) {
	if deep > 2 {
		return
	}
	resp, err := http.Get("http://restapi.amap.com/v3/config/district?key=619b0f4f0f581f7adb245cfd77de9893&subdistrict=1&keywords=" + key)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	v := Body{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		fmt.Println(err)
	}
	for _, item := range v.Districts {
		item.Id, _ = strconv.Atoi(item.Adcode)
		p := item.Id
		if deep == 0 {
			engine.Insert(item)
		}
		for _, row := range item.Children {
			row.Id, _ = strconv.Atoi(row.Adcode)
			row.Parent = p
			engine.Insert(row)
			findChildren(row.Name, deep+1)
		}
	}
	// body, err = json.MarshalIndent(v.Districts, "", "    ")
	// fmt.Println(string(body))
}

func main() {
	engine, _ = xorm.NewEngine("mysql", "root:199337@/deren?charset=utf8")
	engine.ShowSQL(true)
	// engine.Sync2(new(Node))
	findChildren("中国", 0)
}
