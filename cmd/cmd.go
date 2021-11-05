package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func init() {
}

var (
	kubeConfigFlags = genericclioptions.NewConfigFlags(false)

	kubeResouceBuilderFlags = genericclioptions.NewResourceBuilderFlags().
				WithAllNamespaces(false).
				WithLabelSelector("")
)

var rootCmd = &cobra.Command{
	Use:   "kubectl-locality",
	Short: "Plugin to access Kubernetes Pod localities.",
	RunE: func(cmd *cobra.Command, a []string) error {
		resources := "pods,nodes"
		resourceFinder := kubeResouceBuilderFlags.WithAll(true).ToBuilder(kubeConfigFlags, []string{resources})
		nodes := map[string]Locality{}
		pods := []*v1.Pod{}
		hasSubzone := false
		err := resourceFinder.Do().Visit(func(info *resource.Info, e error) error {
			switch info.Object.GetObjectKind().GroupVersionKind().Kind {
			case "Node":
				pm := &v1.Node{}
				if err := runtime.DefaultUnstructuredConverter.FromUnstructured(info.Object.(*unstructured.Unstructured).Object, pm); err != nil {
					return nil
				}
				l := Locality{
					Region:  takeFirst(pm.Labels[v1.LabelTopologyRegion], pm.Labels[v1.LabelFailureDomainBetaRegion]),
					Zone:    takeFirst(pm.Labels[v1.LabelTopologyZone], pm.Labels[v1.LabelFailureDomainBetaZone]),
					Subzone: pm.Labels["topology.istio.io/subzone"],
				}
				if l.Subzone != "" {
					hasSubzone = true
				}
				nodes[pm.Name] = l
			case "Pod":
				pm := &v1.Pod{}
				if err := runtime.DefaultUnstructuredConverter.FromUnstructured(info.Object.(*unstructured.Unstructured).Object, pm); err != nil {
					return nil
				}
				if pm.Spec.NodeName != "" {
					pods = append(pods, pm)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("failed to fetch resources: %v", err)
		}
		tw := new(tabwriter.Writer).Init(os.Stdout, 0, 5, 5, ' ', 0)
		header := "NAMESPACE\tNAME\tREGION\tZONE"
		if hasSubzone {
			header += "\tSUBZONE"
		}
		tw.Write([]byte(header+"\n"))
		for _, p := range pods {
			n := nodes[p.Spec.NodeName]
			d := []string{
				p.Namespace,
				p.Name,
				n.Region,
				n.Zone,
			}
			if hasSubzone {
				d = append(d, n.Subzone)
			}
			tw.Write([]byte(strings.Join(d, "\t") + "\n"))
		}
		tw.Flush()
		return nil
	},
}

func takeFirst(v ...string) string {
	for _, x := range v {
		if x != "" {
			return x
		}
	}
	return ""
}

type Locality struct {
	Region, Zone, Subzone string
}

func Execute() {
	flags := pflag.NewFlagSet("kubectl-resources", pflag.ExitOnError)
	pflag.CommandLine = flags

	kubeConfigFlags.AddFlags(flags)
	kubeResouceBuilderFlags.AddFlags(flags)
	flags.AddFlagSet(rootCmd.PersistentFlags())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
