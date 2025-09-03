#!/usr/bin/env python3
"""
Test script to verify that the macro strategy platform is working correctly.
"""

import requests
import time

def test_backend():
    """Test backend API endpoints"""
    print("Testing backend API...")
    
    # Test health endpoint
    try:
        response = requests.get("http://localhost:8080/api/v1/health")
        if response.status_code == 200:
            print("✅ Backend health check: PASSED")
        else:
            print(f"❌ Backend health check: FAILED (Status: {response.status_code})")
    except Exception as e:
        print(f"❌ Backend health check: FAILED ({e})")
    
    # Test assets endpoint
    try:
        response = requests.get("http://localhost:8080/api/v1/assets")
        if response.status_code == 200:
            print("✅ Backend assets endpoint: PASSED")
        else:
            print(f"❌ Backend assets endpoint: FAILED (Status: {response.status_code})")
    except Exception as e:
        print(f"❌ Backend assets endpoint: FAILED ({e})")

def test_frontend():
    """Test frontend availability"""
    print("\nTesting frontend...")
    
    try:
        response = requests.get("http://localhost:3001")
        if response.status_code == 200:
            print("✅ Frontend availability: PASSED")
        else:
            print(f"❌ Frontend availability: FAILED (Status: {response.status_code})")
    except Exception as e:
        print(f"❌ Frontend availability: FAILED ({e})")

def main():
    print("_MACRO STRATEGY PLATFORM TEST_")
    print("=" * 40)
    
    test_backend()
    test_frontend()
    
    print("\n" + "=" * 40)
    print("Test completed!")
    print("\nSystem status:")
    print("- Backend:  http://localhost:8080")
    print("- Frontend: http://localhost:3001")

if __name__ == "__main__":
    main()