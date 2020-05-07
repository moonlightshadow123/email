package main

import (
	"github.com/astaxie/beego"
	_ "test_proj/email/routers"
	//"test_proj/email/controllers"
)

func main() {
	beego.Run()
}

