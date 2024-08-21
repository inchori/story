package docker

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"text/template"

	e2e "github.com/cometbft/cometbft/test/e2e/pkg"
	"github.com/cometbft/cometbft/test/e2e/pkg/exec"
	"github.com/cometbft/cometbft/test/e2e/pkg/infra"
	cmtdocker "github.com/cometbft/cometbft/test/e2e/pkg/infra/docker"

	"github.com/storyprotocol/iliad/e2e/app/geth"
	"github.com/storyprotocol/iliad/e2e/types"
	"github.com/storyprotocol/iliad/lib/errors"
	"github.com/storyprotocol/iliad/lib/log"

	_ "embed"
)

const ProviderName = "docker"

// composeTmpl is our own custom docker compose template. This differs from cometBFT's.
//
//go:embed compose.yaml.tmpl
var composeTmpl []byte

var _ types.InfraProvider = (*Provider)(nil)

// Provider wraps the cometBFT docker provider, writing a different compose file.
type Provider struct {
	*cmtdocker.Provider
	servicesOnce sync.Once
	testnet      types.Testnet
	iliadTag     string
}

func (*Provider) Clean(ctx context.Context) error {
	log.Info(ctx, "Removing docker containers and networks")

	for _, cmd := range CleanCmds(false, runtime.GOOS == "linux") {
		err := exec.Command(ctx, "bash", "-c", cmd)
		if err != nil {
			return errors.Wrap(err, "remove docker containers")
		}
	}

	return nil
}

// NewProvider returns a new Provider.
func NewProvider(testnet types.Testnet, infd types.InfrastructureData, imgTag string) *Provider {
	return &Provider{
		Provider: &cmtdocker.Provider{
			ProviderData: infra.ProviderData{
				Testnet:            testnet.Testnet,
				InfrastructureData: infd.InfrastructureData,
			},
		},
		testnet:  testnet,
		iliadTag: imgTag,
	}
}

// Setup generates the docker-compose file and write it to disk, erroring if
// any of these operations fail.
func (p *Provider) Setup() error {
	def := ComposeDef{
		Network:     true,
		NetworkName: p.testnet.Name,
		NetworkCIDR: p.testnet.IP.String(),
		BindAll:     false,
		Nodes:       p.testnet.Nodes,
		IliadEVMs:   p.testnet.IliadEVMs,
		Anvils:      p.testnet.AnvilChains,
		Prometheus:  p.testnet.Prometheus,
		IliadTag:    p.iliadTag,
	}

	bz, err := GenerateComposeFile(def)
	if err != nil {
		return errors.Wrap(err, "generate compose file")
	}

	err = os.WriteFile(filepath.Join(p.Testnet.Dir, "docker-compose.yml"), bz, 0o644)
	if err != nil {
		return errors.Wrap(err, "write compose file")
	}

	return nil
}

func (*Provider) Upgrade(context.Context, types.UpgradeConfig) error {
	return errors.New("upgrade not supported for docker provider")
}

func (p *Provider) StartNodes(ctx context.Context, nodes ...*e2e.Node) error {
	var err error
	p.servicesOnce.Do(func() {
		svcs := additionalServices(p.testnet)
		log.Info(ctx, "Starting additional services", "names", svcs)

		err = cmtdocker.ExecCompose(ctx, p.Testnet.Dir, "create") // This fails if containers not available.
		if err != nil {
			err = errors.Wrap(err, "create containers")
			return
		}

		err = cmtdocker.ExecCompose(ctx, p.Testnet.Dir, append([]string{"up", "-d"}, svcs...)...)
		if err != nil {
			err = errors.Wrap(err, "start additional services")
			return
		}
	})
	if err != nil {
		return err
	}

	// if there are no iliad nodes available
	if len(nodes) == 0 {
		panic("no nodes to start")
	}

	// Start all requested nodes (use --no-deps to avoid starting the additional services again).
	nodeNames := make([]string, len(nodes))
	for i, n := range nodes {
		nodeNames[i] = n.Name
	}
	err = cmtdocker.ExecCompose(ctx, p.Testnet.Dir, append([]string{"up", "-d", "--no-deps"}, nodeNames...)...)
	if err != nil {
		return errors.Wrap(err, "start nodes")
	}

	return nil
}

type ComposeDef struct {
	Network     bool
	NetworkName string
	NetworkCIDR string
	BindAll     bool

	Nodes     []*e2e.Node
	IliadEVMs []types.IliadEVM
	Anvils    []types.AnvilChain

	IliadTag   string
	Prometheus bool
}

func (ComposeDef) GethTag() string {
	return geth.Version
}

// NodeIliadEVMs returns a map of node name to IliadEVM instance name; map[node_name]iliad_evm.
func (c ComposeDef) NodeIliadEVMs() map[string]string {
	resp := make(map[string]string)
	for i, node := range c.Nodes {
		evm := c.IliadEVMs[0].InstanceName
		if len(c.IliadEVMs) == len(c.Nodes) {
			evm = c.IliadEVMs[i].InstanceName
		}
		resp[node.Name] = evm
	}

	return resp
}

func GenerateComposeFile(def ComposeDef) ([]byte, error) {
	tmpl, err := template.New("compose").Parse(string(composeTmpl))
	if err != nil {
		return nil, errors.Wrap(err, "parse template")
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, def)
	if err != nil {
		return nil, errors.Wrap(err, "execute template")
	}

	return buf.Bytes(), nil
}

// CleanCmds returns generic docker commands to clean up docker containers and networks.
// This bypasses the need to a specific docker-compose context.
func CleanCmds(sudo bool, isLinux bool) []string {
	// GNU xargs requires the -r flag to not run when input is empty, macOS
	// does this by default. Ugly, but works.
	xargsR := ""
	if isLinux {
		xargsR = "-r"
	}

	// Some environments need sudo to run docker commands.
	perm := ""
	if sudo {
		perm = "sudo"
	}

	return []string{
		fmt.Sprintf("%s docker container ls -qa --filter label=e2e | xargs %v %s docker container rm -f",
			perm, xargsR, perm),
		fmt.Sprintf("%s docker network ls -q --filter label=e2e | xargs %v %s docker network rm",
			perm, xargsR, perm),
	}
}

// additionalServices returns additional (to iliad) docker-compose services to start.
func additionalServices(testnet types.Testnet) []string {
	var resp []string
	if testnet.Prometheus {
		resp = append(resp, "prometheus")
	}

	for _, iliadEVM := range testnet.IliadEVMs {
		resp = append(resp, iliadEVM.InstanceName)
	}
	for _, anvil := range testnet.AnvilChains {
		resp = append(resp, anvil.Chain.Name)
	}

	return resp
}
