package graph_test

import (
	"testing"

	"github.com/planetary-social/go-ssb/fixtures"
	"github.com/planetary-social/go-ssb/refs"
	"github.com/planetary-social/go-ssb/scuttlebutt/graph"
	"github.com/stretchr/testify/require"
)

func TestContacts(t *testing.T) {
	local := fixtures.SomeRefAuthor()

	a := fixtures.SomeRefAuthor()
	b := fixtures.SomeRefAuthor()
	c := fixtures.SomeRefAuthor()

	s := StorageMock{
		contacts: map[string][]refs.Identity{
			local.String(): {
				a,
			},
			a.String(): {
				b,
			},
			b.String(): {
				c,
			},
		},
	}

	g, err := graph.NewSocialGraph(local, graph.MustNewHops(2), s)
	require.NoError(t, err)

	require.Equal(t,
		[]graph.Contact{
			{
				local,
				graph.MustNewHops(0),
			},
			{
				a,
				graph.MustNewHops(1),
			},
			{
				b,
				graph.MustNewHops(2),
			},
		},
		g.Contacts(),
		"social graph should have returned results sorted by distance and filter out nodes that are too far away",
	)

	require.True(t, g.HasContact(local))
	require.True(t, g.HasContact(a))
	require.True(t, g.HasContact(b))
	require.False(t, g.HasContact(c))
}

type StorageMock struct {
	contacts map[string][]refs.Identity
}

func (s StorageMock) GetContacts(node refs.Identity) ([]refs.Identity, error) {
	return s.contacts[node.String()], nil
}
