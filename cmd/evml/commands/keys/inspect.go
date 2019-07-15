package keys

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type outputInspect struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

var showPrivate = false

//AddInspectFlags adds flags to the Inspect command
func AddInspectFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(&showPrivate, "private", false, "include the private key in the output")
	viper.BindPFlags(cmd.Flags())
}

//NewInspectCmd returns the command that inspects an Ethereum keyfile
func NewInspectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "inspect [keyfile]",
		Short: "Inspect a keyfile",
		Long: `
Print various information about the keyfile.

Private key information can be printed by using the --private flag;
make sure to use this feature with great caution!`,
		Args: cobra.ExactArgs(1),
		RunE: inspect,
	}

	AddInspectFlags(cmd)

	return cmd
}

func inspect(cmd *cobra.Command, args []string) error {
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

	// Output all relevant information we can retrieve.
	out := outputInspect{
		Address: key.Address.Hex(),
		PublicKey: hex.EncodeToString(
			crypto.FromECDSAPub(&key.PrivateKey.PublicKey)),
	}
	if showPrivate {
		out.PrivateKey = hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))
	}

	if outputJSON {
		mustPrintJSON(out)
	} else {
		fmt.Println("Address:       ", out.Address)
		fmt.Println("Public key:    ", out.PublicKey)
		if showPrivate {
			fmt.Println("Private key:   ", out.PrivateKey)
		}
	}

	return nil
}
