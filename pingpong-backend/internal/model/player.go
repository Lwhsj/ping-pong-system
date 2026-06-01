package model

type Player struct {
	ID   uint64 `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name string `gorm:"column:name;size:20;not null" json:"name"`
	Sex  string `gorm:"column:sex;size:1;not null" json:"sex"`
	Age  int    `gorm:"column:age;not null" json:"age"`
}

func (Player) TableName() string {
	return "player"
}
