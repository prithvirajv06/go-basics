package main

type GBUser struct {
	GOBID    string   `json:"gobid"`
	Name     string   `json:"name"`
	Password string   `json:"password"`
	Email    string   `json:"email"`
	Role     []GBRole `json:"role"`
}

type GBRole struct {
	ID       string `json:"id"`
	RoleName string `json:"role_name"`
}

type GBCommongResponse struct {
	Token   string `json:"tokren"`
	Message string `json:"message"`
}

type JwtClaims struct {
	Sub  string `json:"sub"`
	Name string `json:"name"`
	Iat  int64  `json:"iat"`
	Exp  int64  `json:"exp"`
}

// Valid implements jwt.Claims.
func (*JwtClaims) Valid() error {
	return nil
}
