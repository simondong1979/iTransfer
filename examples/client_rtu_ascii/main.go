package main

import (
	"fmt"
	"github.com/Transefer/iTransfer/modbus"
	"time"

	"github.com/Transfer/iTransfer/serial"
	//	modbus "github.com/things-go/go-modbus"
)

func main() {
	p := modbus.NewRTUClientProvider(modbus.WithEnableLogger(),
		modbus.WithSerialConfig(serial.Config{
			Address:  "/dev/ttyUSB0",
			BaudRate: 115200,
			DataBits: 8,
			StopBits: 1,
			Parity:   "N",
			Timeout:  modbus.SerialDefaultTimeout,
		}))

	client := modbus.NewClient(p)
	err := client.Connect()
	if err != nil {
		fmt.Println("connect failed, ", err)
		return
	}
	defer client.Close()

	fmt.Println("starting")
	for {
		_, err := client.ReadCoils(3, 0, 10)
		if err != nil {
			fmt.Println(err.Error())
		}

		//	fmt.Printf("ReadDiscreteInputs %#v\r\n", results)

		time.Sleep(time.Second * 2)
	}
}
