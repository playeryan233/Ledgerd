package report

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"ledgerd/internal/domain"
)

const (
	dateLayout     = "2006-01-02"
	floatTolerance = 1e-9
)

// LineItem represents a named amount in a report.
type LineItem struct {
	Name   string
	Amount float64
}

// BalanceSheet summarizes financial position for a given date.
type BalanceSheet struct {
	Date time.Time

	Assets      []LineItem
	Liabilities []LineItem
	Equity      []LineItem

	TotalAssets      float64
	TotalLiabilities float64
	TotalEquity      float64
}

// IncomeStatement summarizes performance for a period.
type IncomeStatement struct {
	StartDate time.Time
	EndDate   time.Time

	Income   []LineItem
	Expenses []LineItem

	TotalIncome  float64
	TotalExpense float64
	NetIncome    float64
}

// GenerateBalanceSheet aggregates journal entries into a balance sheet.
func GenerateBalanceSheet(entries []domain.JournalEntry, asOf time.Time) (BalanceSheet, error) {
	var sheet BalanceSheet
	sheet.Date = asOf

	balances := map[string]*accountBalance{}
	var totalIncome, totalExpenses float64

	for _, entry := range entries {
		entryDate, err := parseDate(entry.Date)
		if err != nil {
			return BalanceSheet{}, err
		}
		if entryDate.After(asOf) {
			continue
		}

		for _, line := range entry.Lines {
			acc := strings.TrimSpace(line.Account)
			if acc == "" {
				continue
			}

			balance := getOrCreateBalance(balances, acc)
			balance.Debit += line.Debit
			balance.Credit += line.Credit
		}
	}

	var assets, liabilities, equity []LineItem

	for account, bal := range balances {
		prefix := accountPrefix(account)
		amount := normalizeAmount(prefix, bal)
		if nearlyZero(amount) {
			continue
		}

		switch prefix {
		case prefixAssets:
			assets = append(assets, LineItem{Name: account, Amount: amount})
			sheet.TotalAssets += amount
		case prefixLiabilities:
			liabilities = append(liabilities, LineItem{Name: account, Amount: amount})
			sheet.TotalLiabilities += amount
		case prefixEquity:
			equity = append(equity, LineItem{Name: account, Amount: amount})
			sheet.TotalEquity += amount
		case prefixIncome:
			totalIncome += amount
		case prefixExpenses:
			totalExpenses += amount
		}
	}

	retained := totalIncome - totalExpenses
	if !nearlyZero(retained) {
		equity = append(equity, LineItem{
			Name:   "Retained Earnings",
			Amount: retained,
		})
		sheet.TotalEquity += retained
	}

	sortLineItems(assets)
	sortLineItems(liabilities)
	sortLineItems(equity)

	sheet.Assets = assets
	sheet.Liabilities = liabilities
	sheet.Equity = equity

	if diff := math.Abs(sheet.TotalAssets - (sheet.TotalLiabilities + sheet.TotalEquity)); diff > floatTolerance {
		return BalanceSheet{}, fmt.Errorf("balance sheet out of balance: assets %.2f, liabilities+equity %.2f", sheet.TotalAssets, sheet.TotalLiabilities+sheet.TotalEquity)
	}

	return sheet, nil
}

// GenerateIncomeStatement aggregates journal entries into an income statement.
func GenerateIncomeStatement(entries []domain.JournalEntry, start, end time.Time) (IncomeStatement, error) {
	if start.After(end) {
		return IncomeStatement{}, fmt.Errorf("start date must be before or equal to end date")
	}

	statement := IncomeStatement{
		StartDate: start,
		EndDate:   end,
	}

	balances := map[string]*accountBalance{}
	for _, entry := range entries {
		entryDate, err := parseDate(entry.Date)
		if err != nil {
			return IncomeStatement{}, err
		}
		if entryDate.Before(start) || entryDate.After(end) {
			continue
		}

		for _, line := range entry.Lines {
			acc := strings.TrimSpace(line.Account)
			if acc == "" {
				continue
			}

			balance := getOrCreateBalance(balances, acc)
			balance.Debit += line.Debit
			balance.Credit += line.Credit
		}
	}

	for account, bal := range balances {
		prefix := accountPrefix(account)
		amount := normalizeAmount(prefix, bal)
		if nearlyZero(amount) {
			continue
		}

		switch prefix {
		case prefixIncome:
			statement.Income = append(statement.Income, LineItem{Name: account, Amount: amount})
			statement.TotalIncome += amount
		case prefixExpenses:
			statement.Expenses = append(statement.Expenses, LineItem{Name: account, Amount: amount})
			statement.TotalExpense += amount
		}
	}

	sortLineItems(statement.Income)
	sortLineItems(statement.Expenses)
	statement.NetIncome = statement.TotalIncome - statement.TotalExpense

	return statement, nil
}

type accountBalance struct {
	Debit  float64
	Credit float64
}

const (
	prefixAssets      = "assets"
	prefixLiabilities = "liabilities"
	prefixEquity      = "equity"
	prefixIncome      = "income"
	prefixExpenses    = "expenses"
)

func accountPrefix(account string) string {
	name := strings.ToLower(account)
	if idx := strings.Index(name, ":"); idx >= 0 {
		name = name[:idx]
	}
	return name
}

func normalizeAmount(prefix string, bal *accountBalance) float64 {
	switch prefix {
	case prefixAssets, prefixExpenses:
		return bal.Debit - bal.Credit
	default:
		return bal.Credit - bal.Debit
	}
}

func getOrCreateBalance(balances map[string]*accountBalance, account string) *accountBalance {
	if existing, ok := balances[account]; ok {
		return existing
	}
	balances[account] = &accountBalance{}
	return balances[account]
}

func parseDate(value string) (time.Time, error) {
	return time.Parse(dateLayout, strings.TrimSpace(value))
}

func sortLineItems(items []LineItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
}

func nearlyZero(value float64) bool {
	return math.Abs(value) <= floatTolerance
}
