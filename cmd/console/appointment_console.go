package main

import (
	"bufio"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"tattoo-studio/internal/appointment"
	"tattoo-studio/internal/client"
	"tattoo-studio/internal/datetime"
	"time"
)

func addAppointment(scanner *bufio.Scanner, clients []client.Client, appointments []appointment.Appointment) (appointment.Appointment, bool) {
	for {
		fmt.Print("Введите ID клиента или 0 для выхода ")
		scanner.Scan()
		id := scanner.Text()
		if id == "0" {
			return appointment.Appointment{}, false
		}
		foundClient, found := client.FindByID(clients, id)
		if !found {
			fmt.Println("Клиент не найден")
			continue
		}
		clientID := foundClient.ID

		fmt.Print("Введите дату записи: ")
		scanner.Scan()
		dateValue := scanner.Text()
		fmt.Print("Введите время записи: ")
		scanner.Scan()
		timeValue := scanner.Text()
		startAt, err := time.ParseInLocation(
			"02.01.2006 15:04",
			dateValue+" "+timeValue,
			datetime.Moscow,
		)
		if err != nil {
			fmt.Println("Некорректный ввод даты")
			continue
		}

		fmt.Print("Введите описание: ")
		scanner.Scan()
		description := strings.TrimSpace(scanner.Text())

		fmt.Print("Введите ожидаему цену: ")
		scanner.Scan()
		plannedPrice, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			fmt.Println("Некорректный ввод планируемой цены")
			continue
		}

		fmt.Print("Введите депозит: ")
		scanner.Scan()
		deposit, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
		if err != nil {
			fmt.Println("Некорректный ввод депозита")
			continue
		}

		input := appointment.CreateInput{
			ClientID:     clientID,
			StartAt:      startAt,
			Description:  description,
			PlannedPrice: plannedPrice,
			Deposit:      deposit,
		}

		createdAppointment, err := appointment.Create(appointments, input)

		switch {
		case errors.Is(err, appointment.ErrStartTimeInPast):
			fmt.Println("Запись не может быть создана в прошлом")
			return appointment.Appointment{}, false

		case errors.Is(err, appointment.ErrTimeTaken):
			fmt.Println("Время уже занято")
			return appointment.Appointment{}, false

		case errors.Is(err, appointment.ErrPlannedPriceLessThanZero):
			fmt.Println("Цена не может быть меньше 0")
			return appointment.Appointment{}, false

		case errors.Is(err, appointment.ErrDepositLessThanZero):
			fmt.Println("Депозит не может быть меньше 0")
			return appointment.Appointment{}, false

		case err != nil:
			fmt.Println("Неизвестная ошибка")
			return appointment.Appointment{}, false
		}

		return createdAppointment, true
	}
}

func printAppointment(appointmentPrint appointment.Appointment, client client.Client) {
	fmt.Println("ID записи:", appointmentPrint.ID)
	fmt.Println(
		"Дата и время:",
		appointmentPrint.StartAt.Format("02.01.2006 15:04"),
	)
	fmt.Printf("Клиент: %s\n", client.Name)
	if !client.IsAdult {
		fmt.Println("Несовершеннолетний клиент")
	}
	if appointmentPrint.Description != "" {
		fmt.Printf("Описание: %s\n", appointmentPrint.Description)
	}
	fmt.Printf("Планируемая стоимость: %d ₽\n", appointmentPrint.PlannedPrice)
	fmt.Printf("Предоплата: %d ₽\n", appointmentPrint.Deposit)
	if appointmentPrint.Status == appointment.StatusCompleted {
		fmt.Printf("Финальная стоимость: %d ₽\n", appointmentPrint.FinalPrice)
	}
	fmt.Printf("Статус: %s\n", appointmentPrint.Status)
	fmt.Println("-------------------")
}

