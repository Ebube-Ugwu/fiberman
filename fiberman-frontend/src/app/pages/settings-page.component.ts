import { CommonModule } from '@angular/common';
import { Component, OnInit, computed, effect, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { PlaygroundSettingsService } from '../core/playground-settings.service';

@Component({
  selector: 'app-settings-page',
  imports: [CommonModule, FormsModule],
  templateUrl: './settings-page.component.html',
  styleUrl: './settings-page.component.css'
})
export class SettingsPageComponent implements OnInit {
  private readonly settingsService = inject(PlaygroundSettingsService);

  protected nodeUrl = '';
  protected authToken = '';
  protected timeoutSeconds = 30;
  protected defaultInvoiceCurrency = '';
  protected statusMessage = signal<string | null>(null);

  protected readonly settings = this.settingsService.settings;
  protected readonly loading = this.settingsService.loading;
  protected readonly saving = this.settingsService.saving;
  protected readonly error = computed(() => this.statusMessage() ?? this.settingsService.error());

  constructor() {
    effect(() => {
      const active = this.settings();
      this.nodeUrl = active.nodeUrl;
      this.authToken = active.authToken;
      this.timeoutSeconds = active.timeoutSeconds;
      this.defaultInvoiceCurrency = active.defaultInvoiceCurrency;
    });
  }

  ngOnInit(): void {
    this.settingsService.load();
  }

  protected save(): void {
    this.statusMessage.set(null);
    this.settingsService.save(
      {
        nodeUrl: this.nodeUrl.trim(),
        authToken: this.authToken.trim(),
        timeoutSeconds: this.timeoutSeconds,
        defaultInvoiceCurrency: this.defaultInvoiceCurrency.trim()
      },
      () => this.statusMessage.set('Settings saved. New RPC calls will use the updated node configuration.')
    );
  }

  protected reset(): void {
    const active = this.settings();
    this.nodeUrl = active.nodeUrl;
    this.authToken = active.authToken;
    this.timeoutSeconds = active.timeoutSeconds;
    this.defaultInvoiceCurrency = active.defaultInvoiceCurrency;
    this.statusMessage.set('Form reset to the currently active runtime settings.');
  }
}
