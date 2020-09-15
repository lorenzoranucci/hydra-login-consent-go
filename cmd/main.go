package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/lorenzoranucci/hydra-login-consent-go/internal/cmd"
)

var version = "dev"

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()

	err := cmd.GetApp(version).Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}
