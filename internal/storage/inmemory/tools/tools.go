package tools

func Paginate[T any](slice []T, limit, offset int) []T {
	if offset > len(slice) {
		return []T{}
	}
	end := offset + limit
	if end > len(slice) {
		end = len(slice)
	}
	return slice[offset:end]
}
