package models

type JsonResponse struct {
	Success bool	`json:"success"`
	Message string	`json:"message"`
	Payload any		`json:"payload"`
}
