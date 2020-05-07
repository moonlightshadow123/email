package controllers

import(
	"fmt"
	//"encoding/json"
	"test_proj/email/utils"
	"test_proj/email/mymsg"
	"test_proj/email/myimap"
	"test_proj/email/mysmtp"
	//"github.com/astaxie/beego"
)

type EmailInstance struct{
	ImapInst *myimap.ImapInstance
	BoxMap map[int]string
	CurPage int
}

func InstLogin(username string, password string)(*EmailInstance, error){
	imapInst, err := myimap.Login(username, password)
	if err != nil{
		return nil, err
	}
	emailinst := &EmailInstance{ImapInst:imapInst, BoxMap:make(map[int]string)}
	return emailinst, nil
}

func (e * EmailInstance)GetBoxs()(*[]byte, error){
	jsonstr,err := e.ImapInst.GetBoxs()
	if err != nil{
		return nil, err
	}
	f, err1 := utils.FromJson(jsonstr)
	if err1 != nil{
		return nil, err
	}
	boxs := f.([]interface{})
	for idx,boxname := range boxs{
		e.BoxMap[idx] = boxname.(string)
	}
	jsonstr1, err2 := utils.ToJson(e.BoxMap)
	if err2 != nil{
		return nil, err2
	}
	return jsonstr1, nil
}

func(e * EmailInstance)Logout(){
	e.ImapInst.Client.Logout()
}

func(e * EmailInstance)Select(idx int)error{
	boxname,has := e.BoxMap[idx]
	if has == false{
		return fmt.Errorf("Box of idx",idx,"not found!")
	}
	err := e.ImapInst.SelectBox(boxname)
	return err
}

func(e * EmailInstance)Page(page int)(*[]byte, error){
	jsonstr, err := e.ImapInst.FetchPageInfo(page)
	if err != nil{
		return nil, err
	}
	e.CurPage = page
	fmt.Println("Page: Curpage = ", e.CurPage)
	return jsonstr, nil
}

func(e * EmailInstance)Item(item int)(*[]byte, error){
	fmt.Println("Item: Curpage = ", e.CurPage)
	if e.CurPage < 1{
		return nil, fmt.Errorf("Select Page first!")
	}
	jsonstr, err := e.ImapInst.FetchItemInfo(e.CurPage, item)
	if err != nil{
		return nil, err
	}
	return jsonstr, nil
}

func(e * EmailInstance)Reply(item int)(*[]byte, error){
	jsonstr, err := e.ImapInst.FetchItemInfo(e.CurPage, item)
	if err != nil{
		return nil, err
	}
	msg,err1 := mymsg.FromJson(jsonstr)
	if err1 != nil{
		return nil, err1
	}
	msg = mysmtp.Reply(msg)
	jsonstr1, err2 := msg.ToJson()
	if err2 != nil{
		return nil, err2
	}
	return jsonstr1, nil
}

func (e * EmailInstance)New()(*[]byte, error){
	msg := mymsg.Message{}
	jsonstr, err := msg.ToJson()
	if err != nil{
		return nil, err
	}
	return jsonstr, nil
}