package fairway

import "github.com/twinj/uuid"

func getUUID() string {
	return uuid.NewV4().String()
}
