package main

import (
	"fmt"
	"log"
	"net/http"
	_ "github.com/hogedigo/shizgo/webapp/handler"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // formパラメータを解析する
	fmt.Fprintf(w, "Hello " + r.FormValue("name") + "!")
}

func main() {
	http.HandleFunc("/hello", helloWorld)       //アクセスのルーティングを設定します。
	err := http.ListenAndServe(":9876", nil) //監視するポートを設定します。
	if err != nil {
		log.Fatal("error occurred: ", err)
	}
}
