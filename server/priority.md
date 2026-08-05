message: 
received message
fatal error: concurrent map read and map write

goroutine 61 [running]:
internal/runtime/maps.fatal({0x2922155?, 0x37f1ce0?})
	/usr/lib/go/src/runtime/panic.go:1181 +0x18
dubmer-bono/app/service/radio.(*Service).GetMessageByMinPriority(0x16e9c51f0c80)
	/home/arnav/vault/dumber-bono/server/app/service/radio/msg.go:15 +0x66
dubmer-bono/app/service/radio.(*Service).radioTheDriver(0x16e9c51f0c80, {0x295c868, 0x3a22be0})
	/home/arnav/vault/dumber-bono/server/app/service/radio/radio.go:65 +0x10c
created by dubmer-bono/app/service/radio.NewService in goroutine 1
	/home/arnav/vault/dumber-bono/server/app/service/radio/radio.go:49 +0x1d0
exit status 2
