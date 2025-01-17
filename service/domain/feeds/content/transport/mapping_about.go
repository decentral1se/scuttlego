package transport

import (
	"encoding/json"
	"strings"

	"github.com/boreq/errors"
	"github.com/hashicorp/go-multierror"
	"github.com/planetary-social/scuttlego/service/domain/feeds/content"
	"github.com/planetary-social/scuttlego/service/domain/refs"
)

var aboutMapping = MessageContentMapping{
	Marshal: func(con content.KnownMessageContent) ([]byte, error) {
		return nil, errors.New("not implemented")
	},
	Unmarshal: func(b []byte) (content.KnownMessageContent, error) {
		var t transportAbout

		if err := json.Unmarshal(b, &t); err != nil {
			return nil, errors.Wrap(err, "json unmarshal failed")
		}

		image, err := unmarshalAboutImage(t.Image)
		if err != nil {
			return nil, errors.Wrap(err, "could not create image ref")
		}

		return content.NewAbout(image)
	},
}

func unmarshalAboutImage(j json.RawMessage) (*refs.Blob, error) {
	if len(j) == 0 {
		return nil, nil
	}

	var blobRefString string
	if err := json.Unmarshal(j, &blobRefString); err == nil {
		if blobRefString == "" {
			return nil, nil
		}

		blob, err := refs.NewBlob(blobRefString)
		if err != nil {
			return nil, errors.Wrap(err, "could not create a blob ref")
		}

		return &blob, nil
	}

	mention, err := unmarshalMention(j)
	if err != nil {
		return nil, errors.Wrap(err, "invalid mention")
	}

	blob, err := refs.NewBlob(mention.Link)
	if err != nil {
		return nil, errors.Wrap(err, "could not create a blob ref")
	}

	return &blob, err
}

type transportAbout struct {
	messageContentType // todo this is stupid

	// this may be a plain string with a blob ref in it or a blobLink object
	Image json.RawMessage `json:"image"`
}

func unmarshalMentions(j json.RawMessage) ([]refs.Blob, error) {
	var returnErr error

	if len(j) == 0 {
		return nil, nil
	}

	var mentionsSlice []json.RawMessage
	if err := json.Unmarshal(j, &mentionsSlice); err != nil {
		returnErr = multierror.Append(returnErr, errors.Wrap(err, "slice unmarshal error"))
	} else {
		return unmarshalMentionsSlice(mentionsSlice)
	}

	var mentionsMap map[string]json.RawMessage
	if err := json.Unmarshal(j, &mentionsMap); err != nil {
		returnErr = multierror.Append(returnErr, errors.Wrap(err, "map unmarshal error"))
	} else {
		return nil, nil
	}

	return nil, returnErr
}

func unmarshalMentionsSlice(slice []json.RawMessage) ([]refs.Blob, error) {
	var blobs []refs.Blob
	for _, rawJSON := range slice {
		mention, err := unmarshalMention(rawJSON)
		if err != nil {
			return nil, errors.Wrap(err, "could not unmarshal a blob link")
		}

		if !strings.HasPrefix(mention.Link, "&") {
			continue
		}

		blob, err := refs.NewBlob(mention.Link)
		if err != nil {
			return nil, errors.Wrap(err, "could not create a blob ref")
		}

		blobs = append(blobs, blob)
	}

	return blobs, nil
}

type mention struct {
	Link string `json:"link"`
}

func unmarshalMention(j json.RawMessage) (mention, error) {
	var m mention
	if err := json.Unmarshal(j, &m); err != nil {
		return m, errors.Wrap(err, "could not unmarshal blob link")
	}

	if m.Link == "" {
		return m, errors.New("link is empty")
	}

	return m, nil
}
