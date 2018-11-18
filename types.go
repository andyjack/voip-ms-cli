package main

import (
	"encoding/json"
)

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
