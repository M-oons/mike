package player

import "github.com/faiface/beep"

type Sound struct {
	Buffer *beep.Buffer
	Format beep.Format
}
