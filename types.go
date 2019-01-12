package main

import (
	"encoding/json"
)

type postParams map[string]string

// BaseResp x
type BaseResp struct {
	Status string `json:"status"`
}

// BalanceResp x
type BalanceResp struct {
	BaseResp
	Balance `json:"balance"`
}

// Balance x
type Balance struct {
	CurrentBalance json.Number `json:"current_balance"`
	SpentTotal     json.Number `json:"spent_total"`
	CallsTotal     json.Number `json:"calls_total"`
	TimeTotal      json.Number `json:"time_total"`
	SpentToday     json.Number `json:"spent_today"`
	CallsToday     json.Number `json:"calls_today"`
	TimeToday      json.Number `json:"time_today"`
}

// SetCallerIDFilterResp x
type SetCallerIDFilterResp struct {
	BaseResp
	Filtering json.Number `json:"filtering"`
}

// GetCallDataRecord x
type GetCallDataRecord struct {
	BaseResp
	CallDataRecords []CallDataRecord `json:"cdr"`
}

// CallDataRecord x
type CallDataRecord struct {
	Date        string      `json:"date"`
	CallerID    string      `json:"callerid"`
	Destination json.Number `json:"destination"`
	Description string      `json:"description"`
	Account     string      `json:"account"`
	Disposition string      `json:"disposition"`
	Duration    string      `json:"duration"`
	Rate        json.Number `json:"rate"`
	Total       json.Number `json:"total"`
	UniqueID    json.Number `json:"uniqueid"`
}
