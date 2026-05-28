package app

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	embedstatic "github.com/ijry/lyshop/internal/embedstatic"
)

func registerEmbeddedSites(r *gin.Engine) {
	if embedstatic.HasFile("web/index.html") {
		r.GET("/", serveEmbeddedFile("web/index.html", "text/html; charset=utf-8"))
		if webAssets, err := embedstatic.SubFileServer("web/assets"); err == nil {
			r.GET("/assets/*filepath", gin.WrapH(http.StripPrefix("/assets/", webAssets)))
		}
	}

	registerEmbeddedSPA(r, "/admin", "admin")
	registerEmbeddedSPA(r, "/h5", "h5")
}

func registerEmbeddedSPA(r *gin.Engine, prefix, root string) {
	indexAsset := root + "/index.html"
	if !embedstatic.HasFile(indexAsset) {
		return
	}

	fsHandler, err := embedstatic.SubFileServer(root)
	if err != nil {
		return
	}

	base := "/" + strings.Trim(strings.TrimSpace(prefix), "/")
	if base == "/" {
		base = ""
	}

	if base != "" {
		r.GET(base, func(c *gin.Context) {
			c.Redirect(http.StatusTemporaryRedirect, base+"/")
		})
	}
	r.GET(base+"/", serveEmbeddedFile(indexAsset, "text/html; charset=utf-8"))

	r.GET(base+"/*filepath", func(c *gin.Context) {
		filepath := strings.TrimPrefix(c.Param("filepath"), "/")
		if filepath == "" {
			serveEmbeddedFile(indexAsset, "text/html; charset=utf-8")(c)
			return
		}

		// API requests should keep normal not-found semantics.
		if base == "/admin" && (filepath == "api" || strings.HasPrefix(filepath, "api/")) {
			c.Status(http.StatusNotFound)
			return
		}

		if embedstatic.HasFile(root + "/" + filepath) {
			c.Request.URL.Path = "/" + filepath
			fsHandler.ServeHTTP(c.Writer, c.Request)
			return
		}

		// SPA fallback.
		serveEmbeddedFile(indexAsset, "text/html; charset=utf-8")(c)
	})
}

func serveEmbeddedFile(assetPath string, contentType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := embedstatic.ReadFile(assetPath)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		if contentType != "" {
			c.Header("Content-Type", contentType)
		}
		c.Writer.WriteHeader(http.StatusOK)
		_, _ = c.Writer.Write(body)
	}
}
