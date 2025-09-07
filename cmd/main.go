package main

import (
	"log"

	"github.com/afonso-borges/t-hub/internal/theme"
	"github.com/afonso-borges/t-hub/internal/utils"
	"github.com/charmbracelet/huh"
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
			NextLabel("Next").WithTheme(theme.Theme()),
		),
	)

	err := welcomeForm.Run()
	if err != nil {
		log.Fatal(err)
	}

	// insert party hunt analyzer screen
	analyzerForm := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title(`Paste here the party hunt analyzer`).
				Description(`use ctrl+shift+v to paste in terminal`).
				Value(&analyzer).WithTheme(theme.Theme()),
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
	}

	// player removal selection screen
	playerOptions := utils.ExtractPlayerNames(players)

	removalForm := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title(`Remove player from loot split?`).
				Description("Select players to exclude from the loot calculation").
				Value(&playersToRemove).
				Options(playerOptions...).WithTheme(theme.Theme()),
		),
	)

	err = removalForm.Run()
	if err != nil {
		log.Fatal(err)
	}
	remainingPlayers := utils.FilterRemainingPlayers(players, playersToRemove)
	split := utils.CalculateGoldSplit(remainingPlayers)

	utils.SaveToClipboard(split)
	utils.DisplayTransfers(split)
}
