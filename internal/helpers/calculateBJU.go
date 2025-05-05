package helpers

import "fmt"

func CalculateBJU(height, weight float64) (string, bool) {
	heightInMeters := height / 100
	bmi := weight / (heightInMeters * heightInMeters)

	calories := 24 * weight
	protein := weight * 1.5
	fat := weight * 0.8
	carbs := (calories - (protein*4 + fat*9)) / 4

	needsWeightLoss := bmi > 25
	weightStatus := "У вас нормальный вес."
	if needsWeightLoss {
		weightStatus = "У вас лишний вес. Нужно худеть."
		calories -= 500
	}

	message := fmt.Sprintf(
		"Твоя примерная дневная норма:\n\nКалории: %.0f ккал\nБелки: %.1f г\nЖиры: %.1f г\nУглеводы: %.1f г\n\nИМТ: %.1f\n%s",
		calories, protein, fat, carbs, bmi, weightStatus)

	return message, needsWeightLoss
}
