package keys

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var newPasswordFile string

//AddUpdateFlags adds flags to the update command
func AddUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&newPasswordFile, "new_passfile", "", "the file containing the new passphrase for the keyfile")
	viper.BindPFlags(cmd.Flags())
}

//NewUpdateCmd returns the command that changes the passphrase of a keyfile
func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [keyfile]",
		Short: "change the passphrase on a keyfile",
		Args:  cobra.ExactArgs(1),
		RunE:  update,
	}

	AddUpdateFlags(cmd)

	return cmd
}

func update(cmd *cobra.Command, args []string) error {
	keyfilepath := args[0]

	// Read key from file.
	keyjson, err := ioutil.ReadFile(keyfilepath)
	if err != nil {
		return fmt.Errorf("Failed to read the keyfile at '%s': %v", keyfilepath, err)
	}

	// Decrypt key with passphrase.
	passphrase, err := getPassphrase()
	if err != nil {
		return err
	}

	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		return fmt.Errorf("Error decrypting key: %v", err)
	}

	// Get a new passphrase.
	fmt.Println("Please provide a new passphrase")
	var newPhrase string
	if newPasswordFile != "" {
		content, err := ioutil.ReadFile(newPasswordFile)
		if err != nil {
			return fmt.Errorf("Failed to read new passphrase file '%s': %v", newPasswordFile, err)
		}
		newPhrase = strings.TrimRight(string(content), "\r\n")
	} else {
		newPhrase, err = promptPassphrase(true)
		if err != nil {
			return err
		}
	}

	// Encrypt the key with the new passphrase.
	newJSON, err := keystore.EncryptKey(key, newPhrase, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return fmt.Errorf("Error encrypting with new passphrase: %v", err)
	}

	// Then write the new keyfile in place of the old one.
	if err := ioutil.WriteFile(keyfilepath, newJSON, 600); err != nil {
		return fmt.Errorf("Error writing new keyfile to disk: %v", err)
	}

	// Don't print anything.  Just return successfully,
	// producing a positive exit code.
	return nil
}
