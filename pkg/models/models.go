package models

// Shared data models for API requests and responses

// GenerateResponse represents the API response when generating a new Hide My Email address
type GenerateResponse struct {
	Success bool   `json:"success"`
	Result  struct {
		HME string `json:"hme"`
	} `json:"result"`
	Error interface{} `json:"error,omitempty"`
}

// ReserveResponse represents the API response when reserving a generated email address
type ReserveResponse struct {
	Success bool        `json:"success"`
	Error   interface{} `json:"error,omitempty"`
}

// Email represents a Hide My Email address with its metadata
type Email struct {
	HME             string `json:"hme"`
	Label           string `json:"label"`
	Note            string `json:"note"`
	IsActive        bool   `json:"isActive"`
	CreateTimestamp int64  `json:"createTimestamp"`
}

// ListResponse represents the API response when listing Hide My Email addresses
type ListResponse struct {
	Success bool `json:"success"`
	Result  struct {
		HMEEmails []Email `json:"hmeEmails"`
	} `json:"result"`
	Error interface{} `json:"error,omitempty"`
}

// APIParams contains the query parameters required for all iCloud API requests
type APIParams struct {
	ClientBuildNumber     string `json:"clientBuildNumber"`
	ClientMasteringNumber string `json:"clientMasteringNumber"`
	ClientID              string `json:"clientId"`
	DSID                  string `json:"dsid"`
}
