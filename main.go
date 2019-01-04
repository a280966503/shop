package main

import (
	_ "mybeego/routers"
	"github.com/astaxie/beego"
	_ "mybeego/models"
)

func main() {
	beego.Run()
}

