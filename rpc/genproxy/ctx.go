package genproxy

import (
	"fmt"
	"strings"
)

type ProxyContext struct {
	Dst     string
	Out     []string
	TypeMap map[string]string
}

func ParseMap(arr []string) map[string]string {
	var mp = make(map[string]string)
	for _, item := range arr {
		if !strings.Contains(item, "=") {
			fmt.Println("can`t parse ", item, " , skip value")
			continue
		}
		values := strings.Split(item, "=")
		if len(values) != 2 {
			fmt.Println("can`t parse ", item, " , skip value")
			continue
		}
		key, val := values[0], values[1]
		mp[key] = val
	}
	return mp
}
