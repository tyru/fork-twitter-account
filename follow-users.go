package main

import (
	"bufio"
	"flag"
	"fork-account/util"
	"os"
	"strconv"
	"strings"
)

func main() {
	os.Exit(Main())
}

func Main() int {
	var writeFile string
	flag.StringVar(&writeFile, "w", "", "this file will be changed (followed users are removed from this list)")
	flag.Parse()

	input := os.Stdin
	if writeFile != "" {
		perm := os.FileMode(0666)
		if fi, err := os.Stat(writeFile); err == nil {
			perm = fi.Mode()
		}
		var err error
		input, err = os.OpenFile(writeFile, os.O_RDWR, perm)
		if err != nil {
			util.LogError(err.Error())
			return 10
		}
		defer input.Close()
	}

	users, err := getUsersFrom(input)
	if err != nil {
		util.LogError(err.Error())
		return 11
	}
	util.LogInfo("Will follow " + strconv.Itoa(len(users)) + " users ...")
	tw, err := util.NewTwitter()
	if err != nil {
		util.LogError(err.Error())
		return 12
	}
	followFail := false
	for len(users) > 0 {
		util.LogInfo("Follow " + users[0])
		if err := tw.FollowUsers(users[0]); err != nil {
			followFail = true
			util.LogError(err.Error())
			break
		}
		users = users[1:]
	}

	// Write to writeFile
	if writeFile != "" {
		if _, err := input.Seek(0, 0); err != nil {
			util.LogError(err.Error())
			return 13
		}
		if err := input.Truncate(0); err != nil {
			util.LogError(err.Error())
			return 14
		}
		if _, err := input.WriteString(strings.Join(users, "\n")); err != nil {
			util.LogError(err.Error())
			return 15
		}
	}

	if followFail {
		return 16
	}
	return 0
}

func getUsersFrom(r *os.File) ([]string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	users := make([]string, 0, 16*1024)
	for scanner.Scan() {
		users = append(users, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
