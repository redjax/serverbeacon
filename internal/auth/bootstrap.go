package auth

import "gorm.io/gorm"

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Token{}, &RegistrationCode{}, &BootstrapState{})
}

// EnsureBootstrap checks if the database has been initialized,
// and does the initialization if not.
func EnsureBootstrap(s *Store) (map[string]string, error) {
	out := map[string]string{}

	// Check if bootstrap indicator exists
	exists, err := s.BootstrapExists()
	if err != nil {
		return nil, err
	}

	if exists {
		return out, nil
	}

	// Create initial 'watcher' user
	user, err := s.CreateUser("watcher")
	if err != nil {
		return nil, err
	}

	// Generate secret for watcher user's token
	watcherToken, err := NewSecret(32)
	if err != nil {
		return nil, err
	}

	// Generate watcher's token
	if _, err := s.CreateToken(user.ID, "watcher", HashSecret(watcherToken), nil); err != nil {
		return nil, err
	}
	out["watcher_token"] = watcherToken

	// Generate secret for admin user's token
	adminToken, err := NewSecret(32)
	if err != nil {
		return nil, err
	}

	// Generate admin's token
	if _, err := s.CreateToken(user.ID, "admin", HashSecret(adminToken), nil); err != nil {
		return nil, err
	}
	out["admin_token"] = adminToken

	// Generate registration token secret
	regCode, err := NewSecret(24)
	if err != nil {
		return nil, err
	}

	// Generate registration token
	if _, err := s.CreateRegistrationCode(HashSecret(regCode), 1, nil); err != nil {
		return nil, err
	}
	out["registration_code"] = regCode

	// Mark database as 'bootstrapped'/initialized
	if err := s.MarkBootstrapped(); err != nil {
		return nil, err
	}

	return out, nil
}
