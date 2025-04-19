package handlers

import (
	"context"
	postgres "diet-bot/internal/store"
	"fmt"
	"log"
	"log/slog"
	"math"
	"strconv"

	"github.com/jackc/pgx/v5"
	tele "gopkg.in/telebot.v4"
)

var userData = make(map[int64]map[string]float64)
var userStates = make(map[int64]string)

func Register(b *tele.Bot) {
	b.Handle("/hello", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Handle("/start", func(c tele.Context) error {
		var userID int64
		err := postgres.Pool.QueryRow(context.Background(),
			"SELECT user_id FROM users WHERE user_id = $1",
			c.Chat().ID,
		).Scan(&userID)

		if err != nil {
			if err == pgx.ErrNoRows {
				var newID int
				err = postgres.Pool.QueryRow(context.Background(),
					"INSERT INTO users (user_id) VALUES ($1) RETURNING id",
					c.Chat().ID,
				).Scan(&newID)

				if err != nil {
					log.Printf("Insert error: %v", err)
					return c.Send("Something went wrong. Please try again later.")
				}

				return c.Send(fmt.Sprintf("Welcome %s!", c.Sender().FirstName))
			} else {
				log.Printf("Query error: %v", err)
				return c.Send("Something went wrong. Please try again later.")
			}
		}

		return c.Send(fmt.Sprintf("Welcome back %s!", c.Sender().FirstName))
	})

	b.Handle("/getInfo", func(c tele.Context) error {
		return c.Send("Hello!")
	})

	b.Handle("/products", func(c tele.Context) error {
		return c.Send("тухум, помидор, зелень, грудинка, оливка, бодринг, мол гошти")
	})

	b.Handle("/help", func(c tele.Context) error {
		ID := c.Chat().ID

		slog.Info("help", "chatID", ID)

		return c.Send("Hello!")
	})

	b.Handle("/options", func(c tele.Context) error {
		userID := c.Chat().ID
		data, ok := userData[userID]
		if !ok || data["height"] == 0 || data["weight"] == 0 {
			userStates[userID] = "awaiting_height"
			return c.Send("Пожалуйста, введи сначала свой рост (в см):")
		}

		bjuBtn := tele.Btn{Unique: "get_bju", Text: "Узнать свой актуальный БЖУ"}
		inlineMenu := &tele.ReplyMarkup{}
		inlineMenu.Inline(
			inlineMenu.Row(inlineMenu.Data(bjuBtn.Text, bjuBtn.Unique)),
		)

		return c.Send("Выберите действие:", inlineMenu)
	})

	bjuBtn := tele.Btn{Unique: "get_bju", Text: "Узнать свой актуальный БЖУ"}
	mealPlanBtn := tele.Btn{Unique: "get_meal_plan", Text: "План блюд для похудения"}

	b.Handle(&bjuBtn, func(c tele.Context) error {
		userID := c.Sender().ID
		data, ok := userData[userID]
		if !ok || data["height"] == 0 || data["weight"] == 0 {
			return c.Send("Пожалуйста, сначала введи свои данные с помощью /start.")
		}

		height := data["height"]
		weight := data["weight"]

		result, needsWeightLoss := calculateBJU(height, weight)

		if needsWeightLoss {
			inlineMenu := &tele.ReplyMarkup{}
			inlineMenu.Inline(
				inlineMenu.Row(mealPlanBtn),
			)
			return c.Send(result, inlineMenu)
		}

		return c.Send(result)
	})

	b.Handle(tele.OnText, func(c tele.Context) error {
		userID := c.Chat().ID
		state := userStates[userID]
		text := c.Text()

		// Инициализируем карту, если её нет
		if _, exists := userData[userID]; !exists {
			userData[userID] = make(map[string]float64)
		}

		switch state {
		case "awaiting_height":
			height, err := strconv.ParseFloat(text, 64)
			if err != nil || height < 50 || height > 300 {
				return c.Send("Пожалуйста, введи корректный рост (в см):")
			}
			userData[userID]["height"] = height
			userStates[userID] = "awaiting_weight"
			return c.Send("Теперь введи свой вес в кг:")
		case "awaiting_weight":
			weight, err := strconv.ParseFloat(text, 64)
			if err != nil || weight < 20 || weight > 300 {
				return c.Send("Пожалуйста, введи корректный вес (в кг):")
			}
			userData[userID]["weight"] = weight
			userStates[userID] = "completed"

			bjuBtn := tele.Btn{Unique: "get_bju", Text: "Узнать свой актуальный БЖУ"}
			inlineMenu := &tele.ReplyMarkup{}
			inlineMenu.Inline(
				inlineMenu.Row(inlineMenu.Data(bjuBtn.Text, bjuBtn.Unique)),
			)
			return c.Send("Выберите действие:", inlineMenu)
		default:
			return c.Send("Напиши /start, чтобы начать.")
		}
	})

	b.Handle(&mealPlanBtn, func(c tele.Context) error {
		userID := c.Sender().ID
		data := userData[userID]

		// Суточные нормы
		calories := 24*data["weight"] - 500
		protein := data["weight"] * 1.5
		fat := data["weight"] * 0.8
		carbs := (calories - (protein*4 + fat*9)) / 4

		// Продукты с БЖУ и реалистичными ограничениями
		products := map[string]struct {
			Protein float64
			Fat     float64
			Carbs   float64
			Max     float64 // Максимальная реалистичная порция
		}{
			"яйцо":     {12.7, 11.5, 0.7, 200}, // Макс 4 яйца (~200г)
			"помидор":  {0.9, 0.2, 3.9, 300},
			"зелень":   {2.0, 0.5, 5.0, 100},
			"грудинка": {22.0, 3.0, 0.0, 250},
			"оливки":   {1.0, 15.0, 6.0, 50}, // ~10 оливок
			"огурец":   {0.8, 0.1, 2.5, 300},
			"говядина": {26.0, 15.0, 0.0, 200},
			"рыба":     {20.0, 5.0, 0.0, 300},
			"масло":    {0.0, 100.0, 0.0, 15}, // 1 ст.л.
		}

		// Умный расчет с ограничениями
		calcPortion := func(product string, targetProtein, targetFat, targetCarbs float64) float64 {
			p := products[product]
			var grams float64

			if p.Protein > 0 {
				grams = targetProtein / (p.Protein / 100)
			}
			if p.Fat > 0 {
				grams = math.Max(grams, targetFat/(p.Fat/100))
			}
			if p.Carbs > 0 {
				grams = math.Max(grams, targetCarbs/(p.Carbs/100))
			}

			// Ограничение по максимальной порции
			return math.Min(math.Round(grams), p.Max)
		}

		// Распределение БЖУ по приемам пищи
		breakfastEggs := calcPortion("яйцо", protein*0.3, fat*0.1, 0)
		breakfastVeggies := calcPortion("помидор", 0, 0, carbs*0.2)

		lunchMeat := calcPortion("грудинка", protein*0.35, fat*0.2, 0)
		lunchSalad := calcPortion("огурец", 0, 0, carbs*0.3)

		dinnerFish := calcPortion("рыба", protein*0.25, fat*0.2, 0)
		dinnerVeggies := calcPortion("помидор", 0, 0, carbs*0.2)

		snackBeef := calcPortion("говядина", protein*0.1, fat*0.1, 0)
		snackOlives := calcPortion("оливки", 0, fat*0.1, 0)

		oil := calcPortion("масло", 0, fat*0.2, 0)

		mealPlan := fmt.Sprintf(`🍳 Реалистичный план (%.0f ккал)

			🌅 Завтрак:
			- %.0fг яиц (%.1fг белка)
			- %.0fг помидоров с зеленью
			- %.1fг оливкового масла

			☀ Перекус:
			- %.0fг говядины
			- %.0fг огурцов
			- %.0fг оливок

			🍗 Обед:
			- %.0fг куриной грудки
			- %.0fг овощного салата
			- %.1fг масла

			🐟 Ужин:
			- %.0fг рыбы
			- %.0fг тушеных овощей
			- Зелень

			📊 Итого:
			⚡ Белки: %.1fг | 🛢 Жиры: %.1fг | 🌾 Углеводы: %.1fг`,
			calories,
			breakfastEggs, breakfastEggs*products["яйцо"].Protein/100,
			breakfastVeggies,
			oil/2,
			snackBeef,
			calcPortion("огурец", 0, 0, carbs*0.1),
			snackOlives,
			lunchMeat,
			lunchSalad,
			oil/2,
			dinnerFish,
			dinnerVeggies,
			protein, fat, carbs)

		return c.Send(mealPlan)
	})

	slog.Info("handlers registered")
}

func calculateBJU(height, weight float64) (string, bool) {
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
