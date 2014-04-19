package main

import (
	"runtime"
	"fmt"
	"log"
	"net/http"
	_ "github.com/hogedigo/shizgo/webapp/handler"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func main() {
	http.HandleFunc("/hello", helloWorld)       //アクセスのルーティングを設定します。
	err := http.ListenAndServe(":9876", nil)
	if err != nil {
		log.Fatal("error occurred: ", err)
	}
}
