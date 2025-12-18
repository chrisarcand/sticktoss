package teamgen

import (
	"errors"
	"math/rand"
	"sort"
	"time"

	"github.com/sticktoss/backend/internal/models"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Team represents a generated team with players
type Team struct {
	Number      int             `json:"number"`
	Players     []models.Player `json:"players"`
	TotalWeight int             `json:"total_weight"`
}

// GenerateBalancedTeams creates balanced teams from a list of players
// lockedPlayers is an array of player ID arrays - each inner array represents players that must be on the same team
// separatedPlayers is an array of player ID arrays - each inner array represents players that must be on different teams
func GenerateBalancedTeams(players []models.Player, numTeams int, lockedPlayers [][]uint, separatedPlayers [][]uint) ([]Team, error) {
	if numTeams < 2 {
		return nil, errors.New("must have at least 2 teams")
	}

	if len(players) < numTeams {
		return nil, errors.New("not enough players for the requested number of teams")
	}

	// Initialize teams
	teams := make([]Team, numTeams)
	for i := range teams {
		teams[i].Number = i + 1
		teams[i].Players = []models.Player{}
		teams[i].TotalWeight = 0
	}

	// Track which players have been assigned
	assignedPlayers := make(map[uint]bool)

	// First, handle locked players
	if len(lockedPlayers) > 0 {
		if len(lockedPlayers) > numTeams {
			return nil, errors.New("cannot have more locked groups than teams")
		}

		// Create a map for quick player lookup
		playerMap := make(map[uint]models.Player)
		for _, p := range players {
			playerMap[p.ID] = p
		}

		// Assign locked groups to teams
		for i, lockedGroup := range lockedPlayers {
			if i >= numTeams {
				break
			}

			for _, playerID := range lockedGroup {
				player, exists := playerMap[playerID]
				if !exists {
					return nil, errors.New("locked player not found in group")
				}

				teams[i].Players = append(teams[i].Players, player)
				teams[i].TotalWeight += player.SkillWeight
				assignedPlayers[playerID] = true
			}
		}
	}

	// Handle separated players (must be on different teams)
	if len(separatedPlayers) > 0 {
		// Create a map for quick player lookup
		playerMap := make(map[uint]models.Player)
		for _, p := range players {
			playerMap[p.ID] = p
		}

		for _, separatedGroup := range separatedPlayers {
			// Validate: can't separate more players than we have teams
			if len(separatedGroup) > numTeams {
				return nil, errors.New("cannot separate more players than the number of teams")
			}

			// Check for conflicts with locked players
			for _, playerID := range separatedGroup {
				if assignedPlayers[playerID] {
					return nil, errors.New("cannot separate a player that is already locked to a team")
				}
			}

			// Assign each player in the separated group to a different team
			// Randomly shuffle which teams they go to for fairness
			teamIndices := make([]int, numTeams)
			for i := range teamIndices {
				teamIndices[i] = i
			}
			rand.Shuffle(len(teamIndices), func(i, j int) {
				teamIndices[i], teamIndices[j] = teamIndices[j], teamIndices[i]
			})

			for i, playerID := range separatedGroup {
				player, exists := playerMap[playerID]
				if !exists {
					return nil, errors.New("separated player not found in group")
				}

				teamIdx := teamIndices[i]
				teams[teamIdx].Players = append(teams[teamIdx].Players, player)
				teams[teamIdx].TotalWeight += player.SkillWeight
				assignedPlayers[playerID] = true
			}
		}
	}

	// Collect remaining players
	remainingPlayers := []models.Player{}
	for _, p := range players {
		if !assignedPlayers[p.ID] {
			remainingPlayers = append(remainingPlayers, p)
		}
	}

	// Shuffle remaining players for randomness
	rand.Shuffle(len(remainingPlayers), func(i, j int) {
		remainingPlayers[i], remainingPlayers[j] = remainingPlayers[j], remainingPlayers[i]
	})

	// Sort remaining players by skill weight (descending) for better balance
	sort.Slice(remainingPlayers, func(i, j int) bool {
		return remainingPlayers[i].SkillWeight > remainingPlayers[j].SkillWeight
	})

	// Assign remaining players using greedy algorithm (assign to team with lowest total weight)
	for _, player := range remainingPlayers {
		// Find all teams with minimum total weight
		minWeight := teams[0].TotalWeight
		for i := 1; i < numTeams; i++ {
			if teams[i].TotalWeight < minWeight {
				minWeight = teams[i].TotalWeight
			}
		}

		// Collect all teams that have the minimum weight
		minTeams := []int{}
		for i := 0; i < numTeams; i++ {
			if teams[i].TotalWeight == minWeight {
				minTeams = append(minTeams, i)
			}
		}

		// Randomly pick one of the teams with minimum weight
		minTeamIdx := minTeams[rand.Intn(len(minTeams))]

		teams[minTeamIdx].Players = append(teams[minTeamIdx].Players, player)
		teams[minTeamIdx].TotalWeight += player.SkillWeight
	}

	return teams, nil
}
