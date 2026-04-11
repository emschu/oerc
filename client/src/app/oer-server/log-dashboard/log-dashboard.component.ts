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
import {Component, OnDestroy, OnInit} from '@angular/core';
import {ApiService} from '../api.service';
import {LogEntryResponse} from '../entities';
import {Subscription} from 'rxjs';
import {NavigationEnd, Router} from '@angular/router';
import {filter} from 'rxjs/operators';

@Component({
    selector: 'app-logdashboard',
    templateUrl: './log-dashboard.component.html',
    styleUrls: ['./log-dashboard.component.scss'],
    standalone: false
})
export class LogDashboardComponent implements OnInit, OnDestroy {
  logEntries: LogEntryResponse | null = null;
  public currentPage = 1;
  public totalPages = 0;
  public totalEntryCount = 0;

  private logEntrySubscription: Subscription | null = null;
  private routerEventsSubscription: Subscription | null = null;

  constructor(private apiService: ApiService, private router: Router) {
  }

  ngOnInit(): void {
    this.routerEventsSubscription = this.router.events.pipe(
      filter(event => event instanceof NavigationEnd)
    ).subscribe((event: any) => {
      const navigationEndEvent = event as NavigationEnd;
      if (navigationEndEvent.urlAfterRedirects === '/log') {
        this.currentPage = 1;
        this.loadLog();
      }
    });
    this.loadLog();
  }

  public goToNextPage(): void {
    this.currentPage++;
    this.loadLog();
  }

  public goToPreviousPage(): void {
    if (this.currentPage > 1) {
      this.currentPage--;
      this.loadLog();
    }
  }

  public goToPage(page: number): void {
    this.currentPage = page;
    this.loadLog();
  }

  private loadLog(): void {
    this.apiService.isLoadingSubject.next(true);
    this.logEntrySubscription = this.apiService.logEntries(this.currentPage - 1).subscribe(value => {
      this.logEntries = value;
      this.totalPages = value.page_count + 1;
      this.totalEntryCount = value.entry_count;
      setTimeout(() => {
        this.apiService.isLoadingSubject.next(false);
      }, 250);
    });
  }

  ngOnDestroy(): void {
    this.logEntrySubscription?.unsubscribe();
    this.routerEventsSubscription?.unsubscribe();
  }
}
