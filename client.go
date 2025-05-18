package plsort

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/schollz/progressbar/v3"
	"github.com/winebarrel/plsort/api"
	"google.golang.org/api/youtube/v3"
)

type Client struct {
	*youtube.Service
}

func NewClient(ctx context.Context, credsPath string) (*Client, error) {
	svc, err := api.NewYoutube(ctx, credsPath)

	if err != nil {
		return nil, err
	}

	client := &Client{
		Service: svc,
	}

	return client, err
}

func (client *Client) getPlaylistItems(ctx context.Context, playlistId string) ([]*youtube.PlaylistItem, error) {
	items := []*youtube.PlaylistItem{}
	call := client.PlaylistItems.List([]string{"snippet"}).
		Context(ctx).PlaylistId(playlistId).MaxResults(50)

	for {
		resp, err := call.Do()

		if err != nil {
			return nil, err
		}

		items = append(items, resp.Items...)

		if resp.NextPageToken == "" {
			break
		}

		call = call.PageToken(resp.NextPageToken)
	}

	return items, nil
}

func (client *Client) Sort(ctx context.Context, options *SortOptions) error {
	items, err := client.getPlaylistItems(ctx, options.PlaylistId)

	if err != nil {
		return fmt.Errorf("unable to get playlist items: %w", err)
	}

	slices.SortFunc(items, func(a, b *youtube.PlaylistItem) int {
		ord := strings.Compare(a.Snippet.Title, b.Snippet.Title)

		if options.Reverse {
			ord = -ord
		}

		return ord
	})

	bar := progressbar.Default(int64(len(items)))

	for i, item := range items {
		newPos := int64(i)

		if item.Snippet.Position != newPos {
			item.Snippet.Position = newPos
			_, err := client.PlaylistItems.Update([]string{"snippet"}, item).
				Context(ctx).Do()

			if err != nil {
				return fmt.Errorf("unable to update playlist items: %w", err)
			}
		}

		bar.Add(1)
	}

	return nil
}
