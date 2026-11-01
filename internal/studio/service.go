package studio

import (
	"tattoo-studio/internal/appointment"
	"tattoo-studio/internal/client"
)

func DeleteClient(clients []client.Client, appointments []appointment.Appointment, clientID string) ([]client.Client, error) {
	_, found := client.FindByID(clients, clientID)
	if !found {
		return clients, client.ErrClientNotFound
	}

	if appointment.HasAppointmentsForClient(appointments, clientID) {
		return clients, ErrClientHasAppointments
	}

	return client.Delete(clients, clientID)
}

func CreateAppointment(clients []client.Client, appointments []appointment.Appointment, input appointment.CreateInput) (appointment.Appointment, error) {
	_, found := client.FindByID(clients, input.ClientID)
	if !found {
		return appointment.Appointment{}, client.ErrClientNotFound
	}

	return appointment.Create(appointments, input)
}
