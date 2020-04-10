package middleware

import "github.com/iris-contrib/middleware/jwt"

func JwtHandler() *jwt.Middleware  {
	var mySecret = []byte("HS2JDFKhu7Y1av7b")
	return jwt.New(jwt.Config{
		ValidationKeyGetter: func(token *jwt.Token) (i interface{}, err error) {
			return mySecret, nil
		},
		ContextKey:          "",
		ErrorHandler:        nil,
		CredentialsOptional: false,
		Extractor:           nil,
		EnableAuthOnOptions: false,
		SigningMethod:       jwt.SigningMethodES256,
		Expiration:          false,
	})
}