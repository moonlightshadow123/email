package main

import(
	"fmt"
	"test_proj/email/utils"
)

func main(){
	themap := make(map[string]interface{})
	themap["hello"] = "hello"
	themap["123"] = 123
	jsonstr, _ :=  utils.ToJson(themap)
	fmt.Println(*jsonstr)
	//var themap2 interface{}
	f, _ :=utils.FromJson(jsonstr)
	themap2 := f.(map[string]interface{})
	themap2["hello"] = "nihao"
	fmt.Println(themap)
	fmt.Println(themap2)
}