package accumulate

const testVersion = 1

func Accumulate(in []string, converter func(string) string) []string {
	out := make([]string, len(in))
	for i, s := range in {
		out[i] = converter(s)
	}
	return out
}
