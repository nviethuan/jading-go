package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID   uint `gorm:"primaryKey"`
	Usn  string `gorm:"unique"`
	Name string
	Age  int
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Tự động tạo bảng (AutoMigrate)
	db.AutoMigrate(&User{})

	// Thêm dữ liệu (Create)
	db.Create(&User{Usn: "Alice", Name: "Alice", Age: 25})
	db.Create(&User{Usn: "Bob", Name: "Bob", Age: 30})

	users := []User{}
	db.Find(&users)

	user := users[0]
	fmt.Println(user.Age)
}
