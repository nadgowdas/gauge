package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/IBM/gauge/pkg/common"
	"github.com/IBM/gauge/pkg/core"
	"github.com/peterbourgon/ff/v3/ffcli"
	"github.com/pkg/errors"
)

//SBOM :
func SBOM() *ffcli.Command {
	var (
		flagset  = flag.NewFlagSet("sbom", flag.ExitOnError)
		sbomfp   = flagset.String("s", "", "sbom filepath")
		format   = flagset.String("o", "", "sbom format (default: cycloneDX)")
		configfp = flagset.String("c", "", "configuration file")
		resultfp = flagset.String("f", "", "result filepath")
		deepscan = flagset.Bool("d", false, "enable deep scan (could be blocked by github API rate-limit)")
	)
	return &ffcli.Command{
		Name:       "sbom",
		ShortUsage: "gauge sbom -s <sbom filepath> -o <sbom format> -f <result filepath>",
		ShortHelp:  `gauge sbom dependencies`,
		LongHelp: `gauge OSS dependencies from SBOM 
EXAMPLES
  # guage sbom 
  gauge sbom -f app-sbom.json -o cyclonedx -c ciso-control.yaml
`,
		FlagSet: flagset,
		Exec: func(ctx context.Context, args []string) error {

			if *sbomfp == "" {
				fmt.Errorf("missing input parameters")
				return errors.New("missing params")
			}
			if os.Getenv(common.GITHUB_API_KEY) == "" {
				fmt.Errorf("please set `GITHUB_API_KEY` env variable")
				return errors.New("missing params")
			}
			if *format == "" {
				*format = "cyclonedx"
			}
			dopts := common.GaugeOpts{}
			dopts.SBOMOptSelected = true
			dopts.SBOMFilepath = *sbomfp
			dopts.ControlFilepath = *configfp
			dopts.SBOMFormat = *format
			dopts.ResultFilepath = *resultfp
			dopts.DeepScanEnabled = *deepscan
			if dopts.ControlFilepath == "" {
				pwd, _ := os.Getwd()
				dopts.ControlFilepath = path.Join(pwd, defaultGaugeConfigFile)
			}
			if err := GuageSBOM(ctx, dopts); err != nil {
				return errors.Wrapf(err, "unpack task for failed")
			}
			return nil
		},
	}
}

//GuageSBOM :
func GuageSBOM(ctx context.Context, dopts common.GaugeOpts) error {
	core.Start(ctx, dopts)
	return nil
}
