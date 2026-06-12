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
