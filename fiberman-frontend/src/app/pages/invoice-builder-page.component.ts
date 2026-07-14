import { CommonModule } from '@angular/common';
import { Component, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { FiberApiService } from '../core/fiber-api.service';
import { FiberCallResponse } from '../core/fiber-types';
import { PlaygroundSettingsService } from '../core/playground-settings.service';

@Component({
  selector: 'app-invoice-builder-page',
  imports: [CommonModule, FormsModule],
  templateUrl: './invoice-builder-page.component.html',
  styleUrl: './invoice-builder-page.component.css'
})
export class InvoiceBuilderPageComponent {
  private readonly api = inject(FiberApiService);
  private readonly settingsService = inject(PlaygroundSettingsService);

  protected description = 'Project Milestone Payment';
  protected amount = 25000;
  protected memo = 'Memo';
  protected expiry = '2026-07-17';
  protected metadata = '{\n  "type": "milestone"\n}';

  protected readonly response = signal<FiberCallResponse | null>(null);
  protected readonly loading = signal(false);
  protected readonly validationMessage = signal<string | null>(null);
  protected readonly defaultInvoiceCurrency = computed(() => this.settingsService.defaultInvoiceCurrency());

  constructor() {
    this.settingsService.load();
  }

  protected buildInvoice(): void {
    const currency = this.defaultInvoiceCurrency();
    if (!currency) {
      this.validationMessage.set('Set a default invoice currency in Settings before generating invoices.');
      return;
    }

    this.validationMessage.set(null);
    this.loading.set(true);
    this.api
      .createInvoice({
        amount: this.amount,
        currency,
        description: `${this.description} ${this.memo ? `| ${this.memo}` : ''}`.trim(),
        expirySeconds: this.expiry ? 3600 : null
      })
      .subscribe({
        next: (response) => {
          this.response.set(response);
          this.loading.set(false);
        },
        error: (error) => {
          this.response.set(error.error ?? null);
          this.loading.set(false);
        }
      });
  }

  protected invoiceValue(): string {
    const result = this.response()?.result as Record<string, unknown> | null;
    return (result?.['invoice_address'] as string) || (result?.['invoice'] as string) || 'Invoice will appear here';
  }

  protected copy(value: string): void {
    navigator.clipboard.writeText(value).catch(() => undefined);
  }
}
