package main

import (
	"github.com/dims/community-images/cmd/community_images/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	cli.InitAndExecute()
}
