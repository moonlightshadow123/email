package myimap

import (
	"io"
	"fmt"
	"log"
	"io/ioutil"
	"encoding/json"
	"test_proj/email/utils"
	"test_proj/email/mymsg"
	netmail "net/mail"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

var(
	serveraddr = "imap.gmail.com:993"
	Maxboxes = 20
	Msgbatch = 10
)

type BoxStatus = imap.MailboxStatus
type ConvertFunc func(*imap.Message)(*mymsg.Message,error)

type ImapInstance struct{
	Client *client.Client
}

// Login returns a Imap instance
func Login(username string, password string)(*ImapInstance,error){
	c, err := client.DialTLS(serveraddr, nil)
	if err != nil{
		log.Println("Error Dial to Server!", err)
		return nil,err
	}
	if err1 := c.Login(username, password); err1 != nil {
		log.Println(err1)
		return nil,err1
	}
	imapIns := &ImapInstance{Client:c}
	return imapIns, nil
}

func (imapIns *ImapInstance)Logout()error{
	err := imapIns.Client.Logout()
	return err
}

// Fill the Imap Instance's MyBoxs
func (imapIns *ImapInstance)GetBoxs()(*[]byte, error){
	if imapIns.Client == nil{
		return nil, fmt.Errorf("Client nil!")
	}
	// Make a chan
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	// Fill the chan
	go func () {
		done <- imapIns.Client.List("", "*", mailboxes)
	}()
	// Get item from the chan
	boxs := make([]string, 0, Maxboxes)
	log.Println("Mailboxes:")
	i:=0
	for  m := range mailboxes {
		_, found := utils.StrInSlice(m.Attributes, "\\Noselect")
		if found{ continue }
		log.Println("* ", m.Name, ", Attr:", m.Attributes)
		boxs = append(boxs, m.Name)
		//imapIns.MyBoxs[i] = m
		i++
	}
	if err := <-done; err != nil {
		return nil, err
	}
	// Convert to json
	jsonstr, err :=json.Marshal(boxs)
	if err != nil{
		return nil, err
	}
	return &jsonstr, nil
}



// Select the imap instalce's box
func(imapIns * ImapInstance)SelectBox(box string)error{
	_, err := imapIns.Client.Select(box, false)
	if err != nil {
		return err
	}
	return nil
}

// return the selected box status, nil if not select.
func(imapIns * ImapInstance)selected()*BoxStatus{
	return imapIns.Client.Mailbox()	
	// return &BoxStatus(boxstat)
}

// return the box name and number
func(imapIns * ImapInstance)BoxNameAndNum()(string, uint32){
	boxsta := imapIns.Client.Mailbox()
	if boxsta == nil{
		return "", 0
	}
	return boxsta.Name, boxsta.Messages
}

func (imapIns * ImapInstance)FetchItems(from uint32, to uint32, items []imap.FetchItem, cfunc ConvertFunc)(*[]*mymsg.Message, error){
	if from > to{
		return nil, fmt.Errorf("From > To!")
	}
	res := make([]*mymsg.Message, 0, int(to-from+1))
	// seqset
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)
	// items
	// items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags}
	//chan
	messages := make(chan *imap.Message, 10)
	//done := make(chan error, 1)
	// Fill the chan
	go func() {
		imapIns.Client.Fetch(seqset, items, messages)
	}()
	var res_err error
	for msg := range messages{
		mymsg,err := cfunc(msg)
		if err !=nil{
			res_err = err
			break
		}
		res = append(res, mymsg)
	}
	return &res, res_err	
}

// Get from to for msg retieval
// Page starts from 1
func getFT(msgnum uint32, page uint32)(uint32,uint32,bool){
	from := int(msgnum) - Msgbatch * int(page) + 1
	to := int(msgnum) - Msgbatch * (int(page)-1)
	var hasnext = true
	if from < 1 && to < 1{ // page too big
		from = 1
		to = 0
		hasnext = false
	}else if from < 1{ // last page
		from = 1
		hasnext = false
	}
	return uint32(from), uint32(to), hasnext
}

func genPageData(page int, hasnext bool, length int)map[string]interface{}{
	themap := make(map[string]interface{})
	themap["page"] = page
	themap["hasnext"] = hasnext
	themap["msglist"] = make([]mymsg.Message, 0, length)
	return themap
}

