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
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ApiService} from '../api.service';
import {LogEntryResponse} from '../entities';
import {Subscription} from 'rxjs';

@Component({
  selector: 'app-logdashboard',
  templateUrl: './log-dashboard.component.html',
  styleUrls: ['./log-dashboard.component.scss']
})
export class LogDashboardComponent implements OnInit, OnDestroy {
  logEntries: LogEntryResponse | null = null;

  private logEntrySubscription: Subscription | null = null;

  constructor(private apiService: ApiService) {
  }

  ngOnInit(): void {
    this.loadLog();
  }

  private loadLog(): void {
    this.apiService.isLoadingSubject.next(true);
    this.logEntrySubscription = this.apiService.logEntries().subscribe(value => {
      this.logEntries = value;
      setTimeout(() => {
        this.apiService.isLoadingSubject.next(false);
      }, 250);
    });
  }

  ngOnDestroy(): void {
    this.logEntrySubscription?.unsubscribe();
  }
}
