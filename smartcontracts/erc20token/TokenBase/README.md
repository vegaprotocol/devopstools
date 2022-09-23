# Base ERC20 Token

Base ERC20 Token used on Vega internal networks.

It is generated with [Open Zeppelin Wizard](https://wizard.openzeppelin.com/).
Apart from basic ERC20 functionality it offers:
- Access Control based on Roles (by Open Zeppelin) - two roles: Owner and Minter,
- Mintable feature (by Open Zeppelin) - can be called by anyone with Minter role,
- Constructor parameters: `token name`, `token symbol`, `decimals`, `total supply` and `default faucet amount`,
- `faucet()` function that can be called by anyone, that mints predefined amount of tokens to the caller,
