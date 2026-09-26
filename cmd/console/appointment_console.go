package main

import (
	"bufio"
	"fmt"
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

		if appointment.IsTimeTaken(appointments, startAt) {
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
			fmt.Printf("%s\n", v.StartAt.Format("15:04"))
			clientAtDate, isExist := client.FindByID(clients, v.ClientID)
			if !isExist {
				fmt.Println("Ошибка поиска клиента") // пока как заглушку оставлю потом надо сделать норм обработку
				continue
			}
			fmt.Printf("Клиент: %s\n", clientAtDate.Name)
			if !clientAtDate.IsAdult {
				fmt.Println("Несовершеннолетний клиент")
			}
			if v.Description != "" {
				fmt.Printf("Описание: %s\n", v.Description)
			}
			fmt.Printf("Планируемая стоимость: %d\n", v.PlannedPrice)
			fmt.Printf("Предоплата: %d\n", v.Deposit)
			if v.FinalPrice != 0 {
				fmt.Printf("Финальная стоимость: %d\n", v.FinalPrice)
			}
			fmt.Printf("Статус: %s\n", v.Status)
			fmt.Println("-------------------")
		}
		return
	}
}
