package messages

import (
	"os"
	"time"
)

type EditorFinishedMsg struct {
	Err  error
	File *os.File
}

type SetAnswerMessage struct {
	Answer         string
	AnswerRendered string
}
type NewAnswerMessage struct {
	Answer string
	Prompt string
}
type NewRenderMessage struct{ Content string }

type HistoryResizedMessage struct {
	NewWidth   int
	TotalWidth int
}

type ViewPortResizeMessage struct {
	Width int
}

type ResponseMsg struct {
	Response string
	Prompt   string
}

type TickMsg time.Time

type ExecutePromptMessage struct {
	Prompt string
}

type ErrMsg struct{ Err error }

func (e ErrMsg) Error() string { return e.Err.Error() }
