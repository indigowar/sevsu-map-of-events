package validators

import "github.com/google/uuid"

func IDExists(array []uuid.UUID, id uuid.UUID) bool {
	for _, v := range array {
		if id == v {
			return true
		}
	}
	return false
}

func StringExists(array []string, id string) bool {
	for _, v := range array {
		if id == v {
			return true
		}
	}
	return false
}
