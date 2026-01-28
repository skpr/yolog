package main

import (
	"fmt"
	"os"
	"time"

	"github.com/skpr/yolog"
)

func main() {
	logger := yolog.NewLogger("example")
	defer logger.Log(os.Stdout)

	logger.SetAttr("event", "user_signup")
	logger.SetErrorCode(500)
	logger.AddError(fmt.Errorf("An error occurred and it was really bad"))

	// Demonstrates how a struct can be utilised.
	logger.SetAttr("user", struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		Name:  "Jane Doe",
		Email: "jane@example.com",
	})

	// Sleep so we can see a good elapsed number.
	time.Sleep(time.Second)
}
