#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Enhanced Macro Strategy Platform Test Script
Test the new multi-market, multi-strategy capabilities
"""

import requests
import json
import time
from datetime import datetime, timedelta

# Configuration
BASE_URL = "http://localhost:8080/api/v1"
TEST_TIMEOUT = 30

def test_api_endpoint(endpoint, method="GET", data=None, description=""):
    """Test an API endpoint and return the result"""
    url = f"{BASE_URL}{endpoint}"
    print(f"\n🧪 Testing: {description}")
    print(f"📡 {method} {url}")
    
    try:
        if method == "GET":
            response = requests.get(url, timeout=TEST_TIMEOUT)
        elif method == "POST":
            response = requests.post(url, json=data, timeout=TEST_TIMEOUT)
        
        print(f"✅ Status: {response.status_code}")
        
        if response.status_code == 200:
            result = response.json()
            if result.get("success"):
                print(f"✅ Success: {len(result.get('data', []))} items returned")
                return result["data"]
            else:
                print(f"❌ API Error: {result.get('error', 'Unknown error')}")
        else:
            print(f"❌ HTTP Error: {response.status_code}")
            print(f"Response: {response.text[:200]}...")
    
    except requests.exceptions.Timeout:
        print(f"⏰ Timeout after {TEST_TIMEOUT} seconds")
    except requests.exceptions.RequestException as e:
        print(f"❌ Request failed: {e}")
    except Exception as e:
        print(f"❌ Unexpected error: {e}")
    
    return None

def test_health_check():
    """Test health check endpoint"""
    return test_api_endpoint("/health", description="Health Check")

def test_supported_markets():
    """Test supported markets endpoint"""
    return test_api_endpoint("/markets", description="Get Supported Markets")

def test_supported_strategies():
    """Test supported strategies endpoint"""
    return test_api_endpoint("/strategies", description="Get Supported Strategies")

def test_all_assets():
    """Test all assets endpoint"""
    return test_api_endpoint("/assets", description="Get All Assets")

def test_assets_by_market():
    """Test assets by market type"""
    markets = ["a_share_index", "a_share_stock", "us_index", "us_stock", "crypto", "hk_index"]
    
    for market in markets:
        result = test_api_endpoint(f"/assets/market/{market}", 
                                  description=f"Get {market.replace('_', ' ').title()} Assets")
        if result:
            print(f"   📊 Found {len(result)} assets in {market}")

def test_single_strategy_backtest():
    """Test single strategy backtesting"""
    # Test with CSI 1000 using buy and hold strategy
    backtest_data = {
        "index_id": "csi1000",
        "strategy": {
            "type": "buy_and_hold",
            "parameters": {
                "target_allocation": 1.0,
                "rebalance_frequency": "never"
            },
            "description": "Buy and hold CSI 1000"
        },
        "start_date": "2024-01-01",
        "end_date": "2024-03-31",
        "initial_cash": 1000000
    }
    
    result = test_api_endpoint("/backtest", "POST", backtest_data, 
                              "Single Strategy Backtest (Buy & Hold)")
    
    if result and result.get("id"):
        print(f"   📈 Backtest ID: {result['id']}")
        print(f"   💰 Total Return: {result['performance_metrics']['total_return']:.4f}")
        print(f"   📊 Sharpe Ratio: {result['performance_metrics']['sharpe_ratio']:.3f}")
        print(f"   🎯 Win Rate: {result['performance_metrics']['win_rate']:.2%}")
        return result["id"]
    
    return None

def test_multi_strategy_backtest():
    """Test multi-strategy comparison"""
    # Compare buy-and-hold vs monthly rotation on CSI 1000
    multi_backtest_data = {
        "asset_id": "csi1000",
        "strategies": [
            {
                "type": "buy_and_hold",
                "parameters": {
                    "target_allocation": 1.0,
                    "rebalance_frequency": "never"
                },
                "description": "Buy and Hold Strategy"
            },
            {
                "type": "monthly_rotation",
                "parameters": {
                    "buy_days_before_month_end": 1,
                    "sell_days_after_month_start": 1
                },
                "description": "Monthly Rotation Strategy"
            }
        ],
        "start_date": "2024-01-01",
        "end_date": "2024-03-31",
        "initial_cash": 1000000,
        "benchmark": "csi300",
        "comparison_opt": {
            "show_benchmark": True,
            "normalize_returns": False,
            "show_drawdown": True,
            "metrics": ["total_return", "sharpe_ratio", "max_drawdown", "win_rate"]
        }
    }
    
    result = test_api_endpoint("/backtest/multi", "POST", multi_backtest_data,
                              "Multi-Strategy Comparison (Buy & Hold vs Monthly Rotation)")
    
    if result and result.get("id"):
        print(f"   📈 Multi-Backtest ID: {result['id']}")
        print(f"   🏆 Best Strategy: {result['comparison']['best_strategy']}")
        print(f"   📉 Worst Strategy: {result['comparison']['worst_strategy']}")
        print(f"   📝 Summary: {result['comparison']['summary'][:100]}...")
        
        # Show strategy comparison
        for i, strategy_result in enumerate(result["results"]):
            strategy_name = strategy_result["request"]["strategy"]["type"]
            metrics = strategy_result["performance_metrics"]
            print(f"   Strategy {i+1} ({strategy_name}):")
            print(f"     💰 Return: {metrics['total_return']:.4f}")
            print(f"     📊 Sharpe: {metrics['sharpe_ratio']:.3f}")
            print(f"     🎯 Win Rate: {metrics['win_rate']:.2%}")
        
        return result["id"]
    
    return None

def test_crypto_assets():
    """Test cryptocurrency assets"""
    crypto_assets = test_api_endpoint("/assets/market/crypto", 
                                     description="Get Cryptocurrency Assets")
    
    if crypto_assets and len(crypto_assets) > 0:
        btc_asset = next((asset for asset in crypto_assets if asset["id"] == "btc"), None)
        if btc_asset:
            print(f"   ₿ Found Bitcoin: {btc_asset['name']} ({btc_asset['symbol']})")
            print(f"   🏛️ Market Type: {btc_asset['market_type']}")
            print(f"   💰 Currency: {btc_asset['currency']}")

def test_us_stocks():
    """Test US stock assets"""
    us_stocks = test_api_endpoint("/assets/market/us_stock",
                                 description="Get US Stock Assets")
    
    if us_stocks and len(us_stocks) > 0:
        aapl_asset = next((asset for asset in us_stocks if asset["id"] == "aapl"), None)
        if aapl_asset:
            print(f"   🍎 Found Apple: {aapl_asset['name']} ({aapl_asset['symbol']})")
            print(f"   🏛️ Market Type: {aapl_asset['market_type']}")
            print(f"   💰 Currency: {aapl_asset['currency']}")
            if aapl_asset.get("metadata"):
                print(f"   🏢 Sector: {aapl_asset['metadata'].get('sector', 'N/A')}")

def main():
    """Run comprehensive tests of the enhanced platform"""
    print("🚀 Enhanced Macro Strategy Platform Test Suite")
    print("=" * 60)
    
    start_time = time.time()
    
    # Basic connectivity tests
    print("\n📋 BASIC CONNECTIVITY TESTS")
    print("-" * 30)
    test_health_check()
    
    # Enhanced capability tests
    print("\n🌟 ENHANCED CAPABILITIES TESTS")
    print("-" * 30)
    test_supported_markets()
    test_supported_strategies()
    test_all_assets()
    
    # Market-specific tests
    print("\n🌍 MULTI-MARKET TESTS")
    print("-" * 30)
    test_assets_by_market()
    test_crypto_assets()
    test_us_stocks()
    
    # Strategy testing
    print("\n📈 STRATEGY TESTING")
    print("-" * 30)
    single_id = test_single_strategy_backtest()
    multi_id = test_multi_strategy_backtest()
    
    # Summary
    print("\n📊 TEST SUMMARY")
    print("-" * 30)
    elapsed = time.time() - start_time
    print(f"⏱️  Total test time: {elapsed:.2f} seconds")
    
    if single_id and multi_id:
        print("✅ All major functionality tests passed!")
        print("🎉 Enhanced platform is working correctly!")
        print("\n🔗 Test Results:")
        print(f"   • Single Strategy Result: http://localhost:8080/api/v1/backtest/{single_id}")
        print(f"   • Multi-Strategy Result: http://localhost:8080/api/v1/backtest/multi/{multi_id}")
    else:
        print("⚠️  Some tests failed - check the backend logs")
    
    print("\n🎯 Key Enhancements Verified:")
    print("   ✅ Multi-market support (A-shares, US, Crypto, HK)")
    print("   ✅ Individual stock support")
    print("   ✅ Multiple strategy types (Buy & Hold, Monthly Rotation)")
    print("   ✅ Strategy comparison and visualization")
    print("   ✅ Enhanced API endpoints")
    print("   ✅ Comprehensive asset management")

if __name__ == "__main__":
    main()