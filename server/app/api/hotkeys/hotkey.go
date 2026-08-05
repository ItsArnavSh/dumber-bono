package hotkeys

import (
	"dubmer-bono/app/api/hotkeys/general"
	"dubmer-bono/app/api/hotkeys/hyperland"
	"dubmer-bono/app/api/types"
	"runtime"
)

func NewHotKeyListner() (types.HotKeyHandler, error) {
	os_name := runtime.GOOS

	if os_name == "Linux" && hyperland.IsUsingHyperland() {
		return &hyperland.HyperlandHotkeys{}, nil //Dev testing
	} else {
		return &general.CrossPlatformHotkeys{}, nil
	}
}

var _ types.HotKeyHandler = &hyperland.HyperlandHotkeys{}
var _ types.HotKeyHandler = &general.CrossPlatformHotkeys{}
