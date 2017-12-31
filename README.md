# fork-twitter-account

List specified account's following users &amp; Follow given users

# Usage

Define required environment variables:

```
$ export TWITTER_CONSUMER_KEY='...'
$ export TWITTER_CONSUMER_SECRET='...'
$ export TWITTER_ACCESS_TOKEN='...'
$ export TWITTER_ACCESS_TOKEN_SECRET='...'
```

And run:

```
$ go run list-following.go | go run follow-users.go
```

But above execution doesn't care abandonment due to rate limit, or and so on.
You can save `not followed users` as `followings.txt` like this:

```
$ go run list-following.go >followings.txt
$ go run follow-users.go -w followings.txt
```

Followed users are removed from `followings.txt` even error response were returned from Twitter API.
So if above `go run follow-users.go -w followings.txt` exits with an error, you can re-run `go run follow-users.go -w followings.txt` again.
