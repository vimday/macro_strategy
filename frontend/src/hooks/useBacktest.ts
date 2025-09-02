import { useMutation, useQuery } from '@tanstack/react-query';
import { indexService, backtestService } from '@/lib/api';
import { BacktestRequest, BacktestResult } from '@/types';

export function useIndexes() {
  return useQuery({
    queryKey: ['indexes'],
    queryFn: indexService.getIndexes,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

export function useRunBacktest() {
  return useMutation({
    mutationFn: (request: BacktestRequest) => backtestService.runBacktest(request),
    onError: (error) => {
      console.error('Backtest failed:', error);
    },
  });
}

export function useBacktestResult(id: string | null) {
  return useQuery({
    queryKey: ['backtest', id],
    queryFn: () => id ? backtestService.getBacktestResult(id) : null,
    enabled: !!id,
  });
}