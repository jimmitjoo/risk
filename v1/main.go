package v1

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Region string
type Gamer int

const (
	NorthAmerica Region = "North America"
	SouthAmerica Region = "South America"
	Europe       Region = "Europe"
	Africa       Region = "Africa"
	Asia         Region = "Asia"
	Oceania      Region = "Oceania"

	Player1 Gamer = 1
	Player2 Gamer = 2
	Player3 Gamer = 3
)

var PlayerNames = map[Gamer]string{
	Player1: "Player 1",
	Player2: "Player 2",
	Player3: "Player 3",
}

type Game struct {
	Regions     map[Region]*RegionStatus
	Players     map[Gamer]*GamerStatus
	CurrentTurn Gamer
}

type RegionStatus struct {
	Name   Region
	Owner  Gamer
	Troops int
}

type GamerStatus struct {
	Alive       bool
	Cards       []*Card
	RegionCount int
}

type Card struct {
	Region Region
}

func NewGame() *Game {
	return &Game{
		Regions: map[Region]*RegionStatus{},
		Players: map[Gamer]*GamerStatus{
			Player1: {Alive: true, Cards: make([]*Card, 0), RegionCount: 0},
			Player2: {Alive: true, Cards: make([]*Card, 0), RegionCount: 0},
			Player3: {Alive: true, Cards: make([]*Card, 0), RegionCount: 0},
		},
		CurrentTurn: Player1,
	}
}

func (g *Game) init() {
	regions := []Region{NorthAmerica, SouthAmerica, Europe, Africa, Asia, Oceania}
	rand.Shuffle(len(regions), func(i, j int) {
		regions[i], regions[j] = regions[j], regions[i]
	})
	for i, region := range regions {
		owner := Gamer(i%len(g.Players) + 1)
		regionStatus := &RegionStatus{
			Name:   region,
			Owner:  owner,
			Troops: rand.Intn(10) + 1,
		}
		g.Regions[region] = regionStatus

		if _, ok := g.Players[owner]; ok {
			g.Players[owner].RegionCount++
		} else {
			fmt.Printf("Player %d does not exist", owner)
		}
	}
}

func (g *Game) PrintStatus() {
	fmt.Println("\nCurrent Game Status:")
	for _, region := range g.Regions {
		fmt.Printf("\tRegion: %s, Owner: %s, Troops: %d\n", region.Name, PlayerNames[region.Owner], region.Troops)
	}
	fmt.Printf("\nIt's %s's turn\n", PlayerNames[g.CurrentTurn])
}

func (g *Game) Attack(from Region, to Region, troops int) error {
	if g.Regions[from].Owner != g.CurrentTurn {
		return errors.New("it's not your turn")
	}
	if g.Regions[from].Troops <= troops {
		return errors.New("not enough troops")
	}
	if g.Regions[to].Owner == g.Regions[from].Owner {
		return errors.New("cannot attack self")
	}
	g.rollDice(from, to, troops)
	return nil
}

func (g *Game) rollDice(from Region, to Region, troops int) {
	attackerRolls := make([]int, troops)
	for i := range attackerRolls {
		attackerRolls[i] = rand.Intn(6) + 1
	}
	defenderRolls := make([]int, g.Regions[to].Troops)
	for i := range defenderRolls {
		defenderRolls[i] = rand.Intn(6) + 1
	}
	for _, attackerRoll := range attackerRolls {
		highestDefenderRoll := 0
		for _, defenderRoll := range defenderRolls {
			if defenderRoll > highestDefenderRoll {
				highestDefenderRoll = defenderRoll
			}
		}
		if attackerRoll > highestDefenderRoll {
			g.Regions[to].Troops--
			g.Regions[from].Troops++
			if g.Regions[to].Troops == 0 {
				g.Regions[to].Owner = g.Regions[from].Owner
				g.Players[g.Regions[from].Owner].RegionCount++
				g.Players[g.Regions[to].Owner].RegionCount--
			}
		} else {
			g.Regions[from].Troops--
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	game := NewGame()
	game.init()
	for {
		game.PrintStatus()
		fmt.Print("Enter attack command in FromRegion,ToRegion,Troops format (e.g., \"North America\",\"South America\",5) or q to quit: ")
		var response string

		_, err := fmt.Scanln(&response)
		if err != nil || response == "q" {
			fmt.Println("Game ended")
			fmt.Println(err)
			os.Exit(0)
		}
		parts := strings.Split(response, ",")
		if len(parts) != 3 {
			fmt.Println("Invalid command, needs to be: FromRegion,ToRegion,Troops")
			continue
		}
		from := strings.Trim(parts[0], "\"")
		to := strings.Trim(parts[1], "\"")
		troops, err := strconv.Atoi(parts[2])
		if err != nil || from == "" || to == "" {
			fmt.Println("Invalid region or troops, they need to be numbers")
			continue
		}
		err = game.Attack(Region(from), Region(to), troops)
		if err != nil {
			fmt.Println("Attack failed:", err)
			continue
		}
		game.CurrentTurn = 4 - game.CurrentTurn
	}
}
