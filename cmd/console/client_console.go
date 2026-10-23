package main

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
	"tattoo-studio/internal/appointment"
	"tattoo-studio/internal/client"
	"tattoo-studio/internal/studio"
)

func addClient(scanner *bufio.Scanner, clients []client.Client) (client.Client, bool, error) {
	for {
		fmt.Print("Введите имя или 0 для выхода: ")
		name, err := readInput(scanner)
		if err != nil {
			return client.Client{}, false, fmt.Errorf("read client name: %w", err)
		}
		if name == "0" {
			return client.Client{}, false, nil
		}

		fmt.Print("Клиент совершеннолетний? Д/Н: ")

		adultInput, err := readInput(scanner)
		if err != nil {
			return client.Client{}, false, fmt.Errorf("read client adult input: %w", err)
		}
		adultInput = strings.ToLower(adultInput)
		isAdult := false
		if adultInput == "д" {
			isAdult = true
		} else if adultInput != "н" {
			fmt.Println("Некорректный ввод")
			continue
		}

		fmt.Print("Введите номер телефона: ")
		phone, err := readInput(scanner)
		if err != nil {
			return client.Client{}, false, fmt.Errorf("read client phone: %w", err)
		}

		fmt.Print("Введите соц. сеть: ")
		social, err := readInput(scanner)
		if err != nil {
			return client.Client{}, false, fmt.Errorf("read client social: %w", err)
		}
		input := client.CreateInput{
			Name:    name,
			IsAdult: isAdult,
			Phone:   phone,
			Social:  social,
		}
		createdClient, err := client.Create(clients, input)

		switch {
		case errors.Is(err, client.ErrNameRequired):
			fmt.Println("Имя не может быть пустым")
			continue
		case errors.Is(err, client.ErrPhoneRequired):
			fmt.Println("Номер не может быть пустым")
			continue
		case errors.Is(err, client.ErrPhoneAlreadyExists):
			fmt.Printf("Номер телефона %s занят пользователем: %s, с ID %s\n", createdClient.Phone, createdClient.Name, createdClient.ID)
			continue
		case err != nil:
			return client.Client{}, false, fmt.Errorf("create client: %w", err)
		}

		fmt.Println("Клиент успешно добавлен")
		return createdClient, true, nil
	}
}

func printClients(clients []client.Client) {
	if len(clients) == 0 {
		fmt.Println("Вы еще не добавили клиентов")
		return
	}
	fmt.Printf("Список из %d клиентов: \n", len(clients))
	for _, client := range clients {
		printClient(client)
	}
}

func printClient(client client.Client) {
	fmt.Println("-------------------")
	fmt.Println("ID:", client.ID)
	fmt.Println("Имя:", client.Name)

	if client.IsAdult {
		fmt.Println("Совершеннолетний")
	} else {
		fmt.Println("Несовершеннолетний")
	}

	fmt.Println("Телефон:", client.Phone)
	fmt.Println("Соц. сеть:", client.Social)
}

func findClient(scanner *bufio.Scanner, clients []client.Client) error {
	fmt.Print("Введите ID клиента или 0 для выхода: ")

	fmt.Println("-------------------")
	id, err := readInput(scanner)
	if err != nil {
		return fmt.Errorf("read client ID: %w", err)
	}
	if id == "0" {
		return nil
	}
	client, isFound := client.FindByID(clients, id)
	if !isFound {
		fmt.Printf("клиент с ID %s не найден\n", id)
		fmt.Println("-------------------")
		return nil
	}
	printClient(client)
	return nil
}

