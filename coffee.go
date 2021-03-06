package buna

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"golang.org/x/crypto/ssh/terminal"
)

type coffee struct {
	name    string
	roaster string
	region  string
	variety string
	method  string
	decaf   bool
}

// Returns the added coffee
func addCoffee(ctx context.Context, db DB) (coffee, error) {
	fmt.Println("Adding new coffee (Enter # to quit):")
	fmt.Print("Enter coffee name: ")
	name, quit := validateStrInput(quitStr, false, nil, nil)
	if quit {
		fmt.Println(quitMsg)
		return coffee{}, nil
	}

	fmt.Print("Enter roaster/producer name: ")
	roaster, quit := validateStrInput(quitStr, false, nil, nil)
	if quit {
		fmt.Println(quitMsg)
		return coffee{}, nil
	}

	fmt.Print("Enter origin/region (Format: Region, Country): ")
	region, quit := validateStrInput(quitStr, true, nil, nil)
	if quit {
		fmt.Println(quitMsg)
		return coffee{}, nil
	}

	fmt.Print("Enter variety (Format: Variety 1, Variety 2, ...): ")
	variety, quit := validateStrInput(quitStr, true, nil, nil)
	if quit {
		fmt.Println(quitMsg)
		return coffee{}, nil
	}

	fmt.Print("Enter processing method: ")
	method, quit := validateStrInput(quitStr, true, nil, nil)
	if quit {
		fmt.Println(quitMsg)
		return coffee{}, nil
	}

	fmt.Print("Is decaf (true or false): ")
	decaf, quit := validateBoolInput(quitStr, true)
	if quit {
		fmt.Println(quitMsg)
		return coffee{}, nil
	}

	newCoffee := coffee{
		name:    name,
		roaster: roaster,
		region:  region,
		variety: variety,
		method:  method,
		decaf:   decaf,
	}

	if err := db.insertCoffee(ctx, newCoffee); err != nil {
		return coffee{}, fmt.Errorf("buna: coffee: failed to insert coffee: %w", err)
	}

	fmt.Println("Added coffee successfully")
	return newCoffee, nil
}

func retrieveCoffee(ctx context.Context, db DB) error {
	options := map[int]string{
		0: "Retrieve coffees ordered by last added",
		// 1: "Retrieve coffee by name",
		// 2: "Retrieve coffees ordered alphabetically",
		// 3: "Retrieve coffees by origin",
		// 4: "Retrieve coffees by roaster",
		// 5: "Retrieve coffees by processing method",
		// 6: "Retrieve decaf coffees ordered by last added",
		// 7: "Retrieve decaf coffees ordered alphabetically",
	}

	fmt.Println("Retrieving coffee (Enter # to quit):")
	if err := displayIntOptions(options); err != nil {
		return fmt.Errorf("buna: coffee: failed to display int options: %w", err)
	}

	selection, quit, err := getIntSelection(options, quitStr)
	if err != nil {
		return fmt.Errorf("buna: coffee: failed to get int selection: %w", err)
	}
	if quit {
		fmt.Println(quitMsg)
		return nil
	}

	if err := runRetrieveCoffeeSelection(ctx, selection, db); err != nil {
		return fmt.Errorf("buna: coffee: failed to run the retrieve selection: %w", err)
	}

	return nil
}

func runRetrieveCoffeeSelection(ctx context.Context, selection int, db DB) error {
	switch selection {
	case 0:
		if err := displayCoffeesByLastAdded(ctx, db); err != nil {
			return fmt.Errorf("buna: coffee: failed to display coffees by last added: %w", err)
		}
	default:
		return errors.New("buna: coffee: invalid retrieve selection")
	}
	return nil
}

// Promts user for an optional limit.
func displayCoffeesByLastAdded(ctx context.Context, db DB) error {
	const defaultDisplayAmount = 15
	const maxDisplayAmount = 60

	fmt.Println("Displaying coffees by last added (Enter # to quit):")
	fmt.Print("Enter a limit for the number of coffees to display: ")
	limit, quit := validateIntInput(quitStr, true, 1, maxDisplayAmount, []int{})
	if quit {
		fmt.Println(quitMsg)
		return nil
	}

	if limit == 0 {
		limit = defaultDisplayAmount
	}

	coffees, err := db.getCoffeesByLastAdded(ctx, limit)
	if err != nil {
		return fmt.Errorf("buna: coffee: failed to get coffees by last added: %w", err)
	}

	t := table.NewWriter()

	t.AppendHeader(table.Row{
		"Name",
		"Roaster",
		"Region/Origin",
		"Variety",
		"Processing method",
		"Decaf",
	})

	for _, coffee := range coffees {
		t.AppendRow(table.Row{
			coffee.name,
			coffee.roaster,
			coffee.region,
			coffee.variety,
			coffee.method,
			coffee.decaf,
		})
		t.AppendSeparator()
	}

	terminalWidth, _, err := terminal.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("buna: coffee: failed to get terminal width: %w", err)
	}
	t.SetAllowedRowLength(terminalWidth)

	t.SetOutputMirror(os.Stdout)
	t.Render()

	return nil
}
