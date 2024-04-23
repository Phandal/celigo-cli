package celigo

import (
	"testing"
)

const programName = "celigo-cli"

var testResource = Resource{
	usage: "test resource usage",
}

func testResourceBuilder(cmd *Command) Resource {
	return testResource
}

func TestNewCommandNoArgs(t *testing.T) {
	var expectedError error = nil
	var expectedResource string = "help"
	cmd, err := NewCommand([]string{programName})

	if err != nil {
		t.Errorf("Expected err to be \"%s\" got \"%s\"", expectedError, err)
	}

	if cmd.ResourceArg != expectedResource {
		t.Errorf("Expected Resource to be \"%s\". Actual \"%s\"", expectedResource, cmd.ResourceArg)
	}
}

func TestNewCommandSingleHelpArg(t *testing.T) {
	var expectedError error = nil
	var expectedResource string = "help"
	firstArgs := []string{"help", "--help", "-h"}

	for _, v := range firstArgs {
		cmd, err := NewCommand([]string{programName, v})

		if err != nil {
			t.Errorf("Expected err to be \"%s\" got \"%s\"", expectedError, err)
		}

		if cmd.ResourceArg != expectedResource {
			t.Errorf("Expected Resource to be \"%s\". Actual \"%s\"", expectedResource, cmd.ResourceArg)
		}
	}
}

func TestNewCommandMissingAction(t *testing.T) {
	var expectedResource string = "resource"
	cmd, err := NewCommand([]string{programName, expectedResource})

	if err.Error() != missingActionErr.Error() {
		t.Errorf("Expected err to be \"%s\". Actual \"%s\"", missingActionErr, err)
	}

	if cmd.ResourceArg != expectedResource {
		t.Errorf("Expected Resource to be \"%s\". Actual \"%s\"", expectedResource, cmd.ResourceArg)
	}
}

func TestNewCommandWorks(t *testing.T) {
	var expectedResource string = "resource"
	var expectedAction string = "action"
	var expectedError error = nil

	cmd, err := NewCommand([]string{programName, expectedResource, expectedAction})

	if err != nil {
		t.Errorf("Expected err to be \"%s\". Actual \"%s\"", expectedError, err)
	}

	if cmd.ResourceArg != expectedResource {
		t.Errorf("Expected Resource to be \"%s\". Actual \"%s\"", expectedResource, cmd.ResourceArg)
	}

	if cmd.ActionArg != expectedAction {
		t.Errorf("Expected Action to be \"%s\". Actual \"%s\"", expectedAction, cmd.ActionArg)
	}
}

func TestNewResource(t *testing.T) {
	var expectedLength = 1
	var expectedName = "first"
	var expectedUsage = "usage"

	cmd, err := NewCommand([]string{programName, "resource", "action"})

	if err != nil {
		t.Errorf("failed to create new command: %s", err)
	}

	cmd.NewResource(expectedName, expectedUsage, testResourceBuilder)

	if len(cmd.mappedResources) != expectedLength {
		t.Errorf("Expected length of cmd.Resources %d. Actual %d", expectedLength, len(cmd.mappedResources))
	}

	res, exists := cmd.mappedResources[expectedName]

	if !exists {
		t.Errorf("Expected Resource with name \"%s\" to exist.", expectedName)
	}

	if res.usage != expectedUsage {
		t.Errorf("Expected resource usage: \"%s\". Actual \"%s\"", expectedUsage, res.usage)
	}
}

func TestCommandExecuteInvalidResource(t *testing.T) {
	var resourceName = "fakeResource"
	var actionName = "fakeAction"
	var expectedErr = invalidResourceErr(resourceName)
	cmd, err := NewCommand([]string{programName, resourceName, actionName})

	if err != nil {
		t.Errorf("failed to create new command: %s", err)
	}

	err = cmd.Execute()

	t.Log(err.Error())
	t.Log(expectedErr.Error())

	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error \"%s\". Actual \"%s\"", expectedErr, err)
	}
}

func TestCommandExecuteInvalidAction(t *testing.T) {
	var resourceName = "fakeResource"
	var actionName = "fakeAction"
	var expectedErr = invalidActionErr(actionName)
	cmd, err := NewCommand([]string{programName, resourceName, actionName})

	if err != nil {
		t.Errorf("failed to create new command: %s", err)
	}

	cmd.NewResource(resourceName, "test usage", testResourceBuilder)

	err = cmd.Execute()

	if err.Error() != expectedErr.Error() {
		t.Errorf("Expected error \"%s\". Actual \"%s\"", expectedErr, err)
	}
}
