package cardcommand

import (
	"github.com/niksmo/gophkeeper/internal/client/command"
	"github.com/spf13/cobra"
)

const (
	SecretKeyFlag  = command.SecreKeyFlag
	NameFlag       = command.NameFlag
	EntryNumFlag   = command.EntryNumFlag
	CardNumFlag    = "number"
	ExpDateFlag    = "exp"
	HolderNameFlag = "holder"
)

const (
	secretKeyShorthand = command.SecretKeyShorthand
	secretKeyDefault   = command.SecretKeyDefault
	secretKeyUsage     = command.SecretKeyUsage

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
	holderNameUsage   = "cardhodler name on the bank card (required)"
)

func New() *command.Command {
	c := &cobra.Command{
		Use:   "card",
		Short: "Use the card command to save your bank cards",
	}
	return &command.Command{Command: c}
}

type AddCmdFlags struct {
	Key, Name, CardNum, Exp, Holder string
}

func NewAdd(h command.GenCmdHandler[AddCmdFlags]) *command.Command {
	var fv AddCmdFlags

	c := &cobra.Command{
		Use: "add",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), fv)
		},
	}

	flagSet := c.Flags()

	flagSet.StringVarP(&fv.Key,
		SecretKeyFlag, secretKeyShorthand, secretKeyDefault, secretKeyUsage)

	flagSet.StringVarP(&fv.Name,
		NameFlag, nameShorthand, nameDefault, nameUsage)

	flagSet.StringVar(&fv.CardNum,
		CardNumFlag, cardNumDefault, cardNumUsage)

	flagSet.StringVar(&fv.Exp,
		ExpDateFlag, expDateDefault, expDateUsage)

	flagSet.StringVar(&fv.Holder,
		HolderNameFlag, holderNameDefault, holderNameUsage)

	c.MarkFlagRequired(SecretKeyFlag)
	c.MarkFlagRequired(NameFlag)
	c.MarkFlagRequired(CardNumFlag)
	c.MarkFlagRequired(ExpDateFlag)
	c.MarkFlagRequired(HolderNameFlag)
	return &command.Command{Command: c}
}

type EditCmdFlags struct {
	Key, Name, CardNum, Exp, Holder string
	EntryNum                        int
}

func NewEdit(h command.GenCmdHandler[EditCmdFlags]) *command.Command {
	var fv EditCmdFlags

	c := &cobra.Command{
		Use: "edit",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), fv)
		},
	}

	flagSet := c.Flags()

	flagSet.StringVarP(&fv.Key,
		SecretKeyFlag, secretKeyShorthand, secretKeyDefault, secretKeyUsage)

	flagSet.IntVarP(&fv.EntryNum,
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage)

	flagSet.StringVarP(&fv.Name,
		NameFlag, nameShorthand, nameDefault, nameUsage)

	flagSet.StringVar(&fv.CardNum,
		CardNumFlag, cardNumDefault, cardNumUsage)

	flagSet.StringVar(&fv.Exp,
		ExpDateFlag, expDateDefault, expDateUsage)

	flagSet.StringVar(&fv.Holder,
		HolderNameFlag, holderNameDefault, holderNameUsage)

	c.MarkFlagRequired(SecretKeyFlag)
	c.MarkFlagRequired(EntryNumFlag)
	c.MarkFlagRequired(NameFlag)
	c.MarkFlagRequired(CardNumFlag)
	c.MarkFlagRequired(ExpDateFlag)
	c.MarkFlagRequired(HolderNameFlag)
	return &command.Command{Command: c}
}

func NewRead(h command.ReadCmdHandler) *command.Command {
	var (
		key      string
		entryNum int
	)

	c := &cobra.Command{
		Use: "read",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), key, entryNum)
		},
	}
	flagSet := c.Flags()

	flagSet.StringVarP(&key,
		SecretKeyFlag, secretKeyShorthand, secretKeyDefault, secretKeyUsage)

	flagSet.IntVarP(&entryNum,
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage)

	c.MarkFlagRequired(SecretKeyFlag)
	c.MarkFlagRequired(EntryNumFlag)

	return &command.Command{Command: c}
}

func NewList(h command.NoFlagsCmdHandler) *command.Command {
	c := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context())
		},
	}
	return &command.Command{Command: c}
}

func NewRemove(h command.RemoveCmdHandler) *command.Command {
	var entryNum int

	c := &cobra.Command{
		Use: "remove",
		Run: func(cmd *cobra.Command, args []string) {
			h.Handle(cmd.Context(), entryNum)
		},
	}

	c.Flags().IntP(
		EntryNumFlag, entryNumShorthand, entryNumDefault, entryNumUsage,
	)

	c.MarkFlagRequired(EntryNumFlag)
	return &command.Command{Command: c}
}
