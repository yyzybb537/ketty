package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	r := &Req{}
	r.Hy = new(Hyaline)
	r.Hy.Val = []int64{1, 2}
	b, _ := json.Marshal(r)
	println(string(b))
	w := &Req{}
	err := json.Unmarshal(b, w)
	if err != nil {
		println("Error:", err.Error())
	} else {
		fmt.Printf("w.Val: %v\n", w.GetHy().GetVal())
	}
}
