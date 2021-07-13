/*
 * oerc, alias oer-collector
 * Copyright (C) 2021 emschu[aet]mailbox.org
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
import {Component, OnDestroy, OnInit} from '@angular/core';
import {Subscription} from 'rxjs';
import {OnPageHidden, OnPageVisible} from 'angular-page-visibility';

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

  @OnPageHidden()
  onPageHidden(): void {
    console.log('on page hidden');
    this.apiService.isWindowOpenedSubject.next(false);
  }

  @OnPageVisible()
  onPageVisible(): void {
    console.log('on page visible');
    this.apiService.isWindowOpenedSubject.next(true);
  }

  ngOnInit(): void {
    this.isLiveSubscription = this.apiService.isLiveSubject.subscribe(value => {
      if (this.inited && !this.isLive && value !== null) {
        window.location.reload();
      }
      if (!this.inited && value !== null) {
        this.inited = true;
      }
      if (value === true || value === false) {
        this.isLive = value;
      }
    });
  }

  ngOnDestroy(): void {
    this.isLiveSubscription?.unsubscribe();
  }
}
