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
