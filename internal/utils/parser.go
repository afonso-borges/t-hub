package utils

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
)

type Party struct {
	SessionData string
	Session     string
	LootType    string
	Loot        int
	Supplies    int
	Balance     int
}

type Player struct {
	Name     string
	Leader   bool
	Loot     int
	Supplies int
	Balance  int
	Damage   int
	Healing  int
}

func parseNumber(s string) int {
	s = strings.ReplaceAll(s, ",", "")
	if s == "" {
		return 0
	}
	val, _ := strconv.Atoi(s)
	return val
}

var (
	LeaderSuffixRX    = regexp.MustCompile(`\s*\(Leader\)\s*$`)
	AnalyzersNumberRX = regexp.MustCompile(`-?\d[\d,]*`)
)

func ExtractPlayerNames(players []Player) []huh.Option[string] {
	options := make([]huh.Option[string], len(players))
	for i, player := range players {
		displayName := player.Name
		if player.Leader {
			displayName += " (Leader)"
		}
		options[i] = huh.NewOption(displayName, player.Name)
	}
	return options
}

func ParseAnalyzer(input string) (Party, []Player, error) {
	var party Party
	var players []Player

	// Parse party header information
	party.SessionData = extractValue(input, "Session data:", "Session:")
	party.Session = extractValue(input, "Session:", "Loot Type:")
	party.LootType = extractValue(input, "Loot Type:", "Loot:")
	party.Loot = parseNumber(extractValue(input, "Loot:", "Supplies:"))
	party.Supplies = parseNumber(extractValue(input, "Supplies:", "Balance:"))
	party.Balance = parseNumber(extractValue(input, "Balance:", getFirstPlayerName(input)))

	// Extract players
	playerNames := extractPlayerNames(input)
	for i, playerName := range playerNames {
		player := Player{
			Name:   playerName,
			Leader: strings.Contains(playerName, "(Leader)"),
		}

		// Clean player name
		if player.Leader {
			player.Name = LeaderSuffixRX.ReplaceAllString(player.Name, "")
		}

		// Find player data section
		playerStart := strings.Index(input, playerName)
		var playerEnd int
		if i+1 < len(playerNames) {
			playerEnd = strings.Index(input, playerNames[i+1])
		} else {
			playerEnd = len(input)
		}

		playerSection := input[playerStart:playerEnd]

		// Extract player stats
		player.Loot = parseNumber(extractPlayerStat(playerSection, "Loot:"))
		player.Supplies = parseNumber(extractPlayerStat(playerSection, "Supplies:"))
		player.Balance = parseNumber(extractPlayerStat(playerSection, "Balance:"))
		player.Damage = parseNumber(extractPlayerStat(playerSection, "Damage:"))
		player.Healing = parseNumber(extractPlayerStat(playerSection, "Healing:"))

		players = append(players, player)
	}

	if len(players) == 0 {
		return party, nil, fmt.Errorf("no players found on party analyzer")
	}

	return party, players, nil
}

// Helper function to extract value between two markers
func extractValue(input, startMarker, endMarker string) string {
	startIdx := strings.Index(input, startMarker)
	if startIdx == -1 {
		return ""
	}
	startIdx += len(startMarker)

	endIdx := strings.Index(input[startIdx:], endMarker)
	if endIdx == -1 {
		return strings.TrimSpace(input[startIdx:])
	}

	return strings.TrimSpace(input[startIdx : startIdx+endIdx])
}

// Helper function to get the first player name
func getFirstPlayerName(input string) string {
	// Find the first occurrence after "Balance:" that doesn't contain ":"
	balanceIdx := strings.Index(input, "Balance:")
	if balanceIdx == -1 {
		return ""
	}

	remaining := input[balanceIdx:]
	// Skip the balance value and find the next non-numeric word
	words := strings.Fields(remaining)
	for i, word := range words {
		if i > 1 && !strings.Contains(word, ":") && !AnalyzersNumberRX.MatchString(word) && !strings.Contains(word, ",") {
			return word
		}
	}
	return ""
}

// Helper function to extract player names from single line format
func extractPlayerNames(input string) []string {
	var names []string

	// Find all potential player names (words that don't contain ":" and appear after balance)
	balanceIdx := strings.Index(input, "Balance:")
	if balanceIdx == -1 {
		return names
	}

	remaining := input[balanceIdx:]
	words := strings.Fields(remaining)

	for i := 0; i < len(words); i++ {
		word := words[i]
		// Skip numbers, colons, and known keywords
		if strings.Contains(word, ":") || AnalyzersNumberRX.MatchString(word) ||
			strings.Contains(word, ",") || word == "Loot" || word == "Supplies" ||
			word == "Balance" || word == "Damage" || word == "Healing" {
			continue
		}

		// Check if this could be a player name by looking ahead for stats
		playerName := word

		// Handle multi-word names and (Leader) suffix
		j := i + 1
		for j < len(words) && !strings.Contains(words[j], ":") {
			if words[j] == "(Leader)" {
				playerName += " " + words[j]
				j++
				break
			} else if !AnalyzersNumberRX.MatchString(words[j]) {
				playerName += " " + words[j]
				j++
			} else {
				break
			}
		}

		// Verify this is a player by checking if "Loot:" follows
		remainingAfterName := strings.Join(words[j:], " ")
		if strings.Contains(remainingAfterName, "Loot:") {
			names = append(names, playerName)
			i = j - 1 // Skip processed words
		}
	}

	return names
}

// Helper function to extract player stat value
func extractPlayerStat(playerSection, statName string) string {
	statIdx := strings.Index(playerSection, statName)
	if statIdx == -1 {
		return "0"
	}

	// Find the number after the stat name
	remaining := playerSection[statIdx+len(statName):]
	match := AnalyzersNumberRX.FindString(remaining)
	if match == "" {
		return "0"
	}

	return match
}

func FilterRemainingPlayers(players []Player, playersToRemove []string) []Player {
	var remainingPlayers []Player

	for _, player := range players {
		shouldRemove := slices.Contains(playersToRemove, player.Name)
		if !shouldRemove {
			remainingPlayers = append(remainingPlayers, player)
		}
	}

	return remainingPlayers
}
