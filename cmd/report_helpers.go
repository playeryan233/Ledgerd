package cmd

import (
	"fmt"
	"io"
	"math"

	"ledgerd/internal/report"
)

const reportDateLayout = "2006-01-02"

func printReportSection(w io.Writer, title string, items []report.LineItem) {
	fmt.Fprintln(w, title)
	if len(items) == 0 {
		fmt.Fprintln(w, "  (none)")
		fmt.Fprintln(w)
		return
	}

	for _, item := range items {
		fmt.Fprintf(w, "  %-20s %10.2f\n", item.Name, item.Amount)
	}
	fmt.Fprintln(w)
}

func printSignedSection(w io.Writer, title string, items []report.LineItem) {
	fmt.Fprintln(w, title)
	if len(items) == 0 {
		fmt.Fprintln(w, "  (none)")
		fmt.Fprintln(w)
		return
	}

	for _, item := range items {
		sign, value := signedParts(item.Amount)
		fmt.Fprintf(w, "  %-20s %s%10.2f\n", item.Name, sign, value)
	}
	fmt.Fprintln(w)
}

func printSignedTotal(w io.Writer, label string, amount float64) {
	sign, value := signedParts(amount)
	fmt.Fprintf(w, "%-20s %s%10.2f\n", label, sign, value)
}

func printUnsignedTotal(w io.Writer, label string, amount float64) {
	fmt.Fprintf(w, "%-20s %10.2f\n", label, amount)
}

func signedParts(amount float64) (string, float64) {
	if amount < 0 {
		return "-", math.Abs(amount)
	}
	return "+", math.Abs(amount)
}
