package messages

import (
	"context"
	"errors"
	"testing"

	mock "github.com/Svoevolin/workshop_1_bot/internal/mocks/messages"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
)

func TestModelSetCurrency(t *testing.T) {
	m := minimock.NewController(t)
	settings := mock.NewConfigGetterMock(m)
	userDB := mock.NewUserDBMock(m)
	sender := mock.NewMessageSenderMock(m)
	model := New(sender, settings, userDB, nil, nil, nil)

	t.Run("return error if currency code unsupported", func(t *testing.T) {
		settings.SupportedCurrencyCodesMock.Expect().Return([]string{"RUB", "EUR"})
		userDB.UserExistsMock.Expect(minimock.AnyContext, int64(123)).Return(false)

		text, err := model.setCurrency(context.TODO(), Message{UserID: 123, Text: "/set_currency USD"})
		assert.NoError(t, err)
		assert.EqualValues(t, text, "Валюта USD не поддерживается, отправьте команду /set_currency с одним из значений [RUB EUR]")
	})

	t.Run("returns an error if an error occurred at the time of saving", func(t *testing.T) {
		settings.SupportedCurrencyCodesMock.Expect().Return([]string{"RUB", "EUR"})
		userDB.UserExistsMock.Expect(minimock.AnyContext, int64(123)).Return(true)
		userDB.ChangeDefaultCurrencyMock.Expect(minimock.AnyContext, int64(123), "EUR").Return(errors.New("error"))
		_, err := model.setCurrency(context.TODO(), Message{UserID: 123, Text: "/set_currency EUR"})

		assert.EqualError(t, err, "failed to change user currency")
	})

	t.Run("a hint on how to use the bot will be returned for new users", func(t *testing.T) {
		settings.SupportedCurrencyCodesMock.Expect().Return([]string{"RUB", "EUR"})
		userDB.UserExistsMock.Expect(minimock.AnyContext, int64(123)).Return(false)
		userDB.ChangeDefaultCurrencyMock.Expect(minimock.AnyContext, int64(123), "EUR").Return(nil)
		text, err := model.setCurrency(context.TODO(), Message{UserID: 123, Text: "/set_currency EUR"})

		assert.NoError(t, err)
		assert.EqualValues(t, helpMessage, text)
	})

	t.Run("notification about change currency will be returned for old users", func(t *testing.T) {
		settings.SupportedCurrencyCodesMock.Expect().Return([]string{"RUB", "EUR"})
		userDB.UserExistsMock.Expect(minimock.AnyContext, int64(123)).Return(true)
		userDB.ChangeDefaultCurrencyMock.Expect(minimock.AnyContext, int64(123), "EUR").Return(nil)
		text, err := model.setCurrency(context.TODO(), Message{UserID: 123, Text: "/set_currency EUR"})

		assert.NoError(t, err)
		assert.EqualValues(t, "Установлена валюта по умолчанию EUR", text)
	})
}

func TestChangeDefaultCurrency(t *testing.T) {
	m := minimock.NewController(t)
	settings := mock.NewConfigGetterMock(m)
	userDB := mock.NewUserDBMock(m)
	sender := mock.NewMessageSenderMock(m)
	model := New(sender, settings, userDB, nil, nil, nil)

	t.Run("return command for changing  default currency", func(t *testing.T) {
		settings.SupportedCurrencyCodesMock.Expect().Return([]string{"RUB", "USD", "EUR"})
		text, buttons := model.changeDefaultCurrency()

		exceptedCurrencies := []map[string]string{
			{
				"RUB": "/set_currency RUB",
				"USD": "/set_currency USD",
				"EUR": "/set_currency EUR",
			},
		}
		assert.EqualValues(t, "Выберите валюту в которой будете производить расходы", text)
		assert.EqualValues(t, exceptedCurrencies, buttons)
	})
}
