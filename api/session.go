package api

// type sessionResponse struct {
// 	ID        int64     `json:"id"`
// 	UserID    int64     `json:"user_id"`
// 	Token     string    `json:"token"`
// 	ExpiresAt time.Time `json:"expires_at"`
// 	CreatedAt time.Time `json:"created_at"`
// }

// // Middleware to check if token is valid and not expired
// func (server *Server) authMiddleware() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authHeader := ctx.GetHeader("Authorization")
// 		if authHeader == "" {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("authorization header is required")))
// 			return
// 		}

// 		// Extract token from Authorization header
// 		tokenString := authHeader
// 		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
// 			tokenString = authHeader[7:]
// 		}

// 		// Verify token is valid and not expired
// 		payload, err := server.tokenMaker.VerifyToken(tokenString)
// 		if err != nil {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
// 			return
// 		}

// 		// Check if token exists in sessions table
// 		session, err := server.store.GetSessionByToken(ctx, tokenString)
// 		if err != nil {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid session")))
// 			return
// 		}

// 		// Check if session is expired
// 		if time.Now().After(session.ExpiresAt) {
// 			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("session expired")))
// 			return
// 		}

// 		// Set user information in context
// 		ctx.Set("user_id", session.UserID)
// 		ctx.Set("public_address", payload.PublicAddress)
// 		ctx.Next()
// 	}
// }

// // Create new session and auth log entry
// func (server *Server) createSession(ctx *gin.Context, userID int64, token string, duration time.Duration) (*db.Session, error) {
// 	// Create session
// 	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
// 		UserID:    userID,
// 		Token:     token,
// 		ExpiresAt: time.Now().Add(duration),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Create auth log entry
// 	_, err = server.store.CreateAuthLog(ctx, db.CreateAuthLogParams{
// 		UserID:    userID,
// 		EventType: "login",
// 		IpAddress: ctx.ClientIP(),
// 	})
// 	if err != nil {
// 		// Log the error but don't fail the login
// 		server.logger.Printf("Failed to create auth log: %v", err)
// 	}

// 	return &session, nil
// }

// // Get user's active sessions
// func (server *Server) getUserSessions(ctx *gin.Context) {
// 	userID, exists := ctx.Get("user_id")
// 	if !exists {
// 		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not authenticated")))
// 		return
// 	}

// 	sessions, err := server.store.GetSessionsByUserID(ctx, userID.(int64))
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, sessions)
// }

// // Logout user (invalidate session)
// func (server *Server) logoutUser(ctx *gin.Context) {
// 	userID, exists := ctx.Get("user_id")
// 	if !exists {
// 		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not authenticated")))
// 		return
// 	}

// 	authHeader := ctx.GetHeader("Authorization")
// 	if authHeader == "" {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("authorization header is required")))
// 		return
// 	}

// 	// Extract token from Authorization header
// 	tokenString := authHeader
// 	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
// 		tokenString = authHeader[7:]
// 	}

// 	// Delete session from database
// 	err := server.store.DeleteSessionByToken(ctx, tokenString)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	// Create auth log entry
// 	_, err = server.store.CreateAuthLog(ctx, db.CreateAuthLogParams{
// 		UserID:    userID.(int64),
// 		EventType: "logout",
// 		IpAddress: ctx.ClientIP(),
// 	})
// 	if err != nil {
// 		// Log the error but don't fail the logout
// 		server.logger.Printf("Failed to create auth log: %v", err)
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
// }

// // Get user's auth logs
// func (server *Server) getUserAuthLogs(ctx *gin.Context) {
// 	userID, exists := ctx.Get("user_id")
// 	if !exists {
// 		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("user not authenticated")))
// 		return
// 	}

// 	logs, err := server.store.GetAuthLogsByUserID(ctx, userID.(int64))
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, logs)
// }

// // Updated loginByWallet to use sessions
// func (server *Server) loginByWallet(ctx *gin.Context) {
// 	var req loginVerifySignature
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, errorResponse(err))
// 		return
// 	}

// 	// Get user from database
// 	user, err := server.store.GetUserByPublicAddress(ctx, req.PublicAddress)
// 	if err != nil {
// 		if errors.Is(err, db.ErrRecordNotFound) {
// 			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user not found")))
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	// Verify the signature
// 	message := []byte("Sign this message to authenticate: " + user.Nonce)
// 	recoveredAddress, err := util.VerifySignature(req.PublicAddress, req.Signature, message)
// 	if err != nil {
// 		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid signature")))
// 		return
// 	}

// 	// Verify the recovered address matches the public address
// 	if recoveredAddress != req.PublicAddress {
// 		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("address mismatch")))
// 		return
// 	}

// 	// Generate JWT token
// 	token, payload, err := server.tokenMaker.CreateToken(user.PublicAddress, server.config.AccessTokenDuration)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	// Create session and auth log
// 	session, err := server.createSession(ctx, user.ID, token, server.config.AccessTokenDuration)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	// Update user's last login time and nonce
// 	err = server.store.UpdateUserAfterLogin(ctx, db.UpdateUserAfterLoginParams{
// 		ID:        user.ID,
// 		Nonce:     util.RandomStringNonce(32),
// 		LastLogin: time.Now(),
// 	})
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"access_token": token,
// 		"token_type":   "bearer",
// 		"user":         payload,
// 		"session":      session,
// 	})
// }
