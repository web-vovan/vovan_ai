package tools

import (
	"encoding/json"
	"os/exec"

	"vovan_ai/client"
)

type Bash struct{}

type Args struct {
	Command string `json:"command"`
}

func NewBashTools() *Bash {
	return &Bash{}
}

func (b *Bash) Name() string {
	return "bash"
}

func (b *Bash) Description() string {
	return "Выполняет команду в bash и возвращает stdout и stderr"
}

func (b *Bash) Schema() client.Tool {
	return client.Tool{
		Type: "function",
		Function: client.FunctionTool{
			Name:        b.Name(),
			Description: b.Description(),
			Parameters: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"command": map[string]interface{}{
						"type":        "string",
						"description": "Команда выполнения в bash",
					},
				},
				"required": []string{"command"},
			},
		},
	}
}

func (b *Bash) Execute(input string) string {
	var args Args

	err := json.Unmarshal([]byte(input), &args)
	if err != nil {
		return "ошибка при разборе аргументов: " + err.Error()
	}

	cmd := exec.Command("bash", "-c", args.Command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "ошибка выполнения команды: " + err.Error() + "\nвывод: " + string(output)
	}

	return string(output)
}
