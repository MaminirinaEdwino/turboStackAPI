package entity


func addChamps(targetSlice []Champs, newField Champs) []Champs {
	targetSlice = append(targetSlice, newField)
	return targetSlice
}

func deleteChamp(targetSlice []Champs, idx int) []Champs {
	for i := range targetSlice {
		if idx == i {
			targetSlice = append(targetSlice[:i-1], targetSlice[i:]...)
		}
	}
	return targetSlice
}