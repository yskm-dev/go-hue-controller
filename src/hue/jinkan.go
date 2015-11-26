package main

import (
	"fmt"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/tty.usbmodem1421")
	sensor := gpio.NewAnalogSensorDriver(firmataAdaptor, "sensor", "0")

	work := func() {
		gobot.On(sensor.Event("data"), func(data interface{}) {
			fmt.Println("Sensor", data)
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{sensor},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
