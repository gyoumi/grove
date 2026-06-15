package gallery

import (
	g "github.com/gyoumi/grove"
	"github.com/gyoumi/grove/ui"
)

func typographySection() *g.Node {
	return section("Typography",
		ui.TypographyH3("The quick brown fox"),
		ui.TypographyP("Grove renders this paragraph with the design system's prose styles, including ", ui.TypographyInlineCode("inline code"), " and links."),
		ui.TypographyBlockquote("Make it work, make it right, make it fast."),
	)
}
