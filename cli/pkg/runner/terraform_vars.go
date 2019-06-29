package runner

import (
	"io"
	"os"
)

// TfVars struct representing the input variables for terraform to create the infrastructure
type TfVars struct {
	PublicKey  string
	AccessCIDR string
}

// NewTfVars creates a TfVars struct with all the defaults
func NewTfVars(publicKey string, accessCIDR string) TfVars {
	return TfVars{
		PublicKey:  publicKey,
		AccessCIDR: accessCIDR,
	}
}

var tmpl = `
access_key={{.PublicKey}}
access_cidr={{.AccessCIDR}}
`

func (tfv *TfVars) String() string {
	return "access_key = \"" + tfv.PublicKey + "\"\n" + "access_cidr = \"" + tfv.AccessCIDR + "\"\n"
}

// EnsureTfVarsFile writes an tfvars file if one hasnt already been made
func EnsureTfVarsFile(tfDir string, publicKey string, accessCIDR string) error {
	filename := tfDir + "/settings/bastion.tfVars"
	exists, err := FileExists(filename)
	if err != nil || exists {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	tfv := NewTfVars(publicKey, accessCIDR)

	_, err = io.WriteString(file, tfv.String())
	if err != nil {
		return err
	}
	return file.Sync()

}
