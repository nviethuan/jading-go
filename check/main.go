package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	var key string
	flag.StringVar(&key, "f", "", "f")
	flag.Parse()
	cmd := exec.Command("./ss3", "-key", key)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Process started, PID:", cmd.Process.Pid)
}
