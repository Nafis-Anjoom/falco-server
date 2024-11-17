package main

import (
	"chat/utils"
	"context"
	"errors"
	"net/http"
	"strings"
)

var (
    missingAuthHeaderError = errors.New("Authorization header missing")
    invalidAuthHeaderError = errors.New("Invalid Authorization header")
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        header := request.Header.Get("Authorization")
        if header == "" {
            utils.WriteErrorResponse(writer, request, http.StatusUnauthorized, missingAuthHeaderError)
            return
        }

        headerParts := strings.Split(header, " ")
        if len(headerParts) != 2 || headerParts[0] != "Bearer" {
            utils.WriteErrorResponse(writer, request, http.StatusUnauthorized, invalidAuthHeaderError)
            return
        }

        tokenString := headerParts[1]
        userId, err := app.authService.VerifyToken(tokenString)
        if err != nil {
            utils.WriteErrorResponse(writer, request, http.StatusUnauthorized, err)
            return
        }

        ctx := context.WithValue(request.Context(), "userId", userId)

        next.ServeHTTP(writer, request.WithContext(ctx))
    })
}
