package tests

import (
	"log"
	"rent-video-game/model"
	"rent-video-game/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return gormDB, mock
}

func TestRegisterUser(t *testing.T) {
	db, mock := NewMockDB()
	repo := repository.NewUserRepository(db)

	testUser := &model.Users{
		Name:     "test",
		Email:    "test@example.com",
		Password: "test_password",
		Address:  "address",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(testUser.UserID))

	mock.ExpectCommit()

	result, err := repo.RegisterUser(testUser)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testUser.Name, result.Name)
	assert.Equal(t, testUser.Email, result.Email)
	assert.Equal(t, testUser.Password, result.Password)
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
