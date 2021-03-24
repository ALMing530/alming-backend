package service

import (
	"fmt"
	"testing"
	"time"
)

func Test_service(t *testing.T) {

	fmt.Println(time.Now())
	info := GetSysInfo()
	fmt.Println(time.Now())
	fmt.Println(info)
}
