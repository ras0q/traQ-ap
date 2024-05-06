package v3

import (
	"fmt"
	"net/http"
	"strings"

	ap "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traPtitech/traQ/router/extension/herror"
)

type WebfingerResponse struct {
	Subject string          `json:"subject"`
	Links   []WebfingerLink `json:"links"`
}

type WebfingerLink struct {
	Rel  string `json:"rel"`
	Type string `json:"type"`
	Href ap.IRI `json:"href"`
}

// GetActivityPubWebfinger GET /ap/webfinger?resource=:resource
// NOTE: this endpoint need to be redirected from /.well-known/webfinger
func (h *Handlers) GetActivityPubWebfinger(c echo.Context) error {
	resource := c.QueryParam("resource")
	if !strings.HasPrefix(resource, "acct:") {
		return herror.BadRequest("resource must start with acct:")
	}

	parts := strings.Split(resource[5:], "@")
	if len(parts) != 2 {
		return herror.BadRequest("resource must have username and domain")
	}

	username := parts[0]
	domain := parts[1]

	user, err := h.Repo.GetUserByName(username, false)
	if err != nil {
		return herror.InternalServerError(err)
	}

	if domain != c.Request().Host {
		return herror.BadRequest("domain does not match")
	}

	apiBaseIRI := ap.IRI(fmt.Sprintf("%s://%s/api/v3", c.Scheme(), c.Request().Host))
	webfinger := WebfingerResponse{
		Subject: resource,
		Links: []WebfingerLink{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: apiBaseIRI.AddPath("ap/users").AddPath(user.GetID().String()),
			},
		},
	}

	return c.JSON(http.StatusOK, webfinger)
}

// GetActivityPubUser GET /ap/users/:userID
func (h *Handlers) GetActivityPubUser(c echo.Context) error {
	userID := uuid.FromStringOrNil(c.Param("userID"))
	user, err := h.Repo.GetUser(userID, false)
	if err != nil {
		return herror.InternalServerError(err)
	}

	apiBaseIRI := ap.IRI(fmt.Sprintf("%s://%s/api/v3", c.Scheme(), c.Request().Host))
	userIRI := apiBaseIRI.AddPath("ap/users").AddPath(user.GetID().String())
	username := user.GetName()

	actor := ap.PersonNew(userIRI)
	actor.Name = ap.DefaultNaturalLanguageValue(username)
	actor.PreferredUsername = ap.DefaultNaturalLanguageValue(username)
	actor.Summary = ap.DefaultNaturalLanguageValue("hello")
	actor.Icon = ap.Image{
		Type:      ap.ImageType,
		MediaType: "image/png",
		URL:       apiBaseIRI.AddPath("public/icon").AddPath(username),
	}
	actor.Inbox = userIRI.AddPath("inbox")
	actor.Outbox = userIRI.AddPath("outbox")
	actor.Following = userIRI.AddPath("following")
	actor.Followers = userIRI.AddPath("followers")

	data, err := jsonld.WithContext(
		jsonld.IRI(ap.ActivityBaseURI),
		jsonld.IRI(ap.PublicNS),
	).Marshal(actor)
	if err != nil {
		return herror.InternalServerError(err)
	}

	c.Response().Header().Set(echo.HeaderContentType, "application/activity+json")

	return c.JSONBlob(http.StatusOK, data)
}

// PostActivityPubInbox POST /ap/users/:userID/inbox
func (h *Handlers) PostActivityPubInbox(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "activitypub is not implemented yet")
}

// GetActivityPubOutbox GET /ap/users/:userID/outbox
func (h *Handlers) GetActivityPubOutbox(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "activitypub is not implemented yet")
}

// GetActivityPubFollowing GET /ap/users/:userID/following
func (h *Handlers) GetActivityPubFollowing(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "activitypub is not implemented yet")
}

// GetActivityPubFollowers GET /ap/users/:userID/followers
func (h *Handlers) GetActivityPubFollowers(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "activitypub is not implemented yet")
}
