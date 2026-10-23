package main

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"tattoo-studio/internal/appointment"
	"tattoo-studio/internal/client"
	"tattoo-studio/internal/datetime"
	"time"
)

func addAppointment(scanner *bufio.Scanner, clients []client.Client, appointments []appointment.Appointment) (appointment.Appointment, bool, error) { // Исправить ux
	for {
		fmt.Print("Введите ID клиента или 0 для выхода ")
		id, err := readInput(scanner)
		if err != nil {
			return appointment.Appointment{}, false, fmt.Errorf("read client ID: %w", err)
		}
		if id == "0" {
			return appointment.Appointment{}, false, nil
		}
		foundClient, found := client.FindByID(clients, id)
		if !found {
			fmt.Println("Клиент не найден")
			continue
		}
		clientID := foundClient.ID

		fmt.Print("Введите дату записи: ")
		dateValue, err := readInput(scanner)
		if err != nil {
			return appointment.Appointment{}, false, fmt.Errorf("read appointment date value: %w", err)
		}
		fmt.Print("Введите время записи: ")
		timeValue, err := readInput(scanner)
		if err != nil {
			return appointment.Appointment{}, false, fmt.Errorf("read appointment time value: %w", err)
		}
		startAt, err := time.ParseInLocation(
			"02.01.2006 15:04",
			dateValue+" "+timeValue,
			datetime.Moscow,
		)
		if err != nil {
			fmt.Println("Ошибка ввода даты")
			continue
		}

		fmt.Print("Введите описание: ")

		description, err := readInput(scanner)
		if err != nil {
			return appointment.Appointment{}, false, fmt.Errorf("read appointment description: %w", err)
		}
		fmt.Print("Введите ожидаему цену: ")

		plannedPriceInput, err := readInput(scanner)
		if err != nil {
			return appointment.Appointment{}, false, fmt.Errorf("read appointment planned price: %w", err)
		}
		plannedPrice, err := strconv.Atoi(plannedPriceInput)
		if err != nil {
			fmt.Println("Ошибка ввода планируемой цены")
			continue
		}

		fmt.Print("Введите депозит: ")
		depositInput, err := readInput(scanner)
		if err != nil {
			return appointment.Appointment{}, false, fmt.Errorf("read appointment deposit: %w", err)
		}
		deposit, err := strconv.Atoi(depositInput)
		if err != nil {
			fmt.Println("Ошибка ввода депозита")
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
			continue
		case errors.Is(err, appointment.ErrTimeTaken):
			fmt.Println("Время уже занято")
			continue
		case errors.Is(err, appointment.ErrPlannedPriceLessThanZero):
			fmt.Println("Цена не может быть меньше 0")
			continue
		case errors.Is(err, appointment.ErrDepositLessThanZero):
			fmt.Println("Депозит не может быть меньше 0")
			continue
		case err != nil:
			return appointment.Appointment{}, false, fmt.Errorf("create appointment: %w", err)
		}

		return createdAppointment, true, nil
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
} // datetime рефакт

func printAppointmentByDate(scanner *bufio.Scanner, appointments []appointment.Appointment, clients []client.Client) error {
	for {
		fmt.Println("Введите дату или 0 для выхода: ")
		date, err := readInput(scanner)
		if err != nil {
			return fmt.Errorf("read appointments date: %w", err)
		}
		if date == "0" {
			return nil
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
		return nil
	}
}

func reschedulingAppointment(scanner *bufio.Scanner, appointments []appointment.Appointment) error {
	for {
		fmt.Println("Введите ID записи или 0 для выхода: ")
		id, err := readInput(scanner)
		if err != nil {
			return fmt.Errorf("read appointment ID: %w", err)
		}

		if id == "0" {
			return nil
		}

		fmt.Println("Введите новую дату записи: ")
		dateValue, err := readInput(scanner)
		if err != nil {
			return fmt.Errorf("read appointment date: %w", err)
		}

		fmt.Println("Введите новое время записи: ")
		timeValue, err := readInput(scanner)
		if err != nil {
			return fmt.Errorf("read appointment time: %w", err)
		}

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
			return fmt.Errorf("reschedule appointment: %w", err)
		}

		fmt.Println("Запись успешно перенесена")

		return nil
	}
}

func completeAppointment(scanner *bufio.Scanner, appointments []appointment.Appointment) error {
	for {
		fmt.Print("Введите ID записи для завершения или 0 для выхода: ")
		id, err := readInput(scanner)
		if err != nil {
			return fmt.Errorf("read appointment ID: %w", err)
		}
		if id == "0" {
			return nil
		}

		fmt.Println("Введите итоговую сумму: ")

		price, err := readInput(scanner)
		if err != nil {
			return fmt.Errorf("read appointment final price: %w", err)
		}

		finalPrice, err := strconv.Atoi(price)
		if err != nil {
			fmt.Println("Ошибка ввода суммы")
			continue
		}

		completedAppointment, err := appointment.Complete(appointments, id, finalPrice)

		switch {
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
		case err != nil:
			return fmt.Errorf("complete appointment: %w", err)
		}
		fmt.Println("Запись успешно завершена")
		fmt.Printf(
			"Планируемая стоимость: %d ₽\n",
			completedAppointment.PlannedPrice,
		)
		fmt.Printf(
			"Итоговая стоимость: %d ₽\n",
			completedAppointment.FinalPrice,
		)
		return nil
	}
}

func cancelAppointment(scanner *bufio.Scanner, appointments []appointment.Appointment) error {
	for {
		fmt.Print("Введите ID записи для отмены или 0 для выхода: ")
		id, err := readInput(scanner)
		if err != nil {
			return fmt.Errorf("read appointment ID: %w", err)
		}
		if id == "0" {
			return nil
		}

		err = appointment.Cancel(appointments, id)

		switch {
		case errors.Is(err, appointment.ErrAppointmentNotFound):
			fmt.Println("Запись не найдена")
			continue
		case errors.Is(err, appointment.ErrAppointmentCancelled):
			fmt.Println("Запись уже отменена")
			continue
		case errors.Is(err, appointment.ErrAppointmentCompleted):
			fmt.Println("Завершённую запись нельзя отменить")
			continue
		case err != nil:
			return fmt.Errorf("cancel appointment: %w", err)
		}
		fmt.Println("Запись успешно отменена")
		return nil
	}
}

func deleteAppointmentChoice(scanner *bufio.Scanner, appointments []appointment.Appointment, clients []client.Client) ([]appointment.Appointment, error) {
	for {
		fmt.Print("Введите ID записи для удаления или 0 для выхода: ")

		id, err := readInput(scanner)
		if err != nil {
			return appointments, fmt.Errorf("read appointment ID: %w", err)
		}

		if id == "0" {
			return appointments, nil
		}

		foundAppointment, found := appointment.FindByID(
			appointments,
			id,
		)
		if !found {
			fmt.Println("Запись не найдена")
			continue
		}

		foundClient, found := client.FindByID(
			clients,
			foundAppointment.ClientID,
		)
		if !found {
			fmt.Println("Клиент записи не найден")
			continue
		}

		printAppointment(foundAppointment, foundClient)

		for {
			fmt.Println()
			fmt.Println("1 — Удалить запись")
			fmt.Println("0 — Выход")
			fmt.Print("Выберите действие: ")

			choice, err := readInput(scanner)
			if err != nil {
				return appointments, fmt.Errorf(
					"read delete appointment choice: %w",
					err,
				)
			}

			switch choice {
			case "1":
				fmt.Println("Вы уверены?")
				fmt.Println("1 — Да")
				fmt.Println("0 — Нет")

				confirmation, err := readInput(scanner)
				if err != nil {
					return appointments, fmt.Errorf(
						"read delete appointment confirmation: %w",
						err,
					)
				}

				switch confirmation {
				case "1":
					afterDelete, err := appointment.Delete(
						appointments,
						id,
					)

					switch {
					case errors.Is(err, appointment.ErrAppointmentNotFound):
						fmt.Println("Запись не найдена")
						return appointments, nil

					case errors.Is(err, appointment.ErrAppointmentCompleted):
						fmt.Println("Нельзя удалить завершённую запись")
						return appointments, nil

					case errors.Is(err, appointment.ErrAppointmentCancelled):
						fmt.Println("Нельзя удалить отменённую запись")
						return appointments, nil

					case err != nil:
						return appointments, fmt.Errorf(
							"delete appointment: %w",
							err,
						)
					}

					fmt.Println("Запись удалена")
					return afterDelete, nil

				case "0":
					fmt.Println("Удаление отменено")
					continue

				default:
					fmt.Println("Неизвестная команда")
				}

			case "0":
				return appointments, nil

			default:
				fmt.Println("Неизвестная команда")
			}
		}
	}
}

func printClientHistory(scanner *bufio.Scanner, appointments []appointment.Appointment, clients []client.Client) error {
	for {
		fmt.Print("Введите ID клиента что бы увидеть историю сеансов или 0 для выхода: ")
		clientID, err := readInput(scanner)
		if err != nil {
			return fmt.Errorf("read client ID: %w", err)
		}

		if clientID == "0" {
			return nil
		}

		foundClient, isFound := client.FindByID(clients, clientID)
		if !isFound {
			fmt.Println("Клиент не найден")
			continue
		}

		clientHistory := appointment.FindByClientID(appointments, foundClient.ID)

		if len(clientHistory) == 0 {
			fmt.Println("У клиента еще нет записей")
			return nil
		}
		for _, v := range clientHistory {
			printAppointment(v, foundClient)
		}
		return nil
	}
}

func printUpcoming(appointments []appointment.Appointment, clients []client.Client) {
	upcomingAppointments := appointment.FindUpcoming(appointments, time.Now())
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
