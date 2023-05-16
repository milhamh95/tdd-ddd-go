package cookies_test

import (
	"context"
	"errors"
	"tdd-ddd-go/cookies"
	"tdd-ddd-go/cookies/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_CookiePurchases(t *testing.T) {
	t.Run(`Given a user tries to purchase a cookie and we have them in stock, "when they tap their card, they get charged and then receive an email receipt a few momenths later"`, func(t *testing.T) {
		emailAddress := "emailAddress"
		cardToken := "token"
		cookiesToBuy := 5
		totalExpectedCost := 250

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mocks.NewMockCookieStockChecker(ctrl)
		c := mocks.NewMockCardCharger(ctrl)
		e := mocks.NewMockEmailSender(ctrl)

		s.EXPECT().AmountInStock(context.Background()).Times(1).Return(cookiesToBuy)
		c.EXPECT().ChargeCard(context.Background(), cardToken, totalExpectedCost).Times(1).Return(nil)
		e.EXPECT().SendEmailReceipt(context.Background(), emailAddress, totalExpectedCost).Times(1).Return(nil)

		cookiesService := cookies.NewCookiesService(s, c, e)

		require.NoError(t, cookiesService.Purchase(context.Background(), cookiesToBuy, cardToken, emailAddress))
	})

	t.Run("Given that a user tries to purchase a cookie and we don’t have any in stock, we return an error to the cashier so they can apologize to the customer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mocks.NewMockCookieStockChecker(ctrl)
		s.EXPECT().AmountInStock(context.Background()).Times(1).Return(0)

		cookiesService := cookies.NewCookiesService(s, nil, nil)

		err := cookiesService.Purchase(context.Background(), 1, "", "")

		require.EqualError(t, errors.New("no cookie in stock"), err.Error())
	})

	t.Run("Given that a user tries to purchase a cookie and we have them in stock, but their card gets declined, we return an error to the cashier so that we can ban the customer from the store", func(t *testing.T) {
		cardToken := "token"
		cookiesToBuy := 5
		totalExpectedCost := 250

		var errorCardIsDeclined = errors.New("card is declined")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mocks.NewMockCookieStockChecker(ctrl)
		c := mocks.NewMockCardCharger(ctrl)

		s.EXPECT().AmountInStock(context.Background()).Times(1).Return(cookiesToBuy)
		c.EXPECT().ChargeCard(context.Background(), cardToken, totalExpectedCost).Times(1).Return(errorCardIsDeclined)

		cookiesService := cookies.NewCookiesService(s, c, nil)

		err := cookiesService.Purchase(context.Background(), cookiesToBuy, cardToken, "")

		require.EqualError(t, errorCardIsDeclined, err.Error())
	})

	t.Run("Given that a user purchases a cookie and we have them in stock, their card is charged successfully, but we fail to send an email, we return a message to the cashier so they can notify the customer that they will not get an email, but the transaction is still considered complete", func(t *testing.T) {
		emailAddress := "emailAddress"
		cardToken := "token"
		cookiesToBuy := 5
		totalExpectedCost := 250

		var errorFailedToSendEmail = errors.New("failed to send email")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mocks.NewMockCookieStockChecker(ctrl)
		c := mocks.NewMockCardCharger(ctrl)
		e := mocks.NewMockEmailSender(ctrl)

		s.EXPECT().AmountInStock(context.Background()).Times(1).Return(cookiesToBuy)
		c.EXPECT().ChargeCard(context.Background(), cardToken, totalExpectedCost).Times(1).Return(nil)
		e.EXPECT().SendEmailReceipt(context.Background(), emailAddress, totalExpectedCost).Times(1).Return(errorFailedToSendEmail)

		cookiesService := cookies.NewCookiesService(s, c, e)

		err := cookiesService.Purchase(context.Background(), cookiesToBuy, cardToken, emailAddress)
		require.EqualError(t, errorFailedToSendEmail, err.Error())
	})

	t.Run("Given someone wants to purchase more cookies than we have in stock we only charge them for the ones we do have", func(t *testing.T) {
		emailAddress := "emailAddress"
		cardToken := "token"
		cookiesToBuy := 5
		totalExpectedCost := 50

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := mocks.NewMockCookieStockChecker(ctrl)
		c := mocks.NewMockCardCharger(ctrl)
		e := mocks.NewMockEmailSender(ctrl)

		s.EXPECT().AmountInStock(context.Background()).Times(1).Return(1)
		c.EXPECT().ChargeCard(context.Background(), cardToken, totalExpectedCost).Times(1).Return(nil)
		e.EXPECT().SendEmailReceipt(context.Background(), emailAddress, totalExpectedCost).Times(1).Return(nil)

		cookiesService := cookies.NewCookiesService(s, c, e)

		require.NoError(t, cookiesService.Purchase(context.Background(), cookiesToBuy, cardToken, emailAddress))
	})
}
