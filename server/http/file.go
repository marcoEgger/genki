package http

import (
	"net/http"
	"os"
	"path"

	"github.com/lukasjarosch/genki/logger"
)

type swaggerHandler struct {
	swaggerFile string
}

func NewSwaggerHandler(swaggerFile string) *swaggerHandler {
	return &swaggerHandler{
		swaggerFile:swaggerFile,
	}
}

// ServeHTTP will serve the swagger file. If the file does not exist, a warning
// is logged and a 404 returned.
func (swagger *swaggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(swagger.swaggerFile); os.IsNotExist(err) {
		logger.Warnf("cannot serve swagger file '%s': %s", swagger.swaggerFile, err)
		http.NotFound(w, r)
		return
	}

	logger.Infof("serving %s", r.URL.Path)
	wd, _ := os.Getwd()
	swaggerFile := path.Join(wd, swagger.swaggerFile)
	http.ServeFile(w, r, swaggerFile)
}
