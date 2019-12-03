// Auth implementation based on JWT

/*
Copyright © 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"context"
	jwt "github.com/dgrijalva/jwt-go"
	u "github.com/redhatinsighs/insights-operator-controller/utils"
	"net/http"
	"os"
	"strings"
)

type contextKey string

const (
	contextKeyUser = contextKey("user")
)

// Token JWT claims struct
type Token struct {
	Login string
	jwt.StandardClaims
}

// JWTAuthentication middleware for checking auth rights
func JWTAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			u.SendForbidden(w, "Missing auth token")
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			u.SendForbidden(w, "Invalid/Malformed auth token")
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			u.SendForbidden(w, "Malformed authentication token")
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			u.SendForbidden(w, "Token is not valid.")
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), contextKeyUser, tk.Login)
		r = r.WithContext(ctx)
		// Proceed to proxy
		next.ServeHTTP(w, r)
	})
}
