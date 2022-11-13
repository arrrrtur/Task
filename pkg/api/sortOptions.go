package api

type Params struct {
	SortByDateOrder   string `json:"sort_by_date_order,omitempty"`
	SortByAmountOrder string `json:"sort_by_amount_order,omitempty"`
	PerPage           int    `json:"per_page,omitempty"`
	Page              int    `json:"page,omitempty"`
}
