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
)

const version = "0.2.0"

var (
	v bool

	host   string
	port   uint
	prefix string
	db     string
	btree  bool
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.BoolVar(&v, "v", false, "show version and exit")

	flag.StringVar(&host, "host", "0.0.0.0", "http listen host")
	flag.UintVar(&port, "port", 7000, "http listen port")

	flag.StringVar(&prefix, "prefix", "", "route prefix")
	flag.StringVar(&db, "db", "data/ip2region.db", "the ip2region.db filepath")
	flag.BoolVar(&btree, "btree", true, "use b-tree algorithm(if no, use memory)")
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
	if db == "" {
		db = os.Getenv("ip_db")
	}
	if prefix == "" {
		prefix = os.Getenv("ip_prefix")
	}
	if os.Getenv("ip_use") == "memory" {
		btree = false
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
	if stat, err := os.Stat(db); err != nil || stat.IsDir() {
		fmt.Println("invalid db")
		return
	}
	if prefix != "" && !strings.HasPrefix(prefix, "/") {
		fmt.Println("prefix must start with /")
		return
	}

	http.HandleFunc(routePattern("/myip"), ipView)
	http.HandleFunc(routePattern("/addr"), addrView)
	http.HandleFunc(routePattern("/rest"), restView)
	listen := fmt.Sprintf("%s:%d", host, port)
	log.Println("HTTP listen on " + listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}
