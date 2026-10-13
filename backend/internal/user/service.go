package user

type UserService struct {
	user *UserStore
}

func NewUserService(user *UserStore)