package generator

import (
	"reflect"
	"testing"
	"time"

	"github.com/xarantolus/rfa-launch-bot/generator/scraper"
)

func Test_extractKeywords(t *testing.T) {
	type args struct {
		title       string
		description string
	}
	tests := []struct {
		args         args
		wantKeywords []string
	}{
		{args{
			title: "RFA Staged combustion event", description: `RFA (Rocket Factory Augsburg) will host the staged combustion event. They might reveal some more info on their NewSpace Launcher vehicle.`,
		}, []string{"RFA", "Rocket Factory Augsburg", "Launcher", "NewSpace", "Augsburg"}},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if gotKeywords := extractKeywords(tt.args.title, tt.args.description); !reflect.DeepEqual(gotKeywords, tt.wantKeywords) {
				t.Errorf("extractKeywords() = %v, want %v", gotKeywords, tt.wantKeywords)
			}
		})
	}
}

func Test_describeLiveStream(t *testing.T) {
	const (
		vidID = "9135813491"
		title = "RFA Staged combustion event"
		desc  = "RFA (Rocket Factory Augsburg) will host the staged combustion event. They might reveal some more info on their Launcher vehicle"
	)
	tests := []struct {
		args scraper.LiveVideo
		want string
	}{
		{
			scraper.LiveVideo{
				VideoID:          vidID,
				Title:            title,
				ShortDescription: desc,
				IsLive:           true,
			},
			`Rocket Factory Augsburg is now live on YouTube:

RFA Staged combustion event

#RFA #RocketFactoryAugsburg #Launcher #Augsburg

https://www.youtube.com/watch?v=9135813491`,
		},
		{
			scraper.LiveVideo{
				VideoID:          vidID,
				Title:            title,
				ShortDescription: desc,
				IsUpcoming:       true,
			},
			`Rocket Factory Augsburg live stream starts soon:

RFA Staged combustion event

#RFA #RocketFactoryAugsburg #Launcher #Augsburg

https://www.youtube.com/watch?v=9135813491`,
		},
		{
			scraper.LiveVideo{
				VideoID:          vidID,
				Title:            title,
				ShortDescription: desc,
				IsUpcoming:       true,
				UpcomingInfo: scraper.LiveBroadcastDetails{
					StartTimestamp: time.Now().Add(1*time.Minute + 10*time.Second),
				},
			},
			`Rocket Factory Augsburg live stream starts in about a minute:

RFA Staged combustion event

#RFA #RocketFactoryAugsburg #Launcher #Augsburg

https://www.youtube.com/watch?v=9135813491`,
		},
		{
			scraper.LiveVideo{
				VideoID:          vidID,
				Title:            title,
				ShortDescription: desc,
				IsUpcoming:       true,
				UpcomingInfo: scraper.LiveBroadcastDetails{
					StartTimestamp: time.Now().Add(3*time.Hour + 10*time.Minute),
				},
			},
			`Rocket Factory Augsburg live stream starts in 3 hours:

RFA Staged combustion event

#RFA #RocketFactoryAugsburg #Launcher #Augsburg

https://www.youtube.com/watch?v=9135813491`,
		},
		{
			scraper.LiveVideo{
				VideoID:          vidID,
				Title:            title,
				ShortDescription: desc,
				IsUpcoming:       true,
				UpcomingInfo: scraper.LiveBroadcastDetails{
					StartTimestamp: time.Now().Add(3*time.Hour + 35*time.Minute),
				},
			},
			`Rocket Factory Augsburg live stream starts in 4 hours:

RFA Staged combustion event

#RFA #RocketFactoryAugsburg #Launcher #Augsburg

https://www.youtube.com/watch?v=9135813491`,
		},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := describeLiveStream(&tt.args); got != tt.want {
				t.Errorf("describeLiveStream() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}
