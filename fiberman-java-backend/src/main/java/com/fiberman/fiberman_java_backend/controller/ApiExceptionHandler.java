package com.fiberman.fiberman_java_backend.controller;

import com.fiberman.sdk.exception.FiberHttpException;
import com.fiberman.sdk.exception.FiberRpcException;
import com.fiberman.sdk.exception.FiberSerializationException;
import com.fiberman.sdk.exception.FiberTimeoutException;
import com.fiberman.sdk.exception.FiberTransportException;
import java.util.Map;
import java.util.stream.Collectors;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@RestControllerAdvice
public class ApiExceptionHandler {

    @ExceptionHandler(MethodArgumentNotValidException.class)
    ResponseEntity<Map<String, Object>> handleValidation(MethodArgumentNotValidException exception) {
        Map<String, String> errors = exception.getBindingResult()
                .getFieldErrors()
                .stream()
                .collect(Collectors.toMap(
                        FieldError::getField,
                        fieldError -> fieldError.getDefaultMessage() == null ? "Invalid value" : fieldError.getDefaultMessage(),
                        (left, right) -> left));

        return ResponseEntity.badRequest().body(Map.of(
                "error", "validation_failed",
                "details", errors));
    }

    @ExceptionHandler(FiberRpcException.class)
    ResponseEntity<Map<String, Object>> handleRpc(FiberRpcException exception) {
        return ResponseEntity.badRequest().body(Map.of(
                "error", "fiber_rpc_error",
                "message", exception.getMessage(),
                "code", exception.getCode(),
                "data", exception.getData()));
    }

    @ExceptionHandler(FiberHttpException.class)
    ResponseEntity<Map<String, Object>> handleHttp(FiberHttpException exception) {
        return ResponseEntity.status(HttpStatus.BAD_GATEWAY).body(Map.of(
                "error", "fiber_http_error",
                "message", exception.getMessage(),
                "status", exception.getStatusCode(),
                "responseBody", exception.getResponseBody()));
    }

    @ExceptionHandler(FiberTimeoutException.class)
    ResponseEntity<Map<String, Object>> handleTimeout(FiberTimeoutException exception) {
        return ResponseEntity.status(HttpStatus.GATEWAY_TIMEOUT).body(Map.of(
                "error", "fiber_timeout",
                "message", exception.getMessage()));
    }

    @ExceptionHandler({FiberTransportException.class, FiberSerializationException.class})
    ResponseEntity<Map<String, Object>> handleFiberFailure(RuntimeException exception) {
        return ResponseEntity.status(HttpStatus.BAD_GATEWAY).body(Map.of(
                "error", "fiber_backend_error",
                "message", exception.getMessage()));
    }
}
