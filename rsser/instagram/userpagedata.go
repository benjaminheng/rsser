package instagram

type EdgeMediaToCaptionEdgeNode struct {
	Text string `json:"text"`
}

type EdgeMediaToCaptionEdge struct {
	Node EdgeMediaToCaptionEdgeNode `json:"node"`
}

type EdgeMediaToCaption struct {
	Edges []EdgeMediaToCaptionEdge `json:"edges"`
}

type EdgeOwnerToTimelineMediaEdgeNode struct {
	EdgeMediaToCaption EdgeMediaToCaption `json:"edge_media_to_caption"`
	Shortcode          string             `json:"shortcode"`
	DisplayURL         string             `json:"display_url"`
	TakenAtTimestamp   int64              `json:"taken_at_timestamp"`
}

type EdgeOwnerToTimelineMediaEdge struct {
	Node EdgeOwnerToTimelineMediaEdgeNode `json:"node"`
}

type EdgeOwnerToTimelineMedia struct {
	Edges []EdgeOwnerToTimelineMediaEdge `json:"edges"`
}

type GraphQLUser struct {
	Biography string `json:"biography"`
	FullName  string `json:"full_name"`
	Username  string `json:"username"`

	EdgeOwnerToTimelineMedia EdgeOwnerToTimelineMedia `json:"edge_owner_to_timeline_media"`
}

type GraphQL struct {
	User GraphQLUser `json:"user"`
}

type ProfilePage struct {
	GraphQL GraphQL `json:"graphql"`
}

type EntryData struct {
	ProfilePages []ProfilePage `json:"ProfilePage"`
}

// UserPageData represents the object Instagram uses to render the page.
type UserPageData struct {
	EntryData `json:"entry_data"`
}
