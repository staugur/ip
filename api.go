/*
   Copyright 2021 Hiroshi.tao

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"ip/ip2region"
	"net"
	"net/http"
	"strings"
)

func cc(s string) string {
	return strings.ReplaceAll(s, "0", "")
}

func isIP(str string) bool {
	return net.ParseIP(str) != nil
}

func search(ip string) (info ip2region.IpInfo, err error) {
	if !isIP(ip) {
		err = errors.New("invalid ip")
		return
	}
	region, err := ip2region.New(db)
	if err != nil {
		return
	}
	defer region.Close()
	if btree {
		return region.BtreeSearch(ip)
	} else {
		return region.MemorySearch(ip)
	}
}

func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		i := strings.IndexAny(ip, ", ")
		if i > 0 {
			return ip[:i]
		}
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	ra, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ra
}

func getIP(r *http.Request) string {
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		ip = realIP(r)
	}
	return ip
}

func getArea(info ip2region.IpInfo) string {
	area := cc(info.Country)
	if area == "" {
		area = cc(info.Region)
	}
	return area
}

func ipView(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, realIP(r))
}

func addrView(w http.ResponseWriter, r *http.Request) {
	ip := getIP(r)
	i, err := search(ip)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
	fmt.Fprintf(
		w, "IP：%s\n国家/地区：%s\n省份：%s\n城市：%s\n运营商：%s\n",
		ip, getArea(i), cc(i.Province), cc(i.City), cc(i.ISP),
	)
}

type restBase struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type restInfo struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	ISP      string `json:"isp"`
}

type resp struct {
	restBase
	restInfo `json:"data"`
}

func apiErrView(w http.ResponseWriter, err error) {
	api := restBase{Message: err.Error()}
	resp, _ := json.Marshal(api)
	w.Write(resp)
}

func restView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ip := getIP(r)
	i, err := search(ip)
	if err != nil {
		apiErrView(w, err)
		return
	}
	api := resp{restBase{true, "ok"}, restInfo{
		ip, getArea(i), cc(i.Province), cc(i.City), cc(i.ISP),
	}}
	resp, err := json.Marshal(api)
	if err != nil {
		apiErrView(w, err)
		return
	}
	w.Write(resp)
}
