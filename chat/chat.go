package chat

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"vovan_ai/client"
	"vovan_ai/terminal"
)

type Chat struct {
    Client *client.AiClient
}

func NewChat(client *client.AiClient) *Chat {
    return &Chat{client}
}

func (c *Chat) Run() {
    onChunk := func(data string) {
		fmt.Print(data)
	}

	reader := bufio.NewReader(os.Stdin)
	messages := make([]client.Message, 0)

	for i := 0; i < 100; i++ {
		fmt.Print(terminal.Color("USER: ", terminal.Green))
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		input = strings.TrimSpace(input)
		fmt.Println()

		msg := client.Message{
			Role:    "user",
			Content: input,
		}

		messages = append(messages, msg)

		fmt.Print(terminal.Color("VOVAN_AI: ", terminal.Cyan))
		resp, err := c.Client.SendStreamChat(messages, onChunk)
		if err != nil {
			fmt.Println(err)
			return
		}

		messages = append(messages, client.Message{
			Role:    "assistant",
			Content: resp,
		})

		fmt.Println()
		fmt.Println()
	}

    fmt.Println("исчерпан лимит диалога")
}