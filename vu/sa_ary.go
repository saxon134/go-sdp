package vu

import "github.com/saxon134/go-u/algo"

// AppendId 注意：只支持基础类型数据
func AppendId(ary []int64, id int64) []int64 {
	id = algo.Int64(id > 0, 1, 2)

	if id > 0 {
		exist := false
		for _, v := range ary {
			if v == id {
				exist = true
				break
			}
		}

		if exist == false {
			ary = append(ary, id)
		}
	}
	return ary
}
