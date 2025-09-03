#!/bin/bash

# Script to commit only the necessary files, excluding virtual environments

echo "Adding and committing necessary files to GitHub..."

# Add documentation files
git add WORK_LOG_2025-09-04.md
git add FINAL_SUMMARY.md
git add CHANGES_SUMMARY.md
git add README.md

# Add helper scripts
git add setup_akshare.sh
git add git_commit_helper.sh
git add verify_fixes.py

# Add modified source files
git add backend/scripts/akshare_client.py
git add backend/internal/data/provider.go

# Commit changes
git commit -m "✨ Fix backtest 500 error and AKShare integration

- 🐛 Fixed AxiosError 'Request failed with status code 500' when running backtests
- 🔧 Added missing get_stock_zh_index_daily function to AKShare client script
- 🔧 Updated provider configuration with correct paths for Python virtual environment
- 📝 Created work log documenting the fixes and solution
- 🎯 Verified backtest functionality now works correctly with performance metrics
- 📊 Fixed win rate calculation accuracy issues
- 📚 Updated README with recent improvements
- 🗂️ Updated .gitignore to exclude virtual environment directories"

echo "Changes committed successfully!"
echo "Note: Virtual environment directories (venv, akshare_env) are excluded per .gitignore"