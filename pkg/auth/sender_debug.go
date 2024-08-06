package auth

import (
	"context"
	"fmt"
)

type DebugSender struct {
}

func (s *DebugSender) SendLoginCode(_ context.Context, _ Email, code LoginCode) error {
	fmt.Println("enter code:", code)
	return nil
}
