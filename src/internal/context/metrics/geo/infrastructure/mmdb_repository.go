package infrastructure

import (
	"context"
	"net/netip"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/geo/domain"
	shared_domain "github.com/aperezgdev/api-snipme/src/internal/context/metrics/shared/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
	"github.com/oschwald/maxminddb-golang/v2"
)

type MMDBRepository struct {
	logger shared_domain_context.Logger
	db maxminddb.Reader
}

func NewMMDBRepository(logger shared_domain_context.Logger, db maxminddb.Reader) *MMDBRepository {
	return &MMDBRepository{
		logger: logger,
		db: db,
	}
}

type geoIPRecord struct {
	Country struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

func (r *MMDBRepository) FindByIp(ctx context.Context, ip shared_domain.Ip) (pkg.Optional[*domain.Country], error) {
	r.logger.Info(ctx, "MMDBRepository - FindByIp - Looking up country for IP", shared_domain_context.NewField("ip", ip))
	var record geoIPRecord
	err := r.db.Lookup(netip.Addr(ip)).Decode(&record)
	if err != nil {
		r.logger.Error(ctx, "MMDBRepository - FindByIp - Error looking up IP in MMDB", shared_domain_context.NewField("ip", ip), shared_domain_context.NewField("error", err.Error()))
		return pkg.EmptyOptional[*domain.Country](), err
	}

	if record.Country.IsoCode == "" {
		r.logger.Info(ctx, "MMDBRepository - FindByIp - No country found for IP", shared_domain_context.NewField("ip", ip))
		return pkg.EmptyOptional[*domain.Country](), nil
	}

	country := &domain.Country{
		ISOCode: shared_domain.CountryCode(record.Country.IsoCode),
	}

	r.logger.Info(ctx, "MMDBRepository - FindByIp - Found country for IP", shared_domain_context.NewField("ip", ip), shared_domain_context.NewField("countryISOCode", country.ISOCode))

	return pkg.Some(country), nil
}