package secretscmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/badie/bdev/internal/core/config"
	"github.com/badie/bdev/internal/core/vault"
	"github.com/badie/bdev/pkg/ui"
)

var v *vault.Vault

func getVault() *vault.Vault {
	if v == nil {
		cfg := config.Get()
		v = vault.New(cfg.VaultFile())
	}
	return v
}

// NewCommand creates the secrets command group
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "secrets",
		Aliases: []string{"vault", "sec"},
		Short:   "Manage encrypted secrets",
		Long:    "Store and retrieve encrypted secrets using AES-256-GCM",
	}

	cmd.AddCommand(initCmd())
	cmd.AddCommand(setCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(lockCmd())

	return cmd
}

// ============================================================
// INIT - Create new vault
// ============================================================

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create a new encrypted vault",
		RunE: func(cmd *cobra.Command, args []string) error {
			vlt := getVault()

			if vlt.Exists() {
				return fmt.Errorf("vault already exists at %s", vlt.FilePath)
			}

			fmt.Print("Enter master password: ")
			password, err := readPassword()
			if err != nil {
				return err
			}

			fmt.Print("\nConfirm password: ")
			confirm, err := readPassword()
			if err != nil {
				return err
			}
			fmt.Println()

			if password != confirm {
				return fmt.Errorf("passwords do not match")
			}

			if len(password) < 8 {
				return fmt.Errorf("password must be at least 8 characters")
			}

			if err := vlt.Create(password); err != nil {
				return err
			}

			fmt.Println(ui.Success("Vault created successfully"))
			fmt.Println(ui.Muted("Store your master password safely - it cannot be recovered!"))
			return nil
		},
	}
}

// ============================================================
// SET - Store a secret
// ============================================================

func setCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> [value]",
		Short: "Store a secret",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			vlt := getVault()

			if err := ensureUnlocked(vlt); err != nil {
				return err
			}

			key := args[0]
			var value string

			if len(args) > 1 {
				value = args[1]
			} else {
				// Read value securely
				fmt.Printf("Enter value for %s: ", key)
				var err error
				value, err = readPasswordString()
				if err != nil {
					return err
				}
				fmt.Println()
			}

			if err := vlt.Set(key, value); err != nil {
				return err
			}

			fmt.Println(ui.Success("Secret stored: ") + ui.Primary(key))
			return nil
		},
	}
}

// ============================================================
// GET - Retrieve a secret
// ============================================================

func getCmd() *cobra.Command {
	var show bool

	cmd := &cobra.Command{
		Use:   "get <key>",
		Short: "Retrieve a secret",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			vlt := getVault()

			if err := ensureUnlocked(vlt); err != nil {
				return err
			}

			key := args[0]
			value, err := vlt.Get(key)
			if err != nil {
				return err
			}

			if show {
				fmt.Println(value)
			} else {
				fmt.Printf("%s %s\n", ui.Primary(key+":"), ui.Muted("[hidden - use --show to reveal]"))
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&show, "show", false, "Show the secret value")
	return cmd
}

// ============================================================
// LIST - List all secrets
// ============================================================

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all stored secrets",
		RunE: func(cmd *cobra.Command, args []string) error {
			vlt := getVault()

			if err := ensureUnlocked(vlt); err != nil {
				return err
			}

			keys, err := vlt.List()
			if err != nil {
				return err
			}

			if len(keys) == 0 {
				fmt.Println(ui.Muted("No secrets stored"))
				return nil
			}

			fmt.Println(ui.Bold(fmt.Sprintf("Secrets (%d)", len(keys))))
			for _, k := range keys {
				fmt.Printf("  %s\n", ui.Primary(k))
			}

			return nil
		},
	}
}

// ============================================================
// DELETE - Remove a secret
// ============================================================

func deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <key>",
		Short: "Delete a secret",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			vlt := getVault()

			if err := ensureUnlocked(vlt); err != nil {
				return err
			}

			key := args[0]
			if err := vlt.Delete(key); err != nil {
				return err
			}

			fmt.Println(ui.Success("Secret deleted: ") + ui.Primary(key))
			return nil
		},
	}
}

// ============================================================
// LOCK - Lock the vault
// ============================================================

func lockCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "lock",
		Short: "Lock the vault",
		Run: func(cmd *cobra.Command, args []string) {
			vlt := getVault()
			vlt.Lock()
			fmt.Println(ui.Success("Vault locked"))
		},
	}
}

// ============================================================
// HELPERS
// ============================================================

func ensureUnlocked(vlt *vault.Vault) error {
	if !vlt.Exists() {
		return fmt.Errorf("vault not initialized. Run: bdev secrets init")
	}

	if vlt.IsUnlocked() {
		return nil
	}

	fmt.Print("Enter master password: ")
	password, err := readPassword()
	if err != nil {
		return err
	}
	fmt.Println()

	return vlt.Unlock(password)
}

func readPassword() (string, error) {
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	return string(password), nil
}

func readPasswordString() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(password), nil
}
