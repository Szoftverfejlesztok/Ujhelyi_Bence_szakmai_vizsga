package types

type Device struct {
	Id     int    `json:"id,omitempty"`
	Device string `json:"device"`
	Date   string `json:"date,omitempty"`
	State  bool   `json:"state"`
}
