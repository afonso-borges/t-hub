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
	SessionData string `json:"session_data,omitempty"`
	Session     string `json:"session,omitempty"`
	LootType    string `json:"loot_type,omitempty"`
	Loot        int    `json:"loot"`
	Supplies    int    `json:"supplies"`
	Balance     int    `json:"balance"`
}

type Player struct {
	Name     string `json:"name"`
	Leader   bool   `json:"leader"`
	Loot     int    `json:"loot"`
	Supplies int    `json:"supplies"`
	Balance  int    `json:"balance"`
	Damage   int    `json:"damage"`
	Healing  int    `json:"healing"`
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
	// Detect format type
	// hasNewlines := strings.Contains(input, "\\n") || strings.Contains(input, "\n")

	return parseSingleLineFormat(input)
	// if hasNewlines {
	// 	fmt.Println("has new lines")
	// 	return parseStructuredFormat(input)
	// } else {
	// 	fmt.Println("hasnt new lines")
	// 	return parseSingleLineFormat(input)
	// }
}

// parseStructuredFormat handles Type 1 format (with \n and \t)
// func parseStructuredFormat(input string) (Party, []Player, error) {
// 	lines := strings.Split(input, "\\n")
// 	var party Party
// 	var players []Player
// 	var currentPlayer *Player
// 	inPlayers := false
//
// 	for _, raw := range lines {
// 		line := strings.TrimSpace(raw)
// 		if line == "" {
// 			continue
// 		}
//
// 		// Check if line starts with tab (player attribute)
// 		isPlayerAttribute := strings.HasPrefix(line, "\\t")
// 		if isPlayerAttribute {
// 			line = strings.TrimPrefix(line, "\\t")
// 		}
//
// 		// If line doesn't have ":" or is a player attribute without current player, skip
// 		if !strings.Contains(line, ":") {
// 			if !isPlayerAttribute {
// 				// This is a player name
// 				inPlayers = true
// 				if currentPlayer != nil {
// 					players = append(players, *currentPlayer)
// 				}
// 				name := line
// 				leader := false
// 				if LeaderSuffixRX.MatchString(name) {
// 					leader = true
// 					name = LeaderSuffixRX.ReplaceAllString(name, "")
// 				}
// 				currentPlayer = &Player{name, leader, 0, 0, 0, 0, 0}
// 			}
// 			continue
// 		}
//
// 		// Handle party header attributes (before players section)
// 		if !inPlayers && !isPlayerAttribute {
// 			switch {
// 			case strings.HasPrefix(line, "Session data:"):
// 				party.SessionData = strings.TrimSpace(strings.TrimPrefix(line, "Session data:"))
// 			case strings.HasPrefix(line, "Session:"):
// 				party.Session = strings.TrimSpace(strings.TrimPrefix(line, "Session:"))
// 			case strings.HasPrefix(line, "Loot Type:"):
// 				party.LootType = strings.TrimSpace(strings.TrimPrefix(line, "Loot Type:"))
// 			case strings.HasPrefix(line, "Loot:"):
// 				party.Loot = parseNumber(AnalyzersNumberRX.FindString(line))
// 			case strings.HasPrefix(line, "Supplies:"):
// 				party.Supplies = parseNumber(AnalyzersNumberRX.FindString(line))
// 			case strings.HasPrefix(line, "Balance:"):
// 				party.Balance = parseNumber(AnalyzersNumberRX.FindString(line))
// 			}
// 			continue
// 		}
//
// 		// Handle player attributes (lines starting with \t)
// 		if isPlayerAttribute && currentPlayer != nil {
// 			switch {
// 			case strings.HasPrefix(line, "Loot:"):
// 				currentPlayer.Loot = parseNumber(AnalyzersNumberRX.FindString(line))
// 			case strings.HasPrefix(line, "Supplies:"):
// 				currentPlayer.Supplies = parseNumber(AnalyzersNumberRX.FindString(line))
// 			case strings.HasPrefix(line, "Balance:"):
// 				currentPlayer.Balance = parseNumber(AnalyzersNumberRX.FindString(line))
// 			case strings.HasPrefix(line, "Damage:"):
// 				currentPlayer.Damage = parseNumber(AnalyzersNumberRX.FindString(line))
// 			case strings.HasPrefix(line, "Healing:"):
// 				currentPlayer.Healing = parseNumber(AnalyzersNumberRX.FindString(line))
// 			}
// 		}
// 	}
//
// 	if currentPlayer != nil {
// 		players = append(players, *currentPlayer)
// 	}
//
// 	if len(players) == 0 {
// 		return party, nil, fmt.Errorf("no players found on party analyzer")
// 	}
//
// 	return party, players, nil
// }

// parseSingleLineFormat handles Type 2 format (single line)
func parseSingleLineFormat(input string) (Party, []Player, error) {
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
