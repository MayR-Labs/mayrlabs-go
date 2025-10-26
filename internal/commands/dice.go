package commands

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// RollDiceCmd rolls dice
var RollDiceCmd = &cobra.Command{
	Use:     "roll-dice [number_of_dice]",
	Aliases: []string{"dice", "rolldice", "rd"},
	Short:   "Roll dice and get random results",
	Long:    "Roll one or more dice (1-6) and display the results",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var numDice int
		var err error

		// Interactive mode or argument provided
		if len(args) == 0 {
			numDiceStr, err := utils.SurveyInput("Enter number of dice to roll:", "1")
			if err != nil {
				return err
			}
			_, err = fmt.Sscanf(numDiceStr, "%d", &numDice)
			if err != nil {
				return fmt.Errorf("invalid number: %s", numDiceStr)
			}
		} else {
			_, err = fmt.Sscanf(args[0], "%d", &numDice)
			if err != nil {
				return fmt.Errorf("invalid number: %s", args[0])
			}
		}

		if numDice < 1 {
			return fmt.Errorf("number of dice must be at least 1")
		}

		if numDice > 100 {
			return fmt.Errorf("number of dice cannot exceed 100")
		}

		fmt.Printf("🎲 Rolling %d dice...\n\n", numDice)

		total := 0
		results := make([]int, numDice)

		for i := 0; i < numDice; i++ {
			roll, err := rollDie()
			if err != nil {
				return err
			}
			results[i] = roll
			total += roll
		}

		// Display results
		if numDice == 1 {
			fmt.Printf("Result: %d\n", results[0])
		} else {
			fmt.Printf("Results: %v\n", results)
			fmt.Printf("Total: %d\n", total)
			fmt.Printf("Average: %.2f\n", float64(total)/float64(numDice))
		}

		return nil
	},
}

func rollDie() (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(6))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + 1, nil
}
