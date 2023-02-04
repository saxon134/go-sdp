package models

const TBNResource = "resource"

type TblResource struct {
	Id      int64
	Type    int     `orm:"tinyint;comment:1-ECS" json:"type" form:"type"`
	Name    string  `orm:"varchar(40)" json:"name" form:"name"`
	Address string  `orm:"40" json:"address" form:"address"`
	Cpu     float32 `orm:"decimal(4,2)" json:"cpu" form:"cpu"`
	Memo    float32 `orm:"decimal(4,2)" json:"memo" form:"memo"`
	Status  int     `orm:"tinyint;comment:-1-已删除 0-未知 2-正常 10-异常" json:"status" form:"status"`
}

func (m *TblResource) TableName() string {
	return TBNResource
}
