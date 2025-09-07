package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/afonso-borges/t-hub/internal/themes"
	"github.com/charmbracelet/lipgloss"
)

type PlayerTransfer struct {
	Player
	TransferAmount int
	FinalBalance   int
	Status         string
}

type DirectTransfer struct {
	From   string
	To     string
	Amount int
}

type GoldSplit struct {
	TotalBalance    int
	EqualShare      int
	PlayerTransfers []PlayerTransfer
	DirectTransfers []DirectTransfer
	Summary         TransferSummary
}

type TransferSummary struct {
	TotalOwed        int
	TotalReceived    int
	PlayersOwing     int
	PlayersReceiving int
	TransferCount    int
}

func CalculateGoldSplit(players []Player) GoldSplit {
	var totalBalance int
	for _, player := range players {
		totalBalance += player.Balance
	}
	playerCount := len(players)
	equalShare := totalBalance / playerCount

	var playerTransfers []PlayerTransfer

	var summary TransferSummary

	// Calculate individual transfer amount
	for _, player := range players {
		transferAmount := player.Balance - equalShare
		finalBalance := equalShare

		var status string
		if transferAmount > 0 {
			status = "owes"
			summary.TotalOwed += transferAmount
			summary.PlayersOwing++
		} else if transferAmount < 0 {
			status = "receives"
			summary.TotalReceived += -transferAmount
			summary.PlayersReceiving++
		} else {
			status = "balanced"
		}

		playerTransfers = append(playerTransfers, PlayerTransfer{
			Player:         player,
			TransferAmount: transferAmount,
			FinalBalance:   finalBalance,
			Status:         status,
		})
	}

	directTransfers := calculateDirectTransfers(playerTransfers)
	summary.TransferCount = len(directTransfers)

	return GoldSplit{
		TotalBalance:    totalBalance,
		EqualShare:      equalShare,
		PlayerTransfers: playerTransfers,
		DirectTransfers: directTransfers,
		Summary:         summary,
	}
}

// calculateDirectTransfers determines who should pay whom to minimize transactions
func calculateDirectTransfers(playerTransfers []PlayerTransfer) []DirectTransfer {
	var debtors []PlayerTransfer   // Players who owe money
	var creditors []PlayerTransfer // Players who should receive money

	// Separate debtors and creditors
	for _, pt := range playerTransfers {
		if pt.TransferAmount > 0 {
			debtors = append(debtors, pt)
		} else if pt.TransferAmount < 0 {
			creditors = append(creditors, pt)
		}
	}

	// Sort debtors by amount owed (descending) and creditors by amount to receive (descending)
	sort.Slice(debtors, func(i, j int) bool {
		return debtors[i].TransferAmount > debtors[j].TransferAmount
	})
	sort.Slice(creditors, func(i, j int) bool {
		return creditors[i].TransferAmount < creditors[j].TransferAmount // More negative first
	})
	var transfers []DirectTransfer

	// Match debtors with creditors
	for len(debtors) > 0 && len(creditors) > 0 {
		debtor := &debtors[0]
		creditor := &creditors[0]
		debt := debtor.TransferAmount
		credit := -creditor.TransferAmount // Make positive
		transferAmount := min(credit, debt)
		transfers = append(transfers, DirectTransfer{
			From:   debtor.Name,
			To:     creditor.Name,
			Amount: transferAmount,
		})

		// Update remaining amounts
		debtor.TransferAmount -= transferAmount
		creditor.TransferAmount += transferAmount
		// Remove settled players
		if debtor.TransferAmount == 0 {
			debtors = debtors[1:]
		}
		if creditor.TransferAmount == 0 {
			creditors = creditors[1:]
		}
	}
	return transfers
}

func DisplayTransfers(split GoldSplit) {
	fmt.Print(FormatTransfers(split))
}

func FormatTransfers(split GoldSplit) string {
	var sb strings.Builder
	kw := func(s string) string {
		return lipgloss.NewStyle().Foreground(themes.Indigo).Render(s)
	}

	dkw := func(s string) string {
		return lipgloss.NewStyle().Foreground(themes.Fuchsia).Render(s)
	}

	// result screen header
	fmt.Fprintf(&sb, "%s\n\n", lipgloss.NewStyle().
		Bold(true).
		Foreground(themes.Indigo).
		Render("Loot split results:"))

	// display transfers
	for _, transfer := range split.DirectTransfers {
		fmt.Fprintf(&sb, "%s %s %s %s\n",
			kw(transfer.From),
			dkw("to pay"),
			kw(transfer.To),
			kw(fmt.Sprintf("%d gp", transfer.Amount)))
	}

	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "%s %s\n",
		dkw("total profit: "),
		kw(fmt.Sprintf("%d gp", split.TotalBalance)))
	fmt.Fprintf(&sb, "%s %s\n",
		dkw("total for each player: "),
		kw(fmt.Sprintf("%d gp", split.EqualShare)))

	return sb.String()
}
