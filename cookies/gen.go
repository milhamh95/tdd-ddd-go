package cookies

import _ "github.com/golang/mock/mockgen/model"

//go:generate mockgen -package mocks -destination mocks/mock_cookies.go tdd-ddd-go/cookies CookieStockChecker,CardCharger,EmailSender
