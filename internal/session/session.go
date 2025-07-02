package session

import (
	"fmt"
	"os"
	"time"

	"github.com/tryoasnafi/gate/internal/crypto"
	"github.com/tryoasnafi/gate/internal/store.go"
)

var (
	cache      []byte
	expiration time.Time
)

func Create(password []byte) {
	cache = password
	expiration = time.Now().Add(5 * time.Minute) // valid for 5 mins
}

func Require() []byte {
	if cache == nil || time.Now().After(expiration) {
		err := crypto.PromptPasswordWithRetry("Enter master password: ", func(pwd []byte) error {
			if err := store.TryDecryptOnly(pwd); err != nil {
				return err
			}
			Create(pwd)
			return nil
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	return cache
}

func Clear() {
	cache = nil
}
