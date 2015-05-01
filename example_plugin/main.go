package main

import (
	"fmt"

	"github.com/natefinch/stdplug"
)

func main() {
	stdplug.Provide(Plugin{})
}

type Plugin struct{}

func (Plugin) SayHi(name string, response *string) error {
	*response = fmt.Sprintf("Hi %s", name)
	return nil
}
