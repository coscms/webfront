package top

import (
	"fmt"

	"github.com/webx-top/com"
)

func CN2Number(cn string) (uint64, error) {
	i := com.ConvertNumberChToAr(cn)
	if i < 0 {
		return 0, fmt.Errorf(`invalid chinese number: %s`, cn)
	}
	return uint64(i), nil
}
