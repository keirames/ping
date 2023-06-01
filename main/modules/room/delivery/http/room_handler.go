package http

import (
	"fmt"
	"main/domain"
	"main/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type roomHandler struct {
	roomUseCase domain.RoomUseCase
}

func NewRoomHandler(routerGroup *gin.RouterGroup, roomUseCase domain.RoomUseCase) {
	handler := &roomHandler{
		roomUseCase,
	}

	routerGroup.GET("/rooms", handler.GetRooms)
	routerGroup.GET("/room/:id", handler.GetRoom)
	routerGroup.POST("/create-room", handler.CreateRoom)
	routerGroup.POST("/send-message", handler.CreateRoom)
	routerGroup.POST("/add-member", handler.CreateRoom)
	routerGroup.POST("/delete-room", handler.CreateRoom)
}

func (rh *roomHandler) GetRooms(c *gin.Context) {

	rooms, err := rh.roomUseCase.GetAll()
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, rooms)
}

func (rh *roomHandler) GetRoom(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": "200"})
}

type Body struct {
	Name string `json:"name" binding:"required"`
}

func validate(c *gin.Context, input interface{}) error {
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return err
	}

	return nil
}

func (rh *roomHandler) CreateRoom(c *gin.Context) {
	createRoomInput := domain.CreateRoomInput{}

	if err := validate(c, &createRoomInput); err != nil {
		return
	}

	if len(createRoomInput.FriendIDs) == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Custom validate on friend's ids
	claims := jwt.GetClaims(c)
	if claims == nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// Add owner id
	var friendIDs []int64
	friendIDs = append(friendIDs, claims.UserID)
	friendIDs = append(friendIDs, createRoomInput.FriendIDs...)

	// Remove any duplicate
	var deduplicateFriendIDs []int64
	dupMap := make(map[int64]bool)
	for _, id := range friendIDs {
		if ok := dupMap[id]; !ok {
			dupMap[id] = true
			deduplicateFriendIDs = append(deduplicateFriendIDs, id)
		}
	}

	room, err := rh.roomUseCase.CreateRoom(
		domain.CreateRoomInput{
			Name:      createRoomInput.Name,
			FriendIDs: deduplicateFriendIDs,
		},
	)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, room)
}

func (rh *roomHandler) SendMessage(c *gin.Context) {
	body := Body{}

	if err := validate(c, &body); err != nil {
		return
	}

	c.JSON(200, gin.H{"ok": ""})
}

func (rh *roomHandler) AddMember(c *gin.Context) {
	body := Body{}

	if err := validate(c, &body); err != nil {
		return
	}

	c.JSON(200, gin.H{"ok": ""})
}

func (rh *roomHandler) DeleteRoom(c *gin.Context) {
	body := Body{}

	if err := validate(c, &body); err != nil {
		return
	}

	c.JSON(200, gin.H{"ok": ""})
}
