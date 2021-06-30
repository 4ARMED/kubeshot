package k8sapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/4armed/kubeshot/internal/config"
	"github.com/kubicorn/kubicorn/pkg/logger"

	// "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetURLs queries the Kubernetes API to get URLs and returns an array
func GetURLs(c *config.Config) ([]string, error) {
	var urls []string
	var clusterConfig *rest.Config
	var err error

	if len(c.KubeConfig) > 0 {
		clusterConfig, err = clientcmd.BuildConfigFromFlags("", c.KubeConfig)
		if err != nil {
			return nil, err
		}
	} else {
		// creates the in-cluster config
		clusterConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(clusterConfig)
	if err != nil {
		panic(err.Error())
	}

	if c.GetK8sPods {
		logger.Info("Fetching pods from Kubernetes API")

	}

	if c.GetK8sSvcs {
		logger.Info("Fetching services from Kubernetes API")
		var urls []string
		svcs, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		logger.Info("Counted %d services", len(svcs.Items))

		for _, service := range svcs.Items {
			logger.Debug("service name: %v", service.Name)
			if service.Spec.Type == "ClusterIP" {
				if service.Spec.ClusterIP != "None" {
					for _, port := range service.Spec.Ports {
						if port.Protocol == "TCP" {
							var url string
							// Try and guess the correct protocol, http or https
							if strings.HasPrefix(port.Name, "https") || port.Port == 443 {
								url = fmt.Sprintf("https://%v:%d", service.Spec.ClusterIP, port.Port)
							} else {
								url = fmt.Sprintf("http://%v:%d", service.Spec.ClusterIP, port.Port)
							}
							logger.Debug("url: %v", url)
							urls = append(urls, url)
						}
					}
				}
			}
		}
		return urls, err
	}

	return urls, nil
}
