package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	//"strings"
)

func main() {
	response, err := http.Get("https://gateway.pinata.cloud/ipns/QmUgddgEc71BH5movhDtLJ91tGy3pKs5iUEgsa69ewCNog")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		leer := string(contents)
		//fmt.Println(strings.Index(leer, "Q"))
		//i := strings.Index(leer, "")
		//chars := leer[:i]
		//arefun := leer[i+1:]
		fmt.Println(leer) //chars)
		//fmt.Println(arefun)
		//fmt.Printf("%s\n", string(contents))

	}
}
