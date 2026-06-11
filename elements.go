package grove

// Element constructors for common HTML tags. Each accepts the same argument
// forms as El: Options, *Node children, strings, []*Node, and nil.

func A(args ...any) *Node          { return El("a", args...) }
func Abbr(args ...any) *Node       { return El("abbr", args...) }
func Address(args ...any) *Node    { return El("address", args...) }
func Article(args ...any) *Node    { return El("article", args...) }
func Aside(args ...any) *Node      { return El("aside", args...) }
func Audio(args ...any) *Node      { return El("audio", args...) }
func B(args ...any) *Node          { return El("b", args...) }
func Blockquote(args ...any) *Node { return El("blockquote", args...) }
func Br(args ...any) *Node         { return El("br", args...) }
func Button(args ...any) *Node     { return El("button", args...) }
func Canvas(args ...any) *Node     { return El("canvas", args...) }
func Caption(args ...any) *Node    { return El("caption", args...) }
func Cite(args ...any) *Node       { return El("cite", args...) }
func Code(args ...any) *Node       { return El("code", args...) }
func Datalist(args ...any) *Node   { return El("datalist", args...) }
func Dd(args ...any) *Node         { return El("dd", args...) }
func Details(args ...any) *Node    { return El("details", args...) }
func Div(args ...any) *Node        { return El("div", args...) }
func Dl(args ...any) *Node         { return El("dl", args...) }
func Dt(args ...any) *Node         { return El("dt", args...) }
func Em(args ...any) *Node         { return El("em", args...) }
func Fieldset(args ...any) *Node   { return El("fieldset", args...) }
func Figcaption(args ...any) *Node { return El("figcaption", args...) }
func Figure(args ...any) *Node     { return El("figure", args...) }
func Footer(args ...any) *Node     { return El("footer", args...) }
func Form(args ...any) *Node       { return El("form", args...) }
func H1(args ...any) *Node         { return El("h1", args...) }
func H2(args ...any) *Node         { return El("h2", args...) }
func H3(args ...any) *Node         { return El("h3", args...) }
func H4(args ...any) *Node         { return El("h4", args...) }
func H5(args ...any) *Node         { return El("h5", args...) }
func H6(args ...any) *Node         { return El("h6", args...) }
func Header(args ...any) *Node     { return El("header", args...) }
func Hr(args ...any) *Node         { return El("hr", args...) }
func I(args ...any) *Node          { return El("i", args...) }
func Iframe(args ...any) *Node     { return El("iframe", args...) }
func Img(args ...any) *Node        { return El("img", args...) }
func Input(args ...any) *Node      { return El("input", args...) }
func Kbd(args ...any) *Node        { return El("kbd", args...) }
func Label(args ...any) *Node      { return El("label", args...) }
func Legend(args ...any) *Node     { return El("legend", args...) }
func Li(args ...any) *Node         { return El("li", args...) }
func Main(args ...any) *Node       { return El("main", args...) }
func Mark(args ...any) *Node       { return El("mark", args...) }
func Nav(args ...any) *Node        { return El("nav", args...) }
func Ol(args ...any) *Node         { return El("ol", args...) }
func Optgroup(args ...any) *Node   { return El("optgroup", args...) }

// OptionEl creates an <option> element (named to avoid clashing with the
// Option interface).
func OptionEl(args ...any) *Node { return El("option", args...) }

func P(args ...any) *Node        { return El("p", args...) }
func Picture(args ...any) *Node  { return El("picture", args...) }
func Pre(args ...any) *Node      { return El("pre", args...) }
func Progress(args ...any) *Node { return El("progress", args...) }
func Q(args ...any) *Node        { return El("q", args...) }
func Section(args ...any) *Node  { return El("section", args...) }
func Select(args ...any) *Node   { return El("select", args...) }
func Small(args ...any) *Node    { return El("small", args...) }
func Span(args ...any) *Node     { return El("span", args...) }
func Strong(args ...any) *Node   { return El("strong", args...) }
func Sub(args ...any) *Node      { return El("sub", args...) }
func Summary(args ...any) *Node  { return El("summary", args...) }
func Sup(args ...any) *Node      { return El("sup", args...) }
func Table(args ...any) *Node    { return El("table", args...) }
func Tbody(args ...any) *Node    { return El("tbody", args...) }
func Td(args ...any) *Node       { return El("td", args...) }
func TextArea(args ...any) *Node { return El("textarea", args...) }
func Tfoot(args ...any) *Node    { return El("tfoot", args...) }
func Th(args ...any) *Node       { return El("th", args...) }
func Thead(args ...any) *Node    { return El("thead", args...) }
func Time(args ...any) *Node     { return El("time", args...) }
func Tr(args ...any) *Node       { return El("tr", args...) }
func U(args ...any) *Node        { return El("u", args...) }
func Ul(args ...any) *Node       { return El("ul", args...) }
func Video(args ...any) *Node    { return El("video", args...) }
