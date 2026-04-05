package report

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"ledgerd/internal/domain"
)

// CashFlowStatement represents the direct-method cash flow report.
type CashFlowStatement struct {
	StartDate time.Time
	EndDate   time.Time

	Operating []LineItem
	Investing []LineItem
	Financing []LineItem

	NetOperating float64
	NetInvesting float64
	NetFinancing float64

	NetIncrease float64

	BeginningCash float64
	EndingCash    float64
}

// GenerateCashFlowStatement builds a cash flow statement using the direct method.
func GenerateCashFlowStatement(entries []domain.JournalEntry, start, end time.Time) (CashFlowStatement, error) {
	if start.After(end) {
		return CashFlowStatement{}, fmt.Errorf("start date must not be after end date")
	}

	beginning, err := computeCashBalance(entries, start, false)
	if err != nil {
		return CashFlowStatement{}, err
	}

	ending, err := computeCashBalance(entries, end, true)
	if err != nil {
		return CashFlowStatement{}, err
	}

	acc := newCashFlowAccumulator()

	for _, entry := range entries {
		entryDate, err := parseDate(entry.Date)
		if err != nil {
			return CashFlowStatement{}, err
		}

		if entryDate.Before(start) || entryDate.After(end) {
			continue
		}

		cashLines, counterpartSegments := classifyEntryLines(entry)
		if len(cashLines) == 0 || len(counterpartSegments) == 0 {
			continue
		}

		if err := acc.allocate(cashLines, counterpartSegments, entry.ID); err != nil {
			return CashFlowStatement{}, err
		}
	}

	statement := CashFlowStatement{
		StartDate: start,
		EndDate:   end,

		Operating: acc.items[categoryOperating],
		Investing: acc.items[categoryInvesting],
		Financing: acc.items[categoryFinancing],

		NetOperating: acc.totals[categoryOperating],
		NetInvesting: acc.totals[categoryInvesting],
		NetFinancing: acc.totals[categoryFinancing],

		BeginningCash: beginning,
		EndingCash:    ending,
	}

	statement.NetIncrease = statement.NetOperating + statement.NetInvesting + statement.NetFinancing

	if diff := math.Abs(statement.BeginningCash + statement.NetIncrease - statement.EndingCash); diff > floatTolerance {
		return CashFlowStatement{}, fmt.Errorf("cash flow statement imbalance: beginning %.2f + net %.2f != ending %.2f",
			statement.BeginningCash, statement.NetIncrease, statement.EndingCash)
	}

	return statement, nil
}

// cashFlowAccumulator keeps per-category detail and totals.
type cashFlowAccumulator struct {
	items  map[flowCategory][]LineItem
	totals map[flowCategory]float64
}

func newCashFlowAccumulator() *cashFlowAccumulator {
	return &cashFlowAccumulator{
		items: map[flowCategory][]LineItem{
			categoryOperating: {},
			categoryInvesting: {},
			categoryFinancing: {},
		},
		totals: map[flowCategory]float64{},
	}
}

func (a *cashFlowAccumulator) allocate(cashLines []cashLine, segments []counterpartSegment, entryID int64) error {
	// Convert line slices into mutable copies to adjust remaining balances.
	remaining := make([]counterpartSegment, len(segments))
	copy(remaining, segments)

	for _, c := range cashLines {
		amountLeft := c.Amount
		sign := 1.0
		if amountLeft < 0 {
			sign = -1.0
			amountLeft = -amountLeft
		}

		if nearlyZero(amountLeft) {
			continue
		}

		for i := range remaining {
			if nearlyZero(amountLeft) {
				break
			}

			if remaining[i].Remaining <= 0 {
				continue
			}

			chunk := remaining[i].Remaining
			if chunk > amountLeft {
				chunk = amountLeft
			}

			remaining[i].Remaining -= chunk
			a.addFlow(remaining[i].Category, remaining[i].Account, chunk*sign)
			amountLeft -= chunk
		}

		if amountLeft > floatTolerance {
			return fmt.Errorf("unable to allocate cash flow for entry %d", entryID)
		}
	}

	// rebuild sorted slices from maps
	for category, lines := range a.items {
		a.items[category] = consolidateAndSort(lines)
	}

	return nil
}

