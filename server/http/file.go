package http

import (
	"net/http"
	"os"

	"github.com/lukasjarosch/genki/logger"
)

type swaggerHandler struct {
	file string
}

func NewFileServer(filePath string) *swaggerHandler {
	return &swaggerHandler{
		file: filePath,
	}
}

// ServeHTTP will serve the file. If the file does not exist, a warning
// is logged and a 404 returned.
func (swagger *swaggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := os.Stat(swagger.file); os.IsNotExist(err) {
		logger.Warnf("cannot serve file '%s': %s", swagger.file, err)
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, swagger.file)
}
