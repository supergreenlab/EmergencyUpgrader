/*
 * Copyright (C) 2022  SuperGreenLab <towelie@supergreenlab.com>
 * Author: Constantin Clauzel <constantin.clauzel@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package input

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SuperGreenLab/EmergencyUpgrader/internal/controller"
	"github.com/SuperGreenLab/EmergencyUpgrader/internal/firmware"
)

func ffs(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func input() {
	var err error
	time.Sleep(1 * time.Second)
	fmt.Println("Welcome to SuperGreenLab emergency upgrader!")
	fmt.Println("---------------------")
	fmt.Println("Make sure your computer is connected to the emoji wifi if it is visible.")
	fmt.Println("Once connected, please enter the controller's ip address then press enter, leave empty for default (when connected to the emoji wifi):")

	reader := bufio.NewReader(os.Stdin)
	ip, err := reader.ReadString('\n')
	ffs(err)
	fmt.Printf(ip)
	if ip == "\n" {
		ip = "192.168.4.1"
	} else {
		ip = ip[0 : len(ip)-1]
	}
	fmt.Printf("Using controller with IP %s\n", ip)

	fmt.Println("Setting up controller for upgrade.")

	baseDir, err := controller.GetStringParameter(ip, "OTA_BASEDIR")
	ffs(err)

	tf := fmt.Sprintf("fs%s/last_timestamp", baseDir)
	timestamp, err := firmware.FirmwareFS.ReadFile(tf)
	ffs(err)

	ctimestamp, err := controller.GetIntParameter(ip, "OTA_TIMESTAMP")
	ffs(err)

	_, err = controller.UploadFile(ip, fmt.Sprintf("fs%s/%s/html_app/app.html", baseDir, string(timestamp[0:len(timestamp)-1])), "/fs/app.html")
	ffs(err)
	fmt.Println("Uploaded app.html.... OK")

	_, err = controller.UploadFile(ip, fmt.Sprintf("fs%s/%s/html_app/config.json", baseDir, string(timestamp[0:len(timestamp)-1])), "/fs/config.json")
	ffs(err)
	fmt.Println("Uploaded config.json.... OK")

	myip, err := controller.GetMyIP(ip)
	ffs(err)
	_, err = controller.SetStringParameter(ip, "OTA_BASEDIR", fmt.Sprintf("/fs%s", baseDir))
	ffs(err)
	_, err = controller.SetStringParameter(ip, "OTA_SERVER_IP", myip)
	ffs(err)
	_, err = controller.SetIntParameter(ip, "OTA_SERVER_PORT", 8081)
	ffs(err)

	fmt.Println("Rebooting controller..")
	_, _ = controller.SetIntParameter(ip, "REBOOT", 1)
	//ffs(err)
	fmt.Println("")
	fmt.Println("Controller should have rebooted.")
	fmt.Println("")
	fmt.Println("Now make sure your computer stays connected to the emoji wifi. If it disconnects while the controller is rebooting, please reconnect it.")
	fmt.Println("The controller will automatically upgrade itself, you can monitor that by looking at the OTA section 'timestamp' parameter.")
	fmt.Printf("It's value should change from %d to %s\n", ctimestamp, timestamp)
	fmt.Println("")
	fmt.Println("Please press any key when it's done.")
	reader.ReadRune()
	fmt.Println("Cleaning up params")
	_, err = controller.SetStringParameter(ip, "OTA_BASEDIR", baseDir)
	ffs(err)
	fmt.Println("Cleaned up params")
	time.Sleep(1 * time.Second)
	os.Exit(0)
}

func Init() {
	go input()
}
