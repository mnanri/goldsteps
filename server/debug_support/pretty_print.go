package common

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func prettyPrint(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
		return
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, data, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(buf.String())
}
