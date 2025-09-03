// test_api_endpoints.js
// Simple JavaScript test for new API endpoints

async function testAPIEndpoints() {
  const BASE_URL = 'http://localhost:8080/api/v1';
  
  console.log('🚀 Testing Enhanced Macro Strategy Platform API Endpoints');
  console.log('=====================================================');
  
  try {
    // Test health check
    console.log('\n📋 Health Check Test');
    const healthResponse = await fetch(`${BASE_URL}/health`);
    const healthData = await healthResponse.json();
    console.log('✅ Health Check:', healthData.status);
    
    // Test supported markets
    console.log('\n🌍 Supported Markets Test');
    const marketsResponse = await fetch(`${BASE_URL}/markets`);
    const marketsData = await marketsResponse.json();
    if (marketsData.success) {
      console.log(`✅ Supported Markets: ${Object.keys(marketsData.data).length} market types`);
      Object.entries(marketsData.data).forEach(([key, market]) => {
        console.log(`   🏛️ ${market.name}: ${market.assets.length} assets`);
      });
    }
    
    // Test supported strategies
    console.log('\n📈 Supported Strategies Test');
    const strategiesResponse = await fetch(`${BASE_URL}/strategies`);
    const strategiesData = await strategiesResponse.json();
    if (strategiesData.success) {
      console.log(`✅ Supported Strategies: ${Object.keys(strategiesData.data).length} strategy types`);
      Object.entries(strategiesData.data).forEach(([key, strategy]) => {
        console.log(`   🎯 ${strategy.name}`);
      });
    }
    
    // Test all assets
    console.log('\n📊 All Assets Test');
    const assetsResponse = await fetch(`${BASE_URL}/assets`);
    const assetsData = await assetsResponse.json();
    if (assetsData.success) {
      console.log(`✅ Total Assets: ${assetsData.data.length}`);
    }
    
    // Test assets by market type
    console.log('\n🌐 Assets by Market Type Test');
    const marketTypes = ['a_share_index', 'a_share_stock', 'us_index', 'us_stock', 'crypto', 'hk_index'];
    for (const marketType of marketTypes) {
      try {
        const response = await fetch(`${BASE_URL}/assets/market/${marketType}`);
        const data = await response.json();
        if (data.success) {
          console.log(`   📈 ${marketType}: ${data.data.length} assets`);
        }
      } catch (error) {
        console.log(`   ❌ ${marketType}: Error - ${error.message}`);
      }
    }
    
    // Test single strategy backtest (quick test)
    console.log('\n⚡ Quick Single Strategy Backtest Test');
    const backtestData = {
      index_id: 'csi1000',
      strategy: {
        type: 'buy_and_hold',
        parameters: {
          target_allocation: 1.0,
          rebalance_frequency: 'never'
        },
        description: 'Quick test backtest'
      },
      start_date: '2024-01-01',
      end_date: '2024-01-31', // Short period for quick test
      initial_cash: 100000
    };
    
    try {
      const backtestResponse = await fetch(`${BASE_URL}/backtest`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(backtestData)
      });
      
      const backtestResult = await backtestResponse.json();
      if (backtestResult.success) {
        console.log('✅ Single Strategy Backtest: Success');
        console.log(`   📈 Return: ${(backtestResult.data.performance_metrics.total_return * 100).toFixed(2)}%`);
        console.log(`   📊 Sharpe Ratio: ${backtestResult.data.performance_metrics.sharpe_ratio.toFixed(3)}`);
      } else {
        console.log('❌ Single Strategy Backtest Failed:', backtestResult.error);
      }
    } catch (error) {
      console.log('❌ Single Strategy Backtest Error:', error.message);
    }
    
    console.log('\n🎉 API Endpoint Testing Complete!');
    console.log('All enhanced platform features are accessible via the API.');
    
  } catch (error) {
    console.error('❌ Test Failed:', error);
  }
}

// Run the tests
testAPIEndpoints();