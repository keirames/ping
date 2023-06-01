package main

import (
	"fmt"

	"main/database"
	"main/middleware"
	"main/modules/room/delivery/http"
	roomRepository "main/modules/room/repository"
	useCase "main/modules/room/use_case"
	userRepository "main/modules/user/repository"
	"main/tools"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	err := tools.LoadConfig(".")
	if err != nil {
		fmt.Println(err)
		panic("Fail to load config")
	}

	err = database.Connect()
	if err != nil {
		fmt.Println(err)
		panic("Fail to load database")
	}

	err = tools.NewNode()
	if err != nil {
		fmt.Println(err)
		panic("Fail to initialize new node snowflake")
	}

	router := gin.Default()

	apiV1 := router.Group("/v1", middleware.AuthRequired())

	apiV1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"1": "!"})
	})

	roomRepository := roomRepository.NewRoomRepository(database.Conn)
	userRepository := userRepository.NewUserRepository(database.Conn)
	roomUseCase := useCase.NewRoomUseCase(roomRepository, userRepository)
	http.NewRoomHandler(apiV1, roomUseCase)

	err = router.Run(tools.Config.ServerAddress)
	if err != nil {
		panic("Fail to start server")
	}
}
