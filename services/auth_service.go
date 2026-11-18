package services

import (
	"context"
	"errors"
	"tiket-event/exception"
	"tiket-event/helpers"
	"tiket-event/models/domain"
	"tiket-event/models/web"
	"tiket-event/repository"

	"github.com/jackc/pgx/v5"
)

type AuthService interface {
	Register(ctx context.Context, request web.RegisterRequest) (web.UserResponse, error)
	Login(ctx context.Context, request web.LoginRequest) (web.LoginResponse, error)
}

type authService struct {
	repository repository.UserRepository
	jwtSecret  string
}

func NewService(repository repository.UserRepository, jwt string) AuthService {
	return &authService{
		repository: repository,
		jwtSecret:  jwt,
	}
}

func (service *authService) Register(ctx context.Context, request web.RegisterRequest) (web.UserResponse, error) {
	//cek apakah email sudah ada didatabase
	_, err := service.repository.FindByEmail(ctx, request.Email)

	// jika err sama dengan nil, artinya ada, berarti email sudah dipakai
	if err == nil {
		// kembalikan user response kosong dan email sudah terdaftar
		return web.UserResponse{}, exception.ErrEmailTaken
	}

	// jika errornya "data tidak ditemukan", berarti databasenya mati
	if !errors.Is(err, pgx.ErrNoRows) {
		return web.UserResponse{}, err
	}

	//hash password
	hashPassword, err := helpers.HashPassword(request.Password)
	//gagal hash password
	if err != nil {
		return web.UserResponse{}, nil
	}

	// siapkan data user yang akan disimpan
	newUser := domain.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hashPassword,
	}

	// simpan ke database lewat repository
	createUser, err := service.repository.Create(ctx, newUser)
	if err != nil {
		return web.UserResponse{}, err
	}

	// siapkan data responsenya yang akan dikirim ke client
	userResponse := web.UserResponse{
		Id:    createUser.Id,
		Name:  createUser.Name,
		Email: createUser.Email,
		Role:  createUser.Role,
	}

	//register berhasil, error di isi nil (artinya tidak error)
	return userResponse, nil
}

func (service *authService) Login(ctx context.Context, request web.LoginRequest) (web.LoginResponse, error) {
	// cari user berdasarkan email yang diketik
	user, err := service.repository.FindByEmail(ctx, request.Email)
	if err != nil {
		//email tidak ada didatabase
		if errors.Is(err, pgx.ErrNoRows) {
			// kembalikan login response kosong dan pesan email/password salah
			return web.LoginResponse{}, exception.ErrInvalidCredential
		}
		return web.LoginResponse{}, err
	}

	// cocokan hash password didatabase dengan password yang diketik
	isMacth := helpers.CheckPassword(user.Password, request.Password)
	// jika tidak sama
	if !isMacth {
		// password salah
		return web.LoginResponse{}, exception.ErrInvalidCredential
	}

	//password benar, bikin token jwt (id dan role disimpan didalamnya)
	token, err := helpers.GenerateToken(service.jwtSecret, user.Id, user.Role)
	if err != nil {
		return web.LoginResponse{}, err
	}

	// siapkan data user yang aman dikirim tanpa password ke client
	userResponse := web.UserResponse{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	// siapkan response akhir: token + data user
	loginResponse := web.LoginResponse{
		Token: token,
		User:  userResponse,
	}

	// login berhasil, error di isi nil (artinya tidak error)
	return loginResponse, nil

}
