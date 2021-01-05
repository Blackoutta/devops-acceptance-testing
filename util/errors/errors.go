package errors

import (
	"encoding/json"
	"log"
)

func HandleError(msg string, err error) {
	if err != nil {
		log.Printf(msg+": %v\n", err)
	}
}

func UnmarshalAndHandleError(resp []byte, respStruct interface{}) {
	err := json.Unmarshal(resp, respStruct)
	if err != nil {
		log.Printf("err unmarshaling json: %v\n", err.Error())
	}
}
