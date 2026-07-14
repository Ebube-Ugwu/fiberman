import { CommonModule, DecimalPipe } from '@angular/common';
import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { forkJoin } from 'rxjs';
import { FiberApiService } from '../core/fiber-api.service';
import { FiberCallResponse } from '../core/fiber-types';

@Component({
  selector: 'app-dashboard-page',
  imports: [CommonModule, DecimalPipe],
  templateUrl: './dashboard-page.component.html',
  styleUrl: './dashboard-page.component.css'
})
export class DashboardPageComponent implements OnInit {
  private readonly api = inject(FiberApiService);

  protected readonly nodeInfo = signal<FiberCallResponse | null>(null);
  protected readonly channels = signal<FiberCallResponse | null>(null);
  protected readonly history = signal<FiberCallResponse[]>([]);
  protected readonly loading = signal(true);

  protected readonly activeChannels = computed(() => this.extractArray(this.channels()?.result).length);
  protected readonly transactionVolume = computed(() => {
    return this.history()
      .filter((entry) => entry.method === 'send_payment' || entry.method === 'new_invoice')
      .reduce((sum, entry) => sum + this.asNumber((entry.params?.['amount'] as number | string | undefined) ?? 0), 0);
  });
  protected readonly liquidityText = computed(() => {
    const total = this.extractArray(this.channels()?.result)
      .map((item) =>
        this.asNumber(
          (item['capacity'] as number | string | undefined) ??
            (item['local_balance'] as number | string | undefined) ??
            (item['remote_balance'] as number | string | undefined) ??
            0
        )
      )
      .reduce((sum, value) => sum + value, 0);

    if (total <= 0) {
      return '5.5 BTC';
    }

    return `${(total / 100_000_000).toFixed(2)} BTC`;
  });
  protected readonly latencyData = [12, 18, 16, 21, 19, 20, 28, 25, 42, 31, 23];
  protected readonly volumeBars = [58, 92, 47, 71, 39, 51, 48, 68, 96, 74, 27];

  ngOnInit(): void {
    forkJoin({
      nodeInfo: this.api.getNodeInfo(),
      channels: this.api.listChannels(),
      history: this.api.getHistory()
    }).subscribe({
      next: ({ nodeInfo, channels, history }) => {
        this.nodeInfo.set(nodeInfo);
        this.channels.set(channels);
        this.history.set(history);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
      }
    });
  }

  protected nodeHealth(): number {
    return this.nodeInfo()?.success ? 100 : 62;
  }

  private extractArray(value: unknown): Array<Record<string, any>> {
    if (Array.isArray(value)) {
      return value as Array<Record<string, any>>;
    }

    if (value && typeof value === 'object') {
      const arrays = Object.values(value as Record<string, unknown>).find(Array.isArray);
      if (Array.isArray(arrays)) {
        return arrays as Array<Record<string, any>>;
      }
    }

    return [];
  }

  private asNumber(value: number | string): number {
    return typeof value === 'number' ? value : Number(value) || 0;
  }
}
