package apiaccess

import (
	"errors"
	"fmt"

	"github.com/newrelic/newrelic-client-go/internal/http"
)

// APIAccessKey represents a New Relic API access ingest or user key.
// type APIAccessKey struct {
// 	ID         string `json:"id,omitempty"`
// 	Key        string `json:"key,omitempty"`
// 	Name       string `json:"name,omitempty"`
// 	Notes      string `json:"notes,omitempty"`
// 	Type       string `json:"type,omitempty"`
// 	AccountID  int    `json:"accountId,omitempty"`
// 	IngestType string `json:"ingestType,omitempty"`
// 	UserID     string `json:"userId,omitempty"`
// }

// apiAccessKeyCreateResponse represents the JSON response returned from creating key(s).
type apiAccessKeyCreateResponse struct {
	APIAccessCreateKeys ApiAccessCreateKeyResponse `json:"apiAccessCreateKeys"`
}

// apiAccessKeyUpdateResponse represents the JSON response returned from updating key(s).
type apiAccessKeyUpdateResponse struct {
	APIAccessUpdateKeys ApiAccessUpdateKeyResponse `json:"apiAccessUpdateKeys"`
}

// apiAccessKeyGetResponse represents the JSON response returned from getting an access key.
type apiAccessKeyGetResponse struct {
	Actor struct {
		APIAccess struct {
			Key ApiAccessKey `json:"key,omitempty"`
		} `json:"apiAccess"`
	} `json:"actor"`
	http.GraphQLErrorResponse
}

type apiAccessKeySearchResponse struct {
	Actor struct {
		APIAccess struct {
			KeySearch ApiAccessKeySearchResult `json:"keySearch,omitempty"`
		} `json:"apiAccess"`
	} `json:"actor"`
	http.GraphQLErrorResponse
}

// apiAccessKeyDeleteResponse represents the JSON response returned from creating key(s).
type apiAccessKeyDeleteResponse struct {
	APIAccessDeleteKeys ApiAccessDeleteKeyResponse `json:"apiAccessDeleteKeys"`
}

const (
	graphqlAPIAccessKeyBaseFields = `
		id
		key
		name
		notes
		type
		... on ApiAccessIngestKey {
			id
			name
			accountId
			ingestType
			key
			notes
			type
		}
		... on ApiAccessUserKey {
			id
			name
			accountId
			key
			notes
			type
			userId
		}
		... on ApiAccessKey {
			id
			name
			key
			notes
			type
		}`

	graphqlAPIAccessCreateKeyFields = `createdKeys {` + graphqlAPIAccessKeyBaseFields + `}`

	graphqlAPIAccessUpdatedKeyFields = `updatedKeys {` + graphqlAPIAccessKeyBaseFields + `}`

	graphqlAPIAccessKeyErrorFields = `errors {
		  message
		  type
		  ... on ApiAccessIngestKeyError {
			id
			ingestErrorType: errorType
			accountId
			ingestType
			message
			type
		  }
		  ... on ApiAccessKeyError {
			message
			type
		  }
		  ... on ApiAccessUserKeyError {
			id
			accountId
			userErrorType: errorType
			message
			type
			userId
		  }
		}
	`

	apiAccessKeyCreateKeys = `mutation($keys: ApiAccessCreateInput!) {
			apiAccessCreateKeys(keys: $keys) {` + graphqlAPIAccessCreateKeyFields + graphqlAPIAccessKeyErrorFields + `
		}}`

	apiAccessKeyGetKey = `query($id: ID!, $keyType: ApiAccessKeyType!) {
		actor {
			apiAccess {
				key(id: $id, keyType: $keyType) {` + graphqlAPIAccessKeyBaseFields + `}}}}`

	apiAccessKeySearch = `query($query: ApiAccessKeySearchQuery!) {
		actor {
			apiAccess {
				keySearch(query: $query) {
					keys {` + graphqlAPIAccessKeyBaseFields + `}
				}}}}`

	apiAccessKeyUpdateKeys = `mutation($keys: ApiAccessUpdateInput!) {
			apiAccessUpdateKeys(keys: $keys) {` + graphqlAPIAccessUpdatedKeyFields + graphqlAPIAccessKeyErrorFields + `
		}}`

	apiAccessKeyDeleteKeys = `mutation($keys: ApiAccessDeleteInput!) {
			apiAccessDeleteKeys(keys: $keys) {
				deletedKeys {
					id
				}` + graphqlAPIAccessKeyErrorFields + `}}`
)

