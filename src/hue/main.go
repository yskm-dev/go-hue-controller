package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
)

const IP_ADDRESS string = "[your api address]"
const USER_NAME string = "[your user name]"
const HUE_INDEX string = "[your hue index]"

func getHueAPI() string {
	return "http://" + IP_ADDRESS + "/api/" + USER_NAME + "/lights/" + HUE_INDEX + "/state"
}

var hueAPI = getHueAPI()

type State struct {
	On bool `json:"on"`
}

func hue(isOn bool) {
	client := &http.Client{}
	data := State{isOn}
	state_json, _ := json.Marshal(data)
	post_body := strings.NewReader(string(state_json))
	request, err := http.NewRequest("PUT", hueAPI, post_body)
	if err != nil {
		fmt.Println(err)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
}

func main() {
	hue(false)
	gbot := gobot.NewGobot()

	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/tty.usbmodem1421")
	sensor := gpio.NewAnalogSensorDriver(firmataAdaptor, "sensor", "0")

	work := func() {
		lastCalled := time.Now().Second()
		sensorState := false
		gobot.On(sensor.Event("data"), func(data interface{}) {
			lastCalled = time.Now().Second()
			if sensorState == false {
				fmt.Println("Sensor On")
				sensorState = true
				hue(sensorState)
			}
		})
		go func() {
			for {
				if time.Now().Second()-lastCalled > 2 {
					if sensorState == true {
						fmt.Println("Sensor Off")
						sensorState = false
						hue(sensorState)
					}
				}
				time.Sleep(time.Millisecond * 100)
			}
		}()
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{sensor},
		work,
	)

	gbot.AddRobot(robot)

	gbot.Start()
}
