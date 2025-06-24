package cardcommand

import (
	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/spf13/cobra"
)

const (
	MasterKeyFlag  = command.MasterKeyFlag
	NameFlag       = command.NameFlag
	EntryNumFlag   = command.EntryNumFlag
	CardNumFlag    = "number"
	ExpDateFlag    = "exp"
	HolderNameFlag = "holder"
)

const (
	masterKeyShorthand = command.MasterKeyShorthand
	masterKeyDefault   = command.MasterKeyDefault
	masterKeyUsage     = command.MasterKeyUsage

	nameShorthand = command.NameShorthand
	nameDefault   = command.NameDefault
	nameUsage     = "title for bank card (required)"

	entryNumShorthand = command.EntryNumShorthand
	entryNumDefault   = command.EntryNumDefault
	entryNumUsage     = "entry number of stored bank card (required)"

	cardNumDefault = ""
	cardNumUsage   = "bank card number (required)"

	expDateDefault = ""
	expDateUsage   = "bank card validity period e.g. 12/2025 (required)"

	holderNameDefault = ""
	hoderNameUsage    = "cardhodler name on the bank card (required)"
)

func New() *command.Command {
	c := &cobra.Command{
		Use:   "card",
		Short: "Use the card command to save your bank cards",
	}
	return &command.Command{Command: c}
}

func NewAdd(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "add",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}

	flagSet := c.Flags()

	flagSet.StringP(
		MasterKeyFlag, masterKeyShorthand, masterKeyDefault, masterKeyUsage,
	)

	flagSet.StringP(NameFlag, nameShorthand, nameDefault, nameUsage)

	flagSet.String(
		CardNumFlag, cardNumDefault, cardNumUsage,
	)

	flagSet.String(
		ExpDateFlag, expDateDefault, expDateUsage,
	)

	flagSet.String(
		HolderNameFlag, holderNameDefault, hoderNameUsage,
	)

	c.MarkFlagRequired(MasterKeyFlag)
	c.MarkFlagRequired(NameFlag)
	c.MarkFlagRequired(CardNumFlag)
	c.MarkFlagRequired(ExpDateFlag)
	c.MarkFlagRequired(HolderNameFlag)

	return &command.Command{Command: c}
}

func NewRead(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "read",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}
	flagSet := c.Flags()

	flagSet.StringP(
		MasterKeyFlag, masterKeyShorthand, masterKeyDefault, masterKeyUsage,
	)

	flagSet.IntP(
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage,
	)

	c.MarkFlagRequired(MasterKeyFlag)
	c.MarkFlagRequired(EntryNumFlag)

	return &command.Command{Command: c}
}

func NewList(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}
	return &command.Command{Command: c}
}

func NewEdit(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "edit",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}

	flagSet := c.Flags()

	flagSet.StringP(
		MasterKeyFlag, masterKeyShorthand, masterKeyDefault, masterKeyUsage,
	)

	flagSet.IntP(
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage,
	)

	flagSet.StringP(NameFlag, nameShorthand, nameDefault, nameUsage)

	flagSet.String(
		CardNumFlag, cardNumDefault, cardNumUsage,
	)

	flagSet.String(
		ExpDateFlag, expDateDefault, expDateUsage,
	)

	flagSet.String(
		HolderNameFlag, holderNameDefault, hoderNameUsage,
	)

	c.MarkFlagRequired(MasterKeyFlag)
	c.MarkFlagRequired(EntryNumFlag)
	c.MarkFlagRequired(NameFlag)
	c.MarkFlagRequired(CardNumFlag)
	c.MarkFlagRequired(ExpDateFlag)
	c.MarkFlagRequired(HolderNameFlag)
	return &command.Command{Command: c}
}

func NewRemove(h command.Handler) *command.Command {
	c := &cobra.Command{
		Use: "remove",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), cmd.Flags())
		},
	}

	c.Flags().IntP(
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage,
	)

	c.MarkFlagRequired(EntryNumFlag)

	return &command.Command{Command: c}
}
