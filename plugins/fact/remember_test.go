package fact

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/velour/catbase/bot"
	"github.com/velour/catbase/bot/msg"
	"github.com/velour/catbase/bot/user"
)

func makeMessage(nick, payload string) msg.Message {
	isCmd := strings.HasPrefix(payload, "!")
	if isCmd {
		payload = payload[1:]
	}
	return msg.Message{
		User:    &user.User{Name: nick},
		Channel: "test",
		Body:    payload,
		Command: isCmd,
	}
}

func makePlugin(t *testing.T) (*RememberPlugin, *Factoid, *bot.MockBot) {
	mb := bot.NewMockBot()
	f := New(mb) // for DB table
	p := NewRemember(mb)
	assert.NotNil(t, p)
	return p, f, mb
}

// Test case
func TestCornerCaseBug(t *testing.T) {
	msgs := []msg.Message{
		makeMessage("user1", "I don’t want to personally touch a horse dick."),
		makeMessage("user3", "idk my bff rose?"),
		makeMessage("user2", "!remember user1 touch"),
	}

	p, _, mb := makePlugin(t)

	for _, m := range msgs {
		p.Message(m)
	}
	assert.Len(t, mb.Messages, 1)
	assert.Contains(t, mb.Messages[0], "horse dick")
	q, err := getSingleFact(mb.DB(), "user1 quotes")
	assert.Nil(t, err)
	assert.Contains(t, q.Tidbit, "horse dick")
}
