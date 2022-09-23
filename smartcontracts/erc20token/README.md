
## Manually compile and create Go-Bindings

1. Select main Solidity file (e.g. `TokenBase.sol`)
2. Check which version of Solidity to use (e.g. `0.8.16`)
3. Generate ABI, compile and generate Go bindings
4. Store: `TokenBase.abi`, `TokenBase.bin`, `TokenBase.go` only.

### Generate ABI and compile
```bash
docker run -v "$(pwd)":/sources -w /sources \
    ethereum/solc:0.8.16 \
        --abi \
        --bin \
        --overwrite \
        --optimize --optimize-runs 200 \
         -o . \
        TokenBase.sol
```

### Generate Go Bindings
```bash
docker run -v "$(pwd)":/sources -w /sources --pull always \
    ethereum/client-go:alltools-latest \
        abigen \
        --abi=TokenBase.abi \
        --bin=TokenBase.bin \
        --pkg=TokenBase \
        --out=TokenBase.go
```