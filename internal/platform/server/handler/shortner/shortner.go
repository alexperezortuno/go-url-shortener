package shortner

import (
	"fmt"
	"github.com/alexperezortuno/go-url-shortner/internal/commons"
	"github.com/alexperezortuno/go-url-shortner/internal/config"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/errors"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/shortener"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/storage/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

type URLCreationRequest struct {
	LongURL string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

func CreateShortURL(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request URLCreationRequest

		if err := ctx.BindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, errors.NewCustomError(errors.BadRequest))
			return
		}

		shortUrl := shortener.GenerateShortURL(request.LongURL, request.UserId)
		store.SaveURLInRedis(shortUrl, request.LongURL)
		ctx.JSON(http.StatusOK, gin.H{
			"short_url": fmt.Sprintf("%s://%s:%d%s/%s/%s", cfg.Protocol,
				cfg.Host,
				cfg.Port,
				cfg.Context,
				commons.ShortenerPath,
				shortUrl),
		})
	}
}

func ReturnLongURL() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shortUrl := ctx.Request.URL.Query().Get("short_url")
		initialUrl := store.RetrieveInitialURLFromRedis(shortUrl)
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"short_url": shortUrl,
			"long_url":  initialUrl,
		})
	}
}

func RedirectURL() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shortUrl := ctx.Param("s")
		initialUrl := store.RetrieveInitialURLFromRedis(shortUrl)
		ctx.Redirect(http.StatusPermanentRedirect, initialUrl)
	}
}
