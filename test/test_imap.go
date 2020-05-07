package main
import(
	"fmt"
	"encoding/json"
	//"test_proj/email/mymsg"
	"test_proj/email/myimap"
)

func main(){
	imapIns, _ := myimap.Login("zhangzach6666@gmail.com", "42236163035")
	jsonstr, _ :=imapIns.GetBoxs()
	var f interface{}
	json.Unmarshal(*jsonstr, &f)
	fmt.Println(f)
	
	imapIns.SelectBox("INBOX")
	jsonstr, _ = imapIns.FetchPageInfo(1)
	json.Unmarshal(*jsonstr, &f)
	fmt.Println(f)

	jsonstr, _ = imapIns.FetchItemInfo(1, 2)
	json.Unmarshal(*jsonstr, &f)
	fmt.Println(f)
	/*
	fmt.Println(imapIns.MyMsgs[2].Body)
	msg := mysmtp.Reply(imapIns.MyMsgs[2])
	fmt.Println(msg.Body)
	msg.Body = "Ok, got it!\n" + msg.Body
	mysmtp.Send("zhangzach6666@gmail.com", "42236163035", msg)
	*/
	imapIns.Logout()
}