package usecase

import (
	"errors"
	"fmt"
	"main/domain"
)

type roomUseCase struct {
	roomRepository domain.RoomRepository
	userRepository domain.UserRepository
}

func NewRoomUseCase(
	roomRepository domain.RoomRepository,
	userRepository domain.UserRepository,
) *roomUseCase {
	return &roomUseCase{
		roomRepository,
		userRepository,
	}
}

func (r *roomUseCase) GetAll() (*[]domain.Room, error) {
	rooms, err := r.roomRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *roomUseCase) CreateRoom(input domain.CreateRoomInput) (*domain.Room, error) {
	userIDs, err := r.userRepository.IsExist(input.FriendIDs)
	if err != nil {
		return nil, err
	}

	if len(userIDs) == 0 {
		return nil, errors.New("friend's ids is empty after remove non-exist - UseCase CreateRoom")
	}

	room, err := r.roomRepository.CreateRoom(input.Name, userIDs)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return room, nil
}
