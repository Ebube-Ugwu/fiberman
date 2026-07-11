package com.fiberman.sdk.exception;

public final class FiberTimeoutException extends FiberTransportException {
    public FiberTimeoutException(String message, Throwable cause) {
        super(message, cause);
    }
}
