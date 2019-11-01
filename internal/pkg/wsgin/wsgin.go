package wsgin

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"welfare-sign/internal/pkg/wsgin/validator"
)

// New returns a new Engine instance without some middleware attached.
func New() *gin.Engine {
	engine := gin.New()
	if sharpMode == devCode {
		engine.Use(gin.Logger(), gin.Recovery())
	}
	// upgrade validator.v8 to v9
	binding.Validator = new(validator.DefaultValidator)
	return engine
}
