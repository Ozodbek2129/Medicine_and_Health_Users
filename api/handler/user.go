package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"user_medic/api/token"
	pb "user_medic/genproto/user"
	"user_medic/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Register a new user
// @Description Register a new user with username and password
// @Tags User
// @Accept json
// @Produce json
// @Param user body user.RegisterUserRequest true "User data"
// @Success 202 {object} user.RegisterUserResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /user/register [post]
func (h Handler) Register(c *gin.Context) {
	req := pb.RegisterUserRequest{}

	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		h.Log.Error(fmt.Sprintf("bodydan malumotlarni olishda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Pasworni hashlashda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	req.Password = string(hashpassword)

	resp, err := h.AuthUser.RegisterUser(context.Background(), &req)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Foydalanuvchi malumotlarni bazga yuborishda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, resp)
}

// @Summary Login user
// @Description Login user with username and password
// @Tags User
// @Accept json
// @Produce json
// @Param user body user.LoginUserRequest true "User data"
// @Success 202 {object} string
// @Failure 400 {object} string
// @Failure 401 {object} string
// @Failure 500 {object} string
// @Router /user/login [post]
func (h Handler) LoginUser(c *gin.Context) {
	req := pb.LoginUserRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.Log.Error(fmt.Sprintf("malumotlarni olishda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	user, err := h.AuthUser.GetByUserEmail(c, &req)
	if err != nil {
		h.Log.Error(fmt.Sprintf("GetbyUserda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		h.Log.Error(fmt.Sprintf("passwordni tekshirishda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := token.GenerateJWT(&model.LoginResponse{
		Id:         user.Id,
		First_name: user.FirstName,
		Gender:     user.Gender,
		Last_name:  user.LastName,
		Email:      user.Email,
		Role:       user.Role,
	})

	_, err = h.AuthUser.StoreRefreshToken(context.Background(), &pb.StoreRefreshTokenReq{
		UserId:       user.Id,
		RefreshToken: token.RefreshToken,
	})

	if err != nil {
		h.Log.Error(fmt.Sprintf("storefreshtokenda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusAccepted, token)
}

// @Summary Refreshe Token
// @Description This endpoint refreshes the signing key and returns a confirmation message.
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200 {object} user.RefreshTokenResponse
// @Failure 500 {object} string
// @Router /user/refresh-token [post]
func (h Handler) RefreshToken(c *gin.Context) {
	req := pb.RefreshTokenRequest{}

	resp, err := h.AuthUser.RefreshToken(c, &req)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Api da malumotlarni olishda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Api da malumotlarni olishda xatolik: " + err.Error(),
		})
		return
	}

	c.JSON(200, resp)
}

// @Summary Get user profile
// @Description Foydalanuvchi profilini olish
// @Tags Users
// @Accept  json
// @Produce  json
// @Param email path string false "Foydalanuvchi email manzili"
// @Success 202 {object} user.GetUserProfileResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /user/profile/{email} [get]
func (h *Handler) GetUserProfile(c *gin.Context) {
	req := pb.GetUserProfileRequest{}

	req.Email = c.Param("email")
	fmt.Println(req.Email)

	resp, err := h.AuthUser.GetUserProfile(context.Background(), &req)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Foydalanuvchi malumotlarni bazadan olishda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, resp)
}

// @Summary Foydalanuvchi profilini yangilash
// @Description Foydalanuvchi profilini yangilash funksiyasi, ma'lumotlar parametrlar orqali olinadi va yangilanadi. Agar parametr bo'sh bo'lsa, eski qiymat saqlanib qoladi.
// @Tags Foydalanuvchi
// @Accept json
// @Produce json
// @Param id path string true "Foydalanuvchi IDsi"
// @Param email path string false "Foydalanuvchi email manzili"
// @Param password path string false "Foydalanuvchi paroli"
// @Param first_name path string false "Foydalanuvchining ismi"
// @Param last_name path string false "Foydalanuvchining familiyasi"
// @Param date_of_birthday path string false "Foydalanuvchining tug'ilgan sanasi (YYYY-MM-DD formatida)"
// @Param gender path string false "Foydalanuvchining jinsi (male/female)"
// @Param role path string false "Foydalanuvchining roli (patient/doctor/admin)"
// @Success 200 {object} string "Muvaffaqiyatli yangilangan profil ma'lumotlari"
// @Failure 400 {object} string "Yaroqsiz so'rov"
// @Failure 500 {object} string "Server xatosi"
// @Router /user/profile/update/{id}/{email}/{password}/{first_name}/{last_name}/{date_of_birthday}/{gender}/{role} [put]
func (h *Handler) UpdateUserProfile(c *gin.Context) {
	req := pb.UpdateUserProfileRequest{}

	req.Id = c.Param("id")
	req.Email = c.Param("email")
	req.Password = c.Param("password")
	req.FirstName = c.Param("first_name")
	req.LastName = c.Param("last_name")
	req.DateOfBirthday = c.Param("date_of_birthday")
	req.Gender = c.Param("gender")
	req.Role = c.Param("role")

	r := pb.LoginUserRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := h.AuthUser.GetByUserEmail(c, &r)
	if err != nil {
		h.Log.Error(fmt.Sprintf("GetByUserda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Password != "" {
		hashpassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			h.Log.Error(fmt.Sprintf("Parolni hashlashda xatolik: %v", err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		req.Password = string(hashpassword)
		user.Password = req.Password
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}

	if req.LastName != "" {
		user.LastName = req.LastName
	}

	if req.DateOfBirthday != "" {
		user.DateOfBirthday = req.DateOfBirthday
	}

	if req.Gender != "" {
		user.Gender = req.Gender
	}

	if req.Role != "" {
		user.Role = req.Role
	}

	req = pb.UpdateUserProfileRequest{
		Id:             req.Id,
		Email:          user.Email,
		Password:       req.Password,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		DateOfBirthday: user.DateOfBirthday,
		Gender:         user.Gender,
		Role:           user.Role,
	}

	resp, err := h.AuthUser.UpdateUserProfile(c, &req)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Profilni update qilishda xatolik: %v", err))
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, resp)
}

// @Summary Logout user
// @Description Foydalanuvchini tizimdan chiqarish
// @Tags Users
// @Accept  json
// @Produce  json
// @Param data body user.LogoutUserRequest true "LogoutUserRequest"
// @Success 200 {object} user.LogoutUserResponse
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /user/logout [put]
func (h *Handler) LogoutUser(c *gin.Context) {
	req := pb.LogoutUserRequest{}

	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		h.Log.Error(fmt.Sprintf("malumotlarni olishda xatolik: %v", err))
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	resp, err := h.AuthUser.LogoutUser(c, &req)
	if err != nil {
		h.Log.Error(fmt.Sprintf("Tokenni logout qilishda xatolik: %v", err))
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, resp)
}
