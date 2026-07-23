package main

import (
	"fmt"
	"log"

	"github.com/Rashomon-code/myblog/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := repository.InitSQL()
	if err != nil {
		log.Fatalln("データベース初始化失敗:", err)
	}
	defer db.Close()

	// userRepo := repository.NewUserRepository(db)

	// err = userRepo.Register("test_user", "123456")
	// if err != nil {
	// 	log.Println("ユーザー登録失敗:", err)
	// }

	// success, err := userRepo.Login("alice", "my_secure_password")
	// if err != nil {
	// 	fmt.Printf("ログインできませんでした: %v\n", err)
	// } else {
	// 	fmt.Printf("ログイン結果: %v\n", success)
	// }

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	fmt.Println("サーバーポート :8080")
	r.Run(":8080")
}
