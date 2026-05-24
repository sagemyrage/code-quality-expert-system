package fuzzy

const recommendationThreshold = 0.5

func buildRecommendations(input Input, scores scores) []string {
	var recommendations []string

	if scores.lineCount < recommendationThreshold {
		recommendations = append(recommendations, "Файл слишком большой: разделите его на несколько по принципу единственной ответственности")
	}
	if scores.commentRatio < recommendationThreshold {
		if input.CommentRatio < 0.1 {
			recommendations = append(recommendations, "Увеличьте долю комментариев: документируйте публичные функции и нетривиальную логику")
		} else {
			recommendations = append(recommendations, "Слишком много комментариев: удалите закомментированный код и избыточные пояснения")
		}
	}
	if scores.functionCount < recommendationThreshold {
		if input.FunctionCount < 2 {
			recommendations = append(recommendations, "Декомпозируйте код: разбейте логику на отдельные функции")
		} else {
			recommendations = append(recommendations, "Файл содержит слишком много функций: разделите его по принципу единственной ответственности")
		}
	}
	if scores.averageFunctionLength < recommendationThreshold {
		recommendations = append(recommendations, "Сократите среднюю длину функций: каждая функция должна делать одно дело")
	}
	if scores.maxFunctionLength < recommendationThreshold {
		recommendations = append(recommendations, "Разбейте самую длинную функцию на несколько меньших")
	}
	if scores.conditionalCount < recommendationThreshold {
		recommendations = append(recommendations, "Упростите ветвление: рассмотрите таблицы решений или паттерн стратегии")
	}
	if scores.loopCount < recommendationThreshold {
		recommendations = append(recommendations, "Уменьшите количество циклов: используйте вспомогательные функции или стандартную библиотеку")
	}
	if scores.maxNestingDepth < recommendationThreshold {
		recommendations = append(recommendations, "Уменьшите глубину вложенности: используйте ранний возврат или выносите условия в отдельные функции")
	}
	if scores.globalVariableCount < recommendationThreshold {
		recommendations = append(recommendations, "Сократите количество глобальных переменных: передавайте зависимости через параметры или структуры")
	}
	if scores.longLineCount < recommendationThreshold {
		recommendations = append(recommendations, "Сократите длинные строки: разбивайте выражения для улучшения читаемости")
	}

	return recommendations
}
