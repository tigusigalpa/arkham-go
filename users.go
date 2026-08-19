package arkham

import (
	"context"
	"net/url"
)

// UserService provides access to user endpoints (alerts, entities, labels).
type UserService struct {
	client *Client
}

// ListAlerts retrieves all alerts for the authenticated user.
// Path: GET /user/alerts
func (s *UserService) ListAlerts(ctx context.Context) ([]Alert, *ResponseMetadata, error) {
	var out []Alert
	meta, err := s.client.get(ctx, "/user/alerts", nil, &out)
	return out, meta, err
}

// CreateAlert creates a new alert.
// Path: POST /user/alerts
func (s *UserService) CreateAlert(ctx context.Context, req *AlertRequest) (*Alert, *ResponseMetadata, error) {
	var out Alert
	meta, err := s.client.post(ctx, "/user/alerts", req, &out)
	return &out, meta, err
}

// UpdateAlert replaces the configuration of an existing alert.
// Path: PUT /user/alerts/{id}
func (s *UserService) UpdateAlert(ctx context.Context, id int, req *AlertRequest) (*ResponseMetadata, error) {
	return s.client.put(ctx, "/user/alerts/"+intToString(id), req, nil)
}

// DeleteAlert deletes an alert by ID.
// Path: DELETE /user/alerts/{id}
func (s *UserService) DeleteAlert(ctx context.Context, id int) (*ResponseMetadata, error) {
	return s.client.delete(ctx, "/user/alerts/"+intToString(id), nil)
}

// ListEntities retrieves all private entities for the authenticated user.
// Path: GET /user/entities
func (s *UserService) ListEntities(ctx context.Context, includeAddresses *bool) ([]UserEntity, *ResponseMetadata, error) {
	q := url.Values{}
	if includeAddresses != nil {
		if *includeAddresses {
			q.Set("includeAddresses", "true")
		} else {
			q.Set("includeAddresses", "false")
		}
	}
	var out []UserEntity
	meta, err := s.client.get(ctx, "/user/entities", q, &out)
	return out, meta, err
}

// ListLabels retrieves all labels for the authenticated user.
// Path: GET /user/labels
func (s *UserService) ListLabels(ctx context.Context) ([]Label, *ResponseMetadata, error) {
	var out []Label
	meta, err := s.client.get(ctx, "/user/labels", nil, &out)
	return out, meta, err
}

// CreateLabels creates one or more labels for the authenticated user.
// Path: POST /user/labels
func (s *UserService) CreateLabels(ctx context.Context, labels []Label) ([]Label, *ResponseMetadata, error) {
	var out []Label
	meta, err := s.client.post(ctx, "/user/labels", labels, &out)
	return out, meta, err
}
