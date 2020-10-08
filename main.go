package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("regexpreprocessor test!")
	js.Global().Set("regextimestamp", Wrapper())
	<-make(chan bool)
}

func regexTimestamp(re, name, input string) (string, error) {
	cfg := RegexTimestampConfig{
		Regex:                 re,
		TS_Match_Name:         name,
		Assume_Local_Timezone: false,
	}
	rts, err := NewRegexTimestampProcessor(cfg)
	if err != nil {
		return "", err
	}

	ent := &Entry{
		Data: []byte(input),
	}
	ret, err := rts.Process(ent)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", ret[0]), nil
}

func Wrapper() js.Func {
	Func := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 3 {
			return "Invalid number of arguments passed"
		}
		re := args[0].String()
		name := args[1].String()
		input := args[2].String()

		out, err := regexTimestamp(re, name, input)
		if err != nil {
			return err.Error()
		}
		return out
	})
	return Func
}
