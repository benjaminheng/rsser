package instagram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/feeds"
)

var userPageDataRegexp = regexp.MustCompile(`(\{.+"entry_data".+\});</script>`)

func GetUserFeed(ctx context.Context, username string) (feed *feeds.Feed, err error) {
	pageData, err := getUserPageData(ctx, username)
	if err != nil {
		return nil, err
	}
	if len(pageData.EntryData.ProfilePages) == 0 {
		return nil, errors.New("malformed page data for profilepages")
	}
	user := pageData.EntryData.ProfilePages[0].GraphQL.User
	feed = &feeds.Feed{
		Title:       fmt.Sprintf("%s (@%s)", user.FullName, user.Username),
		Link:        &feeds.Link{Href: "https://www.instagram.com/" + user.Username},
		Description: user.Biography,
		Author:      &feeds.Author{Name: user.FullName},
	}
	timelineMedia := user.EdgeOwnerToTimelineMedia.Edges
	feed.Items = make([]*feeds.Item, 0, len(timelineMedia))
	for _, media := range timelineMedia {
		if len(media.Node.EdgeMediaToCaption.Edges) == 0 {
			return nil, errors.New("malformed page data for post title")
		}
		description := media.Node.EdgeMediaToCaption.Edges[0].Node.Text
		title := description
		if len(title) > 256 {
			title = title[:256]
		}
		createdAt := time.Unix(media.Node.TakenAtTimestamp, 0)
		item := &feeds.Item{
			Title:       title,
			Link:        &feeds.Link{Href: "https://www.instagram.com/p/" + media.Node.Shortcode},
			Description: description,
			Author:      &feeds.Author{Name: user.FullName},
			Created:     createdAt,
			Enclosure:   &feeds.Enclosure{Url: media.Node.DisplayURL, Type: "image/jpeg", Length: "0"},
		}
		feed.Add(item)
	}
	return feed, nil
}

func getUserPageData(ctx context.Context, username string) (result *UserPageData, err error) {
	url := fmt.Sprintf("https://www.instagram.com/%s/", username)
	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return result, errors.New("got non-200 response when fetching instagram page")
	}
	b, err := ioutil.ReadAll(resp.Body)

	matches := userPageDataRegexp.FindSubmatch(b)
	if len(matches) < 2 {
		return result, errors.New("unable to parse page data from instagram response")
	}

	userPageData := &UserPageData{}
	err = json.Unmarshal(matches[1], userPageData)
	if err != nil {
		return result, err
	}
	return userPageData, nil
}
