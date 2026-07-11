package com.fiberman.fiberman_java_backend.service;

import com.fiberman.fiberman_java_backend.dto.InvoiceQrCodeResponse;
import com.google.zxing.BarcodeFormat;
import com.google.zxing.WriterException;
import com.google.zxing.client.j2se.MatrixToImageWriter;
import com.google.zxing.qrcode.QRCodeWriter;
import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.util.Base64;
import org.springframework.stereotype.Service;

@Service
public class InvoiceQrCodeService {
    private static final int DEFAULT_SIZE = 320;

    public InvoiceQrCodeResponse generate(String value, Integer requestedSize) {
        int size = requestedSize == null ? DEFAULT_SIZE : requestedSize;
        try {
            var matrix = new QRCodeWriter().encode(value, BarcodeFormat.QR_CODE, size, size);
            ByteArrayOutputStream outputStream = new ByteArrayOutputStream();
            MatrixToImageWriter.writeToStream(matrix, "PNG", outputStream);
            String pngBase64 = Base64.getEncoder().encodeToString(outputStream.toByteArray());
            return new InvoiceQrCodeResponse(
                    value,
                    size,
                    pngBase64,
                    "data:image/png;base64," + pngBase64);
        } catch (WriterException | IOException exception) {
            throw new IllegalStateException("Failed to generate invoice QR code", exception);
        }
    }
}
