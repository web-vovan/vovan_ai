package agent

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"vovan_ai/client"
	"vovan_ai/terminal"
	"vovan_ai/tools"
)

type Agent struct {
	Client   *client.AiClient
	Tools    map[string]tools.Tool
	MaxSteps int
}

func NewAgent(client *client.AiClient) *Agent {
	bashTool := tools.NewBashTools()

	return &Agent{
		Client:   client,
		MaxSteps: 100,
		Tools: map[string]tools.Tool{
			bashTool.Name(): bashTool,
		},
	}
}

func (a *Agent) Run() {
	messages := make([]client.Message, 0)

	messages = append(messages, client.Message{
		Role:    "system",
		Content: "Ты — VOVAN_AI, агент для кодинга. Объясняй пользователю, что ты делаешь",
	})

	toolSchemas := make([]client.Tool, 0)
	for _, t := range a.Tools {
		toolSchemas = append(toolSchemas, t.Schema())
	}

	reader := bufio.NewReader(os.Stdin)

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
	fmt.Println("Начинаю работу")
	fmt.Println()

	for i := 1; i <= a.MaxSteps; i++ {
		resp, err := a.Client.SendChat(messages, toolSchemas)
		if err != nil {
			fmt.Printf("ошибка при запросе к модели %s", err.Error())
			return
		}

		messages = append(messages, *resp)

		if len(resp.ToolCalls) > 0 {
			text := fmt.Sprintf("ШАГ %d: ", i)
			fmt.Print(terminal.Color(text, terminal.Cyan))
			fmt.Println(terminal.Color(resp.Reasoning, terminal.Purple))
			fmt.Println()

			for _, tc := range resp.ToolCalls {
				t, ok := a.Tools[tc.Function.Name]

				if !ok {
					fmt.Printf("не найдет инструмент %s", tc.Function.Name)
					return
				}

				r := t.Execute(tc.Function.Arguments)

				if r == "" {
					r = "команда выполнена"
				}

				messages = append(messages, client.Message{
					Role:       "tool",
					Content:    r,
					ToolCallId: tc.ID,
				})
			}
		} else {
			fmt.Println(terminal.Color("Результат:", terminal.Cyan))
			fmt.Println(resp.Content)
			fmt.Println()

			return
		}
	}

	fmt.Println("исчерпан лимит шагов")
}
