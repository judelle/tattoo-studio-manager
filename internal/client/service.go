package client

import (
	"slices"

	"github.com/google/uuid"
)

func FindByID(clients []Client, id string) (Client, bool) {
	for _, client := range clients {
		if client.ID == id {
			return client, true
		}
	}

	return Client{}, false
}

func FindIndexByID(clients []Client, id string) (int, bool) {
	for i, client := range clients {
		if client.ID == id {
			return i, true
		}
	}

	return -1, false
}

func FindByPhone(clients []Client, phone string) (Client, bool) {
	for _, client := range clients {
		if client.Phone == phone {
			return client, true
		}
	}
	return Client{}, false
}

func Create(clients []Client, input CreateInput) (Client, error) {
	if input.Name == "" {
		return Client{}, ErrNameRequired
	}

	if input.Phone == "" {
		return Client{}, ErrPhoneRequired
	}

	foundClient, isFound := FindByPhone(clients, input.Phone)
	if isFound {
		return foundClient, ErrPhoneAlreadyExists
	}

	return Client{
		ID:      uuid.NewString(),
		Name:    input.Name,
		IsAdult: input.IsAdult,
		Phone:   input.Phone,
		Social:  input.Social,
	}, nil
}

func ChangeName(clients []Client, id, newName string) error {
	if newName == "" {
		return ErrNameRequired
	}
	idx, isFound := FindIndexByID(clients, id)
	if !isFound {
		return ErrClientNotFound
	}

	clients[idx].Name = newName
	return nil
}

func ChangePhone(clients []Client, id, newPhone string) error {
	if newPhone == "" {
		return ErrPhoneRequired
	}
	idx, isFound := FindIndexByID(clients, id)
	if !isFound {
		return ErrClientNotFound
	}
	foundClient, isFound := FindByPhone(clients, newPhone)
	if isFound && clients[idx].ID != foundClient.ID {
		return ErrPhoneAlreadyExists
	}
	clients[idx].Phone = newPhone
	return nil
}

func ChangeAdultStatus(clients []Client, id string, isAdult bool) error {
	idx, isFound := FindIndexByID(clients, id)
	if !isFound {
		return ErrClientNotFound
	}

	clients[idx].IsAdult = isAdult
	return nil
}

func ChangeSocial(clients []Client, id, newSocial string) error {
	idx, isFound := FindIndexByID(clients, id)
	if !isFound {
		return ErrClientNotFound
	}

	clients[idx].Social = newSocial
	return nil
}

func Delete(clients []Client, id string) ([]Client, error) {
	idx, isFound := FindIndexByID(clients, id)
	if !isFound {
		return clients, ErrClientNotFound
	}

	clients = slices.Delete(clients, idx, idx+1)
	return clients, nil
}
