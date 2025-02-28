mockgen:
	mockgen -destination=./mocks/mock_user_repository.go -package=mocks rent-video-game/repository IUserRepository \
	&& mockgen -destination=./mocks/mock_user_usecase.go -package=mocks rent-video-game/usecase IUserUsecase \
	&& mockgen -destination=./mocks/mock_user_handler.go -package=mocks rent-video-game/handler IUserHandler \
	&& mockgen -destination=./mocks/mock_topup_history_usecase.go -package=mocks rent-video-game/usecase ITopupHistoryUsecase
	

test:
	go test -cover -v ./...