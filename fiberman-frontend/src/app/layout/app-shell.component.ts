import { CommonModule } from '@angular/common';
import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { FiberApiService } from '../core/fiber-api.service';

interface NavItem {
  label: string;
  icon: string;
  route: string;
  comingSoon?: boolean;
}

@Component({
  selector: 'app-shell',
  imports: [CommonModule, RouterLink, RouterLinkActive, RouterOutlet],
  templateUrl: './app-shell.component.html',
  styleUrl: './app-shell.component.css'
})
export class AppShellComponent implements OnInit {
  private readonly api = inject(FiberApiService);

  protected readonly navItems: NavItem[] = [
    { label: 'Dashboard', icon: '◫', route: '/dashboard' },
    { label: 'RPC Explorer', icon: '⚡', route: '/explorer' },
    { label: 'Invoice Builder', icon: '◫', route: '/invoice' },
    { label: 'Network Topology', icon: '⎔', route: '/topology', comingSoon: true },
    { label: 'Payments History', icon: '↺', route: '/payments' },
    { label: 'Logs', icon: '☰', route: '/logs', comingSoon: true },
    { label: 'Settings', icon: '⚙', route: '/settings' }
  ];

  protected readonly nodeOnline = signal(true);
  protected readonly comingSoonMessage = signal<string | null>(null);
  protected readonly nodeLabel = computed(() =>
    this.nodeOnline() ? 'Node Status: Connected' : 'Node Status: Offline'
  );

  ngOnInit(): void {
    this.api.getNodeInfo().subscribe({
      next: () => this.nodeOnline.set(true),
      error: () => this.nodeOnline.set(false)
    });
  }

  protected openComingSoon(item: NavItem): void {
    this.comingSoonMessage.set(`${item.label} is reserved for the Wails desktop release. The nav stays visible so the IA is clear, but the implementation is intentionally deferred for now.`);
  }

  protected closeComingSoon(): void {
    this.comingSoonMessage.set(null);
  }
}
