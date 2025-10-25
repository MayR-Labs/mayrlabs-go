package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/spf13/cobra"
)

var quotes = []string{
	"Code is like humor. When you have to explain it, it's bad. – Cory House",
	"First, solve the problem. Then, write the code. – John Johnson",
	"Experience is the name everyone gives to their mistakes. – Oscar Wilde",
	"In order to be irreplaceable, one must always be different. – Coco Chanel",
	"Knowledge is power. – Francis Bacon",
	"Sometimes it pays to stay in bed on Monday, rather than spending the rest of the week debugging Monday's code. – Dan Salomon",
	"Perfection is achieved not when there is nothing more to add, but rather when there is nothing more to take away. – Antoine de Saint-Exupery",
	"Code never lies, comments sometimes do. – Ron Jeffries",
	"Simplicity is the soul of efficiency. – Austin Freeman",
	"Make it work, make it right, make it fast. – Kent Beck",
	"Before software can be reusable it first has to be usable. – Ralph Johnson",
	"The best error message is the one that never shows up. – Thomas Fuchs",
	"Walking on water and developing software from a specification are easy if both are frozen. – Edward V. Berard",
	"Talk is cheap. Show me the code. – Linus Torvalds",
	"Programs must be written for people to read, and only incidentally for machines to execute. – Harold Abelson",
}

// QuoteCmd displays a random motivational quote
var QuoteCmd = &cobra.Command{
	Use:   "quote",
	Short: "Display a random motivational quote for developers",
	Long:  "Get inspired with a random motivational quote from famous developers and thinkers",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Using math/rand is acceptable for non-security quote selection
		// #nosec G404
		rand.Seed(time.Now().UnixNano())
		quote := quotes[rand.Intn(len(quotes))]
		fmt.Println("\n✨ " + quote + "\n")
		return nil
	},
}
