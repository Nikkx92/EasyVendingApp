//go:build js && wasm

// storage_wasm.go
package storage

import (
	"errors"
	"fmt"
	"syscall/js"
)

func SetItem(value string) error {
	res := js.Global().Call("_lsSet", "data", value)
	if !res.IsNull() && !res.IsUndefined() && res.Truthy() {
		return fmt.Errorf(res.String())
	}
	return nil
}

func GetItem() (string, error) {
	res := js.Global().Call("_lsGet", "data")
	if !res.Get("ok").Bool() {
		return "", errors.New(res.Get("err").String())
	}
	v := res.Get("val")
	if v.IsNull() || v.IsUndefined() {
		return "", nil
	}
	return v.String(), nil
}
