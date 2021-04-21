package main

import "net/http"

func main() {
	http.HandleFunc("/hello",
		func(writer http.ResponseWriter, r *http.Request) {
			s := r.URL.RawQuery
			writer.Write([]byte("hello,world." + s))
		})
	http.ListenAndServe(":8082", nil)
}
