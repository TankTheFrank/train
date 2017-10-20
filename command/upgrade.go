package trainCommand

import (
	"fmt"
)

const CmdBinPath = "$GOPATH/bin/train"

func Upgrade() {
	bash("go get -u github.com/tankthefrank/train")
	bash("go build -o " + CmdBinPath + " github.com/tankthefrank/train/cmd")
	fmt.Println("Installed latest train command into " + CmdBinPath)
}
