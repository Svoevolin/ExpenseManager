package messages

import (
	"context"
	"errors"
	"time"

	"github.com/Svoevolin/workshop_1_bot/internal/domain"
)

const dateFormat = "02.01.2006"

const (
	newUserMessage = `Привет, я бот - менеджер твоей бухгалтерии. Мне нужно знать валюту которой ты платишь`
	startMessage   = "Привет, я бот - менеджер твоей бухгалтерии.\n\n" + helpMessage
	helpMessage    = `Прочитай команды, чтобы понять что я умею делать:
/change_currency - Изменить валюту расходов
/add сумма; описание; <дата> - Добавь новую трату, если не укажешь дату - будет сегодня

Посмотреть расходы:
/spent - за всё время
/spent_day - за день 
/spent_week - за неделю
/spent_month - за месяц
/spent_year - за год`
	unknownMessage = `Неизвестная команда. Чтобы посмотреть список команд отправь /help`

	FailedMessage               = "Я временно не работаю, повторите попытку позже"
	FailedChangeCurrencyMessage = "Не удалось изменить текущую валюту, повторите попытку позже"
	InvalidCommandMessage       = "Неверный формат команды, исправьте и повторите команду"
	InvalidAmountMessage        = "Неверный формат суммы, исправьте и повторите команду"
	InvalidDateMessage          = "Неверный формат даты, исправьте и повторите команду"
	FailedWriteMessage          = "Не удалось записать расход, повторите попытку позже"
)

var (
	ErrImpossibleToChangeUserCurrency = errors.New("failed to change user currency")
	ErrInvalidCommand                 = errors.New("invalid command")
	ErrInvalidAmount                  = errors.New("invalid amount")
	ErrInvalidDate                    = errors.New("invalid date")
	ErrWriteToDatabase                = errors.New("failed to write to the database")
	ErrGetRecordsInDatabase           = errors.New("failed to get records from the database")
)

type MessageSender interface {
	SendMessage(text string, userID int64, keyboardRows ...map[string]string) error
}

type ConfigGetter interface {
	SupportedCurrencyCodes() []string
	GetBaseCurrency() string
}

type UserDB interface {
	UserExists(ctx context.Context, userID int64) bool
	ChangeDefaultCurrency(ctx context.Context, userID int64, currency string) error
	GetDefaultCurrency(ctx context.Context, userID int64) (string, error)
}

type RateDB interface {
	GetRate(ctx context.Context, code string, date time.Time) *domain.Rate
}

type ExpenseDB interface {
	AddExpense(ctx context.Context, userID int64, kopecks int64, title string, date time.Time) error
	GetExpenses(ctx context.Context, userID int64) ([]domain.Expense, error)
}

type ExchangeRateUpdate interface {
	UpdateCurrency(ctx context.Context, time time.Time) error
}

type Model struct {
	tgClient    MessageSender
	config      ConfigGetter
	userDB      UserDB
	rateDB      RateDB
	expenseDB   ExpenseDB
	rateUpdater ExchangeRateUpdate
}

func New(tgClient MessageSender, config ConfigGetter, userDB UserDB, rateDB RateDB, expenseDB ExpenseDB, rateUpdater ExchangeRateUpdate) *Model {
	return &Model{
		tgClient:    tgClient,
		config:      config,
		userDB:      userDB,
		rateDB:      rateDB,
		expenseDB:   expenseDB,
		rateUpdater: rateUpdater,
	}
}

type Message struct {
	Text   string
	UserID int64
}
