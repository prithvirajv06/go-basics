package main

type GBUser struct {
	GOBID string   `json:"gobid"`
	Name  string   `json:"name"`
	Role  []GBRole `json:"role"`
}

type GBRole struct {
	ID       string `json:"id"`
	RoleName string `json:"role_name"`
}

type GBCommongResponse struct {
	Token string `json:"tokren"`
}
