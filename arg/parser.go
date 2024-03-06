package arg

import "fmt"

type Flag struct {
	name        string
	value       string
	description string
}

type SubCommand struct {
	name  string
	flags []Flag
}

func HelloName(name string) string {
	return fmt.Sprintf("Hello %v!", name)
}
