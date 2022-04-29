package lscmd

import (
	"fmt"
	"strings"

	"github.com/zhiruili/urem/core"
)

// https://docs.unrealengine.com/4.26/en-US/API/Runtime/Projects/EHostType__Type/
var availableModuleTypes = []string{
	"Runtime",
	"RuntimeNoCommandlet",
	"RuntimeAndProgram",
	"CookedOnly",
	"UncookedOnly",
	"Developer",
	"DeveloperTool",
	"Editor",
	"EditorNoCommandlet",
	"EditorAndProgram",
	"Program",
	"ServerOnly",
	"ClientOnly",
	"ClientOnlyNoCommandlet",
}

func GetAvailableModuleTypes() []string {
	return availableModuleTypes
}

func GetFmtAvailableModuleTypes(sep string) string {
	return strings.Join(availableModuleTypes, sep)
}

func IsLegalModuleType(t string) bool {
	return core.StrContains(availableModuleTypes, t)
}

// https://docs.unrealengine.com/4.26/en-US/API/Runtime/Projects/ELoadingPhase__Type/
var availableLoadingPhases = []string{
	"EarliestPossible",
	"PostConfigInit",
	"PostSplashScreen",
	"PreEarlyLoadingScreen",
	"PreLoadingScreen",
	"PreDefault",
	"Default",
	"PostDefault",
	"PostEngineInit",
}

func GetAvailableLoadingPhases() []string {
	return availableLoadingPhases
}

func GetFmtAvailableLoadingPhases(sep string) string {
	return strings.Join(availableLoadingPhases, sep)
}

func IsLegalLoadingPhase(t string) bool {
	return core.StrContains(availableLoadingPhases, t)
}

type Cmd struct {
	Target string `arg:"positional,required" help:"list target modtype/loadphase"`
}

func (cmd *Cmd) Run() error {
	switch strings.ToLower(cmd.Target) {
	case "modtype":
		fmt.Println(GetFmtAvailableModuleTypes("\n"))
	case "loadphase":
		fmt.Println(GetFmtAvailableLoadingPhases("\n"))
	default:
		return fmt.Errorf("missing target: modtype/loadphase")
	}

	return nil
}
