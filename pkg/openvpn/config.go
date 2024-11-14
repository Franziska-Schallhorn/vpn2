// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package openvpn

import (
	"io"
	"net"
	"strings"
	"text/template"

	"github.com/gardener/vpn2/pkg/network"
)

const defaultOpenVPNConfigFile = "/openvpn.config"
const defaultConfigFilePermissions = 0o600

func executeTemplate(name string, w io.Writer, templt string, data any) error {
	cidrMaskFunc :=
		func(n network.CIDR) string {
			mask := net.CIDRMask(n.Mask.Size())
			return net.IPv4(255, 255, 255, 255).Mask(mask).String()
		}
	networksToStringFunc := func(net []network.CIDR) string {
		networks := []string{}
		for _, nw := range net {
			networks = append(networks, nw.String())
		}
		return strings.Join(networks, ",")
	}
	var funcs = map[string]any{"cidrMask": cidrMaskFunc, "networksToString": networksToStringFunc}
	t, err := template.New(name).
		Funcs(funcs).
		Parse(templt)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}
