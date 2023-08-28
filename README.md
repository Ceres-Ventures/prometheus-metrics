# prometheus-metrics
Simple app that collects metrics for prometheus from the Terra2 blockchain

| VARIABLE | DESCRIPTION                                                          |
| ----------- |----------------------------------------------------------------------|
| `BIND_IP` | IP Address the application should listen on. Default `0.0.0.0`       |
| `BIND_PORT` | Port the application should listen on. Default `9292`                |
| `LCD_URL` | Terra API provider                                                   |
| `RPC_ENDPOINT` | RPC endpoint (WIP)                                                   |
| `WALLET_ADDRESSES` | One or more wallet address to query balances for. Separated by a `,` |
| `VALIDATOR_ADDRESSES` | One or more validator to query information about. Separated by a `,`                   |