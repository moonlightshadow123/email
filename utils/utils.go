package utils

import(
	"encoding/json"
)

func ToJson(f interface{})(*[]byte, error){
	jsonstr, err := json.Marshal(f)
	return &jsonstr, err
}


func FromJson(jsonstr *[]byte)(interface{}, error){
	var f interface{}
	err := json.Unmarshal(*jsonstr, &f)
	return f, err
}

func StrInSlice(list []string, a string) (int, bool) {
    for idx, b := range list {
        if b == a {
            return idx, true
        }
    }
    return -1, false
}