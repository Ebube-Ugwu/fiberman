package com.fiberman.sdk.model;

import com.fasterxml.jackson.databind.JsonNode;

public record JsonRpcError(
        int code,
        String message,
        JsonNode data
) {
}
