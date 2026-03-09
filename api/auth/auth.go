package auth

import "github.com/gin-gonic/gin"

func GetAuthHeader(client *gin.Context) string {
	authHeader := client.Request.Header.Get("Authorization")
	return authHeader
}
