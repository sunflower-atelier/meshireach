package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"

	"github.com/gin-gonic/gin"
)

func FirebaseAuth(app *firebase.App) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth, err := app.Auth(context.Background())
		if err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}

		// ヘッダからAuth情報をもらう
		idToken := strings.Replace(ctx.GetHeader("Authorization"), "Bearer ", "", 1)

		/*token*/
		token, err := auth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			ctx.String(http.StatusForbidden, "ごめんなのび太．このページ，会員用なんだ．")
			ctx.Abort()
			return
		}

		// fmt.Printf("Verified ID token: %v\n", token)

		ctx.Set("FirebaseID", token.UID)

		ctx.Next()
	}
}
