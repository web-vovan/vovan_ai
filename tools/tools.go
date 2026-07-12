package tools

import "vovan_ai/client"

type Tool interface {
	Name() string
	Description() string
	Schema() client.Tool
	Execute(argJson string) string
}
