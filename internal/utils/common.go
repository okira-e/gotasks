package utils

func Cond[T any](condition bool, resolve T, reject T) T {
	if condition {
		return resolve
	} else {
		return reject
	}
}
