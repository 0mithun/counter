package display

type Options struct {
	args NewOptionsArgs
}

var ShowHeader bool

type NewOptionsArgs struct {
	ShowWords bool
	ShowBytes bool
	ShowLines bool
}

func NewOptions(args NewOptionsArgs) Options {
	return Options{
		args: args,
	}
}

func (d Options) ShouldShowBytes() bool {
	if !d.args.ShowBytes && !d.args.ShowLines && !d.args.ShowWords {
		return true
	}

	return d.args.ShowBytes
}

func (d Options) ShouldShowLines() bool {
	if !d.args.ShowBytes && !d.args.ShowLines && !d.args.ShowWords {
		return true
	}

	return d.args.ShowLines
}

func (d Options) ShouldShowWords() bool {
	if !d.args.ShowBytes && !d.args.ShowLines && !d.args.ShowWords {
		return true
	}

	return d.args.ShowWords
}
