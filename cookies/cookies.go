package cookies

import (
	"context"
	"errors"
)

type CookieStockChecker interface {
	AmountInStock(ctx context.Context) int
}

type CardCharger interface {
	ChargeCard(ctx context.Context, cardToken string, amountInCents int) error
}

type EmailSender interface {
	SendEmailReceipt(ctx context.Context, emailAddress string, amountInCents int) error
}

type CookiesService struct {
	stockChecker CookieStockChecker
	cardCharger  CardCharger
	emailSender  EmailSender
}

func NewCookiesService(s CookieStockChecker, c CardCharger, e EmailSender) *CookiesService {
	return &CookiesService{
		stockChecker: s,
		cardCharger:  c,
		emailSender:  e,
	}
}

func (cs *CookiesService) Purchase(
	ctx context.Context,
	amountOfCookiesToPurchase int,
	cardToken, emailAddress string) error {
	priceOfCookie := 50

	cookiesInStock := cs.stockChecker.AmountInStock(ctx)
	if cookiesInStock == 0 {
		return errors.New("no cookie in stock")
	}

	if amountOfCookiesToPurchase > cookiesInStock {
		amountOfCookiesToPurchase = cookiesInStock
	}

	cost := priceOfCookie * amountOfCookiesToPurchase
	err := cs.cardCharger.ChargeCard(ctx, cardToken, cost)
	if err != nil {
		return errors.New("card is declined")
	}

	err = cs.emailSender.SendEmailReceipt(ctx, emailAddress, cost)
	if err != nil {
		return errors.New("failed to send email")
	}

	return nil
}
