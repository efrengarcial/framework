package tests

import (
	"context"
	"github.com/efrengarcial/framework/internal/domain"
	"github.com/efrengarcial/framework/internal/platform/auth"
	"github.com/efrengarcial/framework/internal/platform/cache"
	"github.com/efrengarcial/framework/internal/platform/database"
	"github.com/efrengarcial/framework/internal/user/repository"
	"github.com/efrengarcial/framework/internal/user/service"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

// TODO: NewUnit creates a test database inside a Docker container. It creates the
// required table structure but the database is otherwise empty.
//
// It does not return errors as this intended for testing only. Instead it will
// call Fatal on the provided testing.T if anything goes wrong.
//
// It returns the database to use as well as a function to call at the end of
// the test.
func NewUnit(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	db := database.Initialize("sqlite3", ":memory:?_loc=auto")
	t.Log("waiting for database to be ready")

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()
	}

	return db, teardown
}


// Test owns state for running and shutting down tests.
type Test struct {
	DB            *gorm.DB
	Log           *logrus.Logger
	Authenticator *auth.Authenticator
	Cache		  cache.Cache

	t       *testing.T
	cleanup func()
}

// NewIntegration creates a database, seeds it, constructs an authenticator.
func NewIntegration(t *testing.T) *Test {
	t.Helper()

	// Initialize and seed database. Store the cleanup function call later.
	db, cleanup := NewUnit(t)

	// Create the logger to use.
	logger := log.New()
	logger.Out = os.Stdout
	logger.Level = log.InfoLevel
	logger.Formatter = &log.JSONFormatter{}

	// Create RSA keys to enable authentication in our service.
	key := []byte("my-secret-token-to-change-in-production")
	// Build an authenticator using this static key.
	kid := "4754d86b-7a6d-4df5-9c65-224741361492"
	kf := auth.NewSimpleKeyLookupFunc(kid, key)
	authenticator, err := auth.NewAuthenticator(key, kid, "HS512", kf)
	if err != nil {
		t.Fatal(err)
	}
	cache := cache.NewInMemoryCache(time.Hour)

	return &Test{
		DB:            db,
		Log:           logger,
		Authenticator: authenticator,
		Cache:			cache,
		t:             t,
		cleanup:       cleanup,
	}
}

// Teardown releases any resources used for the test.
func (test *Test) Teardown() {
	test.cleanup()
}

// Token generates an authenticated token for a user.
func (test *Test) Token(userName, pass string) string {
	test.t.Helper()
	repo := repository.NewUserGormRepository(test.DB)
	auth := service.NewAuthService(repo, test.Authenticator, test.Log)
	loginVM := &domain.LoginVM{UserName:  userName, Password: pass }
	tkn := &domain.Token{}
	auth.Auth(context.Background(), loginVM, tkn )

	return tkn.Token
}

