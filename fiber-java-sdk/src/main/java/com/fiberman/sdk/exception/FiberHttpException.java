package com.fiberman.sdk.exception;

public final class FiberHttpException extends FiberException {
    private final int statusCode;
    private final String responseBody;

    public FiberHttpException(String message, int statusCode, String responseBody) {
        super(message);
        this.statusCode = statusCode;
        this.responseBody = responseBody;
    }

    public int getStatusCode() {
        return statusCode;
    }

    public String getResponseBody() {
        return responseBody;
    }
}
