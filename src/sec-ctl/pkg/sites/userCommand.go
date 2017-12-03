package sites

import "fmt"

type UserCommandCode string

const (
	CmdArmAway               UserCommandCode = "ArmAway"
	CmdArmStay               UserCommandCode = "ArmStay"
	CmdArmWithPIN            UserCommandCode = "ArmWithPIN"
	CmdArmWithZeroEntryDelay UserCommandCode = "ArmWithZeroEntryDelay"
	CmdDisarm                UserCommandCode = "Disarm"
	CmdPanic                 UserCommandCode = "Panic"
)

const (
	PanicTargetFire      = "1"
	PanicTargetAmbulance = "2"
	PanicTargetPolice    = "3"
)

type UserCommand struct {
	Code        UserCommandCode `binding:"required"`
	PartitionID string          `binding:"required"`
	PIN         string
	PanicTarget string
}

func (cmd UserCommand) Validate() error {

	if cmd.Code != CmdArmAway && cmd.Code != CmdArmStay && cmd.Code != CmdArmWithPIN && cmd.Code != CmdArmWithZeroEntryDelay && cmd.Code != CmdDisarm && cmd.Code != CmdPanic {
		return fmt.Errorf("Invalid command code")
	}

	if (cmd.Code == CmdArmWithPIN || cmd.Code == CmdDisarm) && cmd.PIN == "" {
		return fmt.Errorf("PIN is required")
	}

	if cmd.Code == CmdPanic && cmd.PanicTarget == "" {
		return fmt.Errorf("PanicTarget is required")
	}

	return nil
}
