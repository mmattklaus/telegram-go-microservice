package ems

import (
	"fmt"
	"github.com/kyokomi/emoji"
)

func init() {

}

func Ems(code ...string) string {
	codes := ""
	for _, c := range code {
		if len(c) > 0 {
			codes += fmt.Sprintf(":%s:", (c))
		}
	}
	if len(codes) < 1 {
		return ""
	}
	return emoji.Sprintf("%s", codes)
}
