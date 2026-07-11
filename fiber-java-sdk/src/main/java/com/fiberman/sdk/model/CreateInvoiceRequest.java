package com.fiberman.sdk.model;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.fasterxml.jackson.annotation.JsonInclude;

@JsonInclude(JsonInclude.Include.NON_NULL)
public record CreateInvoiceRequest(
        Long amount,
        String currency,
        String description,
        @JsonProperty("expiry") Long expirySeconds
) {
}
