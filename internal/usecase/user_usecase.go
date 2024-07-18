package usecase

// import (
// 	"context"
// 	"manga/domain"
// 	"manga/domain/dtos"
// 	"time"
// )

// type userUsecase struct {
// 	userRepository domain.UserRepository
// 	contextTimeout time.Duration
// }

// func NewUserUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUsecase {
// 	return &userUsecase{
// 		userRepository: userRepository,
// 		contextTimeout: timeout,
// 	}
// }
// func (uu *userUsecase) GetUserByName(c context.Context, name string, page int64, pageSize int64, rsid string) ([]dtos.UserResponse, dtos.Pagination, error) {
// 	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
// 	defer cancel()

// 	users, totalPages, err := uu.userRepository.FetchByName(ctx, name, "name", 1, page, pageSize, rsid)

// 	var newuser []domain.UserResponse

// 	for _, user := range users {
// 		newuser = append(newuser, domain.UserResponse{
// 			ID: user.ID, Name: user.Name, Email: user.Email, Role: user.Role, Place: user.Place, PlaceName: user.PlaceName, Room: user.Room, RoomName: user.RoomName})
// 	}

// 	if len(newuser) == 0 {

// 		return []domain.UserResponse{}, domain.Pagination{Page: page, Limit: pageSize, Total: totalPages}, nil
// 	}
// 	return newuser, domain.Pagination{Page: page, Limit: pageSize, Total: totalPages}, err
// }

// func (uu *userUsecase) GetUserInfoByID(c context.Context, id string) (domain.UserResponse, error) {
// 	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
// 	defer cancel()
// 	return uu.userRepository.GetUserInfoByID(ctx, id)
// }

// func (uu *userUsecase) DeleteByID(c context.Context, id string) error {
// 	ctx, cancel := context.WithTimeout(c, uu.contextTimeout)
// 	defer cancel()
// 	return uu.userRepository.DeleteByID(ctx, id)
// }
