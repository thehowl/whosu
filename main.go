package main

import (
	"fmt"
	"os"

	"gopkg.in/thehowl/go-osuapi.v1"
)

const format = `%s (uid: %d)
    %.2fpp - #%d (#%d %s)
    Playcount    | %d
    Total score  | %d
    Ranked score | %d
    Level        | %.4f
    Accuracy     | %.2f%%

`

func main() {
	// Get args
	args := os.Args[1:]
	// If we have no args, display the usage of the command.
	if len(args) == 0 {
		fmt.Println("Usage:")
		fmt.Printf("\twhosu [username]\n")
		return
	}

	// make sure the API key is set
	apiKey := os.Getenv("OSU_API_KEY")
	if apiKey == "" {
		fmt.Println("Please set your osu! API key with the environment variable OSU_API_KEY")
		return
	}

	// Start API client
	c := osuapi.NewClient(apiKey)

	// Get user data
	user, err := c.GetUser(osuapi.GetUserOpts{
		Username: args[0],
		Mode:     osuapi.ModeOsu,
	})
	if err != nil {
		fmt.Printf("An error occurred: %v\n", err)
		return
	}
	if user == nil {
		fmt.Println("user is nil! This should never happen! :o")
		return
	}

	// Print user information following format (see above)
	fmt.Printf(
		format,
		user.Username, user.UserID,
		user.PP, user.Rank, user.CountryRank, user.Country,
		user.Playcount,
		user.TotalScore,
		user.RankedScore,
		user.Level,
		user.Accuracy,
	)

	// Get user best scores
	fmt.Print("Waiting for user best...")
	scores, err := c.GetUserBest(osuapi.GetUserScoresOpts{
		UserID: user.UserID,
		Mode:   osuapi.ModeOsu,
		Limit:  15,
	})
	if err != nil {
		fmt.Println(" failed!")
		fmt.Println(err)
		return
	}
	fmt.Println(" ok!")

	// Print all the scores with cool format.
	for _, score := range scores {
		fmt.Printf("    https://osu.ppy.sh/b/%d %v - %.4fpp\n", score.BeatmapID, score.Mods, score.PP)
	}
}
