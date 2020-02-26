package prismaiface

import (
	"io"

	"github.com/CityOfNewYork/prisma-cloud-remediation/api/prisma"
)

// PrismaAPI allow mock test prisma.PrismaClient
type PrismaAPI interface {
	Request(*prisma.PrismaAPIRequestInput) (io.ReadCloser, error)
	ListAlerts(*prisma.ListAlertsInput) (*prisma.Alerts, error)
	DismissAlerts(*prisma.DismissAlertInput) ([]byte, error)
	LoginPrisma(*prisma.LoginPrismaInput) error
	ListAccountGroups() (*prisma.AccountGroups, error)
	ListAccountNames() (*prisma.AccountNames, error)
	RegisterAccount(payload []byte) error
}

var _ PrismaAPI = (*prisma.PrismaClient)(nil)
