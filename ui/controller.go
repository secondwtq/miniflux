// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/miniflux/miniflux/config"
	"github.com/miniflux/miniflux/http/handler"
	"github.com/miniflux/miniflux/model"
	"github.com/miniflux/miniflux/reader/feed"
	"github.com/miniflux/miniflux/reader/opml"
	"github.com/miniflux/miniflux/scheduler"
	"github.com/miniflux/miniflux/storage"
)

type tplParams map[string]interface{}

func (t tplParams) Merge(d tplParams) tplParams {
	for k, v := range d {
		t[k] = v
	}

	return t
}

// Controller contains all HTTP handlers for the user interface.
type Controller struct {
	cfg         *config.Config
	store       *storage.Storage
	pool        *scheduler.WorkerPool
	feedHandler *feed.Handler
	opmlHandler *opml.Handler
}

func (c *Controller) getCommonTemplateArgs(ctx *handler.Context) (tplParams, error) {
	user := ctx.LoggedUser()
	builder := c.store.NewEntryQueryBuilder(user.ID)
	builder.WithStatus(model.EntryStatusUnread)

	countUnread, err := builder.CountEntries()
	if err != nil {
		return nil, err
	}

	params := tplParams{
		"menu":              "",
		"user":              user,
		"countUnread":       countUnread,
		"csrf":              ctx.CSRF(),
		"flashMessage":      ctx.FlashMessage(),
		"flashErrorMessage": ctx.FlashErrorMessage(),
	}
	return params, nil
}

// NewController returns a new Controller.
func NewController(cfg *config.Config, store *storage.Storage, pool *scheduler.WorkerPool, feedHandler *feed.Handler, opmlHandler *opml.Handler) *Controller {
	return &Controller{
		cfg:         cfg,
		store:       store,
		pool:        pool,
		feedHandler: feedHandler,
		opmlHandler: opmlHandler,
	}
}
