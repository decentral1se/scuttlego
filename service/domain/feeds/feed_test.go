package feeds_test

import (
	"testing"

	"github.com/planetary-social/go-ssb/fixtures"
	"github.com/planetary-social/go-ssb/service/domain/feeds"
	msgcontents "github.com/planetary-social/go-ssb/service/domain/feeds/content"
	"github.com/planetary-social/go-ssb/service/domain/feeds/message"
	"github.com/stretchr/testify/require"
)

func TestAppend(t *testing.T) {
	feed := fixtures.SomeRefFeed()
	author := fixtures.SomeRefAuthor()

	msg1 := message.MustNewMessage(
		fixtures.SomeRefMessage(),
		nil,
		message.MustNewSequence(1),
		author,
		feed,
		fixtures.SomeTime(),
		fixtures.SomeContent(),
		fixtures.SomeRawMessage(),
	)

	prevId := msg1.Id()

	msg2 := message.MustNewMessage(
		fixtures.SomeRefMessage(),
		&prevId,
		message.MustNewSequence(2),
		author,
		feed,
		fixtures.SomeTime(),
		fixtures.SomeContent(),
		fixtures.SomeRawMessage(),
	)

	f := feeds.NewFeed(nil)

	err := f.AppendMessage(msg1)
	require.NoError(t, err)

	err = f.AppendMessage(msg2)
	require.NoError(t, err)

	msgs, contacts, pubs := f.PopForPersisting()
	require.Len(t, msgs, 2)
	require.Len(t, contacts, 0)
	require.Len(t, pubs, 0)
}

func TestAppendMessageWithKnownContent(t *testing.T) {
	author := fixtures.SomeRefAuthor()

	someIdentity := fixtures.SomeRefAuthor()

	testCases := []struct {
		Name             string
		Content          msgcontents.KnownMessageContent
		ExpectedContacts []feeds.ContactToSave
		ExpectedPubs     []msgcontents.Pub
	}{
		{
			Name: "contact",
			Content: msgcontents.MustNewContact(
				someIdentity,
				msgcontents.ContactActionFollow,
			),
			ExpectedContacts: []feeds.ContactToSave{
				feeds.NewContactToSave(
					author,
					msgcontents.MustNewContact(
						someIdentity,
						msgcontents.ContactActionFollow,
					),
				),
			},
		},
		{
			Name: "pub",
			Content: msgcontents.MustNewPub(
				someIdentity,
				"host",
				1234,
			),
			ExpectedPubs: []msgcontents.Pub{
				msgcontents.MustNewPub(
					someIdentity,
					"host",
					1234,
				),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			msg := message.MustNewMessage(
				fixtures.SomeRefMessage(),
				nil,
				message.MustNewSequence(1),
				author,
				fixtures.SomeRefFeed(),
				fixtures.SomeTime(),
				testCase.Content,
				fixtures.SomeRawMessage(),
			)

			f := feeds.NewFeed(nil)

			err := f.AppendMessage(msg)
			require.NoError(t, err)

			msgs, contacts, pubs := f.PopForPersisting()
			require.Len(t, msgs, 1)
			require.Equal(t, testCase.ExpectedContacts, contacts)
			require.Equal(t, testCase.ExpectedPubs, pubs)
		})
	}
}
