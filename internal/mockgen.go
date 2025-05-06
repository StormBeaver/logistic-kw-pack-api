package internal

//go:generate mockgen -destination=./mocks/repo_mock.go -package=mocks route255/L-KW-P-API/internal/app/repo EventRepo
//go:generate mockgen -destination=./mocks/sender_mock.go -package=mocks route255/L-KW-P-API/internal/app/sender EventSender
