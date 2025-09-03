#!/usr/bin/env python3
"""
Debug script for analyzing backtest results and identifying accuracy issues
"""

import json
import requests
import sys
from datetime import datetime

def analyze_backtest():
    """Run a detailed analysis of backtest results"""
    
    # Simple backtest request
    request_data = {
        "index_id": "csi300",
        "strategy": {
            "type": "monthly_rotation",
            "parameters": {
                "buy_days_before_month_end": 1,
                "sell_days_after_month_start": 1
            }
        },
        "start_date": "2024-01-01",
        "end_date": "2024-03-31",
        "initial_cash": 1000000
    }
    
    print("=== BACKTEST DEBUG ANALYSIS ===")
    print(f"Request: {json.dumps(request_data, indent=2)}")
    print()
    
    try:
        response = requests.post(
            "http://localhost:8080/api/v1/backtest",
            json=request_data,
            timeout=30
        )
        
        if response.status_code != 200:
            print(f"HTTP Error: {response.status_code}")
            print(f"Response: {response.text}")
            return
            
        result = response.json()
        
        if not result.get("success"):
            print(f"Backtest Error: {result.get('error', 'Unknown error')}")
            return
            
        data = result["data"]
        
        # Analyze basic info
        print("=== BASIC INFO ===")
        print(f"Backtest ID: {data['id']}")
        print(f"Duration: {data['duration']}")
        print(f"Index: {data['request']['index_id']}")
        print(f"Initial Cash: ${data['request']['initial_cash']:,.0f}")
        print()
        
        # Analyze trades
        trades = data["trades"]
        print(f"=== TRADES ANALYSIS ({len(trades)} total) ===")
        
        total_buy_amount = 0
        total_sell_amount = 0
        total_commission = 0
        
        for i, trade in enumerate(trades):
            action = trade["action"]
            date = trade["date"]
            price = trade["price"]
            quantity = trade["quantity"]
            amount = trade["amount"]
            commission = trade["commission"]
            
            total_commission += commission
            
            if action == "buy":
                total_buy_amount += amount
                print(f"Trade {i+1}: {date} BUY {quantity:.0f} @ ${price:.2f} = ${amount:,.2f} (Commission: ${commission:.2f})")
            else:
                total_sell_amount += amount
                print(f"Trade {i+1}: {date} SELL {quantity:.0f} @ ${price:.2f} = ${amount:,.2f} (Commission: ${commission:.2f})")
                
        print(f"Total Buy Amount: ${total_buy_amount:,.2f}")
        print(f"Total Sell Amount: ${total_sell_amount:,.2f}")
        print(f"Total Commission: ${total_commission:,.2f}")
        print(f"Trading P&L (before commission): ${total_sell_amount - total_buy_amount:,.2f}")
        print(f"Net Trading P&L: ${total_sell_amount - total_buy_amount - total_commission:,.2f}")
        print()
        
        # Analyze daily returns
        daily_returns = data["daily_returns"]
        print(f"=== DAILY RETURNS ANALYSIS ({len(daily_returns)} days) ===")
        
        first_day = daily_returns[0]
        last_day = daily_returns[-1]
        
        print(f"First Day: {first_day['date']}")
        print(f"  Portfolio Value: ${first_day['portfolio_value']:,.2f}")
        print(f"  Cash: ${first_day['cash']:,.2f}")
        print(f"  Position Value: ${first_day['position']['market_value']:,.2f}")
        print(f"  Daily Return: {first_day['daily_return']:.4%}")
        print()
        
        print(f"Last Day: {last_day['date']}")
        print(f"  Portfolio Value: ${last_day['portfolio_value']:,.2f}")
        print(f"  Cash: ${last_day['cash']:,.2f}")
        print(f"  Position Value: ${last_day['position']['market_value']:,.2f}")
        print(f"  Daily Return: {last_day['daily_return']:.4%}")
        print(f"  Cumulative Return: {last_day['cumulative_return']:.4%}")
        print()
        
        # Check for potential issues
        print("=== POTENTIAL ISSUES ANALYSIS ===")
        
        # Issue 1: Check if first day daily return is calculated correctly
        if first_day['daily_return'] != 0:
            print(f"❌ ISSUE: First day daily return should be 0, but got {first_day['daily_return']:.4%}")
        else:
            print("✅ First day daily return is correctly 0")
            
        # Issue 2: Portfolio value consistency
        calculated_total = last_day['cash'] + last_day['position']['market_value']
        if abs(calculated_total - last_day['portfolio_value']) > 0.01:
            print(f"❌ ISSUE: Portfolio value mismatch on last day:")
            print(f"   Reported: ${last_day['portfolio_value']:,.2f}")
            print(f"   Calculated (Cash + Position): ${calculated_total:,.2f}")
            print(f"   Difference: ${abs(calculated_total - last_day['portfolio_value']):,.2f}")
        else:
            print("✅ Portfolio value calculation is consistent")
            
        # Issue 3: Total return calculation
        metrics = data["performance_metrics"]
        manual_total_return = (last_day['portfolio_value'] - first_day['portfolio_value']) / first_day['portfolio_value']
        
        if abs(manual_total_return - metrics['total_return']) > 0.0001:
            print(f"❌ ISSUE: Total return calculation mismatch:")
            print(f"   Reported: {metrics['total_return']:.4%}")
            print(f"   Manual calculation: {manual_total_return:.4%}")
            print(f"   Difference: {abs(manual_total_return - metrics['total_return']):.6%}")
        else:
            print("✅ Total return calculation is correct")
            
        # Issue 4: Trade count verification
        buy_trades = len([t for t in trades if t['action'] == 'buy'])
        sell_trades = len([t for t in trades if t['action'] == 'sell'])
        
        print(f"Buy trades: {buy_trades}, Sell trades: {sell_trades}")
        if buy_trades != sell_trades:
            print(f"❌ ISSUE: Unmatched trades - {buy_trades} buys vs {sell_trades} sells")
        else:
            print("✅ Trade pairs are balanced")
            
        # Performance metrics summary
        print()
        print("=== PERFORMANCE METRICS ===")
        print(f"Total Return: {metrics['total_return']:.4%}")
        print(f"Annualized Return: {metrics['annualized_return']:.4%}")
        print(f"Max Drawdown: {metrics['max_drawdown']:.4%}")
        print(f"Sharpe Ratio: {metrics['sharpe_ratio']:.3f}")
        print(f"Win Rate: {metrics['win_rate']:.2%}")
        print(f"Total Trades: {metrics['total_trades']}")
        
    except requests.exceptions.RequestException as e:
        print(f"Request failed: {e}")
    except Exception as e:
        print(f"Analysis failed: {e}")

if __name__ == "__main__":
    analyze_backtest()