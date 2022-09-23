// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

import "./@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./@openzeppelin/contracts/access/AccessControl.sol";

contract TokenBase is ERC20, AccessControl {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

    uint8 private _decimals;
    uint256 public faucetAmount;
    uint256 public faucetCallLimit;
    bool public burnEnabled;
    mapping(address => uint256) lastFaucetCall;

    constructor(string memory name_, string memory symbol_, uint8 decimals_, uint256 totalSupply_) ERC20(name_, symbol_) {
        _decimals = decimals_;
        faucetAmount = 10 ** _decimals;
        faucetCallLimit = 86400;  // in seconds
        burnEnabled = true;
        _mint(msg.sender, totalSupply_);
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(MINTER_ROLE, msg.sender);
    }

    function setFaucetAmount(uint256 faucetAmount_) public onlyRole(MINTER_ROLE) {
        faucetAmount = faucetAmount_;
    }

    function setFaucetCallLimit(uint256 faucetCallLimit_) public onlyRole(MINTER_ROLE) {
        faucetCallLimit = faucetCallLimit_;
    }

    function setBurnEnabled(bool burnEnabled_) public onlyRole(MINTER_ROLE) {
        burnEnabled = burnEnabled_;
    }

    function faucet() public {
        require(faucetAmount > 0, "faucet is disabled");
        require(lastFaucetCall[msg.sender] + faucetCallLimit <= block.timestamp, "must wait faucetCallLimit between faucet calls");
        lastFaucetCall[msg.sender] = block.timestamp;
        _mint(_msgSender(), faucetAmount);
    }

    function mint(address to, uint256 amount) public onlyRole(MINTER_ROLE) {
        _mint(to, amount);
    }

    function burn(uint256 amount) public virtual {
        require(burnEnabled, "burn is disabled");
        _burn(_msgSender(), amount);
    }

    function burnFrom(address account, uint256 amount) public virtual onlyRole(MINTER_ROLE) {
        _burn(account, amount);
    }

    /**
     * @dev See {ERC20-decimals}.
     */
    function decimals() public view virtual override returns (uint8) {
        return _decimals;
    }

}
