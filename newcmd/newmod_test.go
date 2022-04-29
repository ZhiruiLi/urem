package newcmd

import (
	"testing"
)

func TestFormatProjectJsonText(t *testing.T) {
	cases := []struct {
		name   string
		json   string
		expect string
	}{
		{
			name: "simple no module",
			json: `{
	"FileVersion": 3,
	"EngineAssociation": "4.26"
}`,
			expect: `{
	"FileVersion": 3,
	"EngineAssociation": "4.26",
	"Modules": [
		{
			"Name": "NewModule",
			"Type": "Runtime",
			"LoadingPhase": "Default"
		}
	]
}`,
		},
		{
			name: "simple empty module list",
			json: `{
	"FileVersion": 3,
	"EngineAssociation": "4.26",
	"Modules": []
}`,
			expect: `{
	"FileVersion": 3,
	"EngineAssociation": "4.26",
	"Modules": [
		{
			"Name": "NewModule",
			"Type": "Runtime",
			"LoadingPhase": "Default"
		}
	]
}`,
		},
		{
			name: "simple non empty module list",
			json: `{
	"FileVersion": 3,
	"EngineAssociation": "4.26",
	"Modules": [
		{
			"Name": "OtherModule",
			"Type": "Runtime",
			"LoadingPhase": "Default"
		}
	]
}`,
			expect: `{
	"FileVersion": 3,
	"EngineAssociation": "4.26",
	"Modules": [
		{
			"Name": "NewModule",
			"Type": "Runtime",
			"LoadingPhase": "Default"
		},
		{
			"Name": "OtherModule",
			"Type": "Runtime",
			"LoadingPhase": "Default"
		}
	]
}`,
		},
	}

	for i, c := range cases {
		actual := formatProjectJsonText(c.json, "NewModule", "Runtime", "Default")
		if c.expect != actual {
			t.Errorf("%d:%s:\nexpect:\n%s\n\nactual:\n%s", i, c.name, c.expect, actual)
		}
	}
}
