package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	stepTypeDuration := strings.Split(data, ",")
	if len(stepTypeDuration) != 3 {
		return 0, "", 0, errors.New("Неверные данные")
	}

	step, err := strconv.Atoi(stepTypeDuration[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Неверный формат шагов: %w", err)
	}

	if step <= 0 {
		return 0, "", 0, errors.New("Неверное значение шагов")
	}

	duration, err1 := time.ParseDuration(stepTypeDuration[2])
	if err1 != nil {
		return 0, "", 0, fmt.Errorf("Неверное значение времени: %w", err1)
	}

	if duration <= 0 {
		return 0, "", 0, errors.New("Неверное значение времени")
	}

	return step, stepTypeDuration[1], duration, nil

}

func distance(steps int, height float64) float64 {
	return height * stepLengthCoefficient * float64(steps) / float64(mInKm)
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, typeTraining, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}

	switch typeTraining {
	case "Ходьба":
		calories, err1 := WalkingSpentCalories(steps, weight, height, duration)
		if err1 != nil {
			return "", err1
		}
		result1 := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\n", typeTraining, duration.Hours(), distance(steps, height))
		result2 := fmt.Sprintf("Скорость: %.2f км/ч\nСожгли калорий: %.2f\n", meanSpeed(steps, height, duration), calories)
		return result1 + result2, nil

	case "Бег":
		calories, err1 := RunningSpentCalories(steps, weight, height, duration)
		if err1 != nil {
			return "", err1
		}
		result1 := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\n", typeTraining, duration.Hours(), distance(steps, height))
		result2 := fmt.Sprintf("Скорость: %.2f км/ч\nСожгли калорий: %.2f\n", meanSpeed(steps, height, duration), calories)
		return result1 + result2, nil

	default:
		return "", errors.New("неизвестный тип тренировки")

	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Некорректные данные")
	}

	return (duration.Minutes() * meanSpeed(steps, height, duration) * weight) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Некорректные данные")
	}

	s := (duration.Minutes() * meanSpeed(steps, height, duration) * weight) / minInH

	return s * walkingCaloriesCoefficient, nil
}