// CreateAPIAccessKeysMutation create keys. You can create keys for multiple accounts at once.
func (a *APIAccess) CreateAPIAccessKeysMutation(keys ApiAccessCreateInput) ([]ApiAccessKey, error) {
	vars := map[string]interface{}{
		"keys": keys,
	}

	resp := apiAccessKeyCreateResponse{}

	if err := a.client.NerdGraphQuery(apiAccessKeyCreateKeys, vars, &resp); err != nil {
		return nil, err
	}

	if len(resp.APIAccessCreateKeys.Errors) > 0 {
		return nil, errors.New(formatAPIAccessKeyMutationErrors(resp.APIAccessCreateKeys.Errors))
	}

	return resp.APIAccessCreateKeys.CreatedKeys, nil
}

// GetAPIAccessKeyMutation returns a single API access key.
func (a *APIAccess) GetAPIAccessKeyMutation(keyID string, keyType ApiAccessKeyType) (*ApiAccessKey, error) {
	vars := map[string]interface{}{
		"id":      keyID,
		"keyType": keyType,
	}

	resp := apiAccessKeyGetResponse{}

	if err := a.client.NerdGraphQuery(apiAccessKeyGetKey, vars, &resp); err != nil {
		return nil, err
	}

	if resp.Errors != nil {
		return nil, errors.New(resp.Error())
	}

	return &resp.Actor.APIAccess.Key, nil
}

// SearchAPIAccessKeys returns the relevant keys based on search criteria. Returns keys are scoped to the current user.
func (a *APIAccess) SearchAPIAccessKeys(params ApiAccessKeySearchQuery) ([]ApiAccessKey, error) {
	vars := map[string]interface{}{
		// "scope": params.Scope,
		// "types": params.Types,
		"query": params,
	}

	resp := apiAccessKeySearchResponse{}

	if err := a.client.NerdGraphQuery(apiAccessKeySearch, vars, &resp); err != nil {
		return nil, err
	}

	if resp.Errors != nil {
		return nil, errors.New(resp.Error())
	}

	return resp.Actor.APIAccess.KeySearch.Keys, nil
}

// UpdateAPIAccessKeyMutation updates keys. You can update keys for multiple accounts at once.
func (a *APIAccess) UpdateAPIAccessKeyMutation(keys ApiAccessUpdateInput) ([]ApiAccessKey, error) {
	vars := map[string]interface{}{
		"keys": keys,
	}

	resp := apiAccessKeyUpdateResponse{}

	if err := a.client.NerdGraphQuery(apiAccessKeyUpdateKeys, vars, &resp); err != nil {
		return nil, err
	}

	if len(resp.APIAccessUpdateKeys.Errors) > 0 {
		return nil, errors.New(formatAPIAccessKeyMutationErrors(resp.APIAccessUpdateKeys.Errors))
	}

	return resp.APIAccessUpdateKeys.UpdatedKeys, nil
}

// DeleteAPIAccessKeyMutation deletes one or more keys.
func (a *APIAccess) DeleteAPIAccessKeyMutation(keys ApiAccessDeleteInput) ([]ApiAccessDeletedKey, error) {
	vars := map[string]interface{}{
		"keys": keys,
	}

	resp := apiAccessKeyDeleteResponse{}

	if err := a.client.NerdGraphQuery(apiAccessKeyDeleteKeys, vars, &resp); err != nil {
		return nil, err
	}

	if len(resp.APIAccessDeleteKeys.Errors) > 0 {
		return nil, errors.New(formatAPIAccessKeyMutationErrors(resp.APIAccessDeleteKeys.Errors))
	}

	return resp.APIAccessDeleteKeys.DeletedKeys, nil
}

func formatAPIAccessKeyMutationErrors(errors []ApiAccessKeyError) string {
	errorString := ""
	for _, e := range errors {
		errorString += fmt.Sprintf("%v: %v\n", e.Type, e.Message)
	}
	return errorString
}
