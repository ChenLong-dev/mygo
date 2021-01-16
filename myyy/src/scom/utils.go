package scom

import (
	"errors"
	"fmt"
	"github.com/ChenLong-dev/gobase/mbase/mutils"
	"net"
	"reflect"
)

func Signature(region string, rand string, timestamp uint64, password string) string {
	//mutils.SHA1Hex(loginReq.Region + loginReq.Rand + fmt.Sprint(loginReq.Timestamp) + password)
	return mutils.SHA1Hex(region + rand + fmt.Sprint(timestamp) + password)
}

func Duplicate(a interface{}) (ret []interface{}) {
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}
	return ret
}

func RemoveRepeatedElement(arr []*string) (newArr []*string) {
	newArr = make([]*string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return newArr
}

func GetLocalIPv4() (string, error) {
	var ips []string
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addrs, _ := netInterfaces[i].Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					ips = append(ips, ipnet.IP.String())
				}
			}
		}
	}
	if 0 == len(ips) {
		return "", errors.New("No available ip ")
	}
	return ips[0], nil
}
