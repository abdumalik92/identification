package service

import (
	"github.com/abdumalik92/identification/internal/repository"
)

func SendOrderToCFT(product string, clientID string, fileLink string) error {

	return repository.SendOrderToCFT(product, clientID, fileLink)
}
