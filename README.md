# rfa-launch-bot
This is a [Twitter bot](https://twitter.com/wenlauncherbot) that tries to (re)tweet interesting stuff about [Rocket Factory Augsburg](https://www.rfa.space/).

#### Retweets
It searches tweets from the following sources:
* Searches for certain [keywords](collector/search_stream.go)
* Location-tagged tweets around location(s) they use
  * [Augsburg](http://bboxfinder.com/#48.241138,10.656738,48.553887,11.188202)
  * [Andøya Space Center](http://bboxfinder.com/#68.856583,14.864502,69.377411,16.677246)
  * [Estrange Space Center](http://bboxfinder.com/#67.798869,20.041809,68.983031,21.676025)
* Accounts and lists the bot follows

It then checks for a bunch of keywords and retweets matching tweets. If it finds links that is has retweeted within the last 12 hours, it will not retweet them again.

#### Tweets
The bot can also send tweets on its own under certain conditions:
* A YouTube live stream starts on the [official RFA channel](https://www.youtube.com/channel/UC6PsS67tBgDr5w22ZZSgI9w)


-----

If you have any questions or suggestions, please feel free to open an issue or contact the bot account directly.
