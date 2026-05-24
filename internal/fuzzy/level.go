package fuzzy

func classifyLevel(score float64) string {
	if score < 40 {
		return "Плохо"
	}
	if score < 60 {
		return "Удовлетворительно"
	}
	if score < 80 {
		return "Хорошо"
	}
	if score < 95 {
		return "Отлично"
	}

	return "Идеально"
}