func changeClient(scanner *bufio.Scanner, clients []client.Client) error {
	for {
		fmt.Print("Введите ID клиента для изменения или 0 для выхода: ")

		id, err := readInput(scanner)
		if err != nil {
			return fmt.Errorf("read client ID: %w", err)
		}
		if id == "0" {
			return nil
		}

		index, isFound := client.FindIndexByID(clients, id)
		if !isFound {
			fmt.Printf("клиент с ID %s не найден\n", id)
			continue
		}
		for {
			fmt.Println()
			fmt.Println("1 — Изменить имя")
			fmt.Println("2 — Изменить номер")
			fmt.Println("3 — Изменить соц. Сети")
			fmt.Println("4 — Изменить статус совершеннолетнего")
			fmt.Println("0 — Выход")
			fmt.Print("Выберите действие: ")

			infoToChange, err := readInput(scanner)
			if err != nil {
				return fmt.Errorf("read infotochange: %w", err)
			}

			switch infoToChange {
			case "1":
				for {
					fmt.Println("Введите новое имя или 0 для возврата")
					newName, err := readInput(scanner)
					if err != nil {
						return fmt.Errorf("read name: %w", err)
					}
					if newName == "0" {
						break
					}
					err = client.ChangeName(clients, id, newName)

					switch {
					case errors.Is(err, client.ErrNameRequired):
						fmt.Println("Имя не может быть пустым")
						continue

					case errors.Is(err, client.ErrClientNotFound):
						fmt.Println("Клиент не найден")
						return nil

					case err != nil:
						return fmt.Errorf("change client name: %w", err)
					}

					break
				}
			case "2":
				for {
					fmt.Print("Введите номер телефона или 0 для возврата: ")
					newPhone, err := readInput(scanner)
					if err != nil {
						return fmt.Errorf("read phone: %w", err)
					}
					if newPhone == "0" {
						break
					}
					err = client.ChangePhone(clients, id, newPhone)

					switch {
					case errors.Is(err, client.ErrPhoneRequired):
						fmt.Println("Номер телефона не может быть пустым")
						continue

					case errors.Is(err, client.ErrClientNotFound):
						fmt.Println("Клиент не найден")
						return nil

					case errors.Is(err, client.ErrPhoneAlreadyExists):
						fmt.Println("Номер телефона занят другим пользователем")
						continue

					case err != nil:
						return fmt.Errorf("change client phone: %w", err)
					}

					break
				}
			case "3":
				fmt.Println("Введите новую соц. сеть или 0 для возврата")
				newSocial, err := readInput(scanner)
				if err != nil {
					return fmt.Errorf("read social: %w", err)
				}
				if newSocial == "0" {
					break
				}
				err = client.ChangeSocial(clients, id, newSocial)
				if err != nil {
					if errors.Is(err, client.ErrClientNotFound) {
						fmt.Println("Клиент не найден")
						return nil
					}

					return fmt.Errorf("change client social: %w", err)
				}
			case "4":
				for {
					fmt.Print("Клиент совершеннолетний? Д/Н или 0 для возврата: ")

					newIsAdult, err := readInput(scanner)
					if err != nil {
						return fmt.Errorf("read client adult status: %w", err)
					}

					if newIsAdult == "0" {
						break
					}

					newIsAdult = strings.ToLower(newIsAdult)

					var isAdult bool

					switch newIsAdult {
					case "д":
						isAdult = true

					case "н":
						isAdult = false

					default:
						fmt.Println("Некорректный ввод")
						continue
					}

					err = client.ChangeAdultStatus(clients, id, isAdult)
					if err != nil {
						if errors.Is(err, client.ErrClientNotFound) {
							fmt.Println("Клиент не найден")
							break
						}

						return fmt.Errorf("change client adult status: %w", err)
					}

					break
				}
			case "0":
				printClient(clients[index])
				return nil
			default:
				fmt.Println("Неизвестная команда")
			}
		}
	}

}

func deleteClientChoice(scanner *bufio.Scanner, clients []client.Client, appointments []appointment.Appointment) ([]client.Client, error) {
	for {
		fmt.Print("Введите ID клиента для удаления или 0 для выхода: ")

		id, err := readInput(scanner)
		if err != nil {
			return clients, fmt.Errorf("read client ID: %w", err)
		}

		if id == "0" {
			return clients, nil
		}

		foundClient, found := client.FindByID(clients, id)
		if !found {
			fmt.Printf("Клиент с ID %s не найден\n", id)
			continue
		}

		printClient(foundClient)

		for {
			fmt.Println()
			fmt.Println("1 — Удалить клиента")
			fmt.Println("0 — Выход")
			fmt.Print("Выберите действие: ")

			deleteChoice, err := readInput(scanner)
			if err != nil {
				return clients, fmt.Errorf(
					"read delete client choice: %w",
					err,
				)
			}

			switch deleteChoice {
			case "1":
				fmt.Println("Вы уверены?")
				fmt.Println("1 — Да")
				fmt.Println("0 — Нет")

				confirmation, err := readInput(scanner)
				if err != nil {
					return clients, fmt.Errorf(
						"read delete client confirmation: %w",
						err,
					)
				}

				switch confirmation {
				case "1":
					afterDelete, err := studio.DeleteClient(
						clients,
						appointments,
						id,
					)

					switch {
					case errors.Is(err, client.ErrClientNotFound):
						fmt.Println("Клиент не найден")
						return clients, nil

					case errors.Is(err, studio.ErrClientHasAppointments):
						fmt.Println(
							"Нельзя удалить клиента, у него есть история записей",
						)
						return clients, nil

					case err != nil:
						return clients, fmt.Errorf(
							"delete client: %w",
							err,
						)
					}

					fmt.Println("Клиент удалён")
					return afterDelete, nil

				case "0":
					fmt.Println("Удаление отменено")
					continue

				default:
					fmt.Println("Неизвестная команда")
				}

			case "0":
				return clients, nil

			default:
				fmt.Println("Неизвестная команда")
			}
		}
	}
}
