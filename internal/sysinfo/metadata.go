package sysinfo

import (
	"context"
	"errors"
	"fmt"
	"github.com/clarkmcc/cloudcore/internal/agentdb"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"os"
	"regexp"
	"runtime"
)

var (
	alphaNumericExpr = regexp.MustCompile("[a-zA-Z0-9.]+")
)

func BuildSystemMetadata(ctx context.Context, db agentdb.AgentDB, logger *zap.Logger) (*rpc.SystemMetadata, error) {
	s := &rpc.SystemMetadata{}

	// Identifiers
	s.Identifiers = &rpc.SystemMetadata_Identifiers{}
	var err error
	s.Identifiers.AgentIdentifier, err = db.AgentID(ctx)
	if err != nil && !errors.Is(err, agentdb.ErrNoAgentID) {
		return nil, err
	}
	s.Identifiers.HostId, err = host.HostIDWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting host id: %w", err)
	}
	s.Identifiers.HostIdentifier = s.Identifiers.HostId // todo: using system uuid as default, support pluggable identifier
	s.Identifiers.Hostname, err = os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("getting hostname: %w", err)
	}
	s.Identifiers.PublicIpAddress, err = getPublicIpAddress()
	if err != nil {
		logger.Warn("failed to get public ip address", zap.Error(err))
	}
	s.Identifiers.PrivateIpAddress, err = getPrivateIpAddress()
	if err != nil {
		logger.Warn("failed to get private ip address", zap.Error(err))
	}

	info, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	// OS
	s.Os = &rpc.SystemMetadata_OperatingSystem{}
	s.Os.Name = info.Platform
	s.Os.Version = info.PlatformVersion
	s.Os.Family = info.PlatformFamily
	s.Os.Goos = runtime.GOOS
	s.Os.Goarch = runtime.GOARCH

	// Kernel
	s.Kernel = &rpc.SystemMetadata_Kernel{}
	s.Kernel.Version = info.KernelVersion
	s.Kernel.Arch = info.KernelArch

	// CPU
	c, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, err
	}
	s.Cpu = &rpc.SystemMetadata_CPU{}
	for _, v := range c {
		s.Cpu.Cores += v.Cores
		s.Cpu.Model = v.ModelName
	}

	return s, nil
}

// getPublicIpAddress uses an AWS api to return the public IP address
func getPublicIpAddress() (string, error) {
	resp, err := http.Get("https://checkip.amazonaws.com/")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return alphaNumericExpr.FindString(string(bs)), nil
}

func getPrivateIpAddress() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return "", fmt.Errorf("expected local address to be *net.UDPAddr, got %T", conn.LocalAddr())
	}
	return addr.IP.String(), nil
}
