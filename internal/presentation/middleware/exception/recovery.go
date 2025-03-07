package exception

import (
	"github.com/BargheNo/Backend/bootstrap"
	"github.com/gin-gonic/gin"
)

type RecoveryMiddleware struct {
	constants *bootstrap.Constants
}

func NewRecovery(constants *bootstrap.Constants) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		constants: constants,
	}
}

func (recovery RecoveryMiddleware) Recovery(c *gin.Context) {
	defer func() {
		if rec := recover(); rec != nil {
			if err, ok := rec.(error); ok {
				recovery.handleRecoveredError(c, err)
				c.Abort()
			}
		}
	}()

	c.Next()
}

func (recovery RecoveryMiddleware) handleRecoveredError(c *gin.Context, err error) {
	// handle recovered errors
}
