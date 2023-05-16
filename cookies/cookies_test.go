package cookies_test

import (
	"context"
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
}
