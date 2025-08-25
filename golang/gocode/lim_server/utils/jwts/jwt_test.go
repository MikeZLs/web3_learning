package jwts

import (
	"fmt"
	"testing"
)

func TestGenToken(t *testing.T) {
	token, err := GenToken(JwtPayLoad{
		UserID:   1,
		Role:     1,
		Nickname: "lance",
	}, "12345", 8)
	fmt.Println(token, err)
}

func TestParseToken(t *testing.T) {
	payload, err := ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOjEsIm5pY2tuYW1lIjoibGFuY2UiLCJyb2xlIjoxLCJleHAiOjE3NTUxMzYwNjh9.pPfbdqJZjrZ7o9aRa7UpEoGz8Tglaqgw8YiOesljzpk", "12345")
	fmt.Println(payload, err)
}
