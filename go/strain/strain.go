package strain

const testVersion = 1

type Ints []int
type Lists [][]int
type Strings []string

func (ints Ints) Keep(pred func(int) bool) Ints {
	if ints == nil {
		return nil
	}
	res := make([]int, 0, len(ints))
	for _, in := range ints {
		if pred(in) {
			res = append(res, int(in))
		}
	}
	return res
}

func (ints Ints) Discard(pred func(int) bool) Ints {
	if ints == nil {
		return nil
	}
	res := make([]int, 0, len(ints))
	for _, in := range ints {
		if !pred(in) {
			res = append(res, int(in))
		}
	}
	return res
}

func (lists Lists) Keep(pred func([]int) bool) Lists {
	if lists == nil {
		return nil
	}
	res := make([][]int, 0, 0)
	for _, list := range lists {
		if pred(list) {
			res = append(res, list)
		}
	}
	return res
}

func (strs Strings) Keep(pred func(string) bool) Strings {
	res := make([]string, 0, 0)
	for _, str := range strs {
		if pred(str) {
			res = append(res, str)
		}
	}
	return res
}
