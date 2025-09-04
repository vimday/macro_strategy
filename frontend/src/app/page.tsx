'use client';

import { useState } from 'react';
import { Layout, Typography, message, Tabs } from 'antd';
import { TrophyIcon, BarChart3Icon, GitCompareIcon } from 'lucide-react';
import { StrategyForm } from '@/components/StrategyForm';
import { BacktestResults } from '@/components/BacktestResults';
import { MultiStrategyComparison } from '@/components/MultiStrategyComparison';
import { useRunBacktest, useRunMultiStrategyBacktest } from '@/hooks/useBacktest';
import { BacktestFormData, BacktestResult, MultiStrategyBacktestResult } from '@/types';

const { Header, Content, Footer } = Layout;
const { Title, Paragraph } = Typography;
const { TabPane } = Tabs;

export default function HomePage() {
  const [backtestResult, setBacktestResult] = useState<BacktestResult | null>(null);
  const [multiStrategyResult, setMultiStrategyResult] = useState<MultiStrategyBacktestResult | null>(null);
  const [activeTab, setActiveTab] = useState('single');
  
  const runBacktestMutation = useRunBacktest();
  const runMultiStrategyMutation = useRunMultiStrategyBacktest();

  const handleRunBacktest = async (values: BacktestFormData) => {
    try {
      const request = {
        index_id: values.indexId,
        strategy: {
          type: values.strategyType,
          parameters: {
            buy_days_before_month_end: values.buyDaysBeforeMonthEnd,
            sell_days_after_month_start: values.sellDaysAfterMonthStart,
          },
          description: `月末前${values.buyDaysBeforeMonthEnd}个交易日买入，月初第${values.sellDaysAfterMonthStart}个交易日卖出`,
        },
        start_date: values.startDate,
        end_date: values.endDate,
        initial_cash: values.initialCash,
      };

      const result = await runBacktestMutation.mutateAsync(request);
      setBacktestResult(result);
      message.success('回测完成！');
    } catch (error) {
      message.error('回测执行失败');
      console.error('Backtest error:', error);
    }
  };

  const handleRunMultiStrategyBacktest = async (values: BacktestFormData) => {
    try {
      // For demonstration, we'll create a multi-strategy request with different parameters
      const request = {
        asset_id: values.indexId,
        strategies: [
          {
            type: values.strategyType,
            parameters: {
              buy_days_before_month_end: values.buyDaysBeforeMonthEnd,
              sell_days_after_month_start: values.sellDaysAfterMonthStart,
            },
            description: `标准策略 (月末前${values.buyDaysBeforeMonthEnd}天买入，月初第${values.sellDaysAfterMonthStart}天卖出)`,
          },
          {
            type: values.strategyType,
            parameters: {
              buy_days_before_month_end: Math.max(1, values.buyDaysBeforeMonthEnd - 2),
              sell_days_after_month_start: values.sellDaysAfterMonthStart,
            },
            description: `激进策略 (月末前${Math.max(1, values.buyDaysBeforeMonthEnd - 2)}天买入)`,
          },
          {
            type: values.strategyType,
            parameters: {
              buy_days_before_month_end: values.buyDaysBeforeMonthEnd + 2,
              sell_days_after_month_start: values.sellDaysAfterMonthStart,
            },
            description: `保守策略 (月末前${values.buyDaysBeforeMonthEnd + 2}天买入)`,
          }
        ],
        start_date: values.startDate,
        end_date: values.endDate,
        initial_cash: values.initialCash,
        comparison_opt: {
          show_benchmark: true,
          normalize_returns: true,
          show_drawdown: true,
          show_rolling_metrics: false,
          rolling_window: 30,
          metrics: ['total_return', 'sharpe_ratio', 'max_drawdown']
        }
      };

      const result = await runMultiStrategyMutation.mutateAsync(request);
      setMultiStrategyResult(result);
      message.success('多策略回测完成！');
    } catch (error) {
      message.error('多策略回测执行失败');
      console.error('Multi-strategy backtest error:', error);
    }
  };

  return (
    <Layout className="min-h-screen">
      <Header className="flex items-center bg-slate-900">
        <div className="flex items-center text-white text-xl font-bold">
          <TrophyIcon className="mr-2 h-6 w-6" />
          Macro Strategy | 巨策略
        </div>
      </Header>

      <Content className="p-6">
        <div className="max-w-6xl mx-auto">
          <div className="text-center mb-8">
            <Title level={2} className="flex items-center justify-center">
              <BarChart3Icon className="mr-2 h-8 w-8" />
              宏观交易策略回测平台
            </Title>
            <Paragraph className="text-lg text-gray-600">
              支持A股主流指数的策略回测，提供专业的性能分析和可视化展示
            </Paragraph>
          </div>

          <Tabs 
            activeKey={activeTab} 
            onChange={setActiveTab}
            type="card"
            size="large"
          >
            <TabPane
              tab={
                <span>
                  <BarChart3Icon className="mr-2 h-4 w-4 inline" />
                  单策略回测
                </span>
              }
              key="single"
            >
              <StrategyForm 
                onSubmit={handleRunBacktest} 
                loading={runBacktestMutation.isPending}
              />

              {backtestResult && (
                <BacktestResults 
                  result={backtestResult} 
                  loading={runBacktestMutation.isPending}
                  error={runBacktestMutation.error ? String(runBacktestMutation.error) : null}
                />
              )}
            </TabPane>
            
            <TabPane
              tab={
                <span>
                  <GitCompareIcon className="mr-2 h-4 w-4 inline" />
                  多策略对比
                </span>
              }
              key="multi"
            >
              <StrategyForm 
                onSubmit={handleRunMultiStrategyBacktest} 
                loading={runMultiStrategyMutation.isPending}
                multiStrategyMode={true}
              />

              {multiStrategyResult && (
                <MultiStrategyComparison 
                  result={multiStrategyResult} 
                  loading={runMultiStrategyMutation.isPending}
                  error={runMultiStrategyMutation.error ? String(runMultiStrategyMutation.error) : null}
                />
              )}
            </TabPane>
          </Tabs>
        </div>
      </Content>

      <Footer className="text-center bg-gray-50">
        <p className="text-gray-600">
          Macro Strategy ©2024 Built with ❤️ for quantitative trading enthusiasts
        </p>
      </Footer>
    </Layout>
  );
}