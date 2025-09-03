package utils

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
)

func formatForClipboard(split GoldSplit) string {
    var sb strings.Builder

    sb.WriteString("=== LOOT SPLIT RESULTS ===\n\n")

    for _, transfer := range split.DirectTransfers {
        fmt.Fprintf(&sb, "%s to pay %s %d gp\n",
            transfer.From, transfer.To, transfer.Amount)
    }

    fmt.Fprintf(&sb, "\ntotal profit: %d gp\n", split.TotalBalance)
    fmt.Fprintf(&sb, "total for each player: %d gp\n", split.EqualShare)

    return sb.String()
}

func SaveToClipboard(split GoldSplit) error {
    formatted := formatForClipboard(split)
    return clipboard.WriteAll(formatted)
}