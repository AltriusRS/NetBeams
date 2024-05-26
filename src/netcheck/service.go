package netcheck

import (
	_ "embed"
	"io"

	"github.com/altriusrs/netbeams/src/types"
	"github.com/ip2location/ip2proxy-go/v4"
)

//go:embed DB.BIN
var ip6Binary []byte

// Implement db reader interface
type IPDBReader interface {
	io.ReadCloser
	io.ReaderAt
}

type IPDBReaderAt struct {
	inner []byte
}

func (r IPDBReaderAt) ReadAt(p []byte, off int64) (n int, err error) {
	n = copy(p, r.inner[off:])
	return n, nil
}

func (r IPDBReaderAt) Close() error {
	return nil
}

func (r IPDBReaderAt) Read(p []byte) (n int, err error) {
	return r.ReadAt(p, 0)
}

type NetCheckService struct {
	types.Service
	db          *ip2proxy.DB      // IPv6 database
	blockedASNs map[string]string // ASN info for blocked providers (loaded from config on startup)
}

func Service() *NetCheckService {

	svc := &NetCheckService{
		Service:     types.SpinUp("NetCheck"),
		db:          nil,
		blockedASNs: map[string]string{},
	}

	svc.RegisterServiceHooks(svc.Start, svc.Stop, nil)

	return svc
}

func (s *NetCheckService) Check(ip string) (ip2proxy.IP2ProxyRecord, error) {
	return s.db.GetAll(ip)
}

func (s *NetCheckService) Start() (types.Status, error) {
	var err error

	s.db, err = ip2proxy.OpenDBWithReader(IPDBReaderAt{inner: ip6Binary})

	if err != nil {
		return types.StatusErrored, err
	}

	s.Infof("ModuleVersion   : %s", ip2proxy.ModuleVersion())
	s.Infof("Package Version : %s", s.db.PackageVersion())
	s.Infof("Database Version: %s", s.db.DatabaseVersion())

	return types.StatusHealthy, nil
}

func (s *NetCheckService) Stop() (types.Status, error) {
	return types.StatusShutdown, nil
}
