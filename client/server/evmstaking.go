package server

import (
	"github.com/cosmos/cosmos-sdk/types/query"
	"net/http"

	"github.com/piplabs/story/client/server/utils"
	evmstakingtypes "github.com/piplabs/story/client/x/evmstaking/types"
)

func (s *Server) initEvmStakingRoute() {
	s.httpMux.HandleFunc("/evmstaking/params", utils.SimpleWrap(s.aminoCodec, s.GetEvmStakingParams))
	s.httpMux.HandleFunc("/evmstaking/withdrawal-queue", utils.AutoWrap(s.aminoCodec, s.GetEvmStakingWithdrawalQueue))
}

// GetEvmStakingParams queries the parameters of evmstaking module.
func (s *Server) GetEvmStakingParams(r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetEvmStakingKeeper().Params(queryContext, &evmstakingtypes.QueryParamsRequest{})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}

// GetEvmStakingWithdrawalQueue queries the withdrawal queue of evmstaking module.
func (s *Server) GetEvmStakingWithdrawalQueue(req *getEvmStakingWithdrawalQueueRequest, r *http.Request) (resp any, err error) {
	queryContext, err := s.createQueryContextByHeader(r)
	if err != nil {
		return nil, err
	}

	queryResp, err := s.store.GetEvmStakingKeeper().GetWithdrawalQueue(queryContext, &evmstakingtypes.QueryGetWithdrawalQueueRequest{
		Pagination: &query.PageRequest{
			Key:        []byte(req.Pagination.Key),
			Offset:     req.Pagination.Offset,
			Limit:      req.Pagination.Limit,
			CountTotal: req.Pagination.CountTotal,
			Reverse:    req.Pagination.Reverse,
		}})
	if err != nil {
		return nil, err
	}

	return queryResp, nil
}
