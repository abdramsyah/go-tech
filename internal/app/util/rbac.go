package util

import "fmt"

func FormatRbacSubject(adminID uint64) string {
	return fmt.Sprintf("user%d", adminID)
}

func FormatRbacRole(roleID uint64) string {
	return fmt.Sprintf("role%d", roleID)
}
