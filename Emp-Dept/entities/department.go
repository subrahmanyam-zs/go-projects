package entities

type Department struct {
	DeptID   int    `json:"DeptID,omitempty"`
	DeptName string `json:"DeptName,omitempty"`
	FloorNo  int    `json:"FloorNo,omitempty"`
}
