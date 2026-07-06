package main

import (
	"fmt"

	"golang.design/x/hotkey"
)

func main() {
	hk := hotkey.New([]hotkey.Modifier{}, hotkey.KeyA)
	err := hk.Register()
	if err != nil {
		panic(err)
	}
	fmt.Println("Hotkey registered, press A")

	for {
		<-hk.Keydown()
		fmt.Println("PTT: START")
		<-hk.Keyup()
		fmt.Println("PTT: STOP")
	}
}
