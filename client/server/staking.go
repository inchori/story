//nolint:wrapcheck // The api server is our server, so we don't need to wrap it.
package server

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/gorilla/mux"

	"github.com/storyprotocol/iliad/client/server/utils"
)

func (s *Server) initStakingRoute() {
	s.httpMux.HandleFunc("/staking/pool", utils.SimpleWrap(s.aminoCodec, s.GetStakingPool))

	s.httpMux.HandleFunc("/staking/validators", utils.AutoWrap(s.aminoCodec, s.GetValidators))
	s.httpMux.HandleFunc("/staking/validators/{validator_addr}", utils.SimpleWrap(s.aminoCodec, s.GetValidatorByValidatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_addr}/delegations", utils.AutoWrap(s.aminoCodec, s.GetValidatorDelegationsByValidatorAddress))
	s.httpMux.HandleFunc("/staking/validators/{validator_addr}/delegations/{delegator_addr}", utils.SimpleWrap(s.aminoCodec, s.GetDelegationByValidatorAddressDelegatorAddress))

	s.httpMux.HandleFunc("/staking/delegations/{delegator_addr}", utils.AutoWrap(s.aminoCodec, s.GetDelegationsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_addr}/redelegations", utils.AutoWrap(s.aminoCodec, s.GetRedelegationsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_addr}/unbonding_delegations", utils.AutoWrap(s.aminoCodec, s.GetUnbondingDelegationsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_addr}/validators", utils.AutoWrap(s.aminoCodec, s.GetValidatorsByDelegatorAddress))
	s.httpMux.HandleFunc("/staking/delegators/{delegator_addr}/validators/{validator_addr}", utils.SimpleWrap(s.aminoCodec, s.GetValidatorsByDelegatorAddressValidatorAddress))
}

// GetStakingPool queries the staking pool info.
func (s *Server) GetStakingPool(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Pool(queryContext, &stakingtypes.QueryPoolRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetValidators queries all validators that match the given status.
func (s *Server) GetValidators(req *getValidatorsRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Validators(queryContext, &stakingtypes.QueryValidatorsRequest{
		Status: req.Status,
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})

	if err != nil {
		return nil, err
	}

	for _, validator := range queryResp.Validators {
		err = s.prepareUnpackInterfaces(validator)
		if err != nil {
			return nil, err
		}
	}

	return queryResp, nil
}

// GetValidatorByValidatorAddress queries validator info for given validator address.
func (s *Server) GetValidatorByValidatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Validator(queryContext, &stakingtypes.QueryValidatorRequest{
		ValidatorAddr: mux.Vars(r)["validator_addr"],
	})

	if err != nil {
		return nil, err
	}

	err = s.prepareUnpackInterfaces(queryResp.Validator)
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetValidatorDelegationsByValidatorAddress queries delegate info for given validator.
func (s *Server) GetValidatorDelegationsByValidatorAddress(req *getValidatorDelegationsByValidatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).ValidatorDelegations(queryContext, &stakingtypes.QueryValidatorDelegationsRequest{
		ValidatorAddr: mux.Vars(r)["validator_addr"],
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})

	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetDelegationByValidatorAddressDelegatorAddress queries delegate info for given validator delegator pair.
func (s *Server) GetDelegationByValidatorAddressDelegatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	muxVars := mux.Vars(r)
	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Delegation(queryContext, &stakingtypes.QueryDelegationRequest{
		ValidatorAddr: muxVars["validator_addr"],
		DelegatorAddr: muxVars["delegator_addr"],
	})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetDelegationsByDelegatorAddress queries all delegations of a given delegator address.
func (s *Server) GetDelegationsByDelegatorAddress(req *getDelegationsByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorDelegations(queryContext, &stakingtypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: mux.Vars(r)["delegator_addr"],
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})

	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetRedelegationsByDelegatorAddress queries redelegations of given address.
func (s *Server) GetRedelegationsByDelegatorAddress(req *getRedelegationsByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).Redelegations(queryContext, &stakingtypes.QueryRedelegationsRequest{
		DelegatorAddr:    mux.Vars(r)["delegator_addr"],
		SrcValidatorAddr: req.SrcValidatorAddr,
		DstValidatorAddr: req.DstValidatorAddr,
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})

	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetUnbondingDelegationsByDelegatorAddress queries all unbonding delegations of a given delegator address.
func (s *Server) GetUnbondingDelegationsByDelegatorAddress(req *getUnbondingDelegationsByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorUnbondingDelegations(queryContext, &stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
		DelegatorAddr: mux.Vars(r)["delegator_addr"],
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})

	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetValidatorsByDelegatorAddress queries all validators info for given delegator address.
func (s *Server) GetValidatorsByDelegatorAddress(req *getValidatorsByDelegatorAddressRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorValidators(queryContext, &stakingtypes.QueryDelegatorValidatorsRequest{
		DelegatorAddr: mux.Vars(r)["delegator_addr"],
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		},
	})

	if err != nil {
		return nil, err
	}

	for _, validator := range queryResp.Validators {
		err = s.prepareUnpackInterfaces(validator)
		if err != nil {
			return nil, err
		}
	}

	return queryResp, nil
}

// GetValidatorsByDelegatorAddressValidatorAddress queries validator info for given delegator validator pair.
func (s *Server) GetValidatorsByDelegatorAddressValidatorAddress(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	muxVars := mux.Vars(r)
	queryResp, err := keeper.NewQuerier(s.store.GetStakingKeeper()).DelegatorValidator(queryContext, &stakingtypes.QueryDelegatorValidatorRequest{
		DelegatorAddr: muxVars["delegator_addr"],
		ValidatorAddr: muxVars["validator_addr"],
	})

	if err != nil {
		return nil, err
	}

	err = s.prepareUnpackInterfaces(queryResp.Validator)
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
