package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {

	stepAndDuration := strings.Split(data, ",")
	if len(stepAndDuration) != 2 {
		return 0, 0, errors.New("Неверные данные")
	}

	step, err := strconv.Atoi(stepAndDuration[0])
	if err != nil {
		return 0, 0, err
	}

	if step <= 0 {
		return 0, 0, errors.New("Неверное значение шагов")
	}

	duration, err1 := time.ParseDuration(stepAndDuration[1])
	if err1 != nil {
		return 0, 0, err1
	}

	if duration <= 0 {
		return 0, 0, errors.New("Неверное значение вреёмени")
	}

	return step, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distance := float64(steps) * stepLength / mInKm

	calories, err1 := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err1 != nil {
		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)

	return result
}
