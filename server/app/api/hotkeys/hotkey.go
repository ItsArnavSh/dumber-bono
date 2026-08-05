package hotkeys

import (
	"dubmer-bono/app/api/hotkeys/hyperland"
	"dubmer-bono/app/api/types"
	"fmt"
	"runtime"
)

func NewHotKeyListner() (types.HotKeyHandler, error) {
	os_name := runtime.GOOS
	fmt.Println(os_name)
	if os_name == "linux" && hyperland.IsUsingHyperland() {
		fmt.Println("Hyperland")
		return &hyperland.HyperlandHotkeys{}, nil //Dev testing
	} else {

		fmt.Println("General")
		//		return &general.CrossPlatformHotkeys{}, nil
	}
	return &hyperland.HyperlandHotkeys{}, nil
}

var _ types.HotKeyHandler = &hyperland.HyperlandHotkeys{}

//var _ types.HotKeyHandler = &general.CrossPlatformHotkeys{}
