/*
 * oerc, alias oer-collector
 * Copyright (C) 2021-2024 emschu[aet]mailbox.org
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
import {Component, HostListener, OnDestroy, OnInit} from '@angular/core';
import {Subscription} from 'rxjs';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit, OnDestroy {
  isLive = false;
  private isLiveSubscription: Subscription | null = null;
  private inited = false;

  constructor(public apiService: ApiService) {
  }

  @HostListener('document:visibilitychange', ['$event'])
  onPageVisible(): void {
    if (document.hidden) {
      this.apiService.isWindowOpenedSubject.next(false);
    } else {
      this.apiService.isWindowOpenedSubject.next(true);
    }
  }

  ngOnInit(): void {
    this.apiService.init();
    this.isLiveSubscription = this.apiService.isLiveSubject.subscribe(value => {
      if (!this.inited && value !== null) {
        this.inited = true;
      }
      if (value !== null) {
        this.isLive = value;
      }
    });
  }

  ngOnDestroy(): void {
    this.isLiveSubscription?.unsubscribe();
  }
}
