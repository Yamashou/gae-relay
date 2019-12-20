package datastore

import (
	"context"

	"github.com/Yamashou/gae-relay/model"
	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"
)

type Client struct {
	client datastore.Client
}

func NewClient(c datastore.Client) *Client {
	return &Client{client: c}
}

func (c *Client) Save(ctx context.Context, e interface{}) error {
	bm := boom.FromClient(ctx, c.client)
	if _, err := bm.Put(e); err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

func (c *Client) GetUser(ctx context.Context, id model.UserID) (*model.User, error) {
	bm := boom.FromClient(ctx, c.client)
	e := &model.User{ID: id}
	if err := bm.Get(e); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return e, nil
}

func (c *Client) GetEvent(ctx context.Context, id model.EventID) (*model.Event, error) {
	bm := boom.FromClient(ctx, c.client)
	e := &model.Event{ID: id}
	if err := bm.Get(e); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return e, nil
}

func (c *Client) GetUsers(ctx context.Context, limit int, cursor string) (*model.UserConnection, error) {
	// 次のページが存在するか確認するため1件多く取得する
	limitPlusOne := limit + 1
	bm := boom.FromClient(ctx, c.client)
	q := bm.NewQuery(bm.Kind(model.User{})).
		Limit(limitPlusOne)
	if cursor != "" {
		cur, err := c.client.DecodeCursor(cursor)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		q = q.Start(cur)
	}

	var edges []*model.UserEdge
	it := bm.Run(q)
	for {
		var user model.User
		_, err := it.Next(&user)
		if xerrors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		cursor, err := it.Cursor()
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		edge := &model.UserEdge{
			Cursor: cursor.String(),
			Node:   &user,
		}

		edges = append(edges, edge)
	}

	// このページの最初と最後のCursorを返す
	var startCursor, endCursor string
	if len(edges) > 0 {
		startCursor = edges[0].Cursor
		endCursor = edges[len(edges)-1].Cursor
	}

	// 次のページが存在する場合
	var hasNextPage bool
	if len(edges) == limitPlusOne {
		hasNextPage = true
		// 最後の1件は次のページの存在確認用なので除外する
		edges = edges[:len(edges)-1]
	}

	conn := &model.UserConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			StartCursor: &startCursor,
			EndCursor:   &endCursor,
			HasNextPage: hasNextPage,
		},
	}

	return conn, nil
}

func (c *Client) GetEvents(ctx context.Context, userID model.UserID, limit int, cursor string) (*model.EventConnection, error) {
	// 次のページが存在するか確認するため1件多く取得する
	limitPlusOne := limit + 1
	bm := boom.FromClient(ctx, c.client)
	q := bm.NewQuery(bm.Kind(model.Event{})).
		Filter("UserID =", userID).
		Limit(limitPlusOne)
	if cursor != "" {
		cur, err := c.client.DecodeCursor(cursor)
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		q = q.Start(cur)
	}

	var edges []*model.EventEdge
	it := bm.Run(q)
	for {
		var event model.Event
		_, err := it.Next(&event)
		if xerrors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		cursor, err := it.Cursor()
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		edge := &model.EventEdge{
			Cursor: cursor.String(),
			Node:   &event,
		}

		edges = append(edges, edge)
	}

	// このページの最初と最後のCursorを返す
	var startCursor, endCursor string
	if len(edges) > 0 {
		startCursor = edges[0].Cursor
		endCursor = edges[len(edges)-1].Cursor
	}

	// 次のページが存在する場合
	var hasNextPage bool
	if len(edges) == limitPlusOne {
		hasNextPage = true
		// 最後の1件は次のページの存在確認用なので除外する
		edges = edges[:len(edges)-1]
	}

	conn := &model.EventConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			StartCursor: &startCursor,
			EndCursor:   &endCursor,
			HasNextPage: hasNextPage,
		},
	}

	return conn, nil
}
