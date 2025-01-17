package queries

import (
	"github.com/boreq/errors"
	"github.com/planetary-social/scuttlego/internal"
	"github.com/planetary-social/scuttlego/service/app/common"
	"github.com/planetary-social/scuttlego/service/domain/feeds/message"
	"github.com/planetary-social/scuttlego/service/domain/identity"
	"github.com/planetary-social/scuttlego/service/domain/refs"
)

type PublishedLog struct {
	// Only messages with a sequence greater than the sequence of a message
	// pointed to by the given receive log sequence are returned. Given receive
	// log sequence must point to a message published by the current identity.
	// If not such message is known nil should be passed.
	LastSeq *common.ReceiveLogSequence
}

type PublishedLogHandler struct {
	feedRepository       FeedRepository
	receiveLogRepository ReceiveLogRepository
	feed                 refs.Feed
}

func NewPublishedLogHandler(
	feedRepository FeedRepository,
	receiveLogRepository ReceiveLogRepository,
	local identity.Public,
) (*PublishedLogHandler, error) {
	localRef, err := refs.NewIdentityFromPublic(local)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a public identity")
	}

	return &PublishedLogHandler{
		feedRepository:       feedRepository,
		receiveLogRepository: receiveLogRepository,
		feed:                 localRef.MainFeed(),
	}, nil
}

func (h *PublishedLogHandler) Handle(query PublishedLog) ([]LogMessage, error) {
	var startSeq *message.Sequence

	if query.LastSeq != nil {
		msg, err := h.receiveLogRepository.GetMessage(*query.LastSeq)
		if err != nil {
			return nil, errors.Wrap(err, "failed to find a message in the receive log")
		}

		if !msg.Feed().Equal(h.feed) {
			return nil, errors.New("start sequence doesn't point to a message from this feed")
		}

		startSeq = internal.Ptr(msg.Sequence().Next())
	}

	msgs, err := h.feedRepository.GetMessages(h.feed, startSeq, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error getting messages")
	}

	var result []LogMessage
	for _, msg := range msgs {
		receiveLogSequences, err := h.receiveLogRepository.GetSequences(msg.Id())
		if err != nil {
			return nil, errors.Wrap(err, "failed to look up message sequences")
		}

		result = append(result, LogMessage{
			Message:  msg,
			Sequence: receiveLogSequences[0],
		})
	}

	return result, nil
}
