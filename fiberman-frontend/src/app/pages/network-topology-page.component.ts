import { CommonModule } from '@angular/common';
import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { forkJoin } from 'rxjs';
import { FiberApiService } from '../core/fiber-api.service';

@Component({
  selector: 'app-network-topology-page',
  imports: [CommonModule],
  templateUrl: './network-topology-page.component.html',
  styleUrl: './network-topology-page.component.css'
})
export class NetworkTopologyPageComponent implements OnInit {
  private readonly api = inject(FiberApiService);

  protected readonly peers = signal<Array<Record<string, unknown>>>([]);
  protected readonly channels = signal<Array<Record<string, unknown>>>([]);
  protected readonly nodes = computed(() =>
    (this.peers().length ? this.peers() : Array.from({ length: 12 }, (_, index) => ({ id: `peer-${index + 1}` }))).map(
      (peer, index, all) => {
        const peerRecord = peer as Record<string, unknown>;
        const angle = (Math.PI * 2 * index) / all.length;
        return {
          id: String(peerRecord['peer_id'] ?? peerRecord['id'] ?? `peer-${index + 1}`),
          left: `${50 + Math.cos(angle) * 34}%`,
          top: `${50 + Math.sin(angle) * 34}%`,
          color: `hsl(${index * 31} 82% 68%)`
        };
      }
    )
  );

  ngOnInit(): void {
    forkJoin({
      peers: this.api.listPeers(),
      channels: this.api.listChannels()
    }).subscribe({
      next: ({ peers, channels }) => {
        this.peers.set(this.asArray(peers.result));
        this.channels.set(this.asArray(channels.result));
      }
    });
  }

  private asArray(value: unknown): Array<Record<string, unknown>> {
    if (Array.isArray(value)) {
      return value as Array<Record<string, unknown>>;
    }
    if (value && typeof value === 'object') {
      const candidate = Object.values(value as Record<string, unknown>).find(Array.isArray);
      if (Array.isArray(candidate)) {
        return candidate as Array<Record<string, unknown>>;
      }
    }
    return [];
  }
}
