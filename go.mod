module github.com/xarantolus/rfa-launch-bot

go 1.17

require (
	github.com/dghubble/go-twitter v0.0.0-20210609183100-2fdbf421508e
	github.com/dghubble/oauth1 v0.7.0
	github.com/docker/go-units v0.4.0
	github.com/xarantolus/jsonextract v1.5.3
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/andybalholm/cascadia v1.3.1 // indirect
	golang.org/x/net v0.0.0-20210916014120-12bc252f5db8 // indirect
)

require (
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/dghubble/sling v1.3.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/icholy/replace v0.5.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/tdewolff/parse/v2 v2.5.19 // indirect
	golang.org/x/text v0.3.7
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	mvdan.cc/xurls/v2 v2.3.0 // indirect
)

replace github.com/dghubble/go-twitter => ./bot/go-twitter