func (a *cashFlowAccumulator) addFlow(category flowCategory, account string, amount float64) {
	if nearlyZero(amount) {
		return
	}

	a.totals[category] += amount

	a.items[category] = append(a.items[category], LineItem{
		Name:   account,
		Amount: amount,
	})
}

func consolidateAndSort(items []LineItem) []LineItem {
	if len(items) == 0 {
		return nil
	}

	agg := map[string]float64{}
	for _, item := range items {
		agg[item.Name] += item.Amount
	}

	result := make([]LineItem, 0, len(agg))
	for name, amount := range agg {
		if nearlyZero(amount) {
			continue
		}
		result = append(result, LineItem{
			Name:   name,
			Amount: amount,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

type flowCategory string

const (
	categoryOperating flowCategory = "operating"
	categoryInvesting flowCategory = "investing"
	categoryFinancing flowCategory = "financing"
)

type cashLine struct {
	Amount float64
}

type counterpartSegment struct {
	Account   string
	Category  flowCategory
	Remaining float64
}

func classifyEntryLines(entry domain.JournalEntry) ([]cashLine, []counterpartSegment) {
	var cashLines []cashLine
	var counterparts []counterpartSegment

	for _, line := range entry.Lines {
		account := strings.TrimSpace(line.Account)
		if account == "" {
			continue
		}

		if isCashAccount(account) {
			delta := line.Debit - line.Credit
			if nearlyZero(delta) {
				continue
			}
			cashLines = append(cashLines, cashLine{Amount: delta})
			continue
		}

		category, ok := categorizeCounterpart(account)
		if !ok {
			continue
		}

		magnitude := math.Abs(line.Debit - line.Credit)
		if nearlyZero(magnitude) {
			continue
		}

		counterparts = append(counterparts, counterpartSegment{
			Account:   account,
			Category:  category,
			Remaining: magnitude,
		})
	}

	return cashLines, counterparts
}

func categorizeCounterpart(account string) (flowCategory, bool) {
	name := strings.ToLower(strings.TrimSpace(account))
	if name == "" {
		return "", false
	}

	prefix := name
	if idx := strings.Index(prefix, ":"); idx >= 0 {
		prefix = prefix[:idx]
	}

	switch prefix {
	case prefixIncome, prefixExpenses:
		return categoryOperating, true
	case prefixAssets:
		if isCashAccount(account) {
			return "", false
		}
		return categoryInvesting, true
	case prefixLiabilities, prefixEquity:
		return categoryFinancing, true
	default:
		return "", false
	}
}

var cashAccountKinds = map[string]struct{}{
	"cash":   {},
	"bank":   {},
	"alipay": {},
	"wechat": {},
}

func isCashAccount(account string) bool {
	account = strings.TrimSpace(account)
	if account == "" {
		return false
	}

	if !strings.HasPrefix(strings.ToLower(account), prefixAssets+":") {
		return false
	}

	name := account[len("Assets:"):]
	if name == "" {
		return false
	}

	segment := name
	if idx := strings.Index(segment, ":"); idx >= 0 {
		segment = segment[:idx]
	}

	segment = strings.ToLower(strings.TrimSpace(segment))
	_, ok := cashAccountKinds[segment]
	return ok
}

func computeCashBalance(entries []domain.JournalEntry, boundary time.Time, inclusive bool) (float64, error) {
	balances := map[string]*accountBalance{}

	for _, entry := range entries {
		entryDate, err := parseDate(entry.Date)
		if err != nil {
			return 0, err
		}

		if entryDate.After(boundary) || (!inclusive && entryDate.Equal(boundary)) {
			continue
		}

		for _, line := range entry.Lines {
			account := strings.TrimSpace(line.Account)
			if account == "" || !isCashAccount(account) {
				continue
			}

			balance := getOrCreateBalance(balances, account)
			balance.Debit += line.Debit
			balance.Credit += line.Credit
		}
	}

	var total float64
	for _, bal := range balances {
		total += bal.Debit - bal.Credit
	}
	return total, nil
}
