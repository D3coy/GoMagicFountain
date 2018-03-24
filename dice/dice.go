package dice

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

type dieRoll struct {
	DieType int
	Value   int
}

type rolls struct {
	Dice []dieRoll
}

func (r rolls) RollAll() {
	for _, die := range r.Dice {
		rand.Seed(time.Now().Unix())
		die.Value = rand.Intn(die.DieType)
	}
}

func (r *rolls) OnlyTop(num int) {
	r.Sort()
	if num < len(r.Dice) {
		r.Dice = r.Dice[:num]
	}
}

func (r rolls) String() string {
	var returnText string
	for _, die := range r.Dice {
		returnText += fmt.Sprintf("\t d%v: %v\n", die.DieType, die.Value)
	}
	return returnText
}

func (r rolls) Sort() {
	sort.Slice(r.Dice, func(i, j int) bool {
		return r.Dice[i].Value > r.Dice[j].Value
	})
}

func (r rolls) GetType(dieType int) []dieRoll {
	var dice []dieRoll
	for _, die := range r.Dice {
		if die.DieType == dieType {
			dice = append(dice, die)
		}
	}
	sort.Slice(r.Dice, func(i, j int) bool {
		return dice[i].Value > dice[j].Value
	})
	return dice
}

func (r rolls) Total() int {
	total := 0
	for _, die := range r.Dice {
		total += die.Value
	}
	return total
}

// RollDice is used to parse a string for flags and
// dice to roll along with their quantities. The
// flags must be in the beginning, and it uses os.Args
// and golang's flag package to parse for them
func RollDice(args []string) (string, error) {

	var topFlag int
	var sortFlag bool
	var errFlag bool
	var returnString string
	currentRolls := new(rolls)

	diceFlagSet := new(flag.FlagSet)
	diceFlagSet.IntVar(&topFlag, "t", 0, "This will show only the top number of results")
	diceFlagSet.BoolVar(&sortFlag, "s", false, "This will sort the roll results")
	err := diceFlagSet.Parse(args)
	if err != nil {
		return "", errors.New(strings.Split(err.Error(), ":")[0]) // This gets the first part of the error with failed flag
	}
	args = diceFlagSet.Args()

	for len(args) > 0 {
		cmd := args[0]
		if strings.HasPrefix(cmd, "d") {
			rollQ := 1
			if 1 < len(args) {
				if i, err := strconv.Atoi(args[1]); err == nil {
					if i <= 10 {
						rollQ = i
						args = append(args[:1], args[2:]...)
					} else if i > 10 {
						return "", errors.New("you can only have up to 10 rolls at a time")
					}
				} else {
					log.Println("Error converting string to int:", err)
				}
			}

			if i, err := strconv.Atoi(cmd[1:]); err == nil {
				for count := rollQ; count > 0; count-- {
					roll := dieRoll{i, rand.Intn(i) + 1}
					currentRolls.Dice = append(currentRolls.Dice, roll)
				}
			} else {
				log.Println("Error converting string to int:", err)
				return "", fmt.Errorf("%v is an invalid die", cmd)
			}

		} else {
			return "", fmt.Errorf("invalid option: %v", cmd)
		}

		args = args[1:]
	}
	if !errFlag && (len(currentRolls.Dice) != 0) {
		currentRolls.RollAll()
		if topFlag != 0 {
			currentRolls.OnlyTop(topFlag)
		} else if sortFlag {
			currentRolls.Sort()
		}
		returnString += fmt.Sprint(currentRolls)
		returnString += fmt.Sprintln("Total:", currentRolls.Total())
	} else if len(currentRolls.Dice) == 0 {
		return "", errors.New("there are no dice to roll")
	}
	return returnString, nil
}
