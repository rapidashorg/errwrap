package errwrap

func stringPtr(str string) *string {
	return &str
}

func maskFormatterPtr(formatter MaskFormatter) *MaskFormatter {
	return &formatter
}

func messageFormatterPtr(formatter MessageFormatter) *MessageFormatter {
	return &formatter
}
