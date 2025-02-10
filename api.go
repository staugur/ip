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
	"log"
	"net"
	"net/http"
	"strings"

	"mip/third_party/xdb"
)

type ipinfo struct {
	Country  string // 国家
	Region   string // 区域
	Province string // 省份
	City     string // 城市
	ISP      string // 运营商
}

func cc(s string) string {
	return strings.ReplaceAll(s, "0", "")
}

func isIP(str string) bool {
	return net.ParseIP(str) != nil
}

func search(ip string) (i ipinfo, err error) {
	if !isIP(ip) {
		err = errors.New("invalid ip")
		return
	}

	// 用全局的 vIndex 创建带 VectorIndex 缓存的查询对象。
	// 备注：并发使用，全部 goroutine 共享全局的只读 vIndex 缓存，每个 goroutine 创建一个独立的 searcher 对象
	searcher, err := xdb.NewWithVectorIndex(dbpath, vIndex)
	if err != nil {
		log.Printf("failed to create searcher with vector index: %s\n", err)
		return
	}

	txt, err := searcher.SearchByStr(ip)
	if err != nil {
		return
	}
	parts := strings.Split(txt, "|")
	if len(parts) != 5 {
		err = errors.New("invalid format")
		return
	}

	return ipinfo{
		Country:  parts[0],
		Region:   parts[1],
		Province: parts[2],
		City:     parts[3],
		ISP:      parts[4],
	}, nil
}

func realIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		log.Printf("First try to fetch XFF: %s\n", ip)
		i := strings.IndexAny(ip, ", ")
		if i > 0 {
			return ip[:i]
		}
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		log.Printf("Second try to fetch X-Real-IP: %s\n", ip)
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
	log.Printf("query ip: %s\n", ip)
	return ip
}

func getArea(ip ipinfo) string {
	area := cc(ip.Country)
	if area == "" {
		area = cc(ip.Region)
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
