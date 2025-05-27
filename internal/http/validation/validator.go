package validation

import (
	"bytes"
	"context"
	"strings"
)

type Validator interface {
	ValidateStruct(ctx context.Context, s any) error
}

type Errors map[string][]string

func (errs Errors) Error() string {
	buff := bytes.NewBufferString("")

	for field, msgs := range errs {
		buff.WriteString(field + ": " + strings.Join(msgs, ", "))
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}
