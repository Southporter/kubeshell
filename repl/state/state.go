package state

import (
	"errors"
	"os"

	"golang.org/x/term"

	"github.com/spf13/cobra"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientapi "k8s.io/client-go/tools/clientcmd/api"
)

type State struct {
	Namespace  string
	cwd        string
	config     *clientcmd.ClientConfig
	client     *kubernetes.Clientset
	raw        *clientapi.Config
	restConfig *rest.Config
}

func NewState() State {
	return State{
		cwd:       "/",
		Namespace: "",
	}
}

func (s *State) Load(cmd *cobra.Command) error {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

	configPath := cmd.PersistentFlags().Lookup("kubeconfig").Value.String()
	if configPath != "" {
		loadingRules.ExplicitPath = configPath
	}

	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	s.config = &kubeConfig

	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		return errors.New("could not get namespace from config")
	}
	s.Namespace = namespace

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return errors.New("could not get client config from kubeconfig")
	}

	s.restConfig = config

	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return errors.New("could not load raw config")
	}
	s.raw = &rawConfig

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.New("could not create kubernetes client from config file")
	}

	s.client = client

	s.cwd = "/namespaces/" + s.Namespace

	return nil
}

func (s *State) CurrentDirectory() string {
	return s.cwd
}

func (s *State) SetCurrentDirectory(dir string) {
	s.cwd = dir
}

func (s *State) CurrentCluster() string {
	return s.raw.CurrentContext
}

func (s *State) ToDiscoveryClient() (discovery.DiscoveryInterface, error) {
	return discovery.NewDiscoveryClientForConfig(s.restConfig)
}

func (s *State) ToClient() *kubernetes.Clientset {
	return s.client
}

func (s *State) GetTerminalWidth() (int, error) {
	width, _, err := term.GetSize(int(os.Stdin.Fd()))
	return width, err
}
