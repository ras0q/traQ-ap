package v3

import (
	"fmt"
	"net/http"
	"strings"

	ap "github.com/go-ap/activitypub"
	"github.com/go-ap/jsonld"
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
	Href string `json:"href"`
}

// GetActivityPubWebfinger GET /.well-known/webfinger
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

	if domain != c.Request().Host {
		return herror.BadRequest("domain does not match")
	}

	reqURL := fmt.Sprintf("%s://%s/u/%s", c.Scheme(), c.Request().Host, username)

	webfinger := WebfingerResponse{
		Subject: resource,
		Links: []WebfingerLink{
			{
				Rel:  "self",
				Type: "application/activity+json",
				Href: reqURL,
			},
		},
	}

	return c.JSON(http.StatusOK, webfinger)
}

// GetActivityPubUser GET /u/:username
func (h *Handlers) GetActivityPubUser(c echo.Context) error {
	user, err := h.Repo.GetUserByName(c.Param("username"), false)
	if err != nil {
		return herror.InternalServerError(err)
	}

	username := user.GetName()
	reqURL := fmt.Sprintf("%s://%s/u/%s", c.Scheme(), c.Request().Host, username)

	actor := ap.PersonNew(ap.IRI(reqURL))
	actor.Name = ap.DefaultNaturalLanguageValue(username)
	actor.PreferredUsername = ap.DefaultNaturalLanguageValue(username)
	actor.Summary = ap.DefaultNaturalLanguageValue("hello")
	actor.Icon = ap.IRI(ap.IRI(reqURL + "/icon"))
	actor.Inbox = ap.IRI(reqURL + "/inbox")
	actor.Outbox = ap.IRI(reqURL + "/outbox")
	actor.Following = ap.IRI(reqURL + "/following")
	actor.Followers = ap.IRI(reqURL + "/followers")

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

// PostActivityPubInbox POST /u/:username/inbox
func (h *Handlers) PostActivityPubInbox(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "activitypub is not implemented yet")
}

// GetActivityPubOutbox GET /u/:username/outbox
func (h *Handlers) GetActivityPubOutbox(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "activitypub is not implemented yet")
}

// GetActivityPubFollowing GET /u/:username/following
func (h *Handlers) GetActivityPubFollowing(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "activitypub is not implemented yet")
}

// GetActivityPubFollowers GET /u/:username/followers
func (h *Handlers) GetActivityPubFollowers(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented, "activitypub is not implemented yet")
}
