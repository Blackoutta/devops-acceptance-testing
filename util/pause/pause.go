package pause

import (
	"fmt"
	"log"
)

func PauseExecution() {
	log.Println("程序暂停，输入任意字符来执行teardown...")
	var continueInput string
	_, err := fmt.Scanln(&continueInput)
	if err != nil {
		log.Fatalln("err scanning user input:", err)
	}
	log.Println(continueInput)
	log.Println("收到用户输入，继续向下执行teardown")
}
