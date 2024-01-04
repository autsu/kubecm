package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

// ListAuthCommand list cmd struct
type ListAuthCommand struct {
	BaseCommand
}

// Init ListAuthCommand
func (lc *ListAuthCommand) Init() {
	lc.command = &cobra.Command{
		Use:     "list_auth",
		Short:   "List KubeConfig Auth Info",
		Long:    "List KubeConfig Auth Info",
		Aliases: []string{"la"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return lc.runList(cmd, args)
		},
		Example: listAuthExample(),
	}
	lc.command.DisableFlagsInUseLine = true
}

func (lc *ListAuthCommand) runList(command *cobra.Command, args []string) error {
	clusterMessageChan := make(chan *ClusterStatusCheck)
	go func() {
		info, _ := ClusterStatus(2)
		clusterMessageChan <- info
	}()
	config, err := clientcmd.LoadFromFile(cfgFile)
	if err != nil {
		return err
	}
	config = CheckValidContext(false, config)
	outConfig, err := filterArgs(args, config)
	if err != nil {
		return err
	}
	err = PrintTableWithAuth(outConfig)
	if err != nil {
		return err
	}
	clusterMessage := <-clusterMessageChan
	if clusterMessage != nil {
		printString(os.Stdout, "Cluster check succeeded!")
		printString(os.Stdout, "\nKubernetes version ")
		printYellow(os.Stdout, clusterMessage.Version.GitVersion)
		printService(os.Stdout, "\nKubernetes master", clusterMessage.Config.Host)
		err = MoreInfo(clusterMessage.ClientSet, os.Stdout)
		if err != nil {
			fmt.Println("(Error reporting can be ignored and does not affect usage.)")
		}
	}
	return nil
}

func listAuthExample() string {
	return `
# List all the contexts and auth info in your KubeConfig file
kubecm list_auth
# Aliases
kubecm la
# Filter out keywords(Multi-keyword support)
kubecm la kind k3s
`
}
