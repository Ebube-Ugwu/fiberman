package com.fiberman.sdk.exception;

public class FiberException extends RuntimeException {
    public FiberException(String message) {
        super(message);
    }

    public FiberException(String message, Throwable cause) {
        super(message, cause);
    }
}