func getAddr(addrs[]*mail.Address)[]*netmail.Address{
	mailaddrs := make([]*netmail.Address, len(addrs))
	for idx, addr := range addrs{
		//mailaddrs = append(mailaddrs, &mail.Address(*addr))
		netaddr := netmail.Address(*addr)
		mailaddrs[idx] = &netaddr
	}
	return mailaddrs
}

func convAddr(addrs[]*imap.Address)[]*netmail.Address{
	mailaddrs := make([]*netmail.Address, len(addrs))
	for idx, addr := range addrs{
		//mailaddrs = append(mailaddrs, &mail.Address(*addr))
		netaddr := netmail.Address{Name:addr.PersonalName, Address:addr.Address()}
		mailaddrs[idx] = &netaddr
	}
	return mailaddrs
}

func infoConv(msg*imap.Message)(*mymsg.Message, error){
	mymsgIns := mymsg.Message{}
	env := msg.Envelope
	flags := msg.Flags
	mymsgIns.Subject = env.Subject
	mymsgIns.From 	= convAddr(env.From)
	mymsgIns.To 	= convAddr(env.To)
	mymsgIns.Date 	= env.Date
	mymsgIns.Flags 	= flags
	return &mymsgIns, nil
}

func (imapIns * ImapInstance)FetchPageInfo(page int)(*[]byte,error){
	boxsta := imapIns.selected()
	if boxsta == nil{
		return nil, fmt.Errorf("Please select a box First!")
	}
	msgnum := boxsta.Messages
	if msgnum == 0{
		emptymap,_ := json.Marshal(genPageData(1, false, 0))
		return &emptymap,nil
	}
	from, to, hasnext := getFT(msgnum, uint32(page))
	if from > to{
		return nil, fmt.Errorf("No more pages!")
	}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags}
	msgs, err := imapIns.FetchItems(from, to, items, infoConv)
	if err != nil{
		return nil, err
	}
	// Result
	themap := genPageData(page, hasnext, len(*msgs))
	for _, msg := range *msgs{
		themap["msglist"] = append(themap["msglist"].([]mymsg.Message), *msg)	
	}
	// convert to json
	jsonstr, err1 := json.Marshal(themap)
	if err1 != nil{
		return nil, err1
	}
	return &jsonstr, nil
}

func bodyConv(msg *imap.Message)(*mymsg.Message, error){
	var section imap.BodySectionName
	r := msg.GetBody(&section)
	if r == nil {
		return nil, fmt.Errorf("Server didn't returned message body")
	}
	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		return nil, err
	}
	mymsgprt := &mymsg.Message{}

	// Print some info about the message
	header := mr.Header
	if date, err := header.Date(); err == nil {
		mymsgprt.Date = date
	}
	if from, err := header.AddressList("From"); err == nil {
		mymsgprt.From = getAddr(from)
	}
	if to, err := header.AddressList("To"); err == nil {
		mymsgprt.From = getAddr(to)
	}
	if subject, err := header.Subject(); err == nil {
		mymsgprt.Subject = subject
	}
	var body string
	// Process each message's part
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("Error!",err)
			break
		}
		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			// This is the message's text (can be plain-text or HTML)
			b, _ := ioutil.ReadAll(p.Body)
			// log.Println("Got text: %v", string(b))
			body += string(b)
		case *mail.AttachmentHeader:
			// This is an attachment
			filename, _ := h.Filename()
			log.Println("Got attachment: %v", filename)
		}
	}
	mymsgprt.Body  = body
	return mymsgprt,nil
}

func (imapIns * ImapInstance)FetchItemInfo(page int, item int)(*[]byte,error){
	boxsta := imapIns.selected()
	if boxsta == nil{
		return nil, fmt.Errorf("Please select a box First!")
	}
	msgnum := boxsta.Messages
	from, to, _ := getFT(msgnum, uint32(page))
	if from > to{
		return nil, fmt.Errorf("No more pages!")
	}
	idx := from + uint32(item)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{section.FetchItem()}
	msgs, err := imapIns.FetchItems(idx, idx, items, bodyConv)
	if err != nil{
		return nil, err
	}
	// Result
	jsonstr, err1 := json.Marshal((*msgs)[0])
	if err1 != nil{
		return nil, err1
	}
	return &jsonstr, nil
}
