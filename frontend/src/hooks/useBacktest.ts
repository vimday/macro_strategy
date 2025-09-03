'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { indexService, backtestService } from '@/lib/api';
import { Index, BacktestRequest, BacktestResult, MarketType } from '@/types';

/**
 * Hook to fetch all available indexes
 */
export function useIndexes() {
  return useQuery({
    queryKey: ['indexes'],
    queryFn: indexService.getIndexes,
    staleTime: 5 * 60 * 1000, // 5 minutes
    cacheTime: 10 * 60 * 1000, // 10 minutes
  });
}

/**
 * Hook to fetch indexes by market type
 */
export function useIndexesByMarketType(marketType: MarketType) {
  return useQuery({
    queryKey: ['indexes', 'market', marketType],
    queryFn: () => indexService.getIndexesByMarketType(marketType),
    enabled: !!marketType,
    staleTime: 5 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
  });
}

/**
 * Hook to fetch a specific index by ID
 */
export function useIndex(id: string) {
  return useQuery({
    queryKey: ['index', id],
    queryFn: () => indexService.getIndexById(id),
    enabled: !!id,
    staleTime: 5 * 60 * 1000,
    cacheTime: 10 * 60 * 1000,
  });
}

/**
 * Hook to run a backtest
 */
export function useRunBacktest() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (request: BacktestRequest) => backtestService.runBacktest(request),
    onSuccess: (data: BacktestResult) => {
      // Cache the result for potential future retrieval
      queryClient.setQueryData(['backtest', data.id], data);
      
      // Invalidate any related queries if needed
      queryClient.invalidateQueries({ queryKey: ['backtests'] });
    },
    onError: (error) => {
      console.error('Backtest failed:', error);
    },
  });
}

/**
 * Hook to fetch a backtest result by ID
 */
export function useBacktestResult(id: string) {
  return useQuery({
    queryKey: ['backtest', id],
    queryFn: () => backtestService.getBacktestResult(id),
    enabled: !!id,
    staleTime: Infinity, // Backtest results don't change
    cacheTime: 30 * 60 * 1000, // 30 minutes
  });
}

/**
 * Hook to manage backtest state and operations
 */
export function useBacktestManager() {
  const queryClient = useQueryClient();
  const runBacktestMutation = useRunBacktest();

  const clearCache = () => {
    queryClient.invalidateQueries({ queryKey: ['backtests'] });
  };

  const getCachedResult = (id: string): BacktestResult | undefined => {
    return queryClient.getQueryData(['backtest', id]);
  };

  return {
    runBacktest: runBacktestMutation.mutate,
    runBacktestAsync: runBacktestMutation.mutateAsync,
    isRunning: runBacktestMutation.isPending,
    error: runBacktestMutation.error,
    clearCache,
    getCachedResult,
  };
}

/**
 * Hook for managing multiple backtest comparisons
 */
export function useBacktestComparison() {
  const queryClient = useQueryClient();

  const getMultipleResults = (ids: string[]): (BacktestResult | undefined)[] => {
    return ids.map(id => queryClient.getQueryData(['backtest', id]));
  };

  const invalidateResults = (ids: string[]) => {
    ids.forEach(id => {
      queryClient.invalidateQueries({ queryKey: ['backtest', id] });
    });
  };

  return {
    getMultipleResults,
    invalidateResults,
  };
}