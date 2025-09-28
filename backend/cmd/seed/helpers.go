package main

import (
    "log"
    "strings"

    "gorm.io/gorm"
)

// stringPtr returns a pointer to the string value.
func stringPtr(s string) *string {
	return &s
}

// seedWithDuplicateCheck performs a generic seeding operation with duplicate checking
func seedWithDuplicateCheck[T any](db *gorm.DB, items []T, checkQuery func(item T) *gorm.DB, getLogInfo func(item T) string) error {
	for i := range items {
		var existing T
		err := checkQuery(items[i]).First(&existing).Error

		if err == nil {
			log.Printf("%s already exists\n", getLogInfo(items[i]))
			continue
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}

		if err := db.Create(&items[i]).Error; err != nil {
			return err
		}

		log.Printf("Created %s\n", getLogInfo(items[i]))
	}
	return nil
}

// splitAndTrim splits a comma-separated string and trims spaces from each element.
func splitAndTrim(s string) []string {
    parts := strings.Split(s, ",")
    out := make([]string, 0, len(parts))
    for _, p := range parts {
        t := strings.TrimSpace(p)
        if t != "" {
            out = append(out, t)
        }
    }
    return out
}
