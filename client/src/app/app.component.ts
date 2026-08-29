/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2026 emschu[aet]mailbox.org
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public
 * License along with this program.
 * If not, see <https://www.gnu.org/licenses/>.
 */
import {ApiService} from './oer-server/api.service';
import {ChangeDetectionStrategy, Component, effect, HostListener, inject, OnInit} from '@angular/core';
import {NavComponent} from './nav/nav.component';
import {RouterOutlet} from '@angular/router';
import {NgClass} from '@angular/common';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
  changeDetection: ChangeDetectionStrategy.Eager,
  standalone: true,
  imports: [
    NavComponent,
    RouterOutlet,
    NgClass
  ]
})
export class AppComponent implements OnInit {
  isLive = false;
  public apiService = inject(ApiService);
  private inited = false;

  constructor() {
    effect(() => {
      const value = this.apiService.isLive();
      if (!this.inited && value !== null) {
        this.inited = true;
      }
      if (value !== null) {
        this.isLive = value;
      }
    });
  }

  @HostListener('document:visibilitychange', ['$event'])
  onPageVisible(event: Event): void {
    if (document.hidden) {
      this.apiService.isWindowOpenedSubject.set(false);
    } else {
      this.apiService.isWindowOpenedSubject.set(true);
    }
  }

  ngOnInit(): void {
    this.apiService.init();
  }
}
