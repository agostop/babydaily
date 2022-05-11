package resources

import (
	"embed"
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	loggers "management_backend/src/logger"

	"github.com/gin-gonic/gin"
)

const resourcePath = "resources"
const staticPath = "static"

var (
	log = loggers.GetLogger(loggers.ModuleWeb)
)

func NewHtmlHandle(fs embed.FS) *Html {
	return &Html{
		fsResource: fs,
	}
}

type Html struct {
	fsResource embed.FS
}

func (h *Html) HtmlHandle(c *gin.Context) {
	c.Header("content-type", "text/html;charset=utf-8")
	indexContent, err := h.fsResource.ReadFile("resources/index.html")
	if err != nil {
		log.Error("get fs file failed. name: %s", "index.html")
		return
	}
	c.String(200, string(indexContent))
}

func (h *Html) Open(name string) (fs.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path.")
	}
	realPath := filepath.Join(resourcePath, staticPath, path.Clean(name))
	return h.fsResource.Open(realPath)
}
