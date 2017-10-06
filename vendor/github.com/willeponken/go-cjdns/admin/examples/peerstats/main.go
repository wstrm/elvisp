package main

import (
	"fmt"
	"github.com/willeponken/go-cjdns/admin"
	"github.com/kylelemons/godebug/pretty"
)

func main() {
	Pretty := pretty.Config{
		PrintStringers: true,
	}
	conn, err := admin.Connect(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	results, err := conn.InterfaceController_peerStats()
	if err != nil {
		fmt.Println(err)
		return
	}
	Pretty.Print(results)
}
