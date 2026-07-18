// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// Campaign is an auto generated low-level Go binding around an user-defined struct.
type Campaign struct {
	Owner             common.Address
	Title             string
	Description       string
	CurrencyMode      uint8
	Token             common.Address
	GoalEth           *big.Int
	GoalToken         *big.Int
	Deadline          *big.Int
	AmountRaisedEth   *big.Int
	AmountRaisedToken *big.Int
	Withdrawn         bool
}

// CrowdFundingMetaData contains all meta data concerning the CrowdFunding contract.
var CrowdFundingMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"MAX_PAGE_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"campaignCount\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"campaigns\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"title\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"currencyMode\",\"type\":\"uint8\",\"internalType\":\"enumCurrencyMode\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"goalEth\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"goalToken\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountRaisedEth\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountRaisedToken\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"withdrawn\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"closeCampaign\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"contributeEth\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"contributeToken\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"contributionsEth\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"contributionsToken\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createCampaign\",\"inputs\":[{\"name\":\"title\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"currencyMode\",\"type\":\"uint8\",\"internalType\":\"enumCurrencyMode\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"goalEth\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"goalToken\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"durationInSeconds\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getCampaign\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structCampaign\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"title\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"currencyMode\",\"type\":\"uint8\",\"internalType\":\"enumCurrencyMode\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"goalEth\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"goalToken\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountRaisedEth\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountRaisedToken\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"withdrawn\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCampaignStatus\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumCampaignStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCampaigns\",\"inputs\":[{\"name\":\"offset\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"limit\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structCampaign[]\",\"components\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"title\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"description\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"currencyMode\",\"type\":\"uint8\",\"internalType\":\"enumCurrencyMode\"},{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"goalEth\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"goalToken\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountRaisedEth\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amountRaisedToken\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"withdrawn\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getContributionEth\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"contributor\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getContributionToken\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"contributor\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"refund\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"withdraw\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CampaignClosed\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"CampaignCreated\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"currencyMode\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumCurrencyMode\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"goalEth\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"goalToken\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"deadline\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ContributionMade\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"contributor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"ContributionRefunded\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"contributor\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"ethAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FundsWithdrawn\",\"inputs\":[{\"name\":\"campaignId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"ethAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"tokenAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"CampaignDoesNotExist\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CampaignHasEnded\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"CampaignStillActive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ContributionMustBeGreaterThanZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"DurationMustBeGreaterThanZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EthNotAccepted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FundsAlreadyWithdrawn\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GoalAlreadyReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GoalMustBeGreaterThanZero\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GoalNotReached\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoContributionToRefund\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotCampaignOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"TokenAddressRequired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TokenNotAccepted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransferFailed\",\"inputs\":[]}]",
}

// CrowdFundingABI is the input ABI used to generate the binding from.
// Deprecated: Use CrowdFundingMetaData.ABI instead.
var CrowdFundingABI = CrowdFundingMetaData.ABI

// CrowdFunding is an auto generated Go binding around an Ethereum contract.
type CrowdFunding struct {
	CrowdFundingCaller     // Read-only binding to the contract
	CrowdFundingTransactor // Write-only binding to the contract
	CrowdFundingFilterer   // Log filterer for contract events
}

