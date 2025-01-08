package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"yalp_ulab/config"
	"yalp_ulab/internal/entity"
	"yalp_ulab/pkg/etc"
	"yalp_ulab/pkg/hash"
	"yalp_ulab/pkg/jwt"
)

// Login godoc
// @Router /auth/login [post]
// @Summary Login
// @Description Login
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body entity.LoginRequest true "User"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) Login(ctx *gin.Context) {
	var body entity.LoginRequest

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.UseCase.UserRepo.GetSingle(ctx, entity.UserSingleRequest{
		Email: body.Email,
	})
	if h.HandleDbError(ctx, err, "Error getting user") {
		return
	}

	if user.UserType == "user" && body.Platform == "admin" {
		h.ReturnError(ctx, config.ErrorForbidden, "User can't login to admin web", http.StatusBadRequest)
		return
	} else if user.UserType == "admin" && body.Platform != "admin" {
		h.ReturnError(ctx, config.ErrorForbidden, "Admin can only login to admin web", http.StatusBadRequest)
		return
	}

	if !hash.CheckPasswordHash(body.Password, user.Password) {
		h.ReturnError(ctx, config.ErrorInvalidPass, "Incorrect password", http.StatusBadRequest)
		return
	}

	// Create session
	newSession := entity.Session{
		UserID:       user.ID,
		IPAddress:    ctx.ClientIP(),
		ExpiresAt:    time.Now().Add(time.Hour * 999999).Format(time.RFC3339),
		UserAgent:    ctx.Request.UserAgent(),
		IsActive:     true,
		LastActiveAt: time.Now().Format(time.RFC3339),
		Platform:     body.Platform,
	}

	session, err := h.UseCase.SessionRepo.Create(ctx, newSession)
	if h.HandleDbError(ctx, err, "Error while creating new session") {
		return
	}

	// Generate JWT token
	jwtFields := map[string]interface{}{
		"sub":        user.ID,
		"user_role":  user.UserRole,
		"user_type":  user.UserType,
		"platform":   body.Platform,
		"session_id": session.ID,
	}

	user.AccessToken, err = jwt.GenerateJWT(jwtFields, h.Config.JWT.Secret)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Oops, something went wrong!!!", http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":    user,
		"session": session,
	})
}

// Logout godoc
// @Router /auth/logout [post]
// @Summary Logout
// @Description Logout
// @Security BearerAuth
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) Logout(ctx *gin.Context) {
	sessionID := ctx.MustGet("session_id")
	if sessionID == "" {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid session ID", http.StatusBadRequest)
		return
	}

	err := h.UseCase.SessionRepo.Delete(ctx, entity.Id{ID: sessionID.(string)})
	if h.HandleDbError(ctx, err, "Error deleting session") {
		return
	}

	ctx.JSON(http.StatusOK, entity.SuccessResponse{
		Message: "Successfully logged out",
	})
}

// Register godoc
// @Router /auth/register [post]
// @Summary Register
// @Description Register
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body entity.RegisterRequest true "User"
// @Success 200 {object} entity.User
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) Register(ctx *gin.Context) {
	var body entity.RegisterRequest

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.UseCase.UserRepo.GetSingle(ctx, entity.UserSingleRequest{
		Email: body.Email,
	})
	if err == nil {
		h.ReturnError(ctx, config.ErrorConflict, "User already exists", http.StatusBadRequest)
		return
	}

	body.Password, err = hash.HashPassword(body.Password)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Oops, something went wrong!!!", http.StatusInternalServerError)
		return
	}

	user, err = h.UseCase.UserRepo.Create(ctx, entity.User{
		FullName: body.FullName,
		UserType: "user",
		UserRole: "user",
		Email:    body.Email,
		Status:   "inverify",
		Password: body.Password,
	})
	if h.HandleDbError(ctx, err, "Error creating user") {
		return
	}

	// Send verification code to user
	otp := etc.GenerateOTP(6)
	err = h.Redis.Set(ctx, fmt.Sprintf("otp-%s", user.Email), otp, 5*60)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error setting OTP", http.StatusInternalServerError)
		return
	}

	// Send OTP code to user's email
	emailBody, err := etc.GenerateOtpEmailBody(otp)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error generating OTP email body", http.StatusInternalServerError)
		return
	}

	err = etc.SendEmail(h.Config.Gmail.Host, h.Config.Gmail.Port, h.Config.Gmail.Email, h.Config.Gmail.EmailPass, body.Email, emailBody)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error sending OTP", http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusCreated, entity.SuccessResponse{
		Message: "User registered successfully, please verify your email address",
	})
}

// VerifyEmail (continued)
// @Router /auth/verify-email [post]
// @Summary Verify Email
// @Description Verify Email
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body entity.VerifyEmail true "User"
// @Success 200 {object} entity.User
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) VerifyEmail(ctx *gin.Context) {
	var body entity.VerifyEmail

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("otp-%s", body.Email)

	otp, err := h.Redis.Get(ctx, key)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Oops, something went wrong", http.StatusInternalServerError)
		return
	}

	if otp != body.Otp {
		h.ReturnError(ctx, config.ErrorBadRequest, "Incorrect OTP", http.StatusBadRequest)
		return
	}

	user, err := h.UseCase.UserRepo.GetSingle(ctx, entity.UserSingleRequest{
		Email: body.Email,
	})
	if h.HandleDbError(ctx, err, "Error getting user") {
		return
	}

	user.Status = "active"

	_, err = h.UseCase.UserRepo.Update(ctx, user)
	if h.HandleDbError(ctx, err, "Error updating user") {
		return
	}

	// Create session
	newSession := entity.Session{
		UserID:       user.ID,
		IPAddress:    ctx.ClientIP(),
		ExpiresAt:    time.Now().Add(time.Hour * 999999).Format(time.RFC3339),
		UserAgent:    ctx.Request.UserAgent(),
		IsActive:     true,
		LastActiveAt: time.Now().Format(time.RFC3339),
		Platform:     body.Platform,
	}

	session, err := h.UseCase.SessionRepo.Create(ctx, newSession)
	if h.HandleDbError(ctx, err, "Error while creating new session") {
		return
	}

	// Generate JWT token
	jwtFields := map[string]interface{}{
		"sub":        user.ID,
		"user_role":  user.UserRole,
		"user_type":  user.UserType,
		"platform":   body.Platform,
		"session_id": session.ID,
	}

	user.AccessToken, err = jwt.GenerateJWT(jwtFields, h.Config.JWT.Secret)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Oops, something went wrong!!!", http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":    user,
		"session": session,
	})
}
