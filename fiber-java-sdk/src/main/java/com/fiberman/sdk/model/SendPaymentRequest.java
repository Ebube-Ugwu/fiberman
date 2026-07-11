package com.fiberman.sdk.model;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonInclude(JsonInclude.Include.NON_NULL)
public record SendPaymentRequest(
        String invoice,
        Long amount,
        @JsonProperty("timeout") Long timeoutSeconds
) {
}
