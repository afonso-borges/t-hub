package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/afonso-borges/t-hub/internal/utils"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	playersToRemove []string
	analyzer        string = ""
)

func main() {
	// welcome screen
	welcomeForm := huh.NewForm(
		huh.NewGroup(huh.NewNote().
			Title("Welcome to T-HUB").
			Description("loot split in your terminal").
			Next(true).
			NextLabel("Next"),
		),
	)

	err := welcomeForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	// insert party hunt analyzer screen
	analyzerForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(`Paste here the party hunt analyzer`).
				Value(&analyzer),
		),
	)

	err = analyzerForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	// extract players from the input

	_, players, err := utils.ParseAnalyzer(analyzer)
	if err != nil {
		log.Printf("Error parsing analyzer: %v", err)
		log.Fatal("Failed to parse party hunt analyzer")
		os.Exit(1)
	}
	log.Printf("Found %d players", len(players))

	// player removal selection screen
	playerOptions := utils.ExtractPlayerNames(players)

	removalForm := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(`Remove player from loot split?`).
				Description("Select players to exclude from the loot calculation").
				Value(&playersToRemove).
				Options(playerOptions...),
		),
	)

	err = removalForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Players to remove: %v", playersToRemove)
	log.Printf("Remaining players will be included in loot split calculation")

	remainingPlayers := utils.FilterRemainingPlayers(players, playersToRemove)
	split := utils.CalculateGoldSplit(remainingPlayers)

	// Debug logging
	log.Printf("Split calculation results:")
	log.Printf("- Total balance: %d", split.TotalBalance)
	log.Printf("- Equal share: %d", split.EqualShare)
	log.Printf("- Number of remaining players: %d", len(remainingPlayers))
	for i, player := range remainingPlayers {
		log.Printf("- Player %d: %s (balance: %d)", i+1, player.Name, player.Balance)
	}
	log.Printf("- Number of direct transfers: %d", len(split.DirectTransfers))
	for i, transfer := range split.DirectTransfers {
		log.Printf("- Transfer %d: %s -> %s (%d gp)", i+1, transfer.From, transfer.To, transfer.Amount)
	}

	var sb strings.Builder
	kw := func(s string) string {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(s)
	}

	// Result screen
	fmt.Fprintf(&sb, "\n%s\n\n", lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Render("=== LOOT SPLIT RESULTS ==="))

	// Display transfers
	for _, transfer := range split.DirectTransfers {
		fmt.Fprintf(&sb, "%s to pay %s %s\n",
			kw(transfer.From),
			kw(transfer.To),
			kw(fmt.Sprintf("%d gp", transfer.Amount)))
	}

	fmt.Fprintf(&sb, "\n")
	fmt.Fprintf(&sb, "total profit: %s\n", kw(fmt.Sprintf("%d gp", split.TotalBalance)))
	fmt.Fprintf(&sb, "total for each player: %s\n", kw(fmt.Sprintf("%d gp", split.EqualShare)))

	fmt.Print(sb.String())
}
