package state

import (
	"errors"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientapi "k8s.io/client-go/tools/clientcmd/api"
)

type State struct {
	Namespace string
	cwd       string
	config    *clientcmd.ClientConfig
	client    *kubernetes.Clientset
	raw       *clientapi.Config
}

func NewState() State {
	return State{
		cwd:       "/",
		Namespace: "",
		config:    nil,
		client:    nil,
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
		return errors.New("Could not get namespace from config")
	}
	s.Namespace = namespace

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return errors.New("Could not get client config from kubeconfig")
	}

	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return errors.New("Could not load raw config")
	}
	s.raw = &rawConfig

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.New("Could not create kubernetes client from config file")
	}

	s.client = client

	return nil
}

func (s *State) CurrentDirectory() string {
	return s.cwd
}

func (s *State) CurrentCluster() string {
	return s.raw.CurrentContext
}
