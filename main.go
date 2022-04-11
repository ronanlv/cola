package main

import "time"

func main() {
	//config := water.Config{
	//	DeviceType: water.TAP,
	//}
	//config.Name = "test-gw"
	//
	//ifce, err := water.New(config)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//var frame ethernet.Frame
	//
	//for {
	//	frame.Resize(1500)
	//	n, err := ifce.Read([]byte(frame))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	frame = frame[:n]
	//	log.Printf("Dst: %s\n", frame.Destination())
	//	log.Printf("Src: %s\n", frame.Source())
	//	log.Printf("Ethertype: % x\n", frame.Ethertype())
	//	log.Printf("Payload: % x\n", frame.Payload())
	//}

	sx := sx126x{
		rssi:       false,
		crypt:      false,
		address:    0,
		frequency:  433,
		serial_num: 0,
		power:      22,
	}

	SX126xInitialize(sx)
	time.Sleep(1)
	SX126xPrintSettings()
}
