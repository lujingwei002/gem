package userid

type Int64 int64

func (id Int64) Value() any {
	return id
}
