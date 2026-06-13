package ui

import (
	"strconv"
	"strings"

	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/style"
)

type InputOTPProps struct {
	Value    string
	Length   int // number of slots, default 6
	OnChange func(string)
	// Numeric restricts entry to digits (the default); set Alnum for letters.
	Alnum bool
	Class string
}

// InputOTP is a one-time-code field: a row of single-character slots backed by
// a transparent input that captures typing. The value re-syncs from Value and
// is truncated/filtered to Length characters.
func InputOTP(p InputOTPProps) *g.Node {
	length := p.Length
	if length == 0 {
		length = 6
	}

	slots := make([]any, 0, length)
	for i := range length {
		ch := ""
		if i < len(p.Value) {
			ch = string(p.Value[i])
		}
		active := i == len(p.Value) && i < length
		slots = append(slots, g.Div(
			g.Class(style.CN(
				"relative flex h-9 w-9 items-center justify-center border-y border-r border-input text-sm font-medium shadow-sm first:rounded-l-md first:border-l last:rounded-r-md",
				map[string]bool{"z-10 ring-1 ring-ring": active},
			)),
			g.Data("slot", "input-otp-slot"),
			g.Data("index", strconv.Itoa(i)),
			ch,
		))
	}

	input := g.Input(
		g.Class("absolute inset-0 size-full cursor-text opacity-0"),
		g.Data("slot", "input-otp-input"),
		g.Value(p.Value),
		g.Attr("inputmode", inputMode(p.Alnum)),
		g.Attr("maxlength", strconv.Itoa(length)),
		g.Attr("autocomplete", "one-time-code"),
		g.OnInput(func(e *g.Event) {
			if p.OnChange != nil {
				p.OnChange(sanitizeOTP(e.Value(), length, p.Alnum))
			}
		}),
	)

	return g.Div(
		g.Class(style.CN("relative flex w-fit items-center", p.Class)),
		g.Data("slot", "input-otp"),
		input,
		g.Div(append([]any{g.Class("flex items-center")}, slots...)...),
	)
}

func inputMode(alnum bool) string {
	if alnum {
		return "text"
	}
	return "numeric"
}

func sanitizeOTP(s string, length int, alnum bool) string {
	var b strings.Builder
	for _, r := range s {
		ok := r >= '0' && r <= '9'
		if alnum {
			ok = ok || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
		}
		if ok {
			b.WriteRune(r)
		}
		if b.Len() >= length {
			break
		}
	}
	return b.String()
}
