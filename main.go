package main

import (
	"flag"
	"fmt"
	"vovan_ai/chat"
	"vovan_ai/client"
	"vovan_ai/config"
	"vovan_ai/terminal"
)

func main() {
	cnf, err := config.Load()
	if err != nil {
		fmt.Println(terminal.Color(err.Error(), terminal.Red))
		return
	}

	aiClient := client.NewAiClient(cnf)

	mode := flag.String("mode", "chat", "Режим работы vovan_ai (chat/agent)")

	flag.Parse()
	
	switch *mode{
	case "chat":
		chat := chat.NewChat(aiClient)
		chat.Run()
	default:
		fmt.Println(terminal.Color("Укажите режим работы vovan_ai (chat/agent)", terminal.Red))
	}
}
