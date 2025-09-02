import axios, { AxiosResponse } from 'axios';
import { 
  Index, 
  BacktestRequest, 
  BacktestResult, 
  MarketType,
  ApiResponse 
} from '@/types';

// Create axios instance with default config
const apiClient = axios.create({
  baseURL: 'http://localhost:8080',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Response interceptor for handling common response structure
apiClient.interceptors.response.use(
  (response: AxiosResponse<ApiResponse<any>>) => {
    if (response.data.success) {
      return response;
    } else {
      throw new Error(response.data.message || 'API request failed');
    }
  },
  (error) => {
    console.error('API Error:', error);
    throw error;
  }
);

// Index service
export const indexService = {
  async getIndexes(): Promise<Index[]> {
    const response = await apiClient.get<ApiResponse<Index[]>>('/api/v1/indexes');
    return response.data.data;
  },

  async getIndexesByMarketType(marketType: MarketType): Promise<Index[]> {
    const response = await apiClient.get<ApiResponse<Index[]>>(`/api/v1/indexes/market/${marketType}`);
    return response.data.data;
  },

  async getIndexById(id: string): Promise<Index | null> {
    try {
      const indexes = await this.getIndexes();
      return indexes.find(index => index.id === id) || null;
    } catch (error) {
      console.error('Error fetching index:', error);
      return null;
    }
  }
};

// Backtest service
export const backtestService = {
  async runBacktest(request: BacktestRequest): Promise<BacktestResult> {
    const response = await apiClient.post<ApiResponse<BacktestResult>>('/api/v1/backtest', request);
    return response.data.data;
  },

  async getBacktestResult(id: string): Promise<BacktestResult | null> {
    try {
      const response = await apiClient.get<ApiResponse<BacktestResult>>(`/api/v1/backtest/${id}`);
      return response.data.data;
    } catch (error) {
      console.error('Error fetching backtest result:', error);
      return null;
    }
  }
};

// Health service
export const healthService = {
  async checkHealth(): Promise<{ status: string; timestamp: string }> {
    const response = await apiClient.get('/api/v1/health');
    return response.data;
  }
};