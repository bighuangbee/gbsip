package model

import (
	"gosip/data/domain"
)

type Repo struct {
	Device domain.IDeviceRepo
	Channel domain.IChannelRepo
}
