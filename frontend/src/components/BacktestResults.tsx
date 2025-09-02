'use client';

import React from 'react';
import { Card, Divider, Typography } from 'antd';
import { TrendingDownIcon } from 'lucide-react';
import { MetricsDisplay } from './MetricsDisplay';
import { PerformanceChart } from './PerformanceChart';
import { BacktestResult } from '@/types';

const { Title } = Typography;

interface BacktestResultsProps {
  result: BacktestResult;
}

export function BacktestResults({ result }: BacktestResultsProps) {
  return (
    <div className="mt-8">
      <Divider>
        <Title level={3} className="flex items-center">
          <TrendingDownIcon className="mr-2 h-6 w-6" />
          回测结果
        </Title>
      </Divider>
      
      <MetricsDisplay result={result} />
      <PerformanceChart result={result} />
    </div>
  );
}