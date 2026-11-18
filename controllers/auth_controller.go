package controllers

import (
	"errors"
	"log"
	"tiket-event/exception"
	"tiket-event/helpers"
	"tiket-event/models/web"
	"tiket-event/services"

	"github.com/gofiber/fiber/v3"
)

type AuthController interface {
	Register(c fiber.Ctx) error
	Login(c fiber.Ctx) error
}

type authController struct {
	service services.AuthService
}

func NewController(service services.AuthService) AuthController {
	return &authController{
		service: service,
	}
}

func (controller *authController) Register(c fiber.Ctx) error {
	// siapkan kotak kosong untuk menampung isi body request
	var request web.RegisterRequest

	// baca body JSON dari client, lalu isi ke kotak request
	err := c.Bind().Body(&request)
	if err != nil {
		//body bukan JSON yang benar
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "format body tidak valid",
		})
	}

	// cek isian sesuai aturan validate struct
	validationErrors := helpers.ValidateStruct(request)
	if validationErrors != nil {
		// ada isian yanng salah kirim daftar kesalahan
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseApiFailed{
			Error:   "validasi error",
			Details: validationErrors,
		})
	}

	userResponse, err := controller.service.Register(c.Context(), request)
	if err != nil {
		// jika error yang ditamaplikan service adalah exception.ErrEmailTaken, maka email sudah digunakan
		if errors.Is(err, exception.ErrEmailTaken) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// error lain : catat diterminal, jangan bocorkan ke client
		log.Println("error registerasi", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "terjadi kesalahan pada server",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(helpers.ResponseApiSuccess{
		Code:    201,
		Messege: "Berhasil Register",
		Data:    userResponse,
	})
}

func (controller *authController) Login(c fiber.Ctx) error {
	var request web.LoginRequest

	//baca body json
	err := c.Bind().JSON(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "format body tidak valid",
		})
	}

	//cek isian
	validationErrors := helpers.ValidateStruct(request)
	if validationErrors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseApiFailed{
			Error:   "validasi gagal",
			Details: validationErrors,
		})
	}

	loginResponse, err := controller.service.Login(c.Context(), request)
	if err != nil {
		// email atau password salah
		if errors.Is(err, exception.ErrInvalidCredential) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// error lain
		log.Println("error login:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "terjadi kesalahan pada server",
		})
	}

	// berhasil, kirim token dan data user
	return c.Status(fiber.StatusOK).JSON(helpers.ResponseApiSuccess{
		Code:    200,
		Messege: "Berhasil Login",
		Data:    loginResponse,
	})

}
