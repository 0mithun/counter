package display

type Options struct {
	ShowWords bool
	ShowBytes bool
	ShowLines bool
}

var ShowHeader bool

func (d Options) ShouldShowBytes() bool {
	if !d.ShowBytes && !d.ShowLines && !d.ShowWords {
		return true
	}

	return d.ShowBytes
}

func (d Options) ShouldShowLines() bool {
	if !d.ShowBytes && !d.ShowLines && !d.ShowWords {
		return true
	}

	return d.ShowLines
}

func (d Options) ShouldShowWords() bool {
	if !d.ShowBytes && !d.ShowLines && !d.ShowWords {
		return true
	}

	return d.ShowWords
}
