package api

import (
	db "auth_wallet/db/sqlc"
	"auth_wallet/util"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createUserRequest struct {
	PublicAddress string `json:"public_address" binding:"required,eth_address"`
}

func (server *Server) createUserWithNonce(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req.PublicAddress = strings.ToLower(req.PublicAddress)
	getUserOnDB, err := server.store.GetUserBypublic_address(ctx, req.PublicAddress)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			newUser, err := server.store.CreateUser(ctx, db.CreateUserParams{
				PublicAddress: req.PublicAddress,
				Nonce:         util.RandomStringNonce(32),
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"user":    newUser,
				"message": "New user created successfully",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	updatedUser, err := server.store.UpdateUserNonce(ctx, db.UpdateUserNonceParams{
		PublicAddress: getUserOnDB.PublicAddress,
		Nonce:         util.RandomStringNonce(32),
	})

	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":    updatedUser,
		"message": "User nonce updated successfully",
	})
}

type loginVerifySignature struct {
	PublicAddress string `json:"public_address" binding:"required,eth_address"`
	Nonce         string `json:"nonce" binding:"required"`
	Signature     string `json:"signature" binding:"required"`
}

type loginUserResponse struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

type userResponse struct {
	PublicAddress string `json:"public_address"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		PublicAddress: user.PublicAddress,
	}
}

func (server *Server) loginByWallet(ctx *gin.Context) {
	var req loginVerifySignature
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	req.PublicAddress = strings.ToLower(req.PublicAddress)
	// Get user from database
	user, err := server.store.GetUserBypublic_address(ctx, req.PublicAddress)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// // Verify the signature
	message := []byte("Sign this message to authenticate: " + req.Nonce)
	recoveredAddress, err := util.VerifySignature(req.PublicAddress, req.Signature, message)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid signature")))
		return
	}
	recoveredAddress = strings.ToLower(recoveredAddress)
	fmt.Println("recoveredAddress", recoveredAddress)

	// Verify the recovered address matches the public address
	if recoveredAddress != req.PublicAddress {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("address mismatch")))
		return
	}

	// Generate JWT token
	// accessToken
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.PublicAddress,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.PublicAddress,
		server.config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// create Session
	session, err := server.store.CreateSessions(ctx, db.CreateSessionsParams{
		ID:            refreshPayload.ID,
		PublicAddress: refreshPayload.PublicAddress,
		RefreshToken:  refreshToken,
		UserAgent:     ctx.Request.UserAgent(),
		ClientIp:      ctx.ClientIP(),
		IsBlocked:     false,
		ExpiresAt:     refreshPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, rsp)
}
