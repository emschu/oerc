import {Component, OnDestroy, OnInit} from '@angular/core';
import {ApiService} from '../api.service';
import {StatusResponse} from '../entities';
import {Subscription} from 'rxjs';

@Component({
  selector: 'app-oer-status-display',
  templateUrl: './status.component.html',
  styleUrls: ['./status.component.scss']
})
export class StatusComponent implements OnInit, OnDestroy {
  public currentStatus: StatusResponse | null = null;
  statusSubscription: Subscription | null = null;

  constructor(private oerApiService: ApiService) {}

  ngOnInit(): void {
    this.updateStatus();
  }

  ngOnDestroy(): void {
    this.statusSubscription?.unsubscribe();
  }

  updateStatus(): void {
    const statusResponseObservable = this.oerApiService.statusResponse();
    this.statusSubscription = statusResponseObservable.subscribe(statusResponse => {
      this.currentStatus = statusResponse;
    });
  }
}
