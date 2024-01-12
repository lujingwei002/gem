package userid

import "fmt"

type Int64 int64

func (id Int64) Value() any {
	return id
}

func (id Int64) String() string {
	return fmt.Sprintf("%d", id)
}
