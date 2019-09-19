package g

import (
	"os/exec"
	"strconv"
	"unsafe"
)

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToUint64(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}

func StringToFloat32(value string) (float32, error) {
	value_, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, err
	} else {
		return float32(value_), err
	}

}

func FilterStringSlice(list []string, filterStr string) []string {
	list_ := make([]string, 0)
	for _, ele := range list {
		if ele != filterStr {
			list_ = append(list_, ele)
		}
	}
	return list_
}

func DoCmdAndOutPut(cmd string) (string, error) {
	output, err := exec.Command("/bin/bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	return BytesToString(output), nil
}
