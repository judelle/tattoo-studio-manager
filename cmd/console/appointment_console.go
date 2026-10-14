package main

import (
	"bufio"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"tattoo-studio/internal/appointment"
	"tattoo-studio/internal/client"
	"tattoo-studio/internal/datetime"
	"time"

	"github.com/google/uuid"
)

func addAppointment(scanner *bufio.Scanner, clients []client.Client, appointments []appointment.Appointment) (appointment.Appointment, bool) {
	var err error
	clientID := ""
	for {
		fmt.Print("Введите ID клиента для записи или 0 для выхода: ")
		scanner.Scan()
		if scanner.Text() == "0" {
			return appointment.Appointment{}, false
		}
		foundClient, isExist := client.FindByID(clients, scanner.Text())
		if !isExist {
			fmt.Println("Клиента с таким ID не существует")
			continue
		}
		clientID = foundClient.ID
		break
	}
	var startAt time.Time
	for {
		fmt.Println("Введите дату записи: ")
		scanner.Scan()
		date := scanner.Text()
		fmt.Println("Введите Время записи: ")
		scanner.Scan()
		timeValue := scanner.Text()
		startAt, err = time.ParseInLocation(
			"02.01.2006 15:04",
			date+" "+timeValue,
			datetime.Moscow,
		)
		if err != nil {
			fmt.Println("Дата и время введеные некорректно")
			continue
		}
		if startAt.Before(time.Now().In(datetime.Moscow)) {
			fmt.Println("Запись не может быть создана в прошлом")
			continue
		}

		if appointment.IsTimeTaken(appointments, startAt, "") {
			fmt.Println("Время уже занято")
			continue
		}
		break
	}

	fmt.Println("Введите описание: ")
	scanner.Scan()
	description := strings.TrimSpace(scanner.Text())
	var plannedPrice int
	for {
		fmt.Println("Введите планируемую цену или -1 для выхода: ")
		scanner.Scan()
		if scanner.Text() == "-1" {
			return appointment.Appointment{}, false
		}
		plannedPrice, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Ошибка ввода данных")
			continue
		}
		if plannedPrice < 0 {
			fmt.Println("Цена не может быть меньше 0")
			continue
		}
		break
	}
	var deposit int
	for {
		fmt.Println("Введите депозит или -1 для выхода: ")
		scanner.Scan()
		if scanner.Text() == "-1" {
			return appointment.Appointment{}, false
		}
		deposit, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Ошибка ввода данных")
			continue
		}
		if deposit < 0 {
			fmt.Println("Депозит не может быть меньше 0")
			continue
		}
		break
	}

	return appointment.Appointment{
		ID:           uuid.NewString(),
		ClientID:     clientID,
		StartAt:      startAt,
		Description:  description,
		PlannedPrice: plannedPrice,
		Deposit:      deposit,
		Status:       appointment.StatusPlanned,
	}, true
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

func reschedulingAppointment(scanner *bufio.Scanner, appointments []appointment.Appointment) {
	for {
		fmt.Println("Введите ID записи или 0 для выхода: ")
		scanner.Scan()
		id := scanner.Text()

		if id == "0" {
			return
		}

		idx, ok := appointment.FindIndexByID(appointments, id)
		if !ok {
			fmt.Println("Запись не найдена")
			continue
		}

		if appointments[idx].Status == appointment.StatusCompleted {
			fmt.Println("Нельзя перенести выполненную запись")
			continue
		}

		fmt.Println("Введите новую дату записи: ")
		scanner.Scan()
		date := scanner.Text()

		fmt.Println("Введите новое время записи: ")
		scanner.Scan()
		timeValue := scanner.Text()

		startAt, err := time.ParseInLocation(
			"02.01.2006 15:04",
			date+" "+timeValue,
			datetime.Moscow,
		)
		if err != nil {
			fmt.Println("Дата и время введены некорректно")
			continue
		}

		if startAt.Before(time.Now().In(datetime.Moscow)) {
			fmt.Println("Запись не может быть перенесена в прошлое")
			continue
		}

		if appointment.IsTimeTaken(appointments, startAt, appointments[idx].ID) {
			fmt.Println("Время уже занято")
			continue
		}

		if appointments[idx].Status == appointment.StatusCancelled {
			appointments[idx].Status = appointment.StatusPlanned
		}

		appointments[idx].StartAt = startAt

		fmt.Println("Запись успешно перенесена")
		return
	}
}

func completeAppointment(scanner *bufio.Scanner, appointments []appointment.Appointment) {
	for {
		fmt.Print("Введите ID записи для завершения или 0 для выхода: ")
		scanner.Scan()

		if scanner.Text() == "0" {
			return
		}

		idx, found := appointment.FindIndexByID(
			appointments,
			scanner.Text(),
		)
		if !found {
			fmt.Println("Запись не найдена")
			continue
		}

		if appointments[idx].Status == appointment.StatusCancelled {
			fmt.Println("Отменённую запись нельзя завершить")
			continue
		}

		if appointments[idx].Status == appointment.StatusCompleted {
			fmt.Println("Запись уже завершена")
			continue
		}

		fmt.Println("Введите итоговую сумму: ")
		scanner.Scan()

		finalPrice, err := strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Ошибка ввода суммы")
			continue
		}

		if finalPrice < 0 {
			fmt.Println("Итоговая сумма не может быть меньше 0")
			continue
		}

		if finalPrice < appointments[idx].Deposit {
			fmt.Println("Итоговая сумма меньше депозита. Вы уверены? Д/Н")
			scanner.Scan()

			answer := strings.TrimSpace(
				strings.ToLower(scanner.Text()),
			)

			if answer != "д" {
				continue
			}
		}

		appointments[idx].FinalPrice = finalPrice
		appointments[idx].Status = appointment.StatusCompleted

		fmt.Println("Запись успешно завершена")
		fmt.Printf(
			"Планируемая стоимость: %d ₽\n",
			appointments[idx].PlannedPrice,
		)
		fmt.Printf(
			"Итоговая стоимость: %d ₽\n",
			appointments[idx].FinalPrice,
		)

		return
	}
}

func cancelAppointment(scanner *bufio.Scanner, appointments []appointment.Appointment) {
	for {
		fmt.Print("Введите ID записи для отмены или 0 для выхода: ")
		scanner.Scan()

		if scanner.Text() == "0" {
			return
		}

		idx, found := appointment.FindIndexByID(
			appointments,
			scanner.Text(),
		)
		if !found {
			fmt.Println("Запись не найдена")
			continue
		}

		if appointments[idx].Status == appointment.StatusCancelled {
			fmt.Println("Запись уже отменена")
			continue
		}

		if appointments[idx].Status == appointment.StatusCompleted {
			fmt.Println("Завершённую запись нельзя отменить")
			continue
		}

		appointments[idx].Status = appointment.StatusCancelled
		fmt.Println("Запись успешно отменена")
		return
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
