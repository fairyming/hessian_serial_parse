#!/usr/bin/env python3
"""
MCP Server for Hessian Serialization Data Parsing

This MCP server provides a tool to parse Hessian serialized data.
Input: base64 encoded Hessian serialization data
Output: Parsed JSON result

Environment Variables:
    HS_SERIAL_PARSE_PATH: Path to hs_serial_parse executable
"""

import base64
import json
import os
import subprocess
import tempfile
from pathlib import Path

from mcp.server import Server
from mcp.server.stdio import stdio_server
from mcp.types import Tool, TextContent

# Get hs_serial_parse path from environment variable
# Default to looking in the parent directory of this mcp folder
DEFUALT_PARSE_TOOL = Path(__file__).parent.parent / "hs_serial_parse"
PARSE_TOOL_PATH = Path(os.environ.get("HS_SERIAL_PARSE_PATH", str(DEFUALT_PARSE_TOOL)))

app = Server("hessian-mcp")


@app.list_tools()
async def list_tools() -> list[Tool]:
    """List available tools."""
    return [
        Tool(
            name="parse_hessian",
            description="Parse Hessian serialized data from base64 encoded input. "
            "Returns the parsed result as JSON structure showing the deserialized "
            "object hierarchy including class names, maps, and lists.",
            inputSchema={
                "type": "object",
                "properties": {
                    "data": {
                        "type": "string",
                        "description": "Base64 encoded Hessian serialization data"
                    }
                },
                "required": ["data"]
            }
        )
    ]


@app.call_tool()
async def call_tool(name: str, arguments: dict) -> list[TextContent]:
    """Handle tool calls."""
    if name != "parse_hessian":
        return [TextContent(type="text", text=f"Unknown tool: {name}")]

    # Get base64 data
    base64_data = arguments.get("data")
    if not base64_data:
        return [TextContent(type="text", text="Error: 'data' parameter is required")]

    # Validate parse tool exists
    if not PARSE_TOOL_PATH.exists():
        return [TextContent(
            type="text",
            text=f"Error: hs_serial_parse not found at {PARSE_TOOL_PATH}. "
                 f"Please set HS_SERIAL_PARSE_PATH environment variable to the correct path."
        )]

    try:
        # Decode base64 data
        raw_data = base64.b64decode(base64_data)
    except Exception as e:
        return [TextContent(type="text", text=f"Error decoding base64 data: {str(e)}")]

    # Create temporary file for the raw data
    try:
        with tempfile.NamedTemporaryFile(delete=False, suffix=".ser") as tmp_file:
            tmp_file.write(raw_data)
            tmp_path = tmp_file.name

        try:
            # Call hs_serial_parse
            result = subprocess.run(
                [str(PARSE_TOOL_PATH), "--path", tmp_path],
                capture_output=True,
                text=True,
                timeout=30
            )

            if result.returncode != 0:
                error_msg = result.stderr or result.stdout or "Unknown error"
                return [TextContent(
                    type="text",
                    text=f"Error parsing Hessian data: {error_msg}"
                )]

            # Parse and re-format JSON for validation
            try:
                parsed = json.loads(result.stdout)
                formatted = json.dumps(parsed, indent=2, ensure_ascii=False)
                return [TextContent(type="text", text=formatted)]
            except json.JSONDecodeError:
                # Return raw output if not valid JSON
                return [TextContent(type="text", text=result.stdout)]

        finally:
            # Clean up temporary file
            os.unlink(tmp_path)

    except subprocess.TimeoutExpired:
        return [TextContent(type="text", text="Error: Parsing timeout (30s)")]
    except Exception as e:
        return [TextContent(type="text", text=f"Error: {str(e)}")]


async def run_server():
    """Run the MCP server."""
    async with stdio_server() as (read_stream, write_stream):
        await app.run(read_stream, write_stream, app.create_initialization_options())


def main():
    """Entry point for the MCP server."""
    import asyncio
    asyncio.run(run_server())


if __name__ == "__main__":
    main()
