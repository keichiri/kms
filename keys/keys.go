package keys

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os/exec"
)

type KeysManager struct {
}

func NewKeysManager() *KeysManager {
	return &KeysManager{}
}

func (km *KeysManager) CreateNewKeypair(id string) (*KeyPair, error) {
	cmd := exec.Command("divined", "keys", "add", id, "--output", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to execute command to create keypair")
	}

	var output Output
	if err := json.Unmarshal(out, &output); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal command output")
	}

	return &KeyPair{Address: output.Address}, nil
}

type Output struct {
	Address  string `json:"address"`
	Mnemonic string `json:"mnemonic"`
}

type KeyPair struct {
	Address string
}
