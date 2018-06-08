package receiver

import (
	"testing"

	"github.com/spf13/viper"
)

var receiver *Receiver

func TestNewReceiver(t *testing.T) {
	viper.AddConfigPath("../config")
	receiver = NewReceiver()
}
