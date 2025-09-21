package main

// stringPtr returns a pointer to the string value.
func stringPtr(s string) *string {
	return &s
}

// int64Ptr returns a pointer to the int64 value.
func int64Ptr(i int64) *int64 {
	return &i
}

// float64Ptr returns a pointer to the float64 value.
func float64Ptr(f float64) *float64 {
	return &f
}