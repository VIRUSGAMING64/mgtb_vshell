package web

import (
	"fmt"
	"govshell/core"
	"net/http"
)

func Start() {
	fmt.Println("starting webpage")
	http.Handle("/", http.FileServer(http.Dir("web/public")))
	http.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
		if !core.Offline {
			fmt.Println("I am alive!!!")
		}
		s := "{\"alive\": true}"
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(s))
	})
	err := http.ListenAndServe(":80", nil)
	fmt.Println(err)
}
