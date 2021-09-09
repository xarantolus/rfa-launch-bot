package matcher

import "github.com/dghubble/go-twitter/twitter"

// isReplyToOtherUser returns if the given tweet is a reply to another user
func (m *Matcher) isReplyToOtherUser(t *twitter.Tweet) bool {
	// Not a reply at all?
	if t.QuotedStatusID != 0 {
		return false
	}
	if t.User == nil || t.InReplyToUserID == 0 {
		return false
	}

	// Reply to another user?
	if t.User.ID != t.InReplyToUserID {
		return true
	}

	t, _, err := m.client.Statuses.Show(t.InReplyToStatusID, &twitter.StatusShowParams{
		TweetMode: "extended",
	})
	if err != nil {
		// If something goes wrong, we just assume it is a reply;
		// that way the tweet is ignored
		return true
	}

	// Check if there's another user further up in the reply stack
	return m.isReplyToOtherUser(t)
}
