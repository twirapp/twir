package clientinfo

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
)

const (
	geoDbDirName      = "assets/geo-db"
	geoDbFileName     = "dbip-city-lite.mmdb"
	geoDbInfoFileName = "db-info.json"
)

type Opts struct {
	fx.In

	LC     fx.Lifecycle
	Logger *slog.Logger
	Config config.Config
}

type Location struct {
	Country *string
	City    *string
}

type Service struct {
	logger   *slog.Logger
	config   config.Config
	dbPath   string
	infoPath string

	mu     sync.RWMutex
	reader *geoip2.Reader
}

func New(opts Opts) *Service {
	dbDir := resolveDbDir()
	s := &Service{
		logger:   opts.Logger,
		config:   opts.Config,
		dbPath:   filepath.Join(dbDir, geoDbFileName),
		infoPath: filepath.Join(dbDir, geoDbInfoFileName),
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: s.onStart,
			OnStop:  s.onStop,
		},
	)

	return s
}

type ClientInfo struct {
	IP        string
	UserAgent string
}

func (s *Service) GetClientInfo(ctx context.Context) (ClientInfo, error) {
	gCtx, err := gincontext.GetGinContext(ctx)
	if err != nil {
		return ClientInfo{}, err
	}

	ip := strings.TrimSpace(gCtx.ClientIP())
	if ip == "" {
		return ClientInfo{}, errors.New("client ip is empty")
	}

	return ClientInfo{
		IP:        ip,
		UserAgent: gCtx.GetHeader("User-Agent"),
	}, nil
}

func (s *Service) LookupIP(ctx context.Context, ip string) (Location, error) {
	reader, err := s.getReader(ctx)
	if err != nil {
		return Location{}, err
	}

	parsedIP, err := parseIP(ip)
	if err != nil {
		return Location{}, err
	}

	record, err := reader.City(parsedIP)
	if err != nil {
		return Location{}, err
	}

	var country *string
	if record.Country.ISOCode != "" {
		value := record.Country.ISOCode
		country = &value
	}

	var city *string
	if record.City.Names.English != "" {
		value := record.City.Names.English
		city = &value
	}

	return Location{
		Country: country,
		City:    city,
	}, nil
}

func (s *Service) onStart(ctx context.Context) error {
	if err := s.ensureDatabase(ctx, s.config.IsDevelopment()); err != nil {
		if s.config.IsDevelopment() {
			s.logger.WarnContext(ctx, "GeoIP database is not ready", logger.Error(err))
			return nil
		}
		return err
	}

	if err := s.openReader(); err != nil {
		if s.config.IsDevelopment() {
			s.logger.WarnContext(ctx, "Cannot open GeoIP database", logger.Error(err))
			return nil
		}
		return err
	}

	testInfo, err := s.LookupIP(ctx, "194.87.25.102")
	if err != nil {
		if s.config.IsDevelopment() {
			s.logger.WarnContext(ctx, "Cannot test lookup GeoIP database", logger.Error(err))
			return nil
		}
		return err
	}

	s.logger.InfoContext(
		ctx, "GeoIP database initialized",
		slog.Group(
			"test_data",
			slog.String("country", *testInfo.Country),
			slog.String("city", *testInfo.City),
		),
	)

	return nil
}

func (s *Service) onStop(_ context.Context) error {
	s.mu.Lock()
	reader := s.reader
	s.reader = nil
	s.mu.Unlock()

	if reader == nil {
		return nil
	}

	return reader.Close()
}

