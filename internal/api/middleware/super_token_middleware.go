package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// func JwtAuthMiddleware(secret string) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		authHeader := c.Request.Header.Get("Authorization")
// 		t := strings.Split(authHeader, " ")
// 		if len(t) == 2 {
// 			authToken := t[1]
// 			authorized, err := tokenutil.IsAuthorized(authToken, secret)
// 			if authorized {
// 				userID, err := tokenutil.ExtractIDFromToken(authToken, secret)
// 				if err != nil {
// 					c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
// 					c.Abort()
// 					return
// 				}
// 				c.Set("x-user-id", userID)
// 				c.Next()
// 				return
// 			}
// 			c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: err.Error()})
// 			c.Abort()
// 			return
// 		}
// 		c.JSON(http.StatusUnauthorized, domain.ErrorResponse{Message: "Not authorized"})
// 		c.Abort()
// 	}
// }

func SupertokensMiddleware(c *gin.Context) {
	supertokens.Middleware(http.HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			c.Next()
		})).ServeHTTP(c.Writer, c.Request)
	// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
	c.Abort()
}

// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {

// 		c.Header("Access-Control-Allow-Origin", "*")
// 		c.Header("Access-Control-Allow-Credentials", "true")
// 		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
// 		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }
