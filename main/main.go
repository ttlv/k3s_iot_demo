package main

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"

	"github.com/bettercap/gatt"
	"github.com/bettercap/gatt/examples/option"
	"github.com/urfave/cli"
	"k3s_iot_demo"
	"k8s.io/klog"
)

var (
	Version   = "v0.0.0-dev"
	GitCommit = "HEAD"
)

func main() {
	app := cli.NewApp()
	app.Name = "bluetooth-device"
	app.Version = fmt.Sprintf("%s (%s)", Version, GitCommit)
	app.Usage = "Bluetooth device adaptor"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "device_name",
			EnvVar: "DEVICE_NAME",
			Value:  "", // config your decice name in device-deployment.yaml
		},
		cli.StringFlag{
			Name:   "device_mac_address",
			EnvVar: "DEVICE_MAC_ADDRESS",
			Value:  "", // config your device mac address name in device-deployment.yaml
		},
		cli.StringFlag{
			Name:   "device_mqtt",
			EnvVar: "DEVICE_MQTT",
			Value:  "", // config your device mqtt url in device-deployment.yaml
		},
	}
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		klog.Fatal(err)
	}
}

func run(c *cli.Context) {
	klog.Info("Starting bluetooth device adaptor")
	// handle device options
	name := c.String("device_name")
	macAddress := c.String("device_mac_address")
	deviceMQTT := c.String("device_mqtt")

	if len(name) == 0 && len(macAddress) == 0 {
		klog.Fatal("device name or device MAC address is required.")
	}

	if len(deviceMQTT) == 0 {
		klog.Fatal("mqtt config is required for device.")
	}
	mqtt, err := loadMqttConfig(deviceMQTT)
	if err != nil {
		klog.Fatal(err)
	}

	dv := k3s_iot_demo.Device{
		Name:       name,
		MacAddress: macAddress,
		Mqtt:       mqtt,
	}

	cli, err := k3s_iot_demo.ConnectToMQTT("bluetooth-temp", dv.Mqtt)
	if err != nil {
		log.Fatalf("Failed to connect to the mqtt server, err: %s\n", err)
		return
	}
	k3s_iot_demo.MqttCli = cli

	// subscribe to the MQTT device topic
	err = k3s_iot_demo.SubscribeToMQTT(cli, dv.Mqtt)
	if err != nil {
		log.Fatalf("Failed to subscribe the mqtt server, err: %s\n", err)
		return
	}
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}
	for {
		// Register handlers.
		d.Handle(
			gatt.PeripheralDiscovered(dv.OnPeriphDiscovered),
			gatt.PeripheralConnected(dv.OnPeriphConnected),
			gatt.PeripheralDisconnected(dv.OnPeriphDisconnected),
		)
		d.Init(dv.OnStateChanged)
		<-k3s_iot_demo.Done
		fmt.Println("Done")
		fmt.Println("Sleep for 10 seconds")
		time.Sleep(10 * time.Second)
	}

	// disconnect mqtt cli
	if err := cli.Disconnect(); err != nil {
		logrus.Fatalln(err)
	}
}

func loadMqttConfig(mqttStr string) (k3s_iot_demo.Mqtt, error) {
	mqtt := k3s_iot_demo.Mqtt{}
	if len(mqttStr) != 0 {
		err := json.Unmarshal([]byte(mqttStr), &mqtt)
		if err != nil {
			return mqtt, fmt.Errorf("failed to unmarshall mqtt env:%s, err: %s", mqttStr, err.Error())
		}
	}
	return mqtt, nil
}
