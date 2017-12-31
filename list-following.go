package main

import (
	"fmt"
	"os"

	"github.com/tyru/fork-twitter-account/util"
)

func main() {
	os.Exit(Main())
}

func Main() int {
	if len(os.Args) < 2 {
		util.LogError("Please specify screen name which you want to list the following of")
		return 10
	}
	name := os.Args[1]
	tw, err := util.NewTwitter()
	if err != nil {
		util.LogError(err.Error())
		return 11
	}
	util.LogInfo("Getting following users of '" + name + "' ...")
	for result := range tw.GetFollowingUsers(name) {
		if result.Error != nil {
			util.LogError(result.Error.Error())
			return 12
		}
		fmt.Println(result.ScreenName)
	}
	return 0
}