func printAppointmentByDate(scanner *bufio.Scanner, appointments []appointment.Appointment, clients []client.Client) {
	for {
		fmt.Println("Введите дату или 0 для выхода: ")
		scanner.Scan()
		date := scanner.Text()
		if date == "0" {
			return
		}
		parsedDate, err := time.ParseInLocation("02.01.2006", date, datetime.Moscow)
		if err != nil {
			fmt.Println("Некорректная дата")
			continue
		}
		appointmentsByDate := appointment.FindByDate(appointments, parsedDate)
		if len(appointmentsByDate) == 0 {
			fmt.Printf("Записей на: %s не найдено\n", parsedDate.Format("02.01.2006"))
			continue
		}
		fmt.Printf("Записи на %s\n", parsedDate.Format("02.01.2006"))
		fmt.Println()
		for _, v := range appointmentsByDate {
			clientAtDate, isExist := client.FindByID(clients, v.ClientID)
			if !isExist {
				fmt.Println("Ошибка поиска клиента")
				continue
			}
			printAppointment(v, clientAtDate)
		}
		return
	}
}

func reschedulingAppointment(scanner *bufio.Scanner, appointments []appointment.Appointment) { // рефакт
	for {
		fmt.Println("Введите ID записи или 0 для выхода: ")
		scanner.Scan()
		id := scanner.Text()

		if id == "0" {
			return
		}

		fmt.Println("Введите новую дату записи: ")
		scanner.Scan()
		dateValue := scanner.Text()

		fmt.Println("Введите новое время записи: ")
		scanner.Scan()
		timeValue := scanner.Text()

		startAt, err := time.ParseInLocation(
			"02.01.2006 15:04",
			dateValue+" "+timeValue,
			datetime.Moscow,
		)
		if err != nil {
			fmt.Println("Некорректные дата или время")
			continue
		}

		err = appointment.Reschedule(appointments, id, startAt)
		switch {
		case errors.Is(err, appointment.ErrAppointmentNotFound):
			fmt.Println("Запись не найдена")
			continue
		case errors.Is(err, appointment.ErrAppointmentCompleted):
			fmt.Println("Запись уже завершена")
			continue
		case errors.Is(err, appointment.ErrStartTimeInPast):
			fmt.Println("Невозможно перенести запись в прошлое")
			continue
		case errors.Is(err, appointment.ErrTimeTaken):
			fmt.Println("Время уже занято")
			continue
		case err != nil:
			fmt.Println("Неизвестная ошибка")
			continue
		}

		fmt.Println("Запись успешно перенесена")

		return
	}
}

func completeAppointment(scanner *bufio.Scanner, appointments []appointment.Appointment) { // рефакт
	for {
		fmt.Print("Введите ID записи для завершения или 0 для выхода: ")
		scanner.Scan()

		if scanner.Text() == "0" {
			return
		}

		id := scanner.Text()

		fmt.Println("Введите итоговую сумму: ")
		scanner.Scan()

		price := scanner.Text()

		finalPrice, err := strconv.Atoi(price)
		if err != nil {
			fmt.Println("Ошибка ввода суммы")
			continue
		}

		completedAppointment, err := appointment.Complete(appointments, id, finalPrice)

		switch {
		case err == nil:
			fmt.Println("Запись успешно завершена")
		case errors.Is(err, appointment.ErrAppointmentNotFound):
			fmt.Println("Запись не найдена")
			continue
		case errors.Is(err, appointment.ErrAppointmentCompleted):
			fmt.Println("Запись уже завершена")
			continue
		case errors.Is(err, appointment.ErrAppointmentCancelled):
			fmt.Println("Отменённую запись нельзя завершить")
			continue
		case errors.Is(err, appointment.ErrFinalPriceLessThanZero):
			fmt.Println("Итоговая сумма не может быть меньше 0")
			continue
		default:
			fmt.Println("Неизвестная ошибка")
			continue
		}

		fmt.Printf(
			"Планируемая стоимость: %d ₽\n",
			completedAppointment.PlannedPrice,
		)
		fmt.Printf(
			"Итоговая стоимость: %d ₽\n",
			completedAppointment.FinalPrice,
		)
		return
	}
}

