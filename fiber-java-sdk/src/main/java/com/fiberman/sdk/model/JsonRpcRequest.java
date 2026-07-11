package com.fiberman.sdk.model;

public record JsonRpcRequest(
        String jsonrpc,
        String id,
        String method,
        Object params
) {
    public static JsonRpcRequest of(String id, String method, Object params) {
        return new JsonRpcRequest("2.0", id, method, params);
    }
}
