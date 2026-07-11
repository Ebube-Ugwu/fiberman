package com.fiberman.fiberman_java_backend.dto;

import jakarta.validation.constraints.NotBlank;

public record GetChannelApiRequest(@NotBlank String channelId) {
}
