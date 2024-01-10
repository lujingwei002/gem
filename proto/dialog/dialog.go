package dialog

type Dialog struct {
	Source  string
	Content string
}

func New(source string, content string) *Dialog {
	self := &Dialog{
		Source:  source,
		Content: content,
	}
	return self
}
