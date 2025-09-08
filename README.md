# T-Hub

A Terminal User Interface (TUI) application designed for analyzing Party Hunt data from **Tibia**, the classic MMORPG.T-Hub processes party hunt analyzer data and automatically calculates fair gold distribution among players, minimizing the number of required transfers.

## Features

- **Analyzer Processing**: Parses party hunt analyzer data directly from clipboard
- **Player Management**: Select which players to exclude from loot calculations
- **Optimal Split Calculation**: Automatically calculates the most efficient transfer distribution
- **Clipboard Integration**: Copies formatted results back to clipboard for easy sharing
- **Interactive TUI**: Clean, modern terminal interface with intuitive navigation
- **Custom Theming**: Distinctive visual indicators for better user experience

## Installation

### Prerequisites

- Go 1.24.1 or higher
- Linux users may also need `xclip` for clipboard functionality:

  ```bash
  # Ubuntu/Debian
  sudo apt install xclip

  # Arch Linux
  sudo pacman -S xclip

  # Fedora/RHEL
  sudo dnf install xclip
  ```

### Build from Source

```bash
git clone https://github.com/afonso-borges/t-hub.git
cd t-hub
go build -o t-hub cmd/main.go
```

## Usage

1. **Prepare Data**: Copy your party hunt analyzer data to clipboard
2. **Run Application**: Execute `./t-hub` in your terminal
3. **Process Data**: The application will automatically read and parse the analyzer data
4. **Select Players**: Choose any players to exclude from the loot split calculation
5. **View Results**: Review the calculated transfers and copy results to clipboard
6. **Repeat**: Option to process additional analyzer data

### Example Workflow

```bash
# Start the application
./t-hub

# Follow the interactive prompts:
# 1. Welcome screen - Press Enter to start
# 2. Player selection - Choose players to exclude (optional)
# 3. Results display - View calculated transfers
# 4. Copy to clipboard - Results are automatically formatted
# 5. Start over or exit
```

## How It Works

T-Hub uses an optimal algorithm to minimize the number of transfers required to achieve equal profit distribution:

1. **Data Parsing**: Extracts player names, loot values, supplies, and balances from analyzer text
2. **Equal Share Calculation**: Determines fair profit distribution based on total party balance
3. **Transfer Optimization**: Calculates minimal transfers using a greedy matching algorithm
4. **Result Formatting**: Provides both visual display and clipboard-ready text output

## Project Structure

```
t-hub/
├── cmd/
│   └── main.go              # Application entry point and TUI logic
├── internal/
│   ├── themes/
│   │   └── theme.go         # Custom UI theme configuration
│   └── utils/
│       ├── clipboard.go     # Clipboard operations
│       ├── parser.go        # Analyzer data parsing
│       └── transfers.go     # Loot split calculations
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
└── README.md               # Project documentation
```

## Dependencies

- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)**: TUI framework
- **[Huh](https://github.com/charmbracelet/huh)**: Form components
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)**: Style and layout
- **[Clipboard](https://github.com/atotto/clipboard)**: System clipboard access

## Technical Details

### Algorithm

The loot split calculation uses a greedy algorithm that:

- Separates players into debtors (owe money) and creditors (receive money)
- Sorts players by transfer amounts to optimize matching
- Minimizes total number of transactions required
- Ensures mathematical balance across all transfers

### Data Format

T-Hub expects party hunt analyzer data in the standard format(direct from the game) containing:

- Session information
- Loot type and total values
- Individual player statistics (loot, supplies, balance, damage, healing)
- Player roles (leader identification)
