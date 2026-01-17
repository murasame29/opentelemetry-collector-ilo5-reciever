package ilo5receiver

import (
	"context"
	"time"

	"go.uber.org/zap"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/scraper"

	"github.com/murasame29/opentelemetry-collector-ilo5-reciever/internal/ilo"
	"github.com/murasame29/opentelemetry-collector-ilo5-reciever/internal/metadata"
)

type iloScraper struct {
	client   *ilo.Client
	cfg      *Config
	settings receiver.Settings
	mb       *metadata.MetricsBuilder
}

func newScraper(cfg *Config, settings receiver.Settings, client *ilo.Client) (scraper.Metrics, error) {
	s := &iloScraper{
		cfg:      cfg,
		settings: settings,
		client:   client,
		mb:       metadata.NewMetricsBuilder(cfg.MetricsBuilderConfig, settings),
	}
	return scraper.NewMetrics(
		s.scrape,
		scraper.WithStart(s.start),
	)
}

func (s *iloScraper) start(ctx context.Context, _ component.Host) error {
	// Initialize client if not passed (though factory should pass it, but for safety)
	if s.client == nil {
		s.client = ilo.NewClient(ilo.ClientConfig{
			Endpoint:           s.cfg.Endpoint,
			Username:           s.cfg.Username,
			Password:           string(s.cfg.Password),
			InsecureSkipVerify: s.cfg.InsecureSkipVerify,
			Timeout:            10 * time.Second, // Default timeout
		})
	}
	return nil
}

func (s *iloScraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	now := pcommon.NewTimestampFromTime(time.Now())

	var hostName, serialNumber string

	// 1. Get Systems Info
	systems, err := s.client.GetSystems(ctx)
	if err != nil {
		s.settings.Logger.Error("Failed to get systems", zap.Error(err))
	} else {
		for _, sys := range systems {
			systemID := sys.ID

			// Capture hostname and serial from first system
			if hostName == "" {
				hostName = sys.HostName
				serialNumber = sys.SerialNumber
			}

			// Map PowerState: On -> 1, Off/Other -> 0
			powerStateVal := int64(0)
			if sys.PowerState == "On" {
				powerStateVal = 1
			}
			s.mb.RecordIloSystemPowerStateDataPoint(now, powerStateVal, systemID)

			// Map Health: OK -> 1, Warning -> 2, Critical -> 3
			healthVal := int64(0)
			switch sys.Status.Health {
			case "OK":
				healthVal = 1
			case "Warning":
				healthVal = 2
			case "Critical":
				healthVal = 3
			}
			if healthVal > 0 {
				s.mb.RecordIloSystemHealthDataPoint(now, healthVal, systemID)
			}

			// 1.5 Get Storage Drives for this System
			// Note: We need to reconstruct the link if sys.ID is just the ID.
			// Assumption: sys.ID is the short ID (e.g. "1") and the link is standard.
			systemLink := "/redfish/v1/Systems/" + sys.ID
			drivesMap, err := s.client.GetDrives(ctx, systemLink)
			if err != nil {
				// Don't fail everything, just log warn
				s.settings.Logger.Warn("Failed to get drives", zap.String("system", sys.ID), zap.Error(err))
			} else {
				for storageID, drives := range drivesMap {
					for _, drive := range drives {
						healthVal := int64(0)
						switch drive.Status.Health {
						case "OK":
							healthVal = 1
						case "Warning":
							healthVal = 2
						case "Critical":
							healthVal = 3
						}
						if healthVal > 0 {
							s.mb.RecordIloStorageDriveHealthDataPoint(now, healthVal, sys.ID, storageID, drive.ID)
						}
					}
				}
			}
		}
	}

	// 2. Get Chassis Info for Power and Thermal
	chassisIDs, err := s.client.GetChassisIds(ctx)
	if err != nil {
		s.settings.Logger.Error("Failed to get chassis IDs", zap.Error(err))
	} else {
		for _, chassisURI := range chassisIDs {
			// Extract ID from URI (simplification)
			chassisID := chassisURI // simpler for now

			// Power
			power, err := s.client.GetPower(ctx, chassisURI)
			if err != nil {
				s.settings.Logger.Warn("Failed to get power", zap.String("chassis", chassisID), zap.Error(err))
			} else {
				for _, pc := range power.PowerControl {
					s.mb.RecordIloPowerConsumptionDataPoint(now, pc.PowerConsumedWatts, chassisID)
					s.mb.RecordIloPowerCapacityDataPoint(now, pc.PowerCapacityWatts, chassisID)
				}
				for _, v := range power.Voltages {
					s.mb.RecordIloPowerVoltageDataPoint(now, v.ReadingVolts, chassisID, v.MemberID)
				}
				// PSU metrics
				for _, psu := range power.PowerSupplies {
					s.mb.RecordIloPowerPsuOutputDataPoint(now, psu.LastPowerOutputWatts, chassisID, psu.MemberID)
					s.mb.RecordIloPowerPsuInputVoltageDataPoint(now, psu.LineInputVoltage, chassisID, psu.MemberID)
					// Health
					healthVal := int64(0)
					switch psu.Status.Health {
					case "OK":
						healthVal = 1
					case "Warning":
						healthVal = 2
					case "Critical":
						healthVal = 3
					}
					if healthVal > 0 {
						s.mb.RecordIloPowerPsuHealthDataPoint(now, healthVal, chassisID, psu.MemberID)
					}
				}
			}

			// Thermal
			thermal, err := s.client.GetThermal(ctx, chassisURI)
			if err != nil {
				s.settings.Logger.Warn("Failed to get thermal", zap.String("chassis", chassisID), zap.Error(err))
			} else {
				for _, t := range thermal.Temperatures {
					s.mb.RecordIloThermalTemperatureDataPoint(now, t.ReadingCelsius, chassisID, t.Name, t.PhysicalContext, int64(t.Oem.Hpe.LocationXmm), int64(t.Oem.Hpe.LocationYmm))
				}
				for _, f := range thermal.Fans {
					s.mb.RecordIloFanSpeedDataPoint(now, f.Reading, chassisID, f.Name, f.MemberID)
				}
			}
		}
	}

	// Resource Attributes
	rb := s.mb.NewResourceBuilder()
	rb.SetIloEndpoint(s.cfg.Endpoint)
	rb.SetIloModel("iLO 5")
	rb.SetHostName(hostName)
	rb.SetHostSerialNumber(serialNumber)

	return s.mb.Emit(metadata.WithResource(rb.Emit())), nil
}
