package mymsg

import (
	"time"
	"net/mail"
	"encoding/json"
)

var(

)

type Message struct{
	//to Addr
	//from Addr
	To []*mail.Address 	`json:"to"`
	From []*mail.Address	`json:"from"`
	Subject string 		`json:"subject"`
	Body string 		`json:"body"`
	Date time.Time 		`json:"date"`
	Flags []string 		`json:flags`
}

func (msg *Message)ToJson()(*[]byte, error){
	jsonstr, err := json.Marshal(msg)
	return &jsonstr, err
}

func FromJson(jsonstr *[]byte)(* Message, error){
	msg := Message{}
	err := json.Unmarshal(*jsonstr, &msg)
	return &msg, err  
}

