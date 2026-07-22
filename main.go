package main

import (
	"log"
	"myblog/repository"
)

func main() {
	db, err := repository.InitSQL()
	if err != nil {
		log.Fatalln("數據庫初始化失敗:", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)

	err = userRepo.Register("test_user", "123456")
	if err != nil {
		log.Fatalln("建立用戶數據失敗:", err)
	}
}