// CrowdFundingCaller is an auto generated read-only Go binding around an Ethereum contract.
type CrowdFundingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrowdFundingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CrowdFundingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrowdFundingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CrowdFundingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CrowdFundingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CrowdFundingSession struct {
	Contract     *CrowdFunding     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CrowdFundingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CrowdFundingCallerSession struct {
	Contract *CrowdFundingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// CrowdFundingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CrowdFundingTransactorSession struct {
	Contract     *CrowdFundingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// CrowdFundingRaw is an auto generated low-level Go binding around an Ethereum contract.
type CrowdFundingRaw struct {
	Contract *CrowdFunding // Generic contract binding to access the raw methods on
}

// CrowdFundingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CrowdFundingCallerRaw struct {
	Contract *CrowdFundingCaller // Generic read-only contract binding to access the raw methods on
}

// CrowdFundingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CrowdFundingTransactorRaw struct {
	Contract *CrowdFundingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCrowdFunding creates a new instance of CrowdFunding, bound to a specific deployed contract.
func NewCrowdFunding(address common.Address, backend bind.ContractBackend) (*CrowdFunding, error) {
	contract, err := bindCrowdFunding(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CrowdFunding{CrowdFundingCaller: CrowdFundingCaller{contract: contract}, CrowdFundingTransactor: CrowdFundingTransactor{contract: contract}, CrowdFundingFilterer: CrowdFundingFilterer{contract: contract}}, nil
}

// NewCrowdFundingCaller creates a new read-only instance of CrowdFunding, bound to a specific deployed contract.
func NewCrowdFundingCaller(address common.Address, caller bind.ContractCaller) (*CrowdFundingCaller, error) {
	contract, err := bindCrowdFunding(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CrowdFundingCaller{contract: contract}, nil
}

// NewCrowdFundingTransactor creates a new write-only instance of CrowdFunding, bound to a specific deployed contract.
func NewCrowdFundingTransactor(address common.Address, transactor bind.ContractTransactor) (*CrowdFundingTransactor, error) {
	contract, err := bindCrowdFunding(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CrowdFundingTransactor{contract: contract}, nil
}

// NewCrowdFundingFilterer creates a new log filterer instance of CrowdFunding, bound to a specific deployed contract.
func NewCrowdFundingFilterer(address common.Address, filterer bind.ContractFilterer) (*CrowdFundingFilterer, error) {
	contract, err := bindCrowdFunding(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CrowdFundingFilterer{contract: contract}, nil
}

// bindCrowdFunding binds a generic wrapper to an already deployed contract.
func bindCrowdFunding(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CrowdFundingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrowdFunding *CrowdFundingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrowdFunding.Contract.CrowdFundingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrowdFunding *CrowdFundingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrowdFunding.Contract.CrowdFundingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrowdFunding *CrowdFundingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrowdFunding.Contract.CrowdFundingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CrowdFunding *CrowdFundingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CrowdFunding.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CrowdFunding *CrowdFundingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CrowdFunding.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CrowdFunding *CrowdFundingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CrowdFunding.Contract.contract.Transact(opts, method, params...)
}

// MAXPAGESIZE is a free data retrieval call binding the contract method 0x48f4da20.
//
// Solidity: function MAX_PAGE_SIZE() view returns(uint256)
func (_CrowdFunding *CrowdFundingCaller) MAXPAGESIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CrowdFunding.contract.Call(opts, &out, "MAX_PAGE_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXPAGESIZE is a free data retrieval call binding the contract method 0x48f4da20.
//
// Solidity: function MAX_PAGE_SIZE() view returns(uint256)
func (_CrowdFunding *CrowdFundingSession) MAXPAGESIZE() (*big.Int, error) {
	return _CrowdFunding.Contract.MAXPAGESIZE(&_CrowdFunding.CallOpts)
}

// MAXPAGESIZE is a free data retrieval call binding the contract method 0x48f4da20.
//
// Solidity: function MAX_PAGE_SIZE() view returns(uint256)
func (_CrowdFunding *CrowdFundingCallerSession) MAXPAGESIZE() (*big.Int, error) {
	return _CrowdFunding.Contract.MAXPAGESIZE(&_CrowdFunding.CallOpts)
}

// CampaignCount is a free data retrieval call binding the contract method 0x7274e30d.
//
// Solidity: function campaignCount() view returns(uint256)
func (_CrowdFunding *CrowdFundingCaller) CampaignCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CrowdFunding.contract.Call(opts, &out, "campaignCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CampaignCount is a free data retrieval call binding the contract method 0x7274e30d.
//
// Solidity: function campaignCount() view returns(uint256)
func (_CrowdFunding *CrowdFundingSession) CampaignCount() (*big.Int, error) {
	return _CrowdFunding.Contract.CampaignCount(&_CrowdFunding.CallOpts)
}

// CampaignCount is a free data retrieval call binding the contract method 0x7274e30d.
//
// Solidity: function campaignCount() view returns(uint256)
func (_CrowdFunding *CrowdFundingCallerSession) CampaignCount() (*big.Int, error) {
	return _CrowdFunding.Contract.CampaignCount(&_CrowdFunding.CallOpts)
}

// Campaigns is a free data retrieval call binding the contract method 0x141961bc.
//
// Solidity: function campaigns(uint256 ) view returns(address owner, string title, string description, uint8 currencyMode, address token, uint256 goalEth, uint256 goalToken, uint256 deadline, uint256 amountRaisedEth, uint256 amountRaisedToken, bool withdrawn)
func (_CrowdFunding *CrowdFundingCaller) Campaigns(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Owner             common.Address
	Title             string
	Description       string
	CurrencyMode      uint8
	Token             common.Address
	GoalEth           *big.Int
	GoalToken         *big.Int
	Deadline          *big.Int
	AmountRaisedEth   *big.Int
	AmountRaisedToken *big.Int
	Withdrawn         bool
}, error) {
	var out []interface{}
	err := _CrowdFunding.contract.Call(opts, &out, "campaigns", arg0)

	outstruct := new(struct {
		Owner             common.Address
		Title             string
		Description       string
		CurrencyMode      uint8
		Token             common.Address
		GoalEth           *big.Int
		GoalToken         *big.Int
		Deadline          *big.Int
		AmountRaisedEth   *big.Int
		AmountRaisedToken *big.Int
		Withdrawn         bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Owner = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Title = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Description = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.CurrencyMode = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.Token = *abi.ConvertType(out[4], new(common.Address)).(*common.Address)
	outstruct.GoalEth = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.GoalToken = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.Deadline = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.AmountRaisedEth = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.AmountRaisedToken = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)
	outstruct.Withdrawn = *abi.ConvertType(out[10], new(bool)).(*bool)

	return *outstruct, err

}

// Campaigns is a free data retrieval call binding the contract method 0x141961bc.
//
// Solidity: function campaigns(uint256 ) view returns(address owner, string title, string description, uint8 currencyMode, address token, uint256 goalEth, uint256 goalToken, uint256 deadline, uint256 amountRaisedEth, uint256 amountRaisedToken, bool withdrawn)
func (_CrowdFunding *CrowdFundingSession) Campaigns(arg0 *big.Int) (struct {
	Owner             common.Address
	Title             string
	Description       string
	CurrencyMode      uint8
	Token             common.Address
	GoalEth           *big.Int
	GoalToken         *big.Int
	Deadline          *big.Int
	AmountRaisedEth   *big.Int
	AmountRaisedToken *big.Int
	Withdrawn         bool
}, error) {
	return _CrowdFunding.Contract.Campaigns(&_CrowdFunding.CallOpts, arg0)
}

// Campaigns is a free data retrieval call binding the contract method 0x141961bc.
//
// Solidity: function campaigns(uint256 ) view returns(address owner, string title, string description, uint8 currencyMode, address token, uint256 goalEth, uint256 goalToken, uint256 deadline, uint256 amountRaisedEth, uint256 amountRaisedToken, bool withdrawn)
func (_CrowdFunding *CrowdFundingCallerSession) Campaigns(arg0 *big.Int) (struct {
	Owner             common.Address
	Title             string
	Description       string
	CurrencyMode      uint8
	Token             common.Address
	GoalEth           *big.Int
	GoalToken         *big.Int
	Deadline          *big.Int
	AmountRaisedEth   *big.Int
	AmountRaisedToken *big.Int
	Withdrawn         bool
}, error) {
	return _CrowdFunding.Contract.Campaigns(&_CrowdFunding.CallOpts, arg0)
}

// ContributionsEth is a free data retrieval call binding the contract method 0xdac9998e.
//
// Solidity: function contributionsEth(uint256 , address ) view returns(uint256)
func (_CrowdFunding *CrowdFundingCaller) ContributionsEth(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CrowdFunding.contract.Call(opts, &out, "contributionsEth", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ContributionsEth is a free data retrieval call binding the contract method 0xdac9998e.
//
// Solidity: function contributionsEth(uint256 , address ) view returns(uint256)
func (_CrowdFunding *CrowdFundingSession) ContributionsEth(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _CrowdFunding.Contract.ContributionsEth(&_CrowdFunding.CallOpts, arg0, arg1)
}

// ContributionsEth is a free data retrieval call binding the contract method 0xdac9998e.
//
// Solidity: function contributionsEth(uint256 , address ) view returns(uint256)
func (_CrowdFunding *CrowdFundingCallerSession) ContributionsEth(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _CrowdFunding.Contract.ContributionsEth(&_CrowdFunding.CallOpts, arg0, arg1)
}

// ContributionsToken is a free data retrieval call binding the contract method 0x6385fe75.
//
// Solidity: function contributionsToken(uint256 , address ) view returns(uint256)
func (_CrowdFunding *CrowdFundingCaller) ContributionsToken(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CrowdFunding.contract.Call(opts, &out, "contributionsToken", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ContributionsToken is a free data retrieval call binding the contract method 0x6385fe75.
//
// Solidity: function contributionsToken(uint256 , address ) view returns(uint256)
func (_CrowdFunding *CrowdFundingSession) ContributionsToken(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _CrowdFunding.Contract.ContributionsToken(&_CrowdFunding.CallOpts, arg0, arg1)
}

// ContributionsToken is a free data retrieval call binding the contract method 0x6385fe75.
//
// Solidity: function contributionsToken(uint256 , address ) view returns(uint256)
func (_CrowdFunding *CrowdFundingCallerSession) ContributionsToken(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _CrowdFunding.Contract.ContributionsToken(&_CrowdFunding.CallOpts, arg0, arg1)
}

// GetCampaign is a free data retrieval call binding the contract method 0x5598f8cc.
//
// Solidity: function getCampaign(uint256 campaignId) view returns((address,string,string,uint8,address,uint256,uint256,uint256,uint256,uint256,bool))
func (_CrowdFunding *CrowdFundingCaller) GetCampaign(opts *bind.CallOpts, campaignId *big.Int) (Campaign, error) {
	var out []interface{}
	err := _CrowdFunding.contract.Call(opts, &out, "getCampaign", campaignId)

	if err != nil {
		return *new(Campaign), err
	}

	out0 := *abi.ConvertType(out[0], new(Campaign)).(*Campaign)

	return out0, err

}

// GetCampaign is a free data retrieval call binding the contract method 0x5598f8cc.
//
// Solidity: function getCampaign(uint256 campaignId) view returns((address,string,string,uint8,address,uint256,uint256,uint256,uint256,uint256,bool))
func (_CrowdFunding *CrowdFundingSession) GetCampaign(campaignId *big.Int) (Campaign, error) {
	return _CrowdFunding.Contract.GetCampaign(&_CrowdFunding.CallOpts, campaignId)
}

// GetCampaign is a free data retrieval call binding the contract method 0x5598f8cc.
//
// Solidity: function getCampaign(uint256 campaignId) view returns((address,string,string,uint8,address,uint256,uint256,uint256,uint256,uint256,bool))
func (_CrowdFunding *CrowdFundingCallerSession) GetCampaign(campaignId *big.Int) (Campaign, error) {
	return _CrowdFunding.Contract.GetCampaign(&_CrowdFunding.CallOpts, campaignId)
}

// GetCampaignStatus is a free data retrieval call binding the contract method 0x6c19e004.
//
// Solidity: function getCampaignStatus(uint256 campaignId) view returns(uint8)
func (_CrowdFunding *CrowdFundingCaller) GetCampaignStatus(opts *bind.CallOpts, campaignId *big.Int) (uint8, error) {
	var out []interface{}
	err := _CrowdFunding.contract.Call(opts, &out, "getCampaignStatus", campaignId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetCampaignStatus is a free data retrieval call binding the contract method 0x6c19e004.
//
// Solidity: function getCampaignStatus(uint256 campaignId) view returns(uint8)
func (_CrowdFunding *CrowdFundingSession) GetCampaignStatus(campaignId *big.Int) (uint8, error) {
	return _CrowdFunding.Contract.GetCampaignStatus(&_CrowdFunding.CallOpts, campaignId)
}

// GetCampaignStatus is a free data retrieval call binding the contract method 0x6c19e004.
//
// Solidity: function getCampaignStatus(uint256 campaignId) view returns(uint8)
func (_CrowdFunding *CrowdFundingCallerSession) GetCampaignStatus(campaignId *big.Int) (uint8, error) {
	return _CrowdFunding.Contract.GetCampaignStatus(&_CrowdFunding.CallOpts, campaignId)
}

// GetCampaigns is a free data retrieval call binding the contract method 0x09051566.
//
// Solidity: function getCampaigns(uint256 offset, uint256 limit) view returns((address,string,string,uint8,address,uint256,uint256,uint256,uint256,uint256,bool)[])
func (_CrowdFunding *CrowdFundingCaller) GetCampaigns(opts *bind.CallOpts, offset *big.Int, limit *big.Int) ([]Campaign, error) {
	var out []interface{}
	err := _CrowdFunding.contract.Call(opts, &out, "getCampaigns", offset, limit)

	if err != nil {
		return *new([]Campaign), err
	}

	out0 := *abi.ConvertType(out[0], new([]Campaign)).(*[]Campaign)

	return out0, err

}

// GetCampaigns is a free data retrieval call binding the contract method 0x09051566.
//
// Solidity: function getCampaigns(uint256 offset, uint256 limit) view returns((address,string,string,uint8,address,uint256,uint256,uint256,uint256,uint256,bool)[])
func (_CrowdFunding *CrowdFundingSession) GetCampaigns(offset *big.Int, limit *big.Int) ([]Campaign, error) {
	return _CrowdFunding.Contract.GetCampaigns(&_CrowdFunding.CallOpts, offset, limit)
}

// GetCampaigns is a free data retrieval call binding the contract method 0x09051566.
//
// Solidity: function getCampaigns(uint256 offset, uint256 limit) view returns((address,string,string,uint8,address,uint256,uint256,uint256,uint256,uint256,bool)[])
func (_CrowdFunding *CrowdFundingCallerSession) GetCampaigns(offset *big.Int, limit *big.Int) ([]Campaign, error) {
	return _CrowdFunding.Contract.GetCampaigns(&_CrowdFunding.CallOpts, offset, limit)
}

// GetContributionEth is a free data retrieval call binding the contract method 0xa0c1e469.
//
// Solidity: function getContributionEth(uint256 campaignId, address contributor) view returns(uint256)
func (_CrowdFunding *CrowdFundingCaller) GetContributionEth(opts *bind.CallOpts, campaignId *big.Int, contributor common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CrowdFunding.contract.Call(opts, &out, "getContributionEth", campaignId, contributor)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetContributionEth is a free data retrieval call binding the contract method 0xa0c1e469.
//
// Solidity: function getContributionEth(uint256 campaignId, address contributor) view returns(uint256)
func (_CrowdFunding *CrowdFundingSession) GetContributionEth(campaignId *big.Int, contributor common.Address) (*big.Int, error) {
	return _CrowdFunding.Contract.GetContributionEth(&_CrowdFunding.CallOpts, campaignId, contributor)
}

// GetContributionEth is a free data retrieval call binding the contract method 0xa0c1e469.
//
// Solidity: function getContributionEth(uint256 campaignId, address contributor) view returns(uint256)
func (_CrowdFunding *CrowdFundingCallerSession) GetContributionEth(campaignId *big.Int, contributor common.Address) (*big.Int, error) {
	return _CrowdFunding.Contract.GetContributionEth(&_CrowdFunding.CallOpts, campaignId, contributor)
}

// GetContributionToken is a free data retrieval call binding the contract method 0x79873873.
//
// Solidity: function getContributionToken(uint256 campaignId, address contributor) view returns(uint256)
func (_CrowdFunding *CrowdFundingCaller) GetContributionToken(opts *bind.CallOpts, campaignId *big.Int, contributor common.Address) (*big.Int, error) {
	var out []interface{}
	err := _CrowdFunding.contract.Call(opts, &out, "getContributionToken", campaignId, contributor)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetContributionToken is a free data retrieval call binding the contract method 0x79873873.
//
// Solidity: function getContributionToken(uint256 campaignId, address contributor) view returns(uint256)
func (_CrowdFunding *CrowdFundingSession) GetContributionToken(campaignId *big.Int, contributor common.Address) (*big.Int, error) {
	return _CrowdFunding.Contract.GetContributionToken(&_CrowdFunding.CallOpts, campaignId, contributor)
}

// GetContributionToken is a free data retrieval call binding the contract method 0x79873873.
//
// Solidity: function getContributionToken(uint256 campaignId, address contributor) view returns(uint256)
func (_CrowdFunding *CrowdFundingCallerSession) GetContributionToken(campaignId *big.Int, contributor common.Address) (*big.Int, error) {
	return _CrowdFunding.Contract.GetContributionToken(&_CrowdFunding.CallOpts, campaignId, contributor)
}

// CloseCampaign is a paid mutator transaction binding the contract method 0xb0e1c1e1.
//
// Solidity: function closeCampaign(uint256 campaignId) returns()
func (_CrowdFunding *CrowdFundingTransactor) CloseCampaign(opts *bind.TransactOpts, campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.contract.Transact(opts, "closeCampaign", campaignId)
}

// CloseCampaign is a paid mutator transaction binding the contract method 0xb0e1c1e1.
//
// Solidity: function closeCampaign(uint256 campaignId) returns()
func (_CrowdFunding *CrowdFundingSession) CloseCampaign(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.Contract.CloseCampaign(&_CrowdFunding.TransactOpts, campaignId)
}

// CloseCampaign is a paid mutator transaction binding the contract method 0xb0e1c1e1.
//
// Solidity: function closeCampaign(uint256 campaignId) returns()
func (_CrowdFunding *CrowdFundingTransactorSession) CloseCampaign(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.Contract.CloseCampaign(&_CrowdFunding.TransactOpts, campaignId)
}

// ContributeEth is a paid mutator transaction binding the contract method 0x352d722f.
//
// Solidity: function contributeEth(uint256 campaignId) payable returns()
func (_CrowdFunding *CrowdFundingTransactor) ContributeEth(opts *bind.TransactOpts, campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.contract.Transact(opts, "contributeEth", campaignId)
}

// ContributeEth is a paid mutator transaction binding the contract method 0x352d722f.
//
// Solidity: function contributeEth(uint256 campaignId) payable returns()
func (_CrowdFunding *CrowdFundingSession) ContributeEth(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.Contract.ContributeEth(&_CrowdFunding.TransactOpts, campaignId)
}

// ContributeEth is a paid mutator transaction binding the contract method 0x352d722f.
//
// Solidity: function contributeEth(uint256 campaignId) payable returns()
func (_CrowdFunding *CrowdFundingTransactorSession) ContributeEth(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.Contract.ContributeEth(&_CrowdFunding.TransactOpts, campaignId)
}

// ContributeToken is a paid mutator transaction binding the contract method 0x3d925092.
//
// Solidity: function contributeToken(uint256 campaignId, uint256 amount) returns()
func (_CrowdFunding *CrowdFundingTransactor) ContributeToken(opts *bind.TransactOpts, campaignId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.contract.Transact(opts, "contributeToken", campaignId, amount)
}

// ContributeToken is a paid mutator transaction binding the contract method 0x3d925092.
//
// Solidity: function contributeToken(uint256 campaignId, uint256 amount) returns()
func (_CrowdFunding *CrowdFundingSession) ContributeToken(campaignId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.Contract.ContributeToken(&_CrowdFunding.TransactOpts, campaignId, amount)
}

// ContributeToken is a paid mutator transaction binding the contract method 0x3d925092.
//
// Solidity: function contributeToken(uint256 campaignId, uint256 amount) returns()
func (_CrowdFunding *CrowdFundingTransactorSession) ContributeToken(campaignId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.Contract.ContributeToken(&_CrowdFunding.TransactOpts, campaignId, amount)
}

// CreateCampaign is a paid mutator transaction binding the contract method 0x0a55df74.
//
// Solidity: function createCampaign(string title, string description, uint8 currencyMode, address token, uint256 goalEth, uint256 goalToken, uint256 durationInSeconds) returns(uint256 campaignId)
func (_CrowdFunding *CrowdFundingTransactor) CreateCampaign(opts *bind.TransactOpts, title string, description string, currencyMode uint8, token common.Address, goalEth *big.Int, goalToken *big.Int, durationInSeconds *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.contract.Transact(opts, "createCampaign", title, description, currencyMode, token, goalEth, goalToken, durationInSeconds)
}

// CreateCampaign is a paid mutator transaction binding the contract method 0x0a55df74.
//
// Solidity: function createCampaign(string title, string description, uint8 currencyMode, address token, uint256 goalEth, uint256 goalToken, uint256 durationInSeconds) returns(uint256 campaignId)
func (_CrowdFunding *CrowdFundingSession) CreateCampaign(title string, description string, currencyMode uint8, token common.Address, goalEth *big.Int, goalToken *big.Int, durationInSeconds *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.Contract.CreateCampaign(&_CrowdFunding.TransactOpts, title, description, currencyMode, token, goalEth, goalToken, durationInSeconds)
}

// CreateCampaign is a paid mutator transaction binding the contract method 0x0a55df74.
//
// Solidity: function createCampaign(string title, string description, uint8 currencyMode, address token, uint256 goalEth, uint256 goalToken, uint256 durationInSeconds) returns(uint256 campaignId)
func (_CrowdFunding *CrowdFundingTransactorSession) CreateCampaign(title string, description string, currencyMode uint8, token common.Address, goalEth *big.Int, goalToken *big.Int, durationInSeconds *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.Contract.CreateCampaign(&_CrowdFunding.TransactOpts, title, description, currencyMode, token, goalEth, goalToken, durationInSeconds)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 campaignId) returns()
func (_CrowdFunding *CrowdFundingTransactor) Refund(opts *bind.TransactOpts, campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.contract.Transact(opts, "refund", campaignId)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 campaignId) returns()
func (_CrowdFunding *CrowdFundingSession) Refund(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.Contract.Refund(&_CrowdFunding.TransactOpts, campaignId)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 campaignId) returns()
func (_CrowdFunding *CrowdFundingTransactorSession) Refund(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.Contract.Refund(&_CrowdFunding.TransactOpts, campaignId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 campaignId) returns()
func (_CrowdFunding *CrowdFundingTransactor) Withdraw(opts *bind.TransactOpts, campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.contract.Transact(opts, "withdraw", campaignId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 campaignId) returns()
func (_CrowdFunding *CrowdFundingSession) Withdraw(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.Contract.Withdraw(&_CrowdFunding.TransactOpts, campaignId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 campaignId) returns()
func (_CrowdFunding *CrowdFundingTransactorSession) Withdraw(campaignId *big.Int) (*types.Transaction, error) {
	return _CrowdFunding.Contract.Withdraw(&_CrowdFunding.TransactOpts, campaignId)
}

// CrowdFundingCampaignClosedIterator is returned from FilterCampaignClosed and is used to iterate over the raw logs and unpacked data for CampaignClosed events raised by the CrowdFunding contract.
type CrowdFundingCampaignClosedIterator struct {
	Event *CrowdFundingCampaignClosed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CrowdFundingCampaignClosedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdFundingCampaignClosed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CrowdFundingCampaignClosed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CrowdFundingCampaignClosedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdFundingCampaignClosedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdFundingCampaignClosed represents a CampaignClosed event raised by the CrowdFunding contract.
type CrowdFundingCampaignClosed struct {
	CampaignId *big.Int
	Owner      common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCampaignClosed is a free log retrieval operation binding the contract event 0xa78af03f83be4496e7f3344865e5d4a44fa8ed036ba56cabd413e1882b92a79f.
//
// Solidity: event CampaignClosed(uint256 indexed campaignId, address indexed owner)
func (_CrowdFunding *CrowdFundingFilterer) FilterCampaignClosed(opts *bind.FilterOpts, campaignId []*big.Int, owner []common.Address) (*CrowdFundingCampaignClosedIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _CrowdFunding.contract.FilterLogs(opts, "CampaignClosed", campaignIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &CrowdFundingCampaignClosedIterator{contract: _CrowdFunding.contract, event: "CampaignClosed", logs: logs, sub: sub}, nil
}

// WatchCampaignClosed is a free log subscription operation binding the contract event 0xa78af03f83be4496e7f3344865e5d4a44fa8ed036ba56cabd413e1882b92a79f.
//
// Solidity: event CampaignClosed(uint256 indexed campaignId, address indexed owner)
func (_CrowdFunding *CrowdFundingFilterer) WatchCampaignClosed(opts *bind.WatchOpts, sink chan<- *CrowdFundingCampaignClosed, campaignId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _CrowdFunding.contract.WatchLogs(opts, "CampaignClosed", campaignIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdFundingCampaignClosed)
				if err := _CrowdFunding.contract.UnpackLog(event, "CampaignClosed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCampaignClosed is a log parse operation binding the contract event 0xa78af03f83be4496e7f3344865e5d4a44fa8ed036ba56cabd413e1882b92a79f.
//
// Solidity: event CampaignClosed(uint256 indexed campaignId, address indexed owner)
func (_CrowdFunding *CrowdFundingFilterer) ParseCampaignClosed(log types.Log) (*CrowdFundingCampaignClosed, error) {
	event := new(CrowdFundingCampaignClosed)
	if err := _CrowdFunding.contract.UnpackLog(event, "CampaignClosed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrowdFundingCampaignCreatedIterator is returned from FilterCampaignCreated and is used to iterate over the raw logs and unpacked data for CampaignCreated events raised by the CrowdFunding contract.
type CrowdFundingCampaignCreatedIterator struct {
	Event *CrowdFundingCampaignCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CrowdFundingCampaignCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdFundingCampaignCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CrowdFundingCampaignCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CrowdFundingCampaignCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdFundingCampaignCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdFundingCampaignCreated represents a CampaignCreated event raised by the CrowdFunding contract.
type CrowdFundingCampaignCreated struct {
	CampaignId   *big.Int
	Owner        common.Address
	CurrencyMode uint8
	Token        common.Address
	GoalEth      *big.Int
	GoalToken    *big.Int
	Deadline     *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterCampaignCreated is a free log retrieval operation binding the contract event 0xfe08ee610ce61d7fa2fac7befaa25cae823cb93466a011bb8beaac23973bca37.
//
// Solidity: event CampaignCreated(uint256 indexed campaignId, address indexed owner, uint8 currencyMode, address token, uint256 goalEth, uint256 goalToken, uint256 deadline)
func (_CrowdFunding *CrowdFundingFilterer) FilterCampaignCreated(opts *bind.FilterOpts, campaignId []*big.Int, owner []common.Address) (*CrowdFundingCampaignCreatedIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _CrowdFunding.contract.FilterLogs(opts, "CampaignCreated", campaignIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &CrowdFundingCampaignCreatedIterator{contract: _CrowdFunding.contract, event: "CampaignCreated", logs: logs, sub: sub}, nil
}

// WatchCampaignCreated is a free log subscription operation binding the contract event 0xfe08ee610ce61d7fa2fac7befaa25cae823cb93466a011bb8beaac23973bca37.
//
// Solidity: event CampaignCreated(uint256 indexed campaignId, address indexed owner, uint8 currencyMode, address token, uint256 goalEth, uint256 goalToken, uint256 deadline)
func (_CrowdFunding *CrowdFundingFilterer) WatchCampaignCreated(opts *bind.WatchOpts, sink chan<- *CrowdFundingCampaignCreated, campaignId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _CrowdFunding.contract.WatchLogs(opts, "CampaignCreated", campaignIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdFundingCampaignCreated)
				if err := _CrowdFunding.contract.UnpackLog(event, "CampaignCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCampaignCreated is a log parse operation binding the contract event 0xfe08ee610ce61d7fa2fac7befaa25cae823cb93466a011bb8beaac23973bca37.
//
// Solidity: event CampaignCreated(uint256 indexed campaignId, address indexed owner, uint8 currencyMode, address token, uint256 goalEth, uint256 goalToken, uint256 deadline)
func (_CrowdFunding *CrowdFundingFilterer) ParseCampaignCreated(log types.Log) (*CrowdFundingCampaignCreated, error) {
	event := new(CrowdFundingCampaignCreated)
	if err := _CrowdFunding.contract.UnpackLog(event, "CampaignCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrowdFundingContributionMadeIterator is returned from FilterContributionMade and is used to iterate over the raw logs and unpacked data for ContributionMade events raised by the CrowdFunding contract.
type CrowdFundingContributionMadeIterator struct {
	Event *CrowdFundingContributionMade // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CrowdFundingContributionMadeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdFundingContributionMade)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CrowdFundingContributionMade)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CrowdFundingContributionMadeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdFundingContributionMadeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdFundingContributionMade represents a ContributionMade event raised by the CrowdFunding contract.
type CrowdFundingContributionMade struct {
	CampaignId  *big.Int
	Contributor common.Address
	Token       common.Address
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterContributionMade is a free log retrieval operation binding the contract event 0xf6ab15f3ae6cc7f857f32b5a5e1a789bf16059eb626d7aefab39f2ff0447e0df.
//
// Solidity: event ContributionMade(uint256 indexed campaignId, address indexed contributor, address token, uint256 amount)
func (_CrowdFunding *CrowdFundingFilterer) FilterContributionMade(opts *bind.FilterOpts, campaignId []*big.Int, contributor []common.Address) (*CrowdFundingContributionMadeIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var contributorRule []interface{}
	for _, contributorItem := range contributor {
		contributorRule = append(contributorRule, contributorItem)
	}

	logs, sub, err := _CrowdFunding.contract.FilterLogs(opts, "ContributionMade", campaignIdRule, contributorRule)
	if err != nil {
		return nil, err
	}
	return &CrowdFundingContributionMadeIterator{contract: _CrowdFunding.contract, event: "ContributionMade", logs: logs, sub: sub}, nil
}

// WatchContributionMade is a free log subscription operation binding the contract event 0xf6ab15f3ae6cc7f857f32b5a5e1a789bf16059eb626d7aefab39f2ff0447e0df.
//
// Solidity: event ContributionMade(uint256 indexed campaignId, address indexed contributor, address token, uint256 amount)
func (_CrowdFunding *CrowdFundingFilterer) WatchContributionMade(opts *bind.WatchOpts, sink chan<- *CrowdFundingContributionMade, campaignId []*big.Int, contributor []common.Address) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var contributorRule []interface{}
	for _, contributorItem := range contributor {
		contributorRule = append(contributorRule, contributorItem)
	}

	logs, sub, err := _CrowdFunding.contract.WatchLogs(opts, "ContributionMade", campaignIdRule, contributorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdFundingContributionMade)
				if err := _CrowdFunding.contract.UnpackLog(event, "ContributionMade", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseContributionMade is a log parse operation binding the contract event 0xf6ab15f3ae6cc7f857f32b5a5e1a789bf16059eb626d7aefab39f2ff0447e0df.
//
// Solidity: event ContributionMade(uint256 indexed campaignId, address indexed contributor, address token, uint256 amount)
func (_CrowdFunding *CrowdFundingFilterer) ParseContributionMade(log types.Log) (*CrowdFundingContributionMade, error) {
	event := new(CrowdFundingContributionMade)
	if err := _CrowdFunding.contract.UnpackLog(event, "ContributionMade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrowdFundingContributionRefundedIterator is returned from FilterContributionRefunded and is used to iterate over the raw logs and unpacked data for ContributionRefunded events raised by the CrowdFunding contract.
type CrowdFundingContributionRefundedIterator struct {
	Event *CrowdFundingContributionRefunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CrowdFundingContributionRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdFundingContributionRefunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CrowdFundingContributionRefunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CrowdFundingContributionRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdFundingContributionRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdFundingContributionRefunded represents a ContributionRefunded event raised by the CrowdFunding contract.
type CrowdFundingContributionRefunded struct {
	CampaignId  *big.Int
	Contributor common.Address
	EthAmount   *big.Int
	TokenAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterContributionRefunded is a free log retrieval operation binding the contract event 0x94b7d16871fd252f485b98072f622be6a999f952a4cc258fcb4f47942916bed5.
//
// Solidity: event ContributionRefunded(uint256 indexed campaignId, address indexed contributor, uint256 ethAmount, uint256 tokenAmount)
func (_CrowdFunding *CrowdFundingFilterer) FilterContributionRefunded(opts *bind.FilterOpts, campaignId []*big.Int, contributor []common.Address) (*CrowdFundingContributionRefundedIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var contributorRule []interface{}
	for _, contributorItem := range contributor {
		contributorRule = append(contributorRule, contributorItem)
	}

	logs, sub, err := _CrowdFunding.contract.FilterLogs(opts, "ContributionRefunded", campaignIdRule, contributorRule)
	if err != nil {
		return nil, err
	}
	return &CrowdFundingContributionRefundedIterator{contract: _CrowdFunding.contract, event: "ContributionRefunded", logs: logs, sub: sub}, nil
}

// WatchContributionRefunded is a free log subscription operation binding the contract event 0x94b7d16871fd252f485b98072f622be6a999f952a4cc258fcb4f47942916bed5.
//
// Solidity: event ContributionRefunded(uint256 indexed campaignId, address indexed contributor, uint256 ethAmount, uint256 tokenAmount)
func (_CrowdFunding *CrowdFundingFilterer) WatchContributionRefunded(opts *bind.WatchOpts, sink chan<- *CrowdFundingContributionRefunded, campaignId []*big.Int, contributor []common.Address) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var contributorRule []interface{}
	for _, contributorItem := range contributor {
		contributorRule = append(contributorRule, contributorItem)
	}

	logs, sub, err := _CrowdFunding.contract.WatchLogs(opts, "ContributionRefunded", campaignIdRule, contributorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdFundingContributionRefunded)
				if err := _CrowdFunding.contract.UnpackLog(event, "ContributionRefunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseContributionRefunded is a log parse operation binding the contract event 0x94b7d16871fd252f485b98072f622be6a999f952a4cc258fcb4f47942916bed5.
//
// Solidity: event ContributionRefunded(uint256 indexed campaignId, address indexed contributor, uint256 ethAmount, uint256 tokenAmount)
func (_CrowdFunding *CrowdFundingFilterer) ParseContributionRefunded(log types.Log) (*CrowdFundingContributionRefunded, error) {
	event := new(CrowdFundingContributionRefunded)
	if err := _CrowdFunding.contract.UnpackLog(event, "ContributionRefunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CrowdFundingFundsWithdrawnIterator is returned from FilterFundsWithdrawn and is used to iterate over the raw logs and unpacked data for FundsWithdrawn events raised by the CrowdFunding contract.
type CrowdFundingFundsWithdrawnIterator struct {
	Event *CrowdFundingFundsWithdrawn // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *CrowdFundingFundsWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CrowdFundingFundsWithdrawn)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(CrowdFundingFundsWithdrawn)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *CrowdFundingFundsWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CrowdFundingFundsWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CrowdFundingFundsWithdrawn represents a FundsWithdrawn event raised by the CrowdFunding contract.
type CrowdFundingFundsWithdrawn struct {
	CampaignId  *big.Int
	Owner       common.Address
	EthAmount   *big.Int
	TokenAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterFundsWithdrawn is a free log retrieval operation binding the contract event 0x3dbb1a26f29aef1264c12ca1285deb60cca5666d8b17100c7ea1602c0a4014f6.
//
// Solidity: event FundsWithdrawn(uint256 indexed campaignId, address indexed owner, uint256 ethAmount, uint256 tokenAmount)
func (_CrowdFunding *CrowdFundingFilterer) FilterFundsWithdrawn(opts *bind.FilterOpts, campaignId []*big.Int, owner []common.Address) (*CrowdFundingFundsWithdrawnIterator, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _CrowdFunding.contract.FilterLogs(opts, "FundsWithdrawn", campaignIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return &CrowdFundingFundsWithdrawnIterator{contract: _CrowdFunding.contract, event: "FundsWithdrawn", logs: logs, sub: sub}, nil
}

// WatchFundsWithdrawn is a free log subscription operation binding the contract event 0x3dbb1a26f29aef1264c12ca1285deb60cca5666d8b17100c7ea1602c0a4014f6.
//
// Solidity: event FundsWithdrawn(uint256 indexed campaignId, address indexed owner, uint256 ethAmount, uint256 tokenAmount)
func (_CrowdFunding *CrowdFundingFilterer) WatchFundsWithdrawn(opts *bind.WatchOpts, sink chan<- *CrowdFundingFundsWithdrawn, campaignId []*big.Int, owner []common.Address) (event.Subscription, error) {

	var campaignIdRule []interface{}
	for _, campaignIdItem := range campaignId {
		campaignIdRule = append(campaignIdRule, campaignIdItem)
	}
	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}

	logs, sub, err := _CrowdFunding.contract.WatchLogs(opts, "FundsWithdrawn", campaignIdRule, ownerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CrowdFundingFundsWithdrawn)
				if err := _CrowdFunding.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFundsWithdrawn is a log parse operation binding the contract event 0x3dbb1a26f29aef1264c12ca1285deb60cca5666d8b17100c7ea1602c0a4014f6.
//
// Solidity: event FundsWithdrawn(uint256 indexed campaignId, address indexed owner, uint256 ethAmount, uint256 tokenAmount)
func (_CrowdFunding *CrowdFundingFilterer) ParseFundsWithdrawn(log types.Log) (*CrowdFundingFundsWithdrawn, error) {
	event := new(CrowdFundingFundsWithdrawn)
	if err := _CrowdFunding.contract.UnpackLog(event, "FundsWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
