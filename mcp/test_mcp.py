#!/usr/bin/env python3
"""Test script for Hessian MCP server."""

import asyncio
import base64
import json
import sys
import os

# Add parent directory to path for imports
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from hessian_mcp import call_tool, PARSE_TOOL_PATH


async def test_parse():
    # Read a test file
    test_file = os.path.join(os.path.dirname(__file__), "..", "examples", "Rome.ser")
    with open(test_file, "rb") as f:
        raw_data = f.read()

    base64_data = base64.b64encode(raw_data).decode()

    print(f"Parse tool path: {PARSE_TOOL_PATH}")
    print(f"Parse tool exists: {PARSE_TOOL_PATH.exists()}")
    print(f"Base64 data length: {len(base64_data)}")
    print()

    # Call the tool
    result = await call_tool("parse_hessian", {"data": base64_data})

    print("Result:")
    for content in result:
        print(content.text[:2000])  # Print first 2000 chars
        if len(content.text) > 2000:
            print(f"... (truncated, total {len(content.text)} chars)")


if __name__ == "__main__":
    asyncio.run(test_parse())
