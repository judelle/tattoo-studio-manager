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
