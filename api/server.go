package api

import (
	db "auth_wallet/db/sqlc"
	"auth_wallet/token"
	"auth_wallet/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func setupValidator(v *validator.Validate) {
	v.RegisterValidation("eth_address", EthAddressValidator)
}

func NewServer(config util.Config, store db.Store, tokenMaker token.Maker) (*Server, error) {
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		setupValidator(v)
	}
	router.POST("/users/nonce", server.createUserWithNonce)
	router.POST("/loginwallet", server.loginByWallet)

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