func (s *Service) getReader(ctx context.Context) (*geoip2.Reader, error) {
	s.mu.RLock()
	reader := s.reader
	s.mu.RUnlock()

	if reader != nil {
		return reader, nil
	}

	if !s.config.IsDevelopment() {
		return nil, errors.New("geoip database is not initialized")
	}

	if err := s.ensureDatabase(ctx, s.config.IsDevelopment()); err != nil {
		return nil, err
	}

	if err := s.openReader(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	reader = s.reader
	s.mu.RUnlock()

	if reader == nil {
		return nil, errors.New("geoip database is not initialized")
	}

	return reader, nil
}

func (s *Service) openReader() error {
	reader, err := geoip2.Open(s.dbPath)
	if err != nil {
		return fmt.Errorf("cannot open geoip database: %w", err)
	}

	s.mu.Lock()
	if s.reader != nil {
		_ = s.reader.Close()
	}
	s.reader = reader
	s.mu.Unlock()

	return nil
}

func (s *Service) ensureDatabase(ctx context.Context, allowDownload bool) error {
	if _, err := os.Stat(s.dbPath); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cannot stat geoip database: %w", err)
	}

	if !allowDownload {
		return fmt.Errorf("geoip database not found at %s", s.dbPath)
	}

	return s.downloadDatabase(ctx)
}

func (s *Service) downloadDatabase(ctx context.Context) error {
	if err := os.MkdirAll(filepath.Dir(s.dbPath), 0o755); err != nil {
		return fmt.Errorf("cannot create geoip directory: %w", err)
	}

	now := time.Now().UTC()
	url := fmt.Sprintf(
		"https://download.db-ip.com/free/dbip-city-lite-%s-%s.mmdb.gz",
		now.Format("2006"),
		now.Format("01"),
	)

	s.logger.InfoContext(ctx, "Downloading GeoIP database", slog.String("url", url), slog.String("path", s.dbPath))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("cannot build geoip request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot download geoip database: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("geoip download failed: %s", resp.Status)
	}

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read geoip archive: %w", err)
	}
	defer gzReader.Close()

	tmpFile, err := os.CreateTemp(filepath.Dir(s.dbPath), "dbip-city-lite-*.mmdb")
	if err != nil {
		return fmt.Errorf("cannot create geoip temp file: %w", err)
	}
	tmpName := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()

	if _, err := io.Copy(tmpFile, gzReader); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("cannot write geoip database: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("cannot close geoip temp file: %w", err)
	}

	if err := os.Rename(tmpName, s.dbPath); err != nil {
		return fmt.Errorf("cannot move geoip database: %w", err)
	}

	if err := s.writeInfoFile(now); err != nil {
		s.logger.WarnContext(ctx, "Cannot write geoip info file", logger.Error(err))
	}

	return nil
}

func (s *Service) writeInfoFile(now time.Time) error {
	info := struct {
		LastUpdated  string `json:"lastUpdated"`
		Version      string `json:"version"`
		RecordsCount *int   `json:"recordsCount"`
	}{
		LastUpdated:  now.Format(time.RFC3339),
		Version:      "build-" + now.Format("20060102"),
		RecordsCount: nil,
	}

	payload, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("cannot marshal geoip info: %w", err)
	}

	if err := os.WriteFile(s.infoPath, payload, 0o644); err != nil {
		return fmt.Errorf("cannot write geoip info file: %w", err)
	}

	return nil
}

func resolveDbDir() string {
	cwd, err := os.Getwd()
	if err == nil {
		if filepath.Base(cwd) == "api-gql" {
			return filepath.Join(cwd, geoDbDirName)
		}

		appRoot := filepath.Join(cwd, "apps", "api-gql")
		if _, err := os.Stat(appRoot); err == nil {
			return filepath.Join(appRoot, geoDbDirName)
		}
	}

	return geoDbDirName
}

func parseIP(raw string) (netip.Addr, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return netip.Addr{}, errors.New("empty ip")
	}

	if ip := net.ParseIP(raw); ip != nil {
		if addr, ok := netip.AddrFromSlice(ip); ok {
			return addr, nil
		}
	}

	host, _, err := net.SplitHostPort(raw)
	if err == nil {
		if ip := net.ParseIP(host); ip != nil {
			if addr, ok := netip.AddrFromSlice(ip); ok {
				return addr, nil
			}
		}
	}

	return netip.Addr{}, fmt.Errorf("invalid ip: %s", raw)
}
