package usecase

import (
	"gade/srv-goldcard/activations"
)

type activationsUseCase struct {
	aRepo activations.Repository
}

// RegistrationsUseCase represent Registrations Use Case
func RegistrationsUseCase(aRepo activations.Repository) activations.UseCase {
	return &activationsUseCase{aRepo: aRepo}
}
