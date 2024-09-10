package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("should load valid config from file", func(t *testing.T) {
		// Create a temporary directory for the config file
		tmpDir := os.TempDir()
		defer os.RemoveAll(tmpDir)

		// Create a temporary .env file with the config
		configContent := []byte("MONGO_URI=mongodb://localhost:27017\nSERVER_ADDR=localhost:8080")
		err := os.WriteFile(tmpDir+"/app.env", configContent, 0644)
		assert.NoError(t, err)

		// Load the config
		config, err := Load(tmpDir)
		assert.NoError(t, err)

		// Validate that the config was loaded correctly
		assert.Equal(t, "mongodb://localhost:27017", config.DBURI)
		assert.Equal(t, "localhost:8080", config.ServerAddr)
	})

	t.Run("should return error when config file is missing", func(t *testing.T) {
		// Use a non-existing directory
		config, err := Load("/non-existing-directory")

		// Assert that it returns an error
		assert.Error(t, err)
		assert.Equal(t, AppConfig{}, config)
	})
}
