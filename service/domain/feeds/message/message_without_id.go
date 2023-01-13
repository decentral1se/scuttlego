package message

import (
	"fmt"
	"time"

	"github.com/boreq/errors"
	"github.com/planetary-social/scuttlego/service/domain/feeds/content"
	"github.com/planetary-social/scuttlego/service/domain/refs"
)

type MessageWithoutId struct {
	baseMessageFields
	content content.KnownMessageContent

	raw RawMessage
}

func NewMessageWithoutId(
	previous *refs.Message,
	sequence Sequence,
	author refs.Identity,
	feed refs.Feed,
	timestamp time.Time,
	content content.KnownMessageContent,
	raw RawMessage,
) (MessageWithoutId, error) {
	fields, err := newBaseMessageFields(previous, sequence, author, feed, timestamp)
	if err != nil {
		return MessageWithoutId{}, errors.Wrap(err, "could not create base message fields")
	}

	if content == nil {
		return MessageWithoutId{}, errors.Wrap(err, "content is nil")
	}

	if raw.IsZero() {
		return MessageWithoutId{}, errors.New("zero value of raw message")
	}

	return MessageWithoutId{
		baseMessageFields: fields,
		content:           content,
		raw:               raw,
	}, nil
}

func MustNewMessageWithoutId(
	previous *refs.Message,
	sequence Sequence,
	author refs.Identity,
	feed refs.Feed,
	timestamp time.Time,
	content content.KnownMessageContent,
	raw RawMessage,
) MessageWithoutId {
	msg, err := NewMessageWithoutId(previous, sequence, author, feed, timestamp, content, raw)
	if err != nil {
		panic(err)
	}
	return msg
}

func (m MessageWithoutId) Content() content.KnownMessageContent {
	return m.content
}

func (m MessageWithoutId) Raw() RawMessage {
	return m.raw
}

func (m MessageWithoutId) IsZero() bool {
	return m.raw.IsZero()
}

func (m MessageWithoutId) String() string {
	return fmt.Sprintf("feed=%s sequence=%d previous=%s", m.feed.String(), m.sequence.Int(), m.previous)
}