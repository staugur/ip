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
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"mip/third_party/xdb"
)

const version = "0.3.0"

var (
	v bool

	host   string
	port   uint
	prefix string
	dbpath string

	vIndex []byte
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.StringVar(&host, "host", "0.0.0.0", "http listen host")
	flag.UintVar(&port, "port", 7000, "http listen port")
	flag.StringVar(&prefix, "prefix", "", "route prefix")
	flag.StringVar(&dbpath, "db", "data/ip2region.xdb", "the ip2region.xdb filepath")
}

func main() {
	flag.Parse()
	if v {
		fmt.Println(version)
	} else {
		startServer()
	}
}

func routePattern(pattern string) string {
	return strings.TrimSuffix(prefix, "/") + pattern
}

func startServer() {
	if dbpath == "" {
		dbpath = os.Getenv("ip_db")
	}
	if prefix == "" {
		prefix = os.Getenv("ip_prefix")
	}
	envhost := os.Getenv("ip_host")
	envport := os.Getenv("ip_port")
	if envhost != "" {
		host = envhost
	}
	if envport != "" {
		envport, err := strconv.Atoi(envport)
		if err != nil {
			fmt.Println("Invalid environment ip_port")
			return
		}
		port = uint(envport)
	}
	if stat, err := os.Stat(dbpath); err != nil || stat.IsDir() {
		fmt.Println("invalid db file")
		return
	}
	if prefix != "" && !strings.HasPrefix(prefix, "/") {
		fmt.Println("prefix must start with /")
		return
	}

	// 从 dbpath 加载 VectorIndex 缓存，把下述 vIndex 变量全局到内存里面。
	var err error
	vIndex, err = xdb.LoadVectorIndexFromFile(dbpath)
	if err != nil {
		panic(err)
	}

	http.HandleFunc(routePattern("/myip"), ipView)
	http.HandleFunc(routePattern("/addr"), addrView)
	http.HandleFunc(routePattern("/rest"), restView)
	listen := fmt.Sprintf("%s:%d", host, port)
	log.Println("HTTP listen on " + listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}
