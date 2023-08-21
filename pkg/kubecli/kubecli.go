package kubecli

import (
	"bytes"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/delete"
	"k8s.io/kubectl/pkg/cmd/get"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

type Kubecli struct {
	clientConfig *rest.Config
	restClient *kubernetes.Clientset
	dynamicClient dynamic.Interface

	ioStream genericclioptions.IOStreams
	out *bytes.Buffer
	errOut *bytes.Buffer

	configFlag *genericclioptions.ConfigFlags
}

func NewKubecli(config *rest.Config) *Kubecli {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil
	}
	ioStream, _, out, errOut := genericclioptions.NewTestIOStreams()
	configFlag := genericclioptions.NewConfigFlags(true)
	configFlag.WrapConfigFn = func(*rest.Config) *rest.Config {
		return config
	}
	return &Kubecli{
		clientConfig: config,
		restClient: clientset,
		dynamicClient: dynamicClient,
		ioStream: ioStream,
		out: out,
		errOut: errOut,
		configFlag: configFlag}
}

func (cli *Kubecli) Test() error {
	_, err := cli.restClient.ServerVersion()
	return err
}

func (cli *Kubecli) NewFactory(namespace string) cmdutil.Factory {
	if namespace != "" {
		cli.configFlag.Namespace = &namespace
	}
	return cmdutil.NewFactory(cli.configFlag)
}

func (cli *Kubecli) RunCmdWithNormalOutput(cmd *cobra.Command, params []string) {
	o := cli.out
	e := cli.errOut
	cmd.SetOut(o)
	cmd.SetErr(e)
	cmd.Run(cmd, params)

	if e.String() != "" {
		fmt.Println(e.String())
	}
	fmt.Println(o.String())
}

func (cli *Kubecli) Get(resource string, namespace string)  {
	factory := cli.NewFactory(namespace)
	cmd := get.NewCmdGet("kubectl", factory, cli.ioStream)
	cli.RunCmdWithNormalOutput(cmd, []string{resource})
}

func (cli *Kubecli) ApplyFromInput(in io.Reader, namespace string) {
	factory := cli.NewFactory(namespace)
	cmd := apply.NewCmdApply("kubectl", factory, cli.ioStream)
	cmd.SetIn(in)
	cmd.Flags().Set("filename", "-")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		flags := apply.NewApplyFlags(cli.ioStream)
		o, err := flags.ToOptions(factory, cmd, "kubectl", args)
		o.Builder.Stdin()
		cmdutil.CheckErr(err)
		cmdutil.CheckErr(o.Validate())
		cmdutil.CheckErr(o.Run())
	}
	cli.RunCmdWithNormalOutput(cmd, []string{})
}

func (cli *Kubecli) ApplyFromFile(f string, namespace string) {
	factory := cli.NewFactory(namespace)
	cmd := apply.NewCmdApply("kubectl", factory, cli.ioStream)
	cmd.Flags().Set("filename", f)
	cli.RunCmdWithNormalOutput(cmd, []string{})
}

func (cli *Kubecli) DeleteFromFile(f string, namespace string) {
	factory := cli.NewFactory(namespace)
	cmd := delete.NewCmdDelete(factory, cli.ioStream)
	cmd.Flags().Set("filename", f)
	cli.RunCmdWithNormalOutput(cmd, []string{})
}

func (cli *Kubecli) DeleteResource(resourceType, resource, namespace string) {
	factory := cli.NewFactory(namespace)
	cmd := delete.NewCmdDelete(factory, cli.ioStream)
	cli.RunCmdWithNormalOutput(cmd, []string{resourceType, resource})
}