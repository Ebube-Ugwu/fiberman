package com.fiberman.sdk.exception;

import com.fasterxml.jackson.databind.JsonNode;

public final class FiberRpcException extends FiberException {
    private final int code;
    private final JsonNode data;

    public FiberRpcException(String message, int code, JsonNode data) {
        super(message);
        this.code = code;
        this.data = data;
    }

    public int getCode() {
        return code;
    }

    public JsonNode getData() {
        return data;
    }
}
