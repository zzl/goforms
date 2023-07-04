package layouts

func Items[T LayoutItem](items []T) []LayoutItem {
	if items == nil {
		return nil
	}
	result := make([]LayoutItem, len(items))
	for n, item := range items {
		result[n] = item
	}
	return result
}
