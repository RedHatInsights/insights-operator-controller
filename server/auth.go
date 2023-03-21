// Auth implementation based on JWT

/*
Copyright © 2019, 2020, 2021, 2022, 2023 Red Hat, Inc.

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

// Generated documentation is available at:
// https://godoc.org/github.com/RedHatInsights/insights-operator-controller/server
//
// Documentation in literate-programming-style is available at:
// https://redhatinsights.github.io/insights-operator-controller/packages/server/auth.html

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/RedHatInsights/insights-operator-utils/responses"

	jwt "github.com/golang-jwt/jwt/v4"
)

type contextKey string

const (
	// ContextKeyUser is a constant for user authentication token in request
	contextKeyUser = contextKey("user")
)

// Token JWT claims struct
type Token struct {
	Login string
	jwt.StandardClaims
}

// JWTAuthentication middleware for checking auth rights
func (s *Server) JWTAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenHeader := r.Header.Get("Authorization") // Grab the token from the header

		if tokenHeader == "" {
			// Token is missing, returns with error code 403 Unauthorized
			err := responses.SendForbidden(w, "Missing auth token")
			if err != nil {
				log.Println("Error sending response about missing auth token")
			}
			// everything has been handled already
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		// The token normally comes in format `Bearer {token-body}`, we
		// check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			err := responses.SendForbidden(w, "Invalid/Malformed auth token")
			if err != nil {
				log.Println("Error sending response about invalid/malformed auth token")
			}
			// everything has been handled already
			return
		}

		tokenPart := splitted[1] // Grab the token part, what we are truly interested in
		tk := &Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			// malformed token, returns with HTTP code 403 as usual
			err := responses.SendForbidden(w, "Malformed authentication token")
			if err != nil {
				log.Println("Error sending response about malformed authentication token")
			}
			// everything has been handled already
			return
		}

		if !token.Valid {
			// token is invalid, maybe not signed on this server
			err := responses.SendForbidden(w, "Token is not valid.")
			if err != nil {
				log.Println("Error sending response about not valid authentication token")
			}
			// everything has been handled already
			return
		}

		// Everything went well, proceed with the request and set the
		// caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), contextKeyUser, tk.Login)
		r = r.WithContext(ctx)

		// Proceed to proxy
		next.ServeHTTP(w, r)
	})
}
