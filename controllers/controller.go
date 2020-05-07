package controllers

import(
	"fmt"
	"strconv"
	//"net/url"
	//"encoding/json"
	"test_proj/email/utils"
	//"test_proj/email/mymsg"
	"test_proj/email/myimap"
	//"test_proj/email/mysmtp"
	"github.com/astaxie/beego"
)

var(
	sesskey = "emailinst"
)

type EmailController struct{
	beego.Controller
	Imap *myimap.ImapInstance
}

func (ec * EmailController)Error(msg string,err error){
	s := fmt.Sprintf("Msg: %s, Error: %s", msg, err)
	beego.Error(s)
	ec.Data["error"] = s
	ec.TplName = "error.html"
}

func (ec *EmailController)Login(){
	username := ec.Ctx.Request.Form.Get("username")
	password := ec.Ctx.Request.Form.Get("password")
	emailInst,err := InstLogin(username, password)
	if err != nil{
		beego.Error("Email Instance Login Error!", err)
		ec.Redirect("/", 302)
	}
	ec.ClearCookie()
	if err1 := ec.SetSessValue(sesskey, emailInst); err1 != nil{
		beego.Error("Set Sess Value Error!", err1)
		ec.Redirect("/", 302)
	}
	ec.TplName = "index.html"
	boxjson, err2 := emailInst.GetBoxs()
	if err2 != nil{
		ec.Data["boxs"] = make([]string,0, 0)
	}
	f, err3 :=utils.FromJson(boxjson)
	if err3 != nil{
		ec.Data["boxs"] = make([]string,0, 0)
	}
	ec.Data["boxs"] = f.(map[string]interface{})
}

func (ec * EmailController)Select(){
	intf, err := ec.GetSessValue(sesskey)
	if err != nil{
		beego.Error("Set Sess Value Error!", err)
		return
	}
	emailinst := intf.(*EmailInstance)
	idx, err0 := strconv.Atoi(ec.Ctx.Input.Param(":idx"))
	if err0 != nil{
		beego.Error("Index Parse Error!", err0)
		return 
	}
	if err1 := emailinst.Select(idx); err1!= nil{
		beego.Error("Select box error! idx =", idx, "Error:", err1)
		return 
	}
	jsonstr, err2 := emailinst.Page(1)
	if err2 != nil{
		beego.Error("Page error! idx =", idx,"Page = 1", "Error:", err2)
		return 
	}
	jsonobj, err3 := utils.FromJson(jsonstr)
	if err3 != nil{
		beego.Error("Json String Parse Error! Page = 1", "Error:", err3)
	}
	ec.Data["json"] = jsonobj
	ec.ServeJSON()
}

func (ec *EmailController)Page(){
	intf, err := ec.GetSessValue(sesskey)
	if err != nil{
		beego.Error("Set Sess Value Error!", err)
		return
	}
	emailinst := intf.(*EmailInstance)
	page, err1 := strconv.Atoi(ec.Ctx.Input.Param(":page"))
	if err1 != nil{
		beego.Error("Page Para Parse Error!", err1)
		return 
	} 
	jsonstr, err2 := emailinst.Page(page)
	if err2 != nil{
		beego.Error("Page error! Page =", page, "Error:", err2)
		return 
	}
	jsonobj, err3 := utils.FromJson(jsonstr)
	if err3 != nil{
		beego.Error("Json String Parse Error! Page =", page, "Error:", err3)
	}
	ec.Data["json"] = jsonobj
	ec.ServeJSON()
}

func (ec * EmailController)Item(){
	intf, err := ec.GetSessValue(sesskey)
	if err != nil{
		ec.Error("Set Sess Value Error!", err)
		return
	}
	emailinst := intf.(*EmailInstance)
	item, err1 := strconv.Atoi(ec.Ctx.Input.Param(":item"))
	if err1 != nil{
		ec.Error("Item Para Parse Error!", err1)
		return 
	} 
	jsonstr, err2 := emailinst.Item(item)
	if err2 != nil{
		beego.Error("Item error! Item =", item, "Error:", err2)
		return 
	}
	jsonobj, err3 := utils.FromJson(jsonstr)
	if err3 != nil{
		beego.Error("Json String Parse Error! item =", item, "Error:", err3)
	}
	ec.Data["json"] = jsonobj
	ec.ServeJSON()
}

func (ec * EmailController)Reply(){
	intf, err := ec.GetSessValue(sesskey)
	if err != nil{
		ec.Error("Set Sess Value Error!", err)
		return
	}
	emailinst := intf.(*EmailInstance)
	item, err1 := strconv.Atoi(ec.Ctx.Input.Param(":item"))
	if err1 != nil{
		ec.Error("Item Para Parse Error in Reply!", err1)
		return 
	} 
	emailinst.Reply(item)
	ec.Data["name"] = ""
	ec.TplName = "itemedit.html"
}

func (ec * EmailController)New(){
	intf, err := ec.GetSessValue(sesskey)
	if err != nil{
		ec.Error("Set Sess Value Error!", err)
		return
	}
	emailinst := intf.(*EmailInstance)
	emailinst.New()
	ec.Data["name"] = ""
	ec.TplName = "itemedit.html"
}

func(ec * EmailController)Logout(){
	intf, err := ec.GetSessValue(sesskey)
	if err != nil{
		ec.Error("Set Sess Value Error!", err)
		ec.DestroySess()
		ec.Redirect("/", 302)
	}else{
		emailinst := intf.(*EmailInstance)
		emailinst.Logout()
		ec.DestroySess()
		ec.Redirect("/", 302)
	}
}

func (ec * EmailController)Send(){
	intf, err := ec.GetSessValue(sesskey)
	if err != nil{
		ec.Error("Set Sess Value Error!", err)
		return
	}
	emailinst := intf.(*EmailInstance)
	beego.Info(emailinst)

}

