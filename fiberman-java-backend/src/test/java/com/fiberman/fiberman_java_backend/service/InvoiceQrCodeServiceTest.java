package com.fiberman.fiberman_java_backend.service;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class InvoiceQrCodeServiceTest {
    private final InvoiceQrCodeService invoiceQrCodeService = new InvoiceQrCodeService();

    @Test
    void generatesPngDataUrl() {
        var response = invoiceQrCodeService.generate("fiber1demo_invoice", 256);

        assertEquals("fiber1demo_invoice", response.value());
        assertEquals(256, response.size());
        assertNotNull(response.pngBase64());
        assertTrue(response.pngBase64().length() > 100);
        assertTrue(response.dataUrl().startsWith("data:image/png;base64,"));
    }
}
