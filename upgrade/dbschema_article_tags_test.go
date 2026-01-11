package upgrade

import (
	"testing"

	"github.com/coscms/webcore/library/config"
	"github.com/stretchr/testify/require"
)

func TestUpgradeArticleTagsData(t *testing.T) {
	config.FromCLI().Conf = `/home/swh/go/src/github.com/admpub/webx/config/config.yaml`
	err := config.ParseConfig()
	require.NoError(t, err)
	err = upgradeArticleTagsData()
	require.NoError(t, err)
}
