package application

import (
	"context"

	"github.com/aperezgdev/api-snipme/src/internal/context/metrics/link_visit/domain"
	shared_domain_context "github.com/aperezgdev/api-snipme/src/internal/context/shared/domain"
	"github.com/aperezgdev/api-snipme/src/pkg"
)

type LinkVisitProcessor struct {
	logger   shared_domain_context.Logger
	finder   domain.LinkVisitOldsFinder
	remover  domain.LinkVisitBatchRemover
	eventBus shared_domain_context.EventBus
}

func NewLinkVisitProcessor(
	logger shared_domain_context.Logger,
	repo domain.LinkVisitRepository,
	eventBus shared_domain_context.EventBus,
) *LinkVisitProcessor {
	return &LinkVisitProcessor{
		logger:   logger,
		finder:   *domain.NewLinkVisitOldsFinder(logger, repo),
		remover:  *domain.NewLinkVisitBatchRemover(logger, repo),
		eventBus: eventBus,
	}
}

func (p *LinkVisitProcessor) Run(ctx context.Context) error {
	p.logger.Info(ctx, "LinkVisitProcessor - Run - Starting link visit processing")

	linkVisits, err := p.finder.Run(ctx)
	if err != nil {
		p.logger.Error(ctx, "LinkVisitProcessor - Run - Error finding old link visits", shared_domain_context.NewField("error", err.Error()))
		return err
	}

	if len(linkVisits) == 0 {
		p.logger.Info(ctx, "LinkVisitProcessor - Run - No old link visits found to process")
		return nil
	}

	gropedByLink := pkg.GroupBy(linkVisits, func(lv domain.LinkVisit) shared_domain_context.Id {
		return lv.LinkId
	})

	events := p.calculateVisits(ctx, gropedByLink)

	p.logger.Info(ctx, "LinkVisitProcessor - Run - Found old link visits to process", shared_domain_context.NewField("count", len(linkVisits)))
	linkVisitsIds := pkg.Map(linkVisits, func(lv domain.LinkVisit) shared_domain_context.Id {
		return lv.Id
	})

	err = p.remover.Run(ctx, linkVisitsIds)
	if err != nil {
		p.logger.Error(ctx, "LinkVisitProcessor - Run - Error removing old link visits", shared_domain_context.NewField("error", err.Error()))
		return err
	}

	p.logger.Info(ctx, "LinkVisitProcessor - Run - Successfully processed old link visits", shared_domain_context.NewField("count", len(linkVisits)))
	p.eventBus.Publish(ctx, events...)
	return nil
}

func (p *LinkVisitProcessor) calculateVisits(ctx context.Context, linkVisitGrouped map[shared_domain_context.Id][]domain.LinkVisit) []shared_domain_context.DomainEvent {
	var events []shared_domain_context.DomainEvent
	for linkId, visits := range linkVisitGrouped {
		totalViews := uint(len(visits))

		uniqueVisitors := make(map[string]struct{})
		type ipStats struct {
			total    uint
			visitors map[string]struct{}
		}
		ipsStats := make(map[string]*ipStats)

		for _, visit := range visits {
			uniqueVisitors[visit.VisitorId.String()] = struct{}{}
			ipStr := visit.Ip.String()
			if ipsStats[ipStr] == nil {
				ipsStats[ipStr] = &ipStats{visitors: make(map[string]struct{})}
			}
			ipsStats[ipStr].total++
			ipsStats[ipStr].visitors[visit.VisitorId.String()] = struct{}{}
		}
		uniqueViews := uint(len(uniqueVisitors))

		var ipsFrequency []domain.IpsVisits
		for ip, stats := range ipsStats {
			ipsFrequency = append(ipsFrequency, domain.IpsVisits{
				Ip:          ip,
				TotalViews:  stats.total,
				UniqueViews: uint(len(stats.visitors)),
			})
		}

		events = append(events, domain.NewLinkVisitsProcessedDomainEvent(
			linkId.String(),
			totalViews,
			uniqueViews,
			ipsFrequency,
		))

		p.logger.Info(
			ctx, "LinkVisitProcessor - Run - Processed visits for link",
			shared_domain_context.NewField("linkId", linkId.String()),
			shared_domain_context.NewField("totalViews", totalViews),
			shared_domain_context.NewField("uniqueViews", uniqueViews),
		)
	}
	return events
}
