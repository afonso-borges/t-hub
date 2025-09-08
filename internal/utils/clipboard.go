package utils

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
)

func formatFromClipboard(split GoldSplit) string {
	var sb strings.Builder

	sb.WriteString("=== LOOT SPLIT RESULTS ===\n\n")

	for _, transfer := range split.DirectTransfers {
		fmt.Fprintf(&sb, "%s to pay %s %s   |   bank: transfer %d to %s\n\n",
			transfer.From, transfer.To, FormatNumber(transfer.Amount), transfer.Amount, transfer.To)
	}

	fmt.Fprintf(&sb, "\ntotal profit: %s \n", FormatNumber(split.TotalBalance))
	fmt.Fprintf(&sb, "total for each player: %s \n", FormatNumber(split.EqualShare))

	return sb.String()
}

func SaveToClipboard(split GoldSplit) error {
	formatted := formatFromClipboard(split)
	return clipboard.WriteAll(formatted)
}

func CopyFromClipboard() (string, error) {
	i, err := clipboard.ReadAll()
	if err != nil {
		return "", fmt.Errorf("failed to read clipboard: %v", err)
	}
	return i, nil
}
