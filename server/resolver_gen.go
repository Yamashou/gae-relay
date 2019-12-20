package server

import (
	"context"

	"golang.org/x/xerrors"

	"github.com/Yamashou/gae-relay/datastore"

	"github.com/Yamashou/gae-relay/gqlgen"
	"github.com/Yamashou/gae-relay/model"
)

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	DatastoreClient *datastore.Client
}

func NewResolver(c *datastore.Client) *Resolver {
	return &Resolver{DatastoreClient: c}
}

func (r *Resolver) Event() gqlgen.EventResolver {
	return &eventResolver{r}
}

func (r *Resolver) Mutation() gqlgen.MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() gqlgen.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) User() gqlgen.UserResolver {
	return &userResolver{r}
}

type eventResolver struct{ *Resolver }

func (r *eventResolver) ID(_ context.Context, obj *model.Event) (string, error) {
	return string(obj.ID), nil
}

func (r *eventResolver) UserID(_ context.Context, obj *model.Event) (string, error) {
	return string(obj.UserID), nil
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateUser(ctx context.Context, input gqlgen.CreateUserInput) (*gqlgen.CreateUserPayload, error) {
	user := model.NewUser(input.Name)
	if err := r.DatastoreClient.Save(ctx, user); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return &gqlgen.CreateUserPayload{
		User: user,
	}, nil
}

func (r *mutationResolver) CreateEvent(ctx context.Context, input gqlgen.CreateEventInput) (*gqlgen.CreateEventPayload, error) {
	if !model.IsUserID(input.UserID) {
		return nil, xerrors.Errorf("CreateEvent: not a correct UserID")
	}
	event := model.NewEvent(model.UserID(input.UserID), input.Description)
	if err := r.DatastoreClient.Save(ctx, event); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return &gqlgen.CreateEventPayload{
		Event: event,
	}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Node(ctx context.Context, id string) (gqlgen.Node, error) {
	switch {
	case model.IsUserID(id):
		user, err := r.DatastoreClient.GetUser(ctx, model.UserID(id))
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		return user, nil
	case model.IsEventID(id):
		event, err := r.DatastoreClient.GetEvent(ctx, model.EventID(id))
		if err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}

		return event, nil
	default:
		return nil, xerrors.Errorf("Node: not correct ID")
	}
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	if !model.IsUserID(id) {
		return nil, xerrors.Errorf("CreateEvent: not a correct UserID")
	}

	user, err := r.DatastoreClient.GetUser(ctx, model.UserID(id))
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return user, nil
}

func (r *queryResolver) Users(ctx context.Context, first int, after *string) (*model.UserConnection, error) {
	var cursor string
	if after != nil {
		cursor = *after
	}
	userConnection, err := r.DatastoreClient.GetUsers(ctx, first, cursor)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return userConnection, nil
}

type userResolver struct{ *Resolver }

func (r *userResolver) ID(_ context.Context, obj *model.User) (string, error) {
	return string(obj.ID), nil
}

func (r *userResolver) Events(ctx context.Context, obj *model.User, first int, after *string) (*model.EventConnection, error) {
	var cursor string
	if after != nil {
		cursor = *after
	}
	eventConnection, err := r.DatastoreClient.GetEvents(ctx, obj.ID, first, cursor)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return eventConnection, nil
}
