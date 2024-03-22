package main

type User struct {
	Id      int64  `gorm:"primarykey:column:id;type:bigint;not null"`
	Persona string `gorm:"column:persona;type:string;not null"`
}

func (User) TableName() string {
	return "user"
}
