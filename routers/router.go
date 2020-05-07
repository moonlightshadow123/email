package routers

import (
	"github.com/astaxie/beego"
	"test_proj/email/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/index", &controllers.EmailController{}, "post:Login")
	beego.Router("/logout", &controllers.EmailController{}, "get:Logout")
	beego.Router("/select/:idx:int", &controllers.EmailController{}, "get:Select")
	beego.Router("/page/:page:int", &controllers.EmailController{}, "get:Page")
	beego.Router("/item/:item:int", &controllers.EmailController{}, "get:Item")
	beego.Router("/reply/:itme:int", &controllers.EmailController{}, "get:Reply")
	beego.Router("/new", &controllers.EmailController{}, "get:New")
}
