package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/context"

	firebase "firebase.google.com/go"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func FirebaseAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// サーバー起動時に一回やればいい気がする
		opt := option.WithCredentialsFile("./key/otameshi-firebase-adminsdk.json")
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			fmt.Printf("error initializing app: %v", err)
			os.Exit(1)
		}

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

		fmt.Printf("Verified ID token: %v\n", token)

		ctx.Set("FirebaseID", token.UID)

		ctx.Next()
	}
}
