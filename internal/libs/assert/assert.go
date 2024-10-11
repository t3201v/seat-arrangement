package assert

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

func NoError(err error, msg ...string) {
	if err != nil {
		message := strings.Join(msg, " ")
		log.Fatal(err, message)
	}
}

func NoErrorf(err error, msg string, args ...any) {
	if err != nil {
		log.Fatalf("error: %v reason: %v", err, fmt.Sprintf(msg, args))
	}
}

func NotFalse(b bool, msg ...string) {
	if !b {
		message := strings.Join(msg, " ")
		log.Fatal(message)
	}
}

func NotNilf[T any](d *T, msg string, args ...any) {
	if d == nil {
		log.Fatalf("expected not nil %v", fmt.Sprintf(msg, args))
	}
}
