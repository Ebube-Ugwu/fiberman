export interface CodeArtifacts {
  curl: string;
  javaSnippet: string;
  goSnippet: string;
}

export interface InvoiceQrCodeResponse {
  value: string;
  size: number;
  pngBase64: string;
  dataUrl: string;
}

export interface FiberCallResponse {
  historyId: string;
  method: string;
  backendPath: string;
  timestamp: string;
  success: boolean;
  params: Record<string, unknown>;
  result: unknown;
  error: Record<string, unknown> | null;
  codeArtifacts: CodeArtifacts | null;
  invoiceQrCode: InvoiceQrCodeResponse | null;
}

export interface CreateInvoicePayload {
  amount: number;
  currency: string;
  description?: string | null;
  expirySeconds?: number | null;
}

export interface SendPaymentPayload {
  invoice: string;
  amount?: number | null;
  timeoutSeconds?: number | null;
}

export interface GetChannelPayload {
  channelId: string;
}

export interface GetPaymentPayload {
  paymentId: string;
}

export interface InvoiceQrPayload {
  invoice: string;
  size?: number | null;
}

export interface PlaygroundSettings {
  nodeUrl: string;
  authToken: string;
  timeoutSeconds: number;
  defaultInvoiceCurrency: string;
  playgroundBaseUrl: string;
}

export interface UpdatePlaygroundSettingsPayload {
  nodeUrl: string;
  authToken?: string | null;
  timeoutSeconds: number;
  defaultInvoiceCurrency?: string | null;
}
