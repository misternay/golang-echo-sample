package models

type Users struct {
	Id           int      `json:"id,omitempty"`
	Username     string   `json:"username,omitempty"`
	Fullname     string   `json:"fullname,omitempty"`
	ChildLeftId  int      `json:"child_left_id,omitempty"`
	ChildRightId int      `json:"child_right_id,omitempty"`
	Direction    string   `json:"direction,omitempty" sql:"-"`
	Password     string   `json:"-"`
	ChildLeft    []*Users `json:"child_left,omitempty" sql:"-"`
	ChildRight   []*Users `json:"child_right,omitempty" sql:"-"`
}
