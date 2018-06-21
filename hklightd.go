package main

import (
	"github.com/brutella/hc"
	"github.com/brutella/hc/accessory"
	"log"
	"github.com/maznikoff/yeelight"
	"fmt"
	"time"
	"context"
)

func turnLightOn() {
	log.Println("Turn Light On")
}

func turnLightOff() {
	log.Println("Turn Light Off")
}

func setHue(newValue float64) {
	log.Printf("Set hue to %f", newValue)
}

func setSaturation(newValue float64) {
	log.Printf("Set saturation to %f", newValue)
}

func setBrightness(newValue int) {
	log.Printf("Set brightness to %d", newValue)
}

func checkError(err error) {
	if nil != err {
		log.Fatal(err)
	}
}

func main() {
	info := accessory.Info{
		Name:         "Personal Light Bulb",
		Manufacturer: "Matthias",
	}

	acc := accessory.NewLightbulb(info)

	acc.Lightbulb.On.OnValueRemoteUpdate(func(on bool) {
		if on == true {
			turnLightOn()
		} else {
			turnLightOff()
		}
	})

	acc.Lightbulb.Hue.OnValueRemoteUpdate(func(newValue float64) {
		setHue(newValue)
	})
	acc.Lightbulb.Saturation.OnValueRemoteUpdate(func(newValue float64) {
		setSaturation(newValue)
	})
	acc.Lightbulb.Brightness.OnValueRemoteUpdate(func(newValue int) {
		setBrightness(newValue)
	})

	y, err := yeelight.Discover()
	checkError(err)

	on, err := y.GetProp("power")
	checkError(err)
	fmt.Printf("Power is %s", on[0].(string))

	notifications, done, err := y.Listen()
	checkError(err)
	go func() {
		<-time.After(time.Minute * 30)
		done <- struct{}{}
	}()
	for n := range notifications {
		fmt.Println(n)
	}

	context.Background().Done()


	t, err := hc.NewIPTransport(hc.Config{Pin: "32191123"}, acc.Accessory)
	if err != nil {
		log.Fatal(err)
	}

	hc.OnTermination(func() {
		t.Stop()
	})

	t.Start()
}
