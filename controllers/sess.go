package controllers

import (
	"fmt"
	"github.com/astaxie/beego/session"
)

var(
	globalSessions *session.Manager
	cookieName = "myemail"
	secret = "emailsecret"
	lifetime = 3600
	secure = true
)

func (ec * EmailController)ClearCookie(){
	ec.SetSecureCookie(secret, "sessid", "")
}

func Init(){
	sessionConfig := &session.ManagerConfig{
		CookieName: cookieName,
		Gclifetime: int64(lifetime),
		Secure: secure,
	}
	globalSessions, _ = session.NewManager("memory", sessionConfig)
    go globalSessions.GC()
}

// Start a session
func (ec * EmailController)startSess()(string,error){
	if globalSessions == nil{
		Init()
		//return "", fmt.Errorf("GlobalSessions is nil!")
	}
	sess, err1 := globalSessions.SessionStart(ec.Ctx.ResponseWriter, ec.Ctx.Request)
	if err1 != nil{
		return "", err1
	}
	sessid := sess.SessionID()
	ec.SetSecureCookie(secret, "sessid", sessid)
	return sessid, nil
}

func (ec * EmailController)DestroySess(){
	if globalSessions == nil{
		return 
	}
	globalSessions.SessionDestroy(ec.Ctx.ResponseWriter, ec.Ctx.Request)
	ec.SetSecureCookie(secret, "sessid", "")
}

func (ec * EmailController)SetSessValue(key string, val interface{})error{
	sessid, isset := ec.GetSecureCookie(secret, "sessid")
	fmt.Println("Sessid = ", sessid, " Sessid == ''", sessid=="")
	if sessid == "" || isset == false{
		if id, err1 := ec.startSess(); err1 != nil{
			return err1
		}else{
			sessid = id
		}
	}
	sess, err2 := globalSessions.GetSessionStore(sessid)
	if err2 != nil{
		return fmt.Errorf("Get Sess Error! sid = %s error = %s",sessid,err2)
	}
	err3 := sess.Set(key, val)
	if err3 != nil{
		return fmt.Errorf("Get Sess Key Error! sid = %s, key = %s, val = %s, error= %s",sessid, key, val.(string), err3)
	}
	return  nil
}

func (ec * EmailController)GetSessValue(key string)(interface{}, error){
	sessid, isset := ec.GetSecureCookie(secret, "sessid")
	if sessid == "" || isset == false{
		return "", fmt.Errorf("Sess Id is empty!")
	}
	sess, err1 := globalSessions.GetSessionStore(sessid)
	if err1 != nil{
		return nil, fmt.Errorf("Get Sess Error! sid = %s error = %s",sessid,err1)
	}
	val := sess.Get(key)
	if val == nil || val == "" {
		return nil, fmt.Errorf("Val is empty! key = %s, sessid = %s", key ,sessid)
	}
	return val, nil
}