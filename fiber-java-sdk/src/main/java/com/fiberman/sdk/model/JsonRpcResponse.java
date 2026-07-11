package com.fiberman.sdk.model;

import com.fasterxml.jackson.databind.JsonNode;

public record JsonRpcResponse(
        String jsonrpc,
        JsonNode id,
        JsonNode result,
        JsonRpcError error
) {
}
