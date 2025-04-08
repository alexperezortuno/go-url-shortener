package shortner

import (
	"fmt"
	"github.com/alexperezortuno/go-url-shortner/internal/config/environment"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/shortener"
	"github.com/alexperezortuno/go-url-shortner/internal/platform/storage/store"
	"github.com/gin-gonic/gin"
	"net/http"
)

var params = environment.Server()

type URLCreationRequest struct {
	LongURL string `json:"long_url" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
}

func CreateShortURL() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request URLCreationRequest

		if err := ctx.BindJSON(&request); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    "-1000",
				"message": "invalid request",
			})
			return
		}

		shortUrl := shortener.GenerateShortURL(request.LongURL, request.UserId)
		strShortUrl := fmt.Sprintf("%s://%s:%d/%s/%s",
			params.Protocol,
			params.Host,
			params.Port,
			params.Context,
			shortUrl)
		store.SaveURLInRedis(strShortUrl, request.LongURL)
		ctx.JSON(http.StatusOK, gin.H{
			"short_url": strShortUrl,
		})
	}
}

func ReturnLongURL() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		shortUrl := ctx.Request.URL.Query().Get("short_url")
		strShortUrl := fmt.Sprintf("%s://%s:%d/%s/%s",
			params.Protocol,
			params.Host,
			params.Port,
			params.Context,
			shortUrl)
		initialUrl := store.RetrieveInitialURLFromRedis(strShortUrl)
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"short_url": shortUrl,
			"long_url":  initialUrl,
		})
	}
}
