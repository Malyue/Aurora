package arrays

func Distinct[T any](array []T) []T {
	data := make(map[any]struct{})
	for _, v := range array {
		data[v] = struct{}{}
	}

	var result []T
	for k := range data {
		result = append(result, k)
	}
	return result
}

func Concat[T any](array []T, arrays ...[]T) []T {
	for _, arr := range arrays {
		array = append(array, arr...)
	}
	return array
}

func IsContain[T comparable](items []T, item T) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func Paging(pageNo, pageSize, length uint64) (int64, int64) {
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > length {
		pageSize = length
	}
	from, end := (pageNo-1)*pageSize, pageNo*pageSize
	if from > length {
		return -1, -1
	}
	if end > length {
		end = length
	}
	return int64(from), int64(end)
}

// IsArrayContained
// Check if the elements in the `sub` array are a subset of the `array`
// It returns (-1,true) if all elements in `sub` are present int `array`
// Otherwise it returns first elements index in `sub` which is not a subset of the `array` and false
func IsArrayContained[T comparable](array []T, sub []T) (int, bool) {
	if len(sub) == 0 {
		return -1, true
	}
	if len(array) == 0 {
		return 0, false
	}
	arrayMap := make(map[T]struct{})
	for _, item := range array {
		arrayMap[item] = struct{}{}
	}

	for index, item := range sub {
		if _, ok := arrayMap[item]; !ok {
			return index, false
		}
	}

	return -1, true
}
