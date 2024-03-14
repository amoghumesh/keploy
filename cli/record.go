package cli

import (
	"context"

	"github.com/spf13/cobra"
	"go.keploy.io/server/v2/config"
	recordSvc "go.keploy.io/server/v2/pkg/service/record"
	"go.keploy.io/server/v2/utils"
	"go.uber.org/zap"
)

func init() {
	Register("record", Record)
}

func Record(ctx context.Context, logger *zap.Logger, cfg *config.Config, serviceFactory ServiceFactory, cmdConfigurator CmdConfigurator) *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "record",
		Short:   "record the keploy testcases from the API calls",
		Example: `keploy record -c "/path/to/user/app"`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			return cmdConfigurator.ValidateFlags(ctx, cmd, cfg)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			svc, err := serviceFactory.GetService(ctx, cmd.Name(), *cfg)
			if err != nil {
				utils.LogError(logger, err, "failed to get service")
				return nil
			}
			var record recordSvc.Service
			var ok bool
			if record, ok = svc.(recordSvc.Service); !ok {
				utils.LogError(logger, nil, "service doesn't satisfy record service interface")
				return nil
			}
			err = record.Start(ctx)
			if err != nil {
				utils.LogError(logger, err, "failed to record")
				return nil
			}

			return nil
		},
	}

	err := cmdConfigurator.AddFlags(cmd, cfg)
	if err != nil {
		utils.LogError(logger, err, "failed to add record flags")
		return nil
	}
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	return cmd
}