module github.com/vegaprotocol/devopstools

go 1.19

require (
	github.com/spf13/cobra v1.5.0
	go.uber.org/zap v1.23.0
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90
)

require (
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/sys v0.0.0-20220702020025-31831981b65f // indirect
	golang.org/x/term v0.0.0-20220526004731-065cf7ba2467 // indirect
)

replace (
	github.com/fergusstrange/embedded-postgres => github.com/vegaprotocol/embedded-postgres v1.13.1-0.20220607151211-5f2f488de508
	github.com/jackc/pgx/v4 v4.14.1 => github.com/pscott31/pgx/v4 v4.16.2-0.20220531164027-bd666b84b61f
	github.com/shopspring/decimal => github.com/vegaprotocol/decimal v1.3.1-uint256
)
