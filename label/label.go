package label

import "fmt"

type Labelable interface {
	LabelID() string
}

func Label(l Labelable, labelType string) string {
	return fmt.Sprintf("%v(%v)", labelType, l.LabelID())
}
