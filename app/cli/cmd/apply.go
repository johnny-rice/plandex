package cmd

import (
	"plandex-cli/auth"
	"plandex-cli/lib"
	"plandex-cli/plan_exec"
	"plandex-cli/term"
	"plandex-cli/types"

	"github.com/spf13/cobra"
)

var autoCommit, skipCommit, autoExec bool

func init() {
	initApplyFlags(applyCmd, false)
	initExecScriptFlags(applyCmd)
	RootCmd.AddCommand(applyCmd)
}

var applyCmd = &cobra.Command{
	Use:     "apply",
	Aliases: []string{"ap"},
	Short:   "Apply a plan to the project",
	Run:     apply,
}

func apply(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MustResolveProject()
	mustSetPlanExecFlags(cmd)

	if lib.CurrentPlanId == "" {
		term.OutputNoCurrentPlanErrorAndExit()
	}

	applyFlags := types.ApplyFlags{
		AutoConfirm: true,
		AutoCommit:  autoCommit,
		NoCommit:    skipCommit,
		AutoExec:    autoExec,
		NoExec:      noExec,
		AutoDebug:   autoDebug,
	}

	tellFlags := types.TellFlags{
		TellBg:      tellBg,
		TellStop:    tellStop,
		TellNoBuild: tellNoBuild,
		AutoContext: tellAutoContext,
		ExecEnabled: !noExec,
		AutoApply:   tellAutoApply,
	}

	lib.MustApplyPlan(lib.ApplyPlanParams{
		PlanId:     lib.CurrentPlanId,
		Branch:     lib.CurrentBranch,
		ApplyFlags: applyFlags,
		TellFlags:  tellFlags,
		OnExecFail: plan_exec.GetOnApplyExecFail(applyFlags, tellFlags),
	})
}
