package adapter

import (
	"fmt"
	"net/http"
)

func doGet(url string) (r *http.Response, e error) {
	resp, err := http.Get(url)
	if err != nil {
		//fmt.Println(resp.StatusCode)
		fmt.Println(err.Error())
	}
	return resp, err
}