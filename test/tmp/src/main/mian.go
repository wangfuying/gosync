package main 

import (
	"fmt"
	"util"
	"os/exec"
)

func main() {
	fmt.Println("Hello, World!")
	util.Say("Hi")
	exec.Command("PAUSE").Run()
}

