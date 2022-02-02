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

package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/SuperGreenLab/EmergencyUpgrader/internal/firmware"
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

func GetStringParameter(ip, param string) (string, error) {
	url := fmt.Sprintf("http://%s/s?k=%s", ip, param)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	sb := string(body)
	return sb, nil
}

func SetStringParameter(ip, param, value string) (string, error) {
	url := fmt.Sprintf("http://%s/s?k=%s&v=%s", ip, param, value)
	resp, err := client.Post(url, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	sb := string(body)
	return sb, nil
}

func GetIntParameter(ip, param string) (int, error) {
	url := fmt.Sprintf("http://%s/i?k=%s", ip, param)
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	sb := string(body)
	return strconv.Atoi(sb)
}

func SetIntParameter(ip, param string, value int32) (string, error) {
	url := fmt.Sprintf("http://%s/i?k=%s&v=%d", ip, param, value)
	resp, err := client.Post(url, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	sb := string(body)
	return sb, nil
}

func GetMyIP(ip string) (string, error) {
	url := fmt.Sprintf("http://%s/myip", ip)
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	sb := string(body)
	return sb, nil
}

func UploadFile(ip, lpath, rpath string) (string, error) {
	file, err := firmware.FirmwareFS.Open(lpath)

	if err != nil {
		return "", err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return "", err
	}
	// get the size
	size := fi.Size()

	url := fmt.Sprintf("http://%s%s", ip, rpath)
	request, _ := http.NewRequest("POST", url, file)
	request.Header.Add("Content-Type", "application/octet-stream")
	request.ContentLength = size
	response, err := client.Do(request)

	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
