package tool

import (
	"encoding/json"
	"fmt"
)

func Dump(data interface{}) {
	buf, err := json.MarshalIndent(data, "", "   ")
	if err != nil {
		fmt.Println("err = ", err)
		return
	}

	fmt.Printf("\n %c[1;0;36m%s%c[0m\n\n", 0x1B, string(buf), 0x1B)
}
