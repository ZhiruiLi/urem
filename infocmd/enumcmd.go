package infocmd

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

// GetAvailableModuleTypes 获取所有合法的 UE module 类型。
func GetAvailableModuleTypes() []string {
	return availableModuleTypes
}

// GetFmtAvailableModuleTypes 获取所有合法的 UE module 类型的格式化字符串。
func GetFmtAvailableModuleTypes(sep string) string {
	return strings.Join(availableModuleTypes, sep)
}

// IsLegalModuleType 检查一个字符串是否是一个合法的 UE module 类型。
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

// GetAvailableLoadingPhases 获取所有合法的 loading phase 类型。
func GetAvailableLoadingPhases() []string {
	return availableLoadingPhases
}

// GetFmtAvailableLoadingPhases 获取所有合法的 loading phase 类型的格式化字符串。
func GetFmtAvailableLoadingPhases(sep string) string {
	return strings.Join(availableLoadingPhases, sep)
}

// IsLegalLoadingPhase 检查一个字符串是否是一个合法的 loading phase 类型。
func IsLegalLoadingPhase(t string) bool {
	return core.StrContains(availableLoadingPhases, t)
}

// InfoEnumCmd 实现了 enum 查询子命令。
type InfoEnumCmd struct {
	Target string `arg:"positional,required" help:"list target modtype/loadphase"`
}

// Run 执行 info enum 子命令。
func (cmd *InfoEnumCmd) Run() error {
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
