package runner_test

import (
	"github.com/controlplaneio/simulator-standalone/cli/pkg/runner"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func fixture(name string) string {
	return "../../test/fixtures/" + name
}

func readFixture(name string) string {
	file, err := os.Open(fixture(name))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	return string(b)
}

func Test_IsUsable(t *testing.T) {
	tfo := runner.TerraformOutput{}
	assert.False(t, tfo.IsUsable(), "Empty TerraformOutput was usable")
	tfo.BastionPublicIP.Value = "127.0.0.1"
	assert.False(t, tfo.IsUsable(), "TerraformOutput with only bastion was usable")
	tfo.MasterNodesPrivateIP.Value = []string{"127.0.0.1"}
	assert.False(t, tfo.IsUsable(), "TerraformOutput with only master IP was usable")
	tfo.ClusterNodesPrivateIP.Value = []string{"127.0.0.1"}
	assert.False(t, tfo.IsUsable(), "TerraformOutput with only 1 cluster node IPs was usable")

	tfo.ClusterNodesPrivateIP.Value = []string{"127.0.0.1", "127.0.0.2"}
	assert.True(t, tfo.IsUsable(), "Complete TerraformOutput was not usable")
}

func Test_ParseTerraformOutput(t *testing.T) {
	t.Parallel()
	output := readFixture("valid-tf-output.json")

	tfOutput, err := runner.ParseTerraformOutput(output)

	assert.Nil(t, err, "Got an error")
	assert.NotNil(t, tfOutput, "Output was nil")
	assert.Equal(t, tfOutput.BastionPublicIP.Value, "34.244.109.234", "Bastion IP was wrong")
	assert.Equal(t, len(tfOutput.ClusterNodesPrivateIP.Value), 1, "Didn't get 1 node IP")
	assert.Equal(t, tfOutput.ClusterNodesPrivateIP.Value[0], "172.31.2.19")
	assert.Equal(t, len(tfOutput.MasterNodesPrivateIP.Value), 1, "Didn't get 1 master IP")
	assert.Equal(t, tfOutput.MasterNodesPrivateIP.Value[0], "172.31.2.167")
}
