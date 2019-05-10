package keys

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pborman/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultKeyfile = "keyfile.json"

	privateKeyfile string
)

type outputGenerate struct {
	Address      string
	AddressEIP55 string
}

//AddGenerateFlags adds flags to the Generate command
func AddGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&privateKeyfile, "privatekey", "", "file containing a raw private key to encrypt")

	viper.BindPFlags(cmd.Flags())
}

//NewGenerateCmd returns the command that creates a Ethereum keyfile
func NewGenerateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate [keyfile]",
		Short: "Generate a new keyfile",
		Long: `
Generate a new keyfile.

If you want to encrypt an existing private key, it can be specified by setting
--privatekey with the location of the file containing the private key.
`,
		Args: cobra.ExactArgs(1),
		RunE: generate,
	}

	AddGenerateFlags(cmd)

	return cmd
}

func generate(cmd *cobra.Command, args []string) error {
	// Check if keyfile path given and make sure it doesn't already exist.
	keyfilepath := args[0]
	if keyfilepath == "" {
		keyfilepath = defaultKeyfile
	}
	if _, err := os.Stat(keyfilepath); err == nil {
		utils.Fatalf("Keyfile already exists at %s.", keyfilepath)
	} else if !os.IsNotExist(err) {
		utils.Fatalf("Error checking if keyfile exists: %v", err)
	}

	var privateKey *ecdsa.PrivateKey
	var err error
	if file := privateKeyfile; file != "" {
		// Load private key from file.
		privateKey, err = crypto.LoadECDSA(file)
		if err != nil {
			utils.Fatalf("Can't load private key: %v", err)
		}
	} else {
		// If not loaded, generate random.
		privateKey, err = crypto.GenerateKey()
		if err != nil {
			utils.Fatalf("Failed to generate random private key: %v", err)
		}
	}

	// Create the keyfile object with a random UUID.
	id := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}

	// Encrypt key with passphrase.
	passphrase := promptPassphrase(true)
	keyjson, err := keystore.EncryptKey(key, passphrase, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		utils.Fatalf("Error encrypting key: %v", err)
	}

	// Store the file to disk.
	if err := os.MkdirAll(filepath.Dir(keyfilepath), 0700); err != nil {
		utils.Fatalf("Could not create directory %s", filepath.Dir(keyfilepath))
	}
	if err := ioutil.WriteFile(keyfilepath, keyjson, 0600); err != nil {
		utils.Fatalf("Failed to write keyfile to %s: %v", keyfilepath, err)
	}

	// Output some information.
	out := outputGenerate{
		Address: key.Address.Hex(),
	}

	if outputJSON {
		mustPrintJSON(out)
	} else {
		fmt.Println("Address:", out.Address)
	}

	return nil
}
