package style

import "testing"

func TestCN(t *testing.T) {
	cases := []struct {
		name string
		args []any
		want string
	}{
		{"later wins same group", []any{"p-4 p-2"}, "p-2"},
		{"different groups kept", []any{"px-4 py-2"}, "px-4 py-2"},
		{"p overrides px/py", []any{"px-4 py-2 p-3"}, "p-3"},
		{"px does not override p", []any{"p-3 px-4"}, "p-3 px-4"},
		{"bg color conflict", []any{"bg-blue-500 text-white", "bg-red-500"}, "text-white bg-red-500"},
		{"bg size vs color independent", []any{"bg-red-500 bg-cover"}, "bg-red-500 bg-cover"},
		{"variant scoping", []any{"hover:bg-blue-500 bg-blue-500 hover:bg-red-500"}, "bg-blue-500 hover:bg-red-500"},
		{"text size vs color", []any{"text-sm text-red-500 text-lg"}, "text-red-500 text-lg"},
		{"font size kills leading", []any{"leading-7 text-base"}, "text-base"},
		{"border width vs color", []any{"border border-input border-2"}, "border-input border-2"},
		{"border color replaced", []any{"border-input border-destructive"}, "border-destructive"},
		{"rounded sides", []any{"rounded-md rounded-t-lg"}, "rounded-md rounded-t-lg"},
		{"rounded overrides sides", []any{"rounded-t-lg rounded-md"}, "rounded-md"},
		{"display group", []any{"flex hidden"}, "hidden"},
		{"flex-1 is not display", []any{"flex flex-1"}, "flex flex-1"},
		{"flex-direction distinct", []any{"flex flex-col"}, "flex flex-col"},
		{"size overrides w/h", []any{"w-4 h-4 size-9"}, "size-9"},
		{"dedup unknown", []any{"custom custom"}, "custom"},
		{"map and slice args", []any{[]string{"a", "b"}, map[string]bool{"c": true, "d": false}}, "a b c"},
		{"nil skipped", []any{nil, "p-1"}, "p-1"},
		{"negative margins", []any{"-mt-2 mt-4"}, "mt-4"},
		{"important kept separate", []any{"!p-4 p-2"}, "!p-4 p-2"},
		{"arbitrary value width", []any{"w-[100px] w-4"}, "w-4"},
		{"shadow merge", []any{"shadow-sm hover:shadow-md shadow-lg"}, "hover:shadow-md shadow-lg"},
		{"ring width vs color", []any{"ring-1 ring-ring ring-2"}, "ring-ring ring-2"},
		{"component class override", []any{"bg-primary text-primary-foreground h-9 px-4", "bg-destructive h-8"}, "text-primary-foreground px-4 bg-destructive h-8"},
		{"size with line-height modifier", []any{"text-sm/6 text-lg"}, "text-lg"},
		{"size modifier kills leading", []any{"leading-7 text-sm/6"}, "text-sm/6"},
		{"size modifier is not a color", []any{"text-red-500 text-sm/6"}, "text-red-500 text-sm/6"},
		{"color opacity modifier", []any{"bg-black/50 bg-primary"}, "bg-primary"},
		{"col-span distinct from col-start", []any{"col-span-2 col-start-1"}, "col-span-2 col-start-1"},
		{"col-span merges", []any{"col-span-2 col-span-3"}, "col-span-3"},
		{"col-auto resets span and start", []any{"col-span-2 col-start-1 col-auto"}, "col-auto"},
		{"arbitrary property conflicts per property", []any{"[mask-type:luminance] [mask-type:alpha]"}, "[mask-type:alpha]"},
		{"different arbitrary properties kept", []any{"[mask-type:alpha] [paint-order:stroke]"}, "[mask-type:alpha] [paint-order:stroke]"},
		{"empty", []any{""}, ""},

		// v4 trailing important markers
		{"trailing important kept separate", []any{"p-4! p-2"}, "p-4! p-2"},
		{"trailing important merges with leading", []any{"!p-4 p-2!"}, "p-2!"},

		// page breaks vs word breaking are different properties
		{"break families distinct", []any{"break-inside-avoid break-all break-after-page"}, "break-inside-avoid break-all break-after-page"},
		{"word-break merges", []any{"break-all break-keep"}, "break-keep"},
		{"break-inside merges", []any{"break-inside-avoid break-inside-auto"}, "break-inside-auto"},

		// shadow size vs shadow color
		{"shadow color kept with size", []any{"shadow-lg shadow-primary/20"}, "shadow-lg shadow-primary/20"},
		{"shadow sizes merge", []any{"shadow-lg shadow-xl"}, "shadow-xl"},
		{"arbitrary shadow is a size", []any{"shadow-lg shadow-[0_1px_2px_black]"}, "shadow-[0_1px_2px_black]"},
		{"text-shadow distinct from text color", []any{"text-shadow-lg text-red-500"}, "text-shadow-lg text-red-500"},
		{"text-shadow merges", []any{"text-shadow-sm text-shadow-lg"}, "text-shadow-lg"},

		// gradient stop positions vs colors
		{"gradient stop position kept", []any{"from-red-500 from-10%"}, "from-red-500 from-10%"},
		{"gradient colors merge", []any{"from-red-500 from-emerald-400"}, "from-emerald-400"},
		{"gradient positions merge", []any{"via-10% via-[25%]"}, "via-[25%]"},

		// scroll margin/padding mirror the m-/p- side hierarchy
		{"scroll families distinct", []any{"scroll-mt-2 scroll-pb-2 scroll-smooth"}, "scroll-mt-2 scroll-pb-2 scroll-smooth"},
		{"scroll-m overrides sides", []any{"scroll-mt-2 scroll-m-4"}, "scroll-m-4"},
		{"scroll side does not override scroll-m", []any{"scroll-m-4 scroll-mt-2"}, "scroll-m-4 scroll-mt-2"},

		// snap axes vs alignment vs strictness
		{"snap groups distinct", []any{"snap-x snap-start snap-mandatory"}, "snap-x snap-start snap-mandatory"},
		{"snap type merges", []any{"snap-x snap-both"}, "snap-both"},

		// touch resets vs pan axes
		{"touch axes distinct", []any{"touch-pan-x touch-pan-y"}, "touch-pan-x touch-pan-y"},
		{"touch-auto resets axes", []any{"touch-pan-x touch-pinch-zoom touch-auto"}, "touch-auto"},
		{"pan cancels touch reset", []any{"touch-auto touch-pan-left"}, "touch-pan-left"},

		// backdrop filters conflict per filter
		{"backdrop filters distinct", []any{"backdrop-blur-sm backdrop-opacity-50"}, "backdrop-blur-sm backdrop-opacity-50"},
		{"backdrop blur merges", []any{"backdrop-blur-sm backdrop-blur-lg"}, "backdrop-blur-lg"},

		// inset-ring / inset-shadow are not position utilities
		{"inset-ring distinct from inset", []any{"inset-ring-2 inset-4"}, "inset-ring-2 inset-4"},
		{"inset-ring width vs color", []any{"inset-ring-2 inset-ring-primary"}, "inset-ring-2 inset-ring-primary"},
		{"inset-shadow merges", []any{"inset-shadow-2xs inset-shadow-sm"}, "inset-shadow-sm"},
		{"logical inset start", []any{"start-0 inset-x-2"}, "inset-x-2"},

		// font-variant-numeric axes compose; normal-nums resets
		{"fvn axes compose", []any{"ordinal tabular-nums"}, "ordinal tabular-nums"},
		{"fvn spacing merges", []any{"tabular-nums proportional-nums"}, "proportional-nums"},
		{"normal-nums resets fvn", []any{"ordinal tabular-nums normal-nums"}, "normal-nums"},
		{"fvn cancels normal-nums", []any{"normal-nums tabular-nums"}, "tabular-nums"},

		{"divide style vs color", []any{"divide-dashed divide-red-500"}, "divide-dashed divide-red-500"},
		{"divide styles merge", []any{"divide-dashed divide-solid"}, "divide-solid"},
		{"list groups distinct", []any{"list-disc list-inside"}, "list-disc list-inside"},
		{"list type merges", []any{"list-disc list-decimal"}, "list-decimal"},
		{"align-content vs content", []any{"content-center content-none"}, "content-center content-none"},
		{"content property merges", []any{"content-none content-['*']"}, "content-['*']"},
		{"hyphens merges", []any{"hyphens-auto hyphens-none"}, "hyphens-none"},
		{"overscroll axes", []any{"overscroll-x-auto overscroll-contain"}, "overscroll-contain"},
		{"table layout merges", []any{"table-auto table-fixed"}, "table-fixed"},
		{"col shorthand overrides span", []any{"col-span-2 col-auto"}, "col-auto"},
		{"span overrides col shorthand", []any{"col-auto col-span-2"}, "col-span-2"},
		{"border-spacing axes", []any{"border-spacing-x-2 border-spacing-y-2"}, "border-spacing-x-2 border-spacing-y-2"},
		{"border-spacing overrides axes", []any{"border-spacing-x-2 border-spacing-2"}, "border-spacing-2"},
		{"skew axes", []any{"skew-x-3 skew-y-3"}, "skew-x-3 skew-y-3"},
		{"skew overrides axes", []any{"skew-x-3 skew-6"}, "skew-6"},
		{"filters distinct", []any{"grayscale invert sepia hue-rotate-90"}, "grayscale invert sepia hue-rotate-90"},
		{"grayscale merges", []any{"grayscale grayscale-0"}, "grayscale-0"},
		{"mix-blend merges", []any{"mix-blend-multiply mix-blend-screen"}, "mix-blend-screen"},
		{"bg-blend is not bg-color", []any{"bg-blend-multiply bg-red-500"}, "bg-blend-multiply bg-red-500"},
		{"mask families distinct", []any{"mask-clip-border mask-luminance mask-repeat-x"}, "mask-clip-border mask-luminance mask-repeat-x"},
		{"mask image merges", []any{"mask-none mask-[url(m.svg)]"}, "mask-[url(m.svg)]"},
		{"float merges", []any{"float-left float-none"}, "float-none"},
		{"isolation merges", []any{"isolate isolation-auto"}, "isolation-auto"},
		{"box-sizing merges", []any{"box-border box-content"}, "box-content"},
		{"font-smoothing merges", []any{"antialiased subpixel-antialiased"}, "subpixel-antialiased"},
		{"font-stretch distinct from weight", []any{"font-stretch-75% font-bold"}, "font-stretch-75% font-bold"},
		{"transform mode merges", []any{"transform-gpu transform-none"}, "transform-none"},
		{"overflow-wrap merges", []any{"wrap-break-word wrap-anywhere"}, "wrap-anywhere"},
	}
	for _, c := range cases {
		if got := CN(c.args...); got != c.want {
			t.Errorf("%s: CN(%v)\n got: %q\nwant: %q", c.name, c.args, got, c.want)
		}
	}
}

func TestVariants(t *testing.T) {
	v := Variants{
		Base: "inline-flex rounded-md",
		Groups: map[string]map[string]string{
			"variant": {
				"default":     "bg-primary text-primary-foreground",
				"destructive": "bg-destructive text-white",
			},
			"size": {
				"default": "h-9 px-4",
				"sm":      "h-8 px-3",
			},
		},
		Defaults: map[string]string{"variant": "default", "size": "default"},
	}

	if got := v.Class(nil); got != "inline-flex rounded-md h-9 px-4 bg-primary text-primary-foreground" {
		t.Errorf("defaults: %q", got)
	}
	got := v.Class(map[string]string{"variant": "destructive", "size": "sm"})
	want := "inline-flex rounded-md h-8 px-3 bg-destructive text-white"
	if got != want {
		t.Errorf("selected:\n got: %q\nwant: %q", got, want)
	}
	// caller classes override variant classes
	got = v.Class(nil, "bg-accent rounded-full")
	want = "inline-flex h-9 px-4 text-primary-foreground bg-accent rounded-full"
	if got != want {
		t.Errorf("extra:\n got: %q\nwant: %q", got, want)
	}
}
