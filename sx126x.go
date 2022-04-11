package main

import (
	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
	"time"
)

type sx126x struct {
	rssi       bool
	crypt      bool
	address    int
	frequency  int
	serial_num int
	power      int
}

var uartOpts = UARTOptions{
	BitRate:  9600,
	DataBits: 8,
	StopBits: 1,
	Parity:   0,
	Timeout:  0,
}

func SX126xInitialize(sx sx126x) {
	sx126x_cfg := []byte{0xc2, 0x00, 0x09, 0x01, 0x02, 0x03, 0x62, 0x00, 0x12, 0x43, 0x00, 0x00}

	M0 := rpi.GPIO21
	M1 := rpi.GPIO27
	v_low := 0
	v_high := 1
	l_m0, _ := gpiod.RequestLine("gpiochip0", M0, gpiod.AsOutput(v_low))
	l_m1, _ := gpiod.RequestLine("gpiochip0", M1, gpiod.AsOutput(v_high))
	time.Sleep(0.1)

	// Needed in relay mode ONLY
	//low_address = sx.address & 0xff
	//high_address = sx.address >> 8 & 0xff
	//net_id_tmp = net_id & 0xff

	frequency_tmp := sx.frequency - 410
	//start_frequency := 410
	//offset_frequency := frequency_tmp

	air_speed_tmp := 2400
	buffer_sz_tmp := 240
	power_tmp := sx.power

	rssi_tmp := 0x00
	if sx.rssi {
		rssi_tmp = 0x80
	}

	sx126x_cfg[6] = byte(0x60 + air_speed_tmp)
	sx126x_cfg[7] = byte(buffer_sz_tmp + power_tmp + 0x20)
	sx126x_cfg[8] = byte(frequency_tmp)
	sx126x_cfg[9] = byte(0x43 + rssi_tmp)

	crypt_tmp := 0
	if sx.crypt {
		crypt_tmp = 1
	}

	lCrypt := byte(crypt_tmp & 0xff)
	hCrypt := byte(crypt_tmp >> 8 & 0xff)

	sx126x_cfg[10] = hCrypt
	sx126x_cfg[11] = lCrypt

	uartInst, _ := UARTOpen("/dev/serial0", uartOpts)
	for i := 0; i < 2; i++ {
		uartInst.Write(sx126x_cfg)

		time.Sleep(0.2)

		b_available, _ := uartInst.BytesAvailable()
		if b_available > 0 {
			b := []byte{byte(0)}
			uartInst.Read(b)
			if b[0] == 0xc1 {
				break
			}
		}
	}

	l_m0.SetValue(0)
	l_m1.SetValue(0)
	time.Sleep(0.1)

	uartInst.Close()
}

func SX126xPrintSettings() {
	l_m1, _ := gpiod.RequestLine("gpiochip0", 27, gpiod.AsOutput(1))
	time.Sleep(0.1)

	uartInst, _ := UARTOpen("/dev/serial0", uartOpts)

	askConfigCfg := []byte{0xc1, 0x00, 0x09}
	uartInst.Write(askConfigCfg)
	bufferRead := []byte{byte(0)}
	uartInst.Read(bufferRead)

	if bufferRead[0] == 0xc1 && bufferRead[2] == 0x09 {
		print("Frequency is" + string(bufferRead[8]) + ".125MHz.")
		print("Node address is" + string(bufferRead[3]+bufferRead[4]) + ".")
		print("Air speed" + string(bufferRead[6]&0x03) + "bps.")
		print("Power speed" + string(bufferRead[7]&0x03) + "dBm.")
	}

	l_m1.SetValue(0)

	uartInst.Close()
}
