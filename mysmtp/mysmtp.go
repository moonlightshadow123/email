package mysmtp
import(
	//"log"
	"strings"
	"net/smtp"
	"test_proj/email/mymsg"
)

var(
	hostaddr 	= "smtp.gmail.com"
	portaddr 	= "smtp.gmail.com:25"
	format		= "2006-01-02T15:04:05-0700"
)

func Send(username string, password string, msg * mymsg.Message)error{
	auth := smtp.PlainAuth("", username, password, hostaddr)
	body := buildBody(msg)
	err := smtp.SendMail(portaddr, auth, msg.From[0].Address, []string{msg.To[0].Address}, []byte(body))
	return err
}

func buildBody(msg *mymsg.Message)(string){
	var newbody string
	newbody += ("To: " + msg.To[0].Address + "\r\n")
	newbody += ("Subject: " + msg.Subject + "\r\n")
	newbody += "\r\n"
	newbody += msg.Body
	return newbody
}

func BuildReBody(msg * mymsg.Message)(string){
	var newbody string = "\r\n"
	newbody += (">Date: " + msg.Date.Format(format) + "\r\n")
	newbody += (">From: " + msg.From[0].Address + "\r\n")
	newbody += (">To: " + msg.To[0].Address + "\r\n")
	newbody += (">Subject: " + msg.Subject + "\r\n")
	newbody += ">\r\n"
	newbody += (">" + strings.ReplaceAll(msg.Body, "\n", "\n>"))
	return newbody
}

func Reply(msg *mymsg.Message)(*mymsg.Message){
	newmsg := &mymsg.Message{}
	newmsg.To = msg.From
	newmsg.From = msg.To
	newmsg.Subject = "Re:"+msg.Subject
	newmsg.Body = BuildReBody(msg)
	return newmsg
}