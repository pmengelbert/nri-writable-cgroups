package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/containerd/nri/pkg/api"
	"github.com/containerd/nri/pkg/stub"
	"github.com/sirupsen/logrus"
)

const (
	annotationKey = "engelbert.dev/writable-cgroups"
	pluginName    = "writable-cgroups"
)

var log = logrus.StandardLogger()

type plugin struct {
	stub stub.Stub
}

func (p *plugin) CreateContainer(
	_ context.Context,
	pod *api.PodSandbox,
	ctr *api.Container,
) (*api.ContainerAdjustment, []*api.ContainerUpdate, error) {
	log.Infof("CreateContainer: pod=%s container=%s", pod.GetName(), ctr.GetName())

	if !wantsWritableCgroups(pod, ctr) {
		log.Infof("  skipping %s/%s: annotation %s not set to \"true\"",
			pod.GetName(), ctr.GetName(), annotationKey)
		return nil, nil, nil
	}

	adjust := &api.ContainerAdjustment{}

	for _, m := range ctr.GetMounts() {
		if m.GetType() != "cgroup" {
			continue
		}

		log.Infof("  found cgroup mount at %s with options %v",
			m.GetDestination(), m.GetOptions())

		adjust.RemoveMount(m.GetDestination())
		adjust.AddMount(&api.Mount{
			Destination: m.GetDestination(),
			Type:        m.GetType(),
			Source:      m.GetSource(),
			Options:     replaceOption(m.GetOptions(), "ro", "rw"),
		})

		log.Infof("  made cgroup mount at %s writable for %s/%s",
			m.GetDestination(), pod.GetName(), ctr.GetName())
	}

	return adjust, nil, nil
}

// wantsWritableCgroups checks the container annotations first, then falls back
// to pod annotations for the writable-cgroups annotation.
func wantsWritableCgroups(pod *api.PodSandbox, ctr *api.Container) bool {
	if v, ok := ctr.GetAnnotations()[annotationKey]; ok {
		return v == "true"
	}
	if v, ok := pod.GetAnnotations()[annotationKey]; ok {
		return v == "true"
	}
	return false
}

// replaceOption returns a copy of opts with every occurrence of old replaced by
// new. If old is not present, new is appended to ensure it is set.
func replaceOption(opts []string, old, new string) []string {
	out := make([]string, 0, len(opts))
	found := false
	for _, o := range opts {
		if o == old {
			out = append(out, new)
			found = true
		} else {
			out = append(out, o)
		}
	}
	if !found {
		out = append(out, new)
	}
	return out
}

func main() {
	var (
		pluginIdx  string
		socketPath string
	)

	flag.StringVar(&pluginIdx, "idx", "10", "plugin index")
	flag.StringVar(&socketPath, "socket-path", "", "NRI socket path to connect to")
	flag.Parse()

	opts := []stub.Option{
		stub.WithPluginName(pluginName),
		stub.WithPluginIdx(pluginIdx),
	}
	if socketPath != "" {
		opts = append(opts, stub.WithSocketPath(socketPath))
	}

	p := &plugin{}
	var err error
	if p.stub, err = stub.New(p, opts...); err != nil {
		log.Fatalf("failed to create NRI stub: %v", err)
	}

	ctx := context.Background()
	if err := p.stub.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "plugin exited with error: %v\n", err)
		os.Exit(1)
	}
}