func cancelAppointment(scanner *bufio.Scanner, appointments []appointment.Appointment) { // рефакт
	for {
		fmt.Print("Введите ID записи для отмены или 0 для выхода: ")
		scanner.Scan()

		id := scanner.Text()

		if id == "0" {
			return
		}

		err := appointment.Cancel(appointments, id)

		switch {
		case errors.Is(err, appointment.ErrAppointmentNotFound):
			fmt.Println("Запись не найдена")

		case errors.Is(err, appointment.ErrAppointmentCancelled):
			fmt.Println("Запись уже отменена")

		case errors.Is(err, appointment.ErrAppointmentCompleted):
			fmt.Println("Завершённую запись нельзя отменить")

		case err != nil:
			fmt.Println("Неизвестная ошибка")

		default:
			fmt.Println("Запись успешно отменена")
			return
		}
	}
}

func deleteAppointmentChoice(scanner *bufio.Scanner, appointments []appointment.Appointment, clients []client.Client) []appointment.Appointment {
	for {
		fmt.Print("Введите ID записи для удаления или 0 для выхода: ")
		scanner.Scan()

		if scanner.Text() == "0" {
			return appointments
		}

		idx, found := appointment.FindIndexByID(
			appointments,
			scanner.Text(),
		)
		if !found {
			fmt.Println("Запись не найдена")
			continue
		}
		if appointments[idx].Status == appointment.StatusCompleted {
			fmt.Println("Нельзя удалить завершённую запись")
			continue
		}

		if appointments[idx].Status == appointment.StatusCancelled {
			fmt.Println("Нельзя удалить отменённую запись")
			continue
		}
		foundClient, found := client.FindByID(
			clients,
			appointments[idx].ClientID,
		)

		if !found {
			fmt.Println("Клиент записи не найден")
			continue
		}

		printAppointment(appointments[idx], foundClient)
		for {
			fmt.Println()
			fmt.Println("1 — Удалить запись")
			fmt.Println("0 — Выход")
			fmt.Print("Выберите действие: ")

			scanner.Scan()
			deleteChoice := strings.TrimSpace(strings.ToLower(scanner.Text()))
			fmt.Println("-------------------")

			switch deleteChoice {
			case "1":
				fmt.Println("Вы уверены?")
				fmt.Println("1 — Да")
				fmt.Println("0 — Нет")
				scanner.Scan()
				if strings.TrimSpace(scanner.Text()) == "1" {
					appointments = slices.Delete(appointments, idx, idx+1)
					fmt.Println("Запись удалена")
					return appointments
				}
			case "0":
				return appointments
			default:
				fmt.Println("Неизвестная команда")
			}
		}
	}
}

func printClientHistory(scanner *bufio.Scanner, appointments []appointment.Appointment, clients []client.Client) {
	for {
		fmt.Print("Введите ID клиента что бы увидеть историю сеансов или 0 для выхода: ")
		scanner.Scan()

		if scanner.Text() == "0" {
			return
		}
		clientId := scanner.Text()

		client, isFound := client.FindByID(clients, clientId)
		if !isFound {
			fmt.Println("Клиент не найден")
			continue
		}

		clientHistory := appointment.FindByClientID(appointments, client.ID)

		if len(clientHistory) == 0 {
			fmt.Println("У клиента еще нет записей")
		}
		for _, v := range clientHistory {
			printAppointment(v, client)
		}
		return
	}
}

func printUpcoming(appointments []appointment.Appointment, clients []client.Client) {
	upcomingAppointments := appointment.FindUpcoming(appointments, time.Now().In(datetime.Moscow))
	if len(upcomingAppointments) == 0 {
		fmt.Println("Ближайших записей нет")
		return
	}
	for _, upcomingAppointment := range upcomingAppointments {
		foundClient, found := client.FindByID(clients, upcomingAppointment.ClientID)
		if !found {
			fmt.Println("Ошибка поиска клиента")
			continue
		}
		printAppointment(upcomingAppointment, foundClient)
	}

}
