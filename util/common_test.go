package util

import "testing"

func TestHashTagText(t *testing.T) {
	tests := []struct {
		args []string
		want string
	}{
		{[]string{"RFA", "Rocket Factory Augsburg"}, "#RFA #RocketFactoryAugsburg"},
		{[]string{"RFA Launcher"}, "#RFALauncher"},
	}
	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			if got := HashTagText(tt.args); got != tt.want {
				t.Errorf("HashTagText() = %v, want %v", got, tt.want)
			}
		})
	}
}
