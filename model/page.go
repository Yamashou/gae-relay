package model

type PageInfo struct {
	HasNextPage     bool    `json:"hasNextPage"`
	StartCursor     *string `json:"startCursor"`
	EndCursor       *string `json:"endCursor"`
}
